package game

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oapi-codegen/nullable"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/api"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/taskqueue"
)

type playerClientState int

const (
	playerClientStateWaitingEntries playerClientState = iota
	playerClientStateEntried
	playerClientStateReady
)

type gameHub struct {
	ctx               context.Context
	game              *game
	q                 *db.Queries
	taskQueue         *taskqueue.Queue
	players           map[*playerClient]playerClientState
	registerPlayer    chan *playerClient
	unregisterPlayer  chan *playerClient
	playerC2SMessages chan *playerMessageC2SWithClient
	watchers          map[*watcherClient]bool
	registerWatcher   chan *watcherClient
	unregisterWatcher chan *watcherClient
	taskResults       chan taskqueue.TaskResult
}

func newGameHub(ctx context.Context, game *game, q *db.Queries, taskQueue *taskqueue.Queue) *gameHub {
	return &gameHub{
		ctx:               ctx,
		game:              game,
		q:                 q,
		taskQueue:         taskQueue,
		players:           make(map[*playerClient]playerClientState),
		registerPlayer:    make(chan *playerClient),
		unregisterPlayer:  make(chan *playerClient),
		playerC2SMessages: make(chan *playerMessageC2SWithClient),
		watchers:          make(map[*watcherClient]bool),
		registerWatcher:   make(chan *watcherClient),
		unregisterWatcher: make(chan *watcherClient),
		taskResults:       make(chan taskqueue.TaskResult),
	}
}

func (hub *gameHub) run() {
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case player := <-hub.registerPlayer:
			hub.players[player] = playerClientStateWaitingEntries
		case player := <-hub.unregisterPlayer:
			if _, ok := hub.players[player]; ok {
				hub.closePlayerClient(player)
			}
		case watcher := <-hub.registerWatcher:
			hub.watchers[watcher] = true
		case watcher := <-hub.unregisterWatcher:
			if _, ok := hub.watchers[watcher]; ok {
				hub.closeWatcherClient(watcher)
			}
		case message := <-hub.playerC2SMessages:
			switch msg := message.message.(type) {
			case *playerMessageC2SEntry:
				log.Printf("entry: %v", message.message)
				// TODO: assert state is waiting_entries
				hub.players[message.client] = playerClientStateEntried
				entriedPlayerCount := 0
				for _, state := range hub.players {
					if playerClientStateEntried <= state {
						entriedPlayerCount++
					}
				}
				if entriedPlayerCount == hub.game.playerCount {
					err := hub.q.UpdateGameState(hub.ctx, db.UpdateGameStateParams{
						GameID: int32(hub.game.gameID),
						State:  string(gameStateWaitingStart),
					})
					if err != nil {
						log.Fatalf("failed to set game state: %v", err)
					}
					hub.game.state = gameStateWaitingStart
				}
			case *playerMessageC2SReady:
				log.Printf("ready: %v", message.message)
				// TODO: assert state is prepare
				hub.players[message.client] = playerClientStateReady
				readyPlayerCount := 0
				for _, state := range hub.players {
					if playerClientStateReady <= state {
						readyPlayerCount++
					}
				}
				if readyPlayerCount == hub.game.playerCount {
					startAt := time.Now().Add(11 * time.Second).UTC()
					for player := range hub.players {
						player.s2cMessages <- &playerMessageS2CStart{
							Type: playerMessageTypeS2CStart,
							Data: playerMessageS2CStartPayload{
								StartAt: int(startAt.Unix()),
							},
						}
					}
					hub.broadcastToWatchers(&watcherMessageS2CStart{
						Type: watcherMessageTypeS2CStart,
						Data: watcherMessageS2CStartPayload{
							StartAt: int(startAt.Unix()),
						},
					})
					err := hub.q.UpdateGameStartedAt(hub.ctx, db.UpdateGameStartedAtParams{
						GameID: int32(hub.game.gameID),
						StartedAt: pgtype.Timestamp{
							Time:             startAt,
							InfinityModifier: pgtype.Finite,
							Valid:            true,
						},
					})
					if err != nil {
						log.Fatalf("failed to set game state: %v", err)
					}
					hub.game.startedAt = &startAt
					err = hub.q.UpdateGameState(hub.ctx, db.UpdateGameStateParams{
						GameID: int32(hub.game.gameID),
						State:  string(gameStateStarting),
					})
					if err != nil {
						log.Fatalf("failed to set game state: %v", err)
					}
					hub.game.state = gameStateStarting
				}
			case *playerMessageC2SCode:
				// TODO: assert game state is gaming
				log.Printf("code: %v", message.message)
				code := msg.Data.Code
				hub.broadcastToWatchers(&watcherMessageS2CCode{
					Type: watcherMessageTypeS2CCode,
					Data: watcherMessageS2CCodePayload{
						PlayerID: message.client.playerID,
						Code:     code,
					},
				})
			case *playerMessageC2SSubmit:
				// TODO: assert game state is gaming
				log.Printf("submit: %v", message.message)
				code := msg.Data.Code
				codeSize := len(code) // TODO: exclude whitespaces.
				codeHash := calcHash(code)
				if err := hub.taskQueue.EnqueueTaskCreateSubmissionRecord(
					hub.game.gameID,
					message.client.playerID,
					code,
					codeSize,
					taskqueue.MD5HexHash(codeHash),
				); err != nil {
					// TODO: notify failure to player
					log.Fatalf("failed to enqueue task: %v", err)
				}
				hub.broadcastToWatchers(&watcherMessageS2CSubmit{
					Type: watcherMessageTypeS2CSubmit,
					Data: watcherMessageS2CSubmitPayload{
						PlayerID:         message.client.playerID,
						PreliminaryScore: codeSize,
					},
				})
			default:
				log.Printf("unexpected message type: %T", message.message)
			}
		case <-ticker.C:
			if hub.game.state == gameStateStarting {
				if time.Now().After(*hub.game.startedAt) {
					err := hub.q.UpdateGameState(hub.ctx, db.UpdateGameStateParams{
						GameID: int32(hub.game.gameID),
						State:  string(gameStateGaming),
					})
					if err != nil {
						log.Fatalf("failed to set game state: %v", err)
					}
					hub.game.state = gameStateGaming
				}
			} else if hub.game.state == gameStateGaming {
				if time.Now().After(hub.game.startedAt.Add(time.Duration(hub.game.durationSeconds) * time.Second)) {
					err := hub.q.UpdateGameState(hub.ctx, db.UpdateGameStateParams{
						GameID: int32(hub.game.gameID),
						State:  string(gameStateFinished),
					})
					if err != nil {
						log.Fatalf("failed to set game state: %v", err)
					}
					hub.game.state = gameStateFinished
					hub.close()
					return
				}
			}
		}
	}
}

