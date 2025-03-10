// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const aggregateTestcaseResults = `-- name: AggregateTestcaseResults :one
SELECT
    CASE
        WHEN COUNT(CASE WHEN r.status IS NULL            THEN 1 END) > 0 THEN 'running'
        WHEN COUNT(CASE WHEN r.status = 'internal_error' THEN 1 END) > 0 THEN 'internal_error'
        WHEN COUNT(CASE WHEN r.status = 'timeout'        THEN 1 END) > 0 THEN 'timeout'
        WHEN COUNT(CASE WHEN r.status = 'runtime_error'  THEN 1 END) > 0 THEN 'runtime_error'
        WHEN COUNT(CASE WHEN r.status = 'wrong_answer'   THEN 1 END) > 0 THEN 'wrong_answer'
        ELSE 'success'
    END AS status
FROM testcases
LEFT JOIN testcase_results AS r ON testcases.testcase_id = r.testcase_id
WHERE r.submission_id = $1
`

func (q *Queries) AggregateTestcaseResults(ctx context.Context, submissionID int32) (string, error) {
	row := q.db.QueryRow(ctx, aggregateTestcaseResults, submissionID)
	var status string
	err := row.Scan(&status)
	return status, err
}

const createSubmission = `-- name: CreateSubmission :one
INSERT INTO submissions (game_id, user_id, code, code_size, status)
VALUES ($1, $2, $3, $4, 'running')
RETURNING submission_id
`

type CreateSubmissionParams struct {
	GameID   int32
	UserID   int32
	Code     string
	CodeSize int32
}

func (q *Queries) CreateSubmission(ctx context.Context, arg CreateSubmissionParams) (int32, error) {
	row := q.db.QueryRow(ctx, createSubmission,
		arg.GameID,
		arg.UserID,
		arg.Code,
		arg.CodeSize,
	)
	var submission_id int32
	err := row.Scan(&submission_id)
	return submission_id, err
}

const createTestcaseResult = `-- name: CreateTestcaseResult :exec
INSERT INTO testcase_results (submission_id, testcase_id, status, stdout, stderr)
VALUES ($1, $2, $3, $4, $5)
`

type CreateTestcaseResultParams struct {
	SubmissionID int32
	TestcaseID   int32
	Status       string
	Stdout       string
	Stderr       string
}

func (q *Queries) CreateTestcaseResult(ctx context.Context, arg CreateTestcaseResultParams) error {
	_, err := q.db.Exec(ctx, createTestcaseResult,
		arg.SubmissionID,
		arg.TestcaseID,
		arg.Status,
		arg.Stdout,
		arg.Stderr,
	)
	return err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, display_name, is_admin)
VALUES ($1, $1, false)
RETURNING user_id
`

func (q *Queries) CreateUser(ctx context.Context, username string) (int32, error) {
	row := q.db.QueryRow(ctx, createUser, username)
	var user_id int32
	err := row.Scan(&user_id)
	return user_id, err
}

const createUserAuth = `-- name: CreateUserAuth :exec
INSERT INTO user_auths (user_id, auth_type)
VALUES ($1, $2)
`

type CreateUserAuthParams struct {
	UserID   int32
	AuthType string
}

func (q *Queries) CreateUserAuth(ctx context.Context, arg CreateUserAuthParams) error {
	_, err := q.db.Exec(ctx, createUserAuth, arg.UserID, arg.AuthType)
	return err
}

const getGameByID = `-- name: GetGameByID :one
SELECT game_id, game_type, is_public, display_name, duration_seconds, created_at, started_at, games.problem_id, problems.problem_id, title, description, sample_code FROM games
JOIN problems ON games.problem_id = problems.problem_id
WHERE games.game_id = $1
LIMIT 1
`

type GetGameByIDRow struct {
	GameID          int32
	GameType        string
	IsPublic        bool
	DisplayName     string
	DurationSeconds int32
	CreatedAt       pgtype.Timestamp
	StartedAt       pgtype.Timestamp
	ProblemID       int32
	ProblemID_2     int32
	Title           string
	Description     string
	SampleCode      string
}

