package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/nullable"

	"github.com/nsfisis/phperkaigi-2025-albatross/backend/auth"
	"github.com/nsfisis/phperkaigi-2025-albatross/backend/db"
)

type Handler struct {
	q   *db.Queries
	hub GameHubInterface
}

type GameHubInterface interface {
	CalcCodeSize(code string) int
	EnqueueTestTasks(ctx context.Context, submissionID, gameID, userID int, code string) error
}

func (h *Handler) PostLogin(ctx context.Context, request PostLoginRequestObject) (PostLoginResponseObject, error) {
	username := request.Body.Username
	password := request.Body.Password
	userID, err := auth.Login(ctx, h.q, username, password)
	if err != nil {
		log.Printf("login failed: %v", err)
		var msg string
		if errors.Is(err, auth.ErrForteeLoginTimeout) {
			msg = "ログインに失敗しました"
		} else {
			msg = "ユーザー名またはパスワードが誤っています"
		}
		return PostLogin401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: msg,
			},
		}, nil
	}

	user, err := h.q.GetUserByID(ctx, int32(userID))
	if err != nil {
		return PostLogin401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "ログインに失敗しました",
			},
		}, nil
	}

	jwt, err := auth.NewJWT(&user)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return PostLogin200JSONResponse{
		Token: jwt,
	}, nil
}

func (h *Handler) GetGames(ctx context.Context, _ GetGamesRequestObject, _ *auth.JWTClaims) (GetGamesResponseObject, error) {
	gameRows, err := h.q.ListPublicGames(ctx)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	games := make([]Game, len(gameRows))
	gameIDs := make([]int32, len(gameRows))
	gameID2Index := make(map[int32]int, len(gameRows))
	for i, row := range gameRows {
		var startedAt *int64
		if row.StartedAt.Valid {
			startedAtTimestamp := row.StartedAt.Time.Unix()
			startedAt = &startedAtTimestamp
		}
		games[i] = Game{
			GameID:          int(row.GameID),
			GameType:        GameGameType(row.GameType),
			IsPublic:        row.IsPublic,
			DisplayName:     row.DisplayName,
			DurationSeconds: int(row.DurationSeconds),
			StartedAt:       startedAt,
			Problem: Problem{
				ProblemID:   int(row.ProblemID),
				Title:       row.Title,
				Description: row.Description,
				SampleCode:  row.SampleCode,
			},
		}
		gameIDs[i] = row.GameID
		gameID2Index[row.GameID] = i
	}
	mainPlayerRows, err := h.q.ListMainPlayers(ctx, gameIDs)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, row := range mainPlayerRows {
		idx := gameID2Index[row.GameID]
		game := &games[idx]
		game.MainPlayers = append(game.MainPlayers, User{
			UserID:      int(row.UserID),
			Username:    row.Username,
			DisplayName: row.DisplayName,
			IconPath:    row.IconPath,
			IsAdmin:     row.IsAdmin,
		})
	}
	return GetGames200JSONResponse{
		Games: games,
	}, nil
}