func (hub *gameHub) broadcastToWatchers(msg watcherMessageS2C) {
	for watcher := range hub.watchers {
		watcher.s2cMessages <- msg
	}
}

type codeSubmissionError struct {
	Status string
	Stdout string
	Stderr string
}

func (err *codeSubmissionError) Error() string {
	return err.Stderr
}

func (hub *gameHub) processTaskResults() {
	for taskResult := range hub.taskResults {
		switch taskResult := taskResult.(type) {
		case *taskqueue.TaskResultCreateSubmissionRecord:
			err := hub.processTaskResultCreateSubmissionRecord(taskResult)
			if err != nil {
				for player := range hub.players {
					if player.playerID != taskResult.TaskPayload.UserID() {
						continue
					}
					player.s2cMessages <- &playerMessageS2CExecResult{
						Type: playerMessageTypeS2CExecResult,
						Data: playerMessageS2CExecResultPayload{
							Score:  nil,
							Status: api.GamePlayerMessageS2CExecResultPayloadStatus(err.Status),
						},
					}
				}
				hub.broadcastToWatchers(&watcherMessageS2CSubmitResult{
					Type: watcherMessageTypeS2CSubmitResult,
					Data: watcherMessageS2CSubmitResultPayload{
						PlayerID: taskResult.TaskPayload.UserID(),
						Status:   api.GameWatcherMessageS2CSubmitResultPayloadStatus(err.Status),
					},
				})
			}
		case *taskqueue.TaskResultCompileSwiftToWasm:
			err := hub.processTaskResultCompileSwiftToWasm(taskResult)
			if err != nil {
				for player := range hub.players {
					if player.playerID != taskResult.TaskPayload.UserID() {
						continue
					}
					player.s2cMessages <- &playerMessageS2CExecResult{
						Type: playerMessageTypeS2CExecResult,
						Data: playerMessageS2CExecResultPayload{
							Score:  nil,
							Status: api.GamePlayerMessageS2CExecResultPayloadStatus(err.Status),
						},
					}
				}
				hub.broadcastToWatchers(&watcherMessageS2CExecResult{
					Type: watcherMessageTypeS2CExecResult,
					Data: watcherMessageS2CExecResultPayload{
						PlayerID: taskResult.TaskPayload.UserID(),
						Status:   api.GameWatcherMessageS2CExecResultPayloadStatus(err.Status),
						Stdout:   err.Stdout,
						Stderr:   err.Stderr,
					},
				})
				hub.broadcastToWatchers(&watcherMessageS2CSubmitResult{
					Type: watcherMessageTypeS2CSubmitResult,
					Data: watcherMessageS2CSubmitResultPayload{
						PlayerID: taskResult.TaskPayload.UserID(),
						Status:   api.GameWatcherMessageS2CSubmitResultPayloadStatus(err.Status),
					},
				})
			}
		case *taskqueue.TaskResultCompileWasmToNativeExecutable:
			err := hub.processTaskResultCompileWasmToNativeExecutable(taskResult)
			if err != nil {
				for player := range hub.players {
					if player.playerID != taskResult.TaskPayload.UserID() {
						continue
					}
					player.s2cMessages <- &playerMessageS2CExecResult{
						Type: playerMessageTypeS2CExecResult,
						Data: playerMessageS2CExecResultPayload{
							Score:  nil,
							Status: api.GamePlayerMessageS2CExecResultPayloadStatus(err.Status),
						},
					}
				}
				hub.broadcastToWatchers(&watcherMessageS2CExecResult{
					Type: watcherMessageTypeS2CExecResult,
					Data: watcherMessageS2CExecResultPayload{
						PlayerID: taskResult.TaskPayload.UserID(),
						Status:   api.GameWatcherMessageS2CExecResultPayloadStatus(err.Status),
						Stdout:   err.Stdout,
						Stderr:   err.Stderr,
					},
				})
				hub.broadcastToWatchers(&watcherMessageS2CSubmitResult{
					Type: watcherMessageTypeS2CSubmitResult,
					Data: watcherMessageS2CSubmitResultPayload{
						PlayerID: taskResult.TaskPayload.UserID(),
						Status:   api.GameWatcherMessageS2CSubmitResultPayloadStatus(err.Status),
					},
				})
			} else {
				hub.broadcastToWatchers(&watcherMessageS2CExecResult{
					Type: watcherMessageTypeS2CExecResult,
					Data: watcherMessageS2CExecResultPayload{
						PlayerID: taskResult.TaskPayload.UserID(),
						Status:   api.GameWatcherMessageS2CExecResultPayloadStatus("success"),
						// TODO: inherit the command stdout/stderr.
						Stdout: "Successfully compiled",
						Stderr: "",
					},
				})
			}
		case *taskqueue.TaskResultRunTestcase:
			// FIXME: error handling
			var err error
			err1 := hub.processTaskResultRunTestcase(taskResult)
			_ = err // TODO: handle err?
			aggregatedStatus, err := hub.q.AggregateTestcaseResults(hub.ctx, int32(taskResult.TaskPayload.SubmissionID))
			_ = err // TODO: handle err?
			err = hub.q.CreateSubmissionResult(hub.ctx, db.CreateSubmissionResultParams{
				SubmissionID: int32(taskResult.TaskPayload.SubmissionID),
				Status:       aggregatedStatus,
				Stdout:       "",
				Stderr:       "",
			})
			if err != nil {
				for player := range hub.players {
					if player.playerID != taskResult.TaskPayload.UserID() {
						continue
					}
					player.s2cMessages <- &playerMessageS2CExecResult{
						Type: playerMessageTypeS2CExecResult,
						Data: playerMessageS2CExecResultPayload{
							Score:  nil,
							Status: api.GamePlayerMessageS2CExecResultPayloadStatus("internal_error"),
						},
					}
				}
				hub.broadcastToWatchers(&watcherMessageS2CExecResult{
					Type: watcherMessageTypeS2CExecResult,
					Data: watcherMessageS2CExecResultPayload{
						PlayerID:   taskResult.TaskPayload.UserID(),
						TestcaseID: nullable.NewNullableWithValue(int(taskResult.TaskPayload.TestcaseID)),
						Status:     api.GameWatcherMessageS2CExecResultPayloadStatus("internal_error"),
						// TODO: inherit the command stdout/stderr?
						Stdout: "",
						Stderr: "internal error",
					},
				})
				hub.broadcastToWatchers(&watcherMessageS2CSubmitResult{
					Type: watcherMessageTypeS2CSubmitResult,
					Data: watcherMessageS2CSubmitResultPayload{
						PlayerID: taskResult.TaskPayload.UserID(),
						Status:   api.GameWatcherMessageS2CSubmitResultPayloadStatus("internal_error"),
					},
				})
				continue
			}
			for player := range hub.players {
				if player.playerID != taskResult.TaskPayload.UserID() {
					continue
				}
				player.s2cMessages <- &playerMessageS2CExecResult{
					Type: playerMessageTypeS2CExecResult,
					Data: playerMessageS2CExecResultPayload{
						Score:  nil,
						Status: api.GamePlayerMessageS2CExecResultPayloadStatus(aggregatedStatus),
					},
				}
			}
			if err1 != nil {
				hub.broadcastToWatchers(&watcherMessageS2CExecResult{
					Type: watcherMessageTypeS2CExecResult,
					Data: watcherMessageS2CExecResultPayload{
						PlayerID:   taskResult.TaskPayload.UserID(),
						TestcaseID: nullable.NewNullableWithValue(int(taskResult.TaskPayload.TestcaseID)),
						Status:     api.GameWatcherMessageS2CExecResultPayloadStatus(err1.Status),
						Stdout:     err1.Stdout,
						Stderr:     err1.Stderr,
					},
				})
			} else {
				hub.broadcastToWatchers(&watcherMessageS2CExecResult{
					Type: watcherMessageTypeS2CExecResult,
					Data: watcherMessageS2CExecResultPayload{
						PlayerID:   taskResult.TaskPayload.UserID(),
						TestcaseID: nullable.NewNullableWithValue(int(taskResult.TaskPayload.TestcaseID)),
						Status:     api.GameWatcherMessageS2CExecResultPayloadStatus("success"),
						// TODO: inherit the command stdout/stderr?
						Stdout: "Testcase passed",
						Stderr: "",
					},
				})
			}
			if aggregatedStatus != "running" {
				hub.broadcastToWatchers(&watcherMessageS2CSubmitResult{
					Type: watcherMessageTypeS2CSubmitResult,
					Data: watcherMessageS2CSubmitResultPayload{
						PlayerID: taskResult.TaskPayload.UserID(),
						Status:   api.GameWatcherMessageS2CSubmitResultPayloadStatus(aggregatedStatus),
					},
				})
			}
		default:
			panic("unexpected task result type")
		}
	}
}