func (q *Queries) GetGameByID(ctx context.Context, gameID int32) (GetGameByIDRow, error) {
	row := q.db.QueryRow(ctx, getGameByID, gameID)
	var i GetGameByIDRow
	err := row.Scan(
		&i.GameID,
		&i.GameType,
		&i.IsPublic,
		&i.DisplayName,
		&i.DurationSeconds,
		&i.CreatedAt,
		&i.StartedAt,
		&i.ProblemID,
		&i.ProblemID_2,
		&i.Title,
		&i.Description,
		&i.SampleCode,
	)
	return i, err
}

const getLatestState = `-- name: GetLatestState :one
SELECT game_states.game_id, game_states.user_id, game_states.code, game_states.status, best_score_submission_id, submission_id, submissions.game_id, submissions.user_id, submissions.code, code_size, submissions.status, created_at FROM game_states
LEFT JOIN submissions ON game_states.best_score_submission_id = submissions.submission_id
WHERE game_states.game_id = $1 AND game_states.user_id = $2
LIMIT 1
`

type GetLatestStateParams struct {
	GameID int32
	UserID int32
}

type GetLatestStateRow struct {
	GameID                int32
	UserID                int32
	Code                  string
	Status                string
	BestScoreSubmissionID *int32
	SubmissionID          *int32
	GameID_2              *int32
	UserID_2              *int32
	Code_2                *string
	CodeSize              *int32
	Status_2              *string
	CreatedAt             pgtype.Timestamp
}

func (q *Queries) GetLatestState(ctx context.Context, arg GetLatestStateParams) (GetLatestStateRow, error) {
	row := q.db.QueryRow(ctx, getLatestState, arg.GameID, arg.UserID)
	var i GetLatestStateRow
	err := row.Scan(
		&i.GameID,
		&i.UserID,
		&i.Code,
		&i.Status,
		&i.BestScoreSubmissionID,
		&i.SubmissionID,
		&i.GameID_2,
		&i.UserID_2,
		&i.Code_2,
		&i.CodeSize,
		&i.Status_2,
		&i.CreatedAt,
	)
	return i, err
}

const getLatestStatesOfMainPlayers = `-- name: GetLatestStatesOfMainPlayers :many
SELECT game_main_players.game_id, game_main_players.user_id, game_states.game_id, game_states.user_id, game_states.code, game_states.status, best_score_submission_id, submission_id, submissions.game_id, submissions.user_id, submissions.code, code_size, submissions.status, created_at FROM game_main_players
LEFT JOIN game_states ON game_main_players.game_id = game_states.game_id AND game_main_players.user_id = game_states.user_id
LEFT JOIN submissions ON game_states.best_score_submission_id = submissions.submission_id
WHERE game_main_players.game_id = $1
`

type GetLatestStatesOfMainPlayersRow struct {
	GameID                int32
	UserID                int32
	GameID_2              *int32
	UserID_2              *int32
	Code                  *string
	Status                *string
	BestScoreSubmissionID *int32
	SubmissionID          *int32
	GameID_3              *int32
	UserID_3              *int32
	Code_2                *string
	CodeSize              *int32
	Status_2              *string
	CreatedAt             pgtype.Timestamp
}