func (h *Handler) GetGame(ctx context.Context, request GetGameRequestObject, _ *auth.JWTClaims) (GetGameResponseObject, error) {
	gameID := request.GameID
	row, err := h.q.GetGameByID(ctx, int32(gameID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetGame404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{
					Message: "Game not found",
				},
			}, nil
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if !row.IsPublic {
		return GetGame404JSONResponse{
			NotFoundJSONResponse: NotFoundJSONResponse{
				Message: "Game not found",
			},
		}, nil
	}
	var startedAt *int64
	if row.StartedAt.Valid {
		startedAtTimestamp := row.StartedAt.Time.Unix()
		startedAt = &startedAtTimestamp
	}
	mainPlayerRows, err := h.q.ListMainPlayers(ctx, []int32{int32(gameID)})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	mainPlayers := make([]User, len(mainPlayerRows))
	for i, playerRow := range mainPlayerRows {
		mainPlayers[i] = User{
			UserID:      int(playerRow.UserID),
			Username:    playerRow.Username,
			DisplayName: playerRow.DisplayName,
			IconPath:    playerRow.IconPath,
			IsAdmin:     playerRow.IsAdmin,
		}
	}
	game := Game{
		GameID:          int(row.GameID),
		GameType:        GameGameType(row.GameType),
		IsPublic:        row.IsPublic,
		DisplayName:     row.DisplayName,
		DurationSeconds: int(row.DurationSeconds),
		StartedAt:       startedAt,
		Problem: Problem{
			ProblemID:   int(row.ProblemID),
			Title:       row.Title,
			Description: row.Description,
			SampleCode:  row.SampleCode,
		},
		MainPlayers: mainPlayers,
	}
	return GetGame200JSONResponse{
		Game: game,
	}, nil
}

func (h *Handler) GetGamePlayLatestState(ctx context.Context, request GetGamePlayLatestStateRequestObject, user *auth.JWTClaims) (GetGamePlayLatestStateResponseObject, error) {
	gameID := request.GameID
	userID := user.UserID
	row, err := h.q.GetLatestState(ctx, db.GetLatestStateParams{
		GameID: int32(gameID),
		UserID: int32(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetGamePlayLatestState200JSONResponse{
				State: LatestGameState{
					Code:   "",
					Score:  nullable.NewNullNullable[int](),
					Status: None,
				},
			}, nil
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	var score nullable.Nullable[int]
	if row.CodeSize != nil {
		score = nullable.NewNullableWithValue(int(*row.CodeSize))
	} else {
		score = nullable.NewNullNullable[int]()
	}
	return GetGamePlayLatestState200JSONResponse{
		State: LatestGameState{
			Code:   row.Code,
			Score:  score,
			Status: ExecutionStatus(row.Status),
		},
	}, nil
}

func (h *Handler) GetGameWatchLatestStates(ctx context.Context, request GetGameWatchLatestStatesRequestObject, user *auth.JWTClaims) (GetGameWatchLatestStatesResponseObject, error) {
	gameID := request.GameID
	rows, err := h.q.GetLatestStatesOfMainPlayers(ctx, int32(gameID))
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	states := make(map[string]LatestGameState, len(rows))
	for _, row := range rows {
		var code string
		if row.Code != nil {
			code = *row.Code
		}
		var score nullable.Nullable[int]
		if row.CodeSize != nil {
			score = nullable.NewNullableWithValue(int(*row.CodeSize))
		} else {
			score = nullable.NewNullNullable[int]()
		}
		var status ExecutionStatus
		if row.Status != nil {
			status = ExecutionStatus(*row.Status)
		} else {
			status = None
		}
		states[string(row.UserID)] = LatestGameState{
			Code:   code,
			Score:  score,
			Status: status,
		}

		if int(row.UserID) == user.UserID {
			return GetGameWatchLatestStates403JSONResponse{
				ForbiddenJSONResponse: ForbiddenJSONResponse{
					Message: "You are one of the main players of this game",
				},
			}, nil
		}
	}
	return GetGameWatchLatestStates200JSONResponse{
		States: states,
	}, nil
}

func (h *Handler) GetGameWatchRanking(ctx context.Context, request GetGameWatchRankingRequestObject, _ *auth.JWTClaims) (GetGameWatchRankingResponseObject, error) {
	gameID := request.GameID
	rows, err := h.q.GetRanking(ctx, int32(gameID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetGameWatchRanking200JSONResponse{}, nil
		}
	}
	ranking := make([]RankingEntry, len(rows))
	for i, row := range rows {
		ranking[i] = RankingEntry{
			Player: User{
				UserID:      int(row.UserID),
				Username:    row.Username,
				DisplayName: row.DisplayName,
				IconPath:    row.IconPath,
				IsAdmin:     row.IsAdmin,
			},
			Score: int(row.CodeSize),
		}
	}
	return GetGameWatchRanking200JSONResponse{
		Ranking: ranking,
	}, nil
}

func (h *Handler) PostGamePlayCode(ctx context.Context, request PostGamePlayCodeRequestObject, user *auth.JWTClaims) (PostGamePlayCodeResponseObject, error) {
	gameID := request.GameID
	userID := user.UserID
	err := h.q.UpdateCode(ctx, db.UpdateCodeParams{
		GameID: int32(gameID),
		UserID: int32(userID),
		Code:   request.Body.Code,
		Status: "none",
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return PostGamePlayCode200Response{}, nil
}

func (h *Handler) PostGamePlaySubmit(ctx context.Context, request PostGamePlaySubmitRequestObject, user *auth.JWTClaims) (PostGamePlaySubmitResponseObject, error) {
	gameID := request.GameID
	userID := user.UserID
	code := request.Body.Code
	codeSize := h.hub.CalcCodeSize(code)
	// TODO: transaction
	err := h.q.UpdateCode(ctx, db.UpdateCodeParams{
		GameID: int32(gameID),
		UserID: int32(userID),
		Code:   code,
		Status: "running",
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	submissionID, err := h.q.CreateSubmission(ctx, db.CreateSubmissionParams{
		GameID:   int32(gameID),
		UserID:   int32(userID),
		Code:     code,
		CodeSize: int32(codeSize),
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = h.hub.EnqueueTestTasks(ctx, int(submissionID), gameID, userID, code)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return PostGamePlaySubmit200Response{}, nil
}