func (hub *gameHub) processTaskResultCreateSubmissionRecord(
	taskResult *taskqueue.TaskResultCreateSubmissionRecord,
) *codeSubmissionError {
	if taskResult.Err != nil {
		return &codeSubmissionError{
			Status: "internal_error",
			Stderr: taskResult.Err.Error(),
		}
	}

	if err := hub.taskQueue.EnqueueTaskCompileSwiftToWasm(
		taskResult.TaskPayload.GameID(),
		taskResult.TaskPayload.UserID(),
		taskResult.TaskPayload.Code,
		taskResult.TaskPayload.CodeHash(),
		taskResult.SubmissionID,
	); err != nil {
		return &codeSubmissionError{
			Status: "internal_error",
			Stderr: err.Error(),
		}
	}
	return nil
}

func (hub *gameHub) processTaskResultCompileSwiftToWasm(
	taskResult *taskqueue.TaskResultCompileSwiftToWasm,
) *codeSubmissionError {
	if taskResult.Err != nil {
		return &codeSubmissionError{
			Status: "internal_error",
			Stderr: taskResult.Err.Error(),
		}
	}

	if taskResult.Status != "success" {
		if err := hub.q.CreateSubmissionResult(hub.ctx, db.CreateSubmissionResultParams{
			SubmissionID: int32(taskResult.TaskPayload.SubmissionID),
			Status:       taskResult.Status,
			Stdout:       taskResult.Stdout,
			Stderr:       taskResult.Stderr,
		}); err != nil {
			return &codeSubmissionError{
				Status: "internal_error",
				Stderr: err.Error(),
			}
		}
		return &codeSubmissionError{
			Status: taskResult.Status,
			Stdout: taskResult.Stdout,
			Stderr: taskResult.Stderr,
		}
	}
	if err := hub.taskQueue.EnqueueTaskCompileWasmToNativeExecutable(
		taskResult.TaskPayload.GameID(),
		taskResult.TaskPayload.UserID(),
		taskResult.TaskPayload.CodeHash(),
		taskResult.TaskPayload.SubmissionID,
	); err != nil {
		return &codeSubmissionError{
			Status: "internal_error",
			Stderr: err.Error(),
		}
	}
	return nil
}