func (q *Queries) GetLatestStatesOfMainPlayers(ctx context.Context, gameID int32) ([]GetLatestStatesOfMainPlayersRow, error) {
	rows, err := q.db.Query(ctx, getLatestStatesOfMainPlayers, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLatestStatesOfMainPlayersRow
	for rows.Next() {
		var i GetLatestStatesOfMainPlayersRow
		if err := rows.Scan(
			&i.GameID,
			&i.UserID,
			&i.GameID_2,
			&i.UserID_2,
			&i.Code,
			&i.Status,
			&i.BestScoreSubmissionID,
			&i.SubmissionID,
			&i.GameID_3,
			&i.UserID_3,
			&i.Code_2,
			&i.CodeSize,
			&i.Status_2,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRanking = `-- name: GetRanking :many
SELECT game_states.game_id, game_states.user_id, game_states.code, game_states.status, best_score_submission_id, users.user_id, username, display_name, icon_path, is_admin, label, users.created_at, submission_id, submissions.game_id, submissions.user_id, submissions.code, code_size, submissions.status, submissions.created_at FROM game_states
JOIN users ON game_states.user_id = users.user_id
JOIN submissions ON game_states.best_score_submission_id = submissions.submission_id
WHERE game_states.game_id = $1
ORDER BY submissions.code_size ASC, submissions.created_at ASC
`

type GetRankingRow struct {
	GameID                int32
	UserID                int32
	Code                  string
	Status                string
	BestScoreSubmissionID *int32
	UserID_2              int32
	Username              string
	DisplayName           string
	IconPath              *string
	IsAdmin               bool
	Label                 *string
	CreatedAt             pgtype.Timestamp
	SubmissionID          int32
	GameID_2              int32
	UserID_3              int32
	Code_2                string
	CodeSize              int32
	Status_2              string
	CreatedAt_2           pgtype.Timestamp
}

func (q *Queries) GetRanking(ctx context.Context, gameID int32) ([]GetRankingRow, error) {
	rows, err := q.db.Query(ctx, getRanking, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRankingRow
	for rows.Next() {
		var i GetRankingRow
		if err := rows.Scan(
			&i.GameID,
			&i.UserID,
			&i.Code,
			&i.Status,
			&i.BestScoreSubmissionID,
			&i.UserID_2,
			&i.Username,
			&i.DisplayName,
			&i.IconPath,
			&i.IsAdmin,
			&i.Label,
			&i.CreatedAt,
			&i.SubmissionID,
			&i.GameID_2,
			&i.UserID_3,
			&i.Code_2,
			&i.CodeSize,
			&i.Status_2,
			&i.CreatedAt_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserAuthByUsername = `-- name: GetUserAuthByUsername :one
SELECT users.user_id, username, display_name, icon_path, is_admin, label, created_at, user_auth_id, user_auths.user_id, auth_type, password_hash FROM users
JOIN user_auths ON users.user_id = user_auths.user_id
WHERE users.username = $1
LIMIT 1
`

type GetUserAuthByUsernameRow struct {
	UserID       int32
	Username     string
	DisplayName  string
	IconPath     *string
	IsAdmin      bool
	Label        *string
	CreatedAt    pgtype.Timestamp
	UserAuthID   int32
	UserID_2     int32
	AuthType     string
	PasswordHash *string
}

func (q *Queries) GetUserAuthByUsername(ctx context.Context, username string) (GetUserAuthByUsernameRow, error) {
	row := q.db.QueryRow(ctx, getUserAuthByUsername, username)
	var i GetUserAuthByUsernameRow
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.DisplayName,
		&i.IconPath,
		&i.IsAdmin,
		&i.Label,
		&i.CreatedAt,
		&i.UserAuthID,
		&i.UserID_2,
		&i.AuthType,
		&i.PasswordHash,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT user_id, username, display_name, icon_path, is_admin, label, created_at FROM users
WHERE users.user_id = $1
LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, userID int32) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.DisplayName,
		&i.IconPath,
		&i.IsAdmin,
		&i.Label,
		&i.CreatedAt,
	)
	return i, err
}

const getUserIDByUsername = `-- name: GetUserIDByUsername :one
SELECT user_id FROM users
WHERE users.username = $1
LIMIT 1
`

func (q *Queries) GetUserIDByUsername(ctx context.Context, username string) (int32, error) {
	row := q.db.QueryRow(ctx, getUserIDByUsername, username)
	var user_id int32
	err := row.Scan(&user_id)
	return user_id, err
}

const listAllGames = `-- name: ListAllGames :many
SELECT game_id, game_type, is_public, display_name, duration_seconds, created_at, started_at, problem_id FROM games
ORDER BY games.game_id
`

func (q *Queries) ListAllGames(ctx context.Context) ([]Game, error) {
	rows, err := q.db.Query(ctx, listAllGames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Game
	for rows.Next() {
		var i Game
		if err := rows.Scan(
			&i.GameID,
			&i.GameType,
			&i.IsPublic,
			&i.DisplayName,
			&i.DurationSeconds,
			&i.CreatedAt,
			&i.StartedAt,
			&i.ProblemID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMainPlayers = `-- name: ListMainPlayers :many
SELECT game_id, game_main_players.user_id, users.user_id, username, display_name, icon_path, is_admin, label, created_at FROM game_main_players
JOIN users ON game_main_players.user_id = users.user_id
WHERE game_main_players.game_id = ANY($1::INT[])
ORDER BY game_main_players.user_id
`

type ListMainPlayersRow struct {
	GameID      int32
	UserID      int32
	UserID_2    int32
	Username    string
	DisplayName string
	IconPath    *string
	IsAdmin     bool
	Label       *string
	CreatedAt   pgtype.Timestamp
}

func (q *Queries) ListMainPlayers(ctx context.Context, dollar_1 []int32) ([]ListMainPlayersRow, error) {
	rows, err := q.db.Query(ctx, listMainPlayers, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListMainPlayersRow
	for rows.Next() {
		var i ListMainPlayersRow
		if err := rows.Scan(
			&i.GameID,
			&i.UserID,
			&i.UserID_2,
			&i.Username,
			&i.DisplayName,
			&i.IconPath,
			&i.IsAdmin,
			&i.Label,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPublicGames = `-- name: ListPublicGames :many
SELECT game_id, game_type, is_public, display_name, duration_seconds, created_at, started_at, games.problem_id, problems.problem_id, title, description, sample_code FROM games
JOIN problems ON games.problem_id = problems.problem_id
WHERE is_public = true
ORDER BY games.game_id
`

type ListPublicGamesRow struct {
	GameID          int32
	GameType        string
	IsPublic        bool
	DisplayName     string
	DurationSeconds int32
	CreatedAt       pgtype.Timestamp
	StartedAt       pgtype.Timestamp
	ProblemID       int32
	ProblemID_2     int32
	Title           string
	Description     string
	SampleCode      string
}

func (q *Queries) ListPublicGames(ctx context.Context) ([]ListPublicGamesRow, error) {
	rows, err := q.db.Query(ctx, listPublicGames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPublicGamesRow
	for rows.Next() {
		var i ListPublicGamesRow
		if err := rows.Scan(
			&i.GameID,
			&i.GameType,
			&i.IsPublic,
			&i.DisplayName,
			&i.DurationSeconds,
			&i.CreatedAt,
			&i.StartedAt,
			&i.ProblemID,
			&i.ProblemID_2,
			&i.Title,
			&i.Description,
			&i.SampleCode,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTestcasesByGameID = `-- name: ListTestcasesByGameID :many
SELECT testcase_id, problem_id, stdin, stdout FROM testcases
WHERE testcases.problem_id = (SELECT problem_id FROM games WHERE game_id = $1)
ORDER BY testcases.testcase_id
`

func (q *Queries) ListTestcasesByGameID(ctx context.Context, gameID int32) ([]Testcase, error) {
	rows, err := q.db.Query(ctx, listTestcasesByGameID, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Testcase
	for rows.Next() {
		var i Testcase
		if err := rows.Scan(
			&i.TestcaseID,
			&i.ProblemID,
			&i.Stdin,
			&i.Stdout,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
SELECT user_id, username, display_name, icon_path, is_admin, label, created_at FROM users
ORDER BY users.user_id
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.UserID,
			&i.Username,
			&i.DisplayName,
			&i.IconPath,
			&i.IsAdmin,
			&i.Label,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const syncGameStateBestScoreSubmission = `-- name: SyncGameStateBestScoreSubmission :exec
UPDATE game_states
SET best_score_submission_id = (
    SELECT submission_id FROM submissions AS s
    WHERE s.game_id = $1 AND s.user_id = $2 AND s.status = 'success'
    ORDER BY s.code_size ASC, s.created_at ASC
    LIMIT 1
)
WHERE game_id = $1 AND user_id = $2
`

type SyncGameStateBestScoreSubmissionParams struct {
	GameID int32
	UserID int32
}

func (q *Queries) SyncGameStateBestScoreSubmission(ctx context.Context, arg SyncGameStateBestScoreSubmissionParams) error {
	_, err := q.db.Exec(ctx, syncGameStateBestScoreSubmission, arg.GameID, arg.UserID)
	return err
}

const updateCode = `-- name: UpdateCode :exec
INSERT INTO game_states (game_id, user_id, code, status)
VALUES ($1, $2, $3, $4)
ON CONFLICT (game_id, user_id)
DO UPDATE SET code = EXCLUDED.code
`

type UpdateCodeParams struct {
	GameID int32
	UserID int32
	Code   string
	Status string
}

func (q *Queries) UpdateCode(ctx context.Context, arg UpdateCodeParams) error {
	_, err := q.db.Exec(ctx, updateCode,
		arg.GameID,
		arg.UserID,
		arg.Code,
		arg.Status,
	)
	return err
}

const updateGame = `-- name: UpdateGame :exec
UPDATE games
SET
    game_type = $2,
    is_public = $3,
    display_name = $4,
    duration_seconds = $5,
    started_at = $6,
    problem_id = $7
WHERE game_id = $1
`

type UpdateGameParams struct {
	GameID          int32
	GameType        string
	IsPublic        bool
	DisplayName     string
	DurationSeconds int32
	StartedAt       pgtype.Timestamp
	ProblemID       int32
}

func (q *Queries) UpdateGame(ctx context.Context, arg UpdateGameParams) error {
	_, err := q.db.Exec(ctx, updateGame,
		arg.GameID,
		arg.GameType,
		arg.IsPublic,
		arg.DisplayName,
		arg.DurationSeconds,
		arg.StartedAt,
		arg.ProblemID,
	)
	return err
}

const updateGameStartedAt = `-- name: UpdateGameStartedAt :exec
UPDATE games
SET started_at = $2
WHERE game_id = $1
`

type UpdateGameStartedAtParams struct {
	GameID    int32
	StartedAt pgtype.Timestamp
}

func (q *Queries) UpdateGameStartedAt(ctx context.Context, arg UpdateGameStartedAtParams) error {
	_, err := q.db.Exec(ctx, updateGameStartedAt, arg.GameID, arg.StartedAt)
	return err
}

const updateGameStateStatus = `-- name: UpdateGameStateStatus :exec
UPDATE game_states
SET status = $3
WHERE game_id = $1 AND user_id = $2
`

type UpdateGameStateStatusParams struct {
	GameID int32
	UserID int32
	Status string
}

func (q *Queries) UpdateGameStateStatus(ctx context.Context, arg UpdateGameStateStatusParams) error {
	_, err := q.db.Exec(ctx, updateGameStateStatus, arg.GameID, arg.UserID, arg.Status)
	return err
}

const updateSubmissionStatus = `-- name: UpdateSubmissionStatus :exec
UPDATE submissions
SET status = $2
WHERE submission_id = $1
`

type UpdateSubmissionStatusParams struct {
	SubmissionID int32
	Status       string
}

func (q *Queries) UpdateSubmissionStatus(ctx context.Context, arg UpdateSubmissionStatusParams) error {
	_, err := q.db.Exec(ctx, updateSubmissionStatus, arg.SubmissionID, arg.Status)
	return err
}

const updateUserIconPath = `-- name: UpdateUserIconPath :exec
UPDATE users
SET icon_path = $2
WHERE user_id = $1
`

type UpdateUserIconPathParams struct {
	UserID   int32
	IconPath *string
}

func (q *Queries) UpdateUserIconPath(ctx context.Context, arg UpdateUserIconPathParams) error {
	_, err := q.db.Exec(ctx, updateUserIconPath, arg.UserID, arg.IconPath)
	return err
}
