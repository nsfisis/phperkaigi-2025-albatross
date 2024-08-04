// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getGameByID = `-- name: GetGameByID :one
SELECT game_id, game_type, state, display_name, duration_seconds, created_at, started_at, games.problem_id, problems.problem_id, title, description FROM games
LEFT JOIN problems ON games.problem_id = problems.problem_id
WHERE games.game_id = $1
LIMIT 1
`

type GetGameByIDRow struct {
	GameID          int32
	GameType        string
	State           string
	DisplayName     string
	DurationSeconds int32
	CreatedAt       pgtype.Timestamp
	StartedAt       pgtype.Timestamp
	ProblemID       *int32
	ProblemID_2     *int32
	Title           *string
	Description     *string
}

func (q *Queries) GetGameByID(ctx context.Context, gameID int32) (GetGameByIDRow, error) {
	row := q.db.QueryRow(ctx, getGameByID, gameID)
	var i GetGameByIDRow
	err := row.Scan(
		&i.GameID,
		&i.GameType,
		&i.State,
		&i.DisplayName,
		&i.DurationSeconds,
		&i.CreatedAt,
		&i.StartedAt,
		&i.ProblemID,
		&i.ProblemID_2,
		&i.Title,
		&i.Description,
	)
	return i, err
}

const getUserAuthByUsername = `-- name: GetUserAuthByUsername :one
SELECT users.user_id, username, display_name, icon_path, is_admin, created_at, user_auth_id, user_auths.user_id, auth_type, password_hash FROM users
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
		&i.CreatedAt,
		&i.UserAuthID,
		&i.UserID_2,
		&i.AuthType,
		&i.PasswordHash,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT user_id, username, display_name, icon_path, is_admin, created_at FROM users
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
		&i.CreatedAt,
	)
	return i, err
}

const listGames = `-- name: ListGames :many
SELECT game_id, game_type, state, display_name, duration_seconds, created_at, started_at, games.problem_id, problems.problem_id, title, description FROM games
LEFT JOIN problems ON games.problem_id = problems.problem_id
`

type ListGamesRow struct {
	GameID          int32
	GameType        string
	State           string
	DisplayName     string
	DurationSeconds int32
	CreatedAt       pgtype.Timestamp
	StartedAt       pgtype.Timestamp
	ProblemID       *int32
	ProblemID_2     *int32
	Title           *string
	Description     *string
}

func (q *Queries) ListGames(ctx context.Context) ([]ListGamesRow, error) {
	rows, err := q.db.Query(ctx, listGames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListGamesRow
	for rows.Next() {
		var i ListGamesRow
		if err := rows.Scan(
			&i.GameID,
			&i.GameType,
			&i.State,
			&i.DisplayName,
			&i.DurationSeconds,
			&i.CreatedAt,
			&i.StartedAt,
			&i.ProblemID,
			&i.ProblemID_2,
			&i.Title,
			&i.Description,
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

const listGamesForPlayer = `-- name: ListGamesForPlayer :many
SELECT games.game_id, game_type, state, display_name, duration_seconds, created_at, started_at, games.problem_id, problems.problem_id, title, description, game_players.game_id, user_id FROM games
LEFT JOIN problems ON games.problem_id = problems.problem_id
JOIN game_players ON games.game_id = game_players.game_id
WHERE game_players.user_id = $1
`

type ListGamesForPlayerRow struct {
	GameID          int32
	GameType        string
	State           string
	DisplayName     string
	DurationSeconds int32
	CreatedAt       pgtype.Timestamp
	StartedAt       pgtype.Timestamp
	ProblemID       *int32
	ProblemID_2     *int32
	Title           *string
	Description     *string
	GameID_2        int32
	UserID          int32
}

func (q *Queries) ListGamesForPlayer(ctx context.Context, userID int32) ([]ListGamesForPlayerRow, error) {
	rows, err := q.db.Query(ctx, listGamesForPlayer, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListGamesForPlayerRow
	for rows.Next() {
		var i ListGamesForPlayerRow
		if err := rows.Scan(
			&i.GameID,
			&i.GameType,
			&i.State,
			&i.DisplayName,
			&i.DurationSeconds,
			&i.CreatedAt,
			&i.StartedAt,
			&i.ProblemID,
			&i.ProblemID_2,
			&i.Title,
			&i.Description,
			&i.GameID_2,
			&i.UserID,
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
SELECT user_id, username, display_name, icon_path, is_admin, created_at FROM users
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

const updateGame = `-- name: UpdateGame :exec
UPDATE games
SET
    state = $2,
    display_name = $3,
    duration_seconds = $4,
    started_at = $5,
    problem_id = $6
WHERE game_id = $1
`

type UpdateGameParams struct {
	GameID          int32
	State           string
	DisplayName     string
	DurationSeconds int32
	StartedAt       pgtype.Timestamp
	ProblemID       *int32
}

func (q *Queries) UpdateGame(ctx context.Context, arg UpdateGameParams) error {
	_, err := q.db.Exec(ctx, updateGame,
		arg.GameID,
		arg.State,
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

const updateGameState = `-- name: UpdateGameState :exec
UPDATE games
SET state = $2
WHERE game_id = $1
`

type UpdateGameStateParams struct {
	GameID int32
	State  string
}

func (q *Queries) UpdateGameState(ctx context.Context, arg UpdateGameStateParams) error {
	_, err := q.db.Exec(ctx, updateGameState, arg.GameID, arg.State)
	return err
}