func (hub *gameHub) processTaskResultCompileWasmToNativeExecutable(
	taskResult *taskqueue.TaskResultCompileWasmToNativeExecutable,
) *codeSubmissionError {
	if taskResult.Err != nil {
		return &codeSubmissionError{
			Status: "internal_error",
			Stderr: taskResult.Err.Error(),
		}
	}

	if taskResult.Status != "success" {
		if err := hub.q.CreateSubmissionResult(hub.ctx, db.CreateSubmissionResultParams{
			SubmissionID: int32(taskResult.TaskPayload.SubmissionID),
			Status:       taskResult.Status,
			Stdout:       taskResult.Stdout,
			Stderr:       taskResult.Stderr,
		}); err != nil {
			return &codeSubmissionError{
				Status: "internal_error",
				Stderr: err.Error(),
			}
		}
		return &codeSubmissionError{
			Status: taskResult.Status,
			Stdout: taskResult.Stdout,
			Stderr: taskResult.Stderr,
		}
	}

	testcases, err := hub.q.ListTestcasesByGameID(hub.ctx, int32(taskResult.TaskPayload.GameID()))
	if err != nil {
		return &codeSubmissionError{
			Status: "internal_error",
			Stderr: err.Error(),
		}
	}
	if len(testcases) == 0 {
		return &codeSubmissionError{
			Status: "internal_error",
			Stderr: "no testcases found",
		}
	}

	for _, testcase := range testcases {
		if err := hub.taskQueue.EnqueueTaskRunTestcase(
			taskResult.TaskPayload.GameID(),
			taskResult.TaskPayload.UserID(),
			taskResult.TaskPayload.CodeHash(),
			taskResult.TaskPayload.SubmissionID,
			int(testcase.TestcaseID),
			testcase.Stdin,
			testcase.Stdout,
		); err != nil {
			return &codeSubmissionError{
				Status: "internal_error",
				Stderr: err.Error(),
			}
		}
	}
	return nil
}

func (hub *gameHub) processTaskResultRunTestcase(
	taskResult *taskqueue.TaskResultRunTestcase,
) *codeSubmissionError {
	if taskResult.Err != nil {
		return &codeSubmissionError{
			Status: "internal_error",
			Stderr: taskResult.Err.Error(),
		}
	}

	if taskResult.Status != "success" {
		if err := hub.q.CreateTestcaseResult(hub.ctx, db.CreateTestcaseResultParams{
			SubmissionID: int32(taskResult.TaskPayload.SubmissionID),
			TestcaseID:   int32(taskResult.TaskPayload.TestcaseID),
			Status:       taskResult.Status,
			Stdout:       taskResult.Stdout,
			Stderr:       taskResult.Stderr,
		}); err != nil {
			return &codeSubmissionError{
				Status: "internal_error",
				Stderr: err.Error(),
			}
		}
		return &codeSubmissionError{
			Status: taskResult.Status,
			Stdout: taskResult.Stdout,
			Stderr: taskResult.Stderr,
		}
	}
	if !isTestcaseResultCorrect(taskResult.TaskPayload.Stdout, taskResult.Stdout) {
		if err := hub.q.CreateTestcaseResult(hub.ctx, db.CreateTestcaseResultParams{
			SubmissionID: int32(taskResult.TaskPayload.SubmissionID),
			TestcaseID:   int32(taskResult.TaskPayload.TestcaseID),
			Status:       "wrong_answer",
			Stdout:       taskResult.Stdout,
			Stderr:       taskResult.Stderr,
		}); err != nil {
			return &codeSubmissionError{
				Status: "internal_error",
				Stderr: err.Error(),
			}
		}
		return &codeSubmissionError{
			Status: "wrong_answer",
			Stdout: taskResult.Stdout,
			Stderr: taskResult.Stderr,
		}
	}
	return nil
}

func (hub *gameHub) startGame() error {
	for player := range hub.players {
		player.s2cMessages <- &playerMessageS2CPrepare{
			Type: playerMessageTypeS2CPrepare,
		}
	}

	err := hub.q.UpdateGameState(hub.ctx, db.UpdateGameStateParams{
		GameID: int32(hub.game.gameID),
		State:  string(gameStatePrepare),
	})
	if err != nil {
		return err
	}
	hub.game.state = gameStatePrepare
	return nil
}

func (hub *gameHub) close() {
	for client := range hub.players {
		hub.closePlayerClient(client)
	}
	close(hub.registerPlayer)
	close(hub.unregisterPlayer)
	close(hub.playerC2SMessages)
	for watcher := range hub.watchers {
		hub.closeWatcherClient(watcher)
	}
	close(hub.registerWatcher)
	close(hub.unregisterWatcher)
}

func (hub *gameHub) closePlayerClient(player *playerClient) {
	delete(hub.players, player)
	close(player.s2cMessages)
}

func (hub *gameHub) closeWatcherClient(watcher *watcherClient) {
	delete(hub.watchers, watcher)
	close(watcher.s2cMessages)
}

type Hubs struct {
	hubs        map[int]*gameHub
	q           *db.Queries
	taskQueue   *taskqueue.Queue
	taskResults chan taskqueue.TaskResult
}

func NewGameHubs(q *db.Queries, taskQueue *taskqueue.Queue, taskResults chan taskqueue.TaskResult) *Hubs {
	return &Hubs{
		hubs:        make(map[int]*gameHub),
		q:           q,
		taskQueue:   taskQueue,
		taskResults: taskResults,
	}
}

func (hubs *Hubs) Close() {
	log.Println("closing all game hubs")
	for _, hub := range hubs.hubs {
		hub.close()
	}
}

func (hubs *Hubs) getHub(gameID int) *gameHub {
	return hubs.hubs[gameID]
}

func (hubs *Hubs) RestoreFromDB(ctx context.Context) error {
	games, err := hubs.q.ListGames(ctx)
	if err != nil {
		return err
	}
	for _, row := range games {
		var startedAt *time.Time
		if row.StartedAt.Valid {
			startedAt = &row.StartedAt.Time
		}
		pr := &problem{
			problemID:   int(row.ProblemID),
			title:       row.Title,
			description: row.Description,
		}
		// TODO: N+1
		playerRows, err := hubs.q.ListGamePlayers(ctx, int32(row.GameID))
		if err != nil {
			return err
		}
		hubs.hubs[int(row.GameID)] = newGameHub(ctx, &game{
			gameID:          int(row.GameID),
			gameType:        gameType(row.GameType),
			durationSeconds: int(row.DurationSeconds),
			state:           gameState(row.State),
			displayName:     row.DisplayName,
			startedAt:       startedAt,
			problem:         pr,
			playerCount:     len(playerRows),
		}, hubs.q, hubs.taskQueue)
	}
	return nil
}

func (hubs *Hubs) Run() {
	for _, hub := range hubs.hubs {
		go hub.run()
		go hub.processTaskResults()
	}

	for taskResult := range hubs.taskResults {
		hub := hubs.getHub(taskResult.GameID())
		if hub == nil {
			log.Printf("no such game: %d", taskResult.GameID())
			continue
		}
		hub.taskResults <- taskResult
	}
}

func (hubs *Hubs) SockHandler() *SockHandler {
	return newSockHandler(hubs)
}

func (hubs *Hubs) StartGame(gameID int) error {
	hub := hubs.getHub(gameID)
	if hub == nil {
		return errors.New("no such game")
	}
	return hub.startGame()
}

func isTestcaseResultCorrect(expectedStdout, actualStdout string) bool {
	expectedStdout = strings.TrimSpace(expectedStdout)
	actualStdout = strings.TrimSpace(actualStdout)
	return actualStdout == expectedStdout
}

func calcHash(code string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(code)))
}
