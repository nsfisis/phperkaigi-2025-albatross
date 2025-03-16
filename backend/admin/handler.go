package admin

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"github.com/nsfisis/phperkaigi-2025-albatross/backend/account"
	"github.com/nsfisis/phperkaigi-2025-albatross/backend/auth"
	"github.com/nsfisis/phperkaigi-2025-albatross/backend/db"
)

const (
	basePath = "/phperkaigi/2025/code-battle"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

type Handler struct {
	q *db.Queries
}

func NewHandler(q *db.Queries) *Handler {
	return &Handler{q: q}
}

func newAdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			jwt, err := c.Cookie("albatross_token")
			if err != nil {
				return c.Redirect(http.StatusSeeOther, basePath+"/login")
			}
			claims, err := auth.ParseJWT(jwt.Value)
			if err != nil {
				return c.Redirect(http.StatusSeeOther, basePath+"/login")
			}
			if !claims.IsAdmin {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			return next(c)
		}
	}
}

func (h *Handler) RegisterHandlers(g *echo.Group) {
	g.Use(newAssetsMiddleware())
	g.Use(newAdminMiddleware())

	g.GET("/dashboard", h.getDashboard)
	g.GET("/users", h.getUsers)
	g.GET("/users/:userID", h.getUserEdit)
	g.POST("/users/:userID", h.postUserEdit)
	g.POST("/users/:userID/fetch-icon", h.postUserFetchIcon)
	g.GET("/games", h.getGames)
	g.GET("/games/:gameID", h.getGameEdit)
	g.POST("/games/:gameID", h.postGameEdit)
	g.POST("/games/:gameID/start", h.postGameStart)
}

func (h *Handler) getDashboard(c echo.Context) error {
	return c.Render(http.StatusOK, "dashboard", echo.Map{
		"BasePath": basePath,
		"Title":    "Dashboard",
	})
}

func (h *Handler) getUsers(c echo.Context) error {
	rows, err := h.q.ListUsers(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	users := make([]echo.Map, len(rows))
	for i, u := range rows {
		users[i] = echo.Map{
			"UserID":      u.UserID,
			"Username":    u.Username,
			"DisplayName": u.DisplayName,
			"IconPath":    u.IconPath,
			"IsAdmin":     u.IsAdmin,
		}
	}

	return c.Render(http.StatusOK, "users", echo.Map{
		"BasePath": basePath,
		"Title":    "Users",
		"Users":    users,
	})
}

func (h *Handler) getUserEdit(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
	}
	row, err := h.q.GetUserByID(c.Request().Context(), int32(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "user_edit", echo.Map{
		"BasePath": basePath,
		"Title":    "User Edit",
		"User": echo.Map{
			"UserID":      row.UserID,
			"Username":    row.Username,
			"DisplayName": row.DisplayName,
			"IconPath":    row.IconPath,
			"IsAdmin":     row.IsAdmin,
			"Label":       row.Label,
		},
	})
}

func (h *Handler) postUserEdit(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user_id")
	}

	displayName := c.FormValue("display_name")
	iconPathRaw := c.FormValue("icon_path")
	isAdmin := (c.FormValue("is_admin") != "")
	labelRaw := c.FormValue("label")

	var iconPath *string
	if iconPathRaw != "" {
		iconPath = &iconPathRaw
	}
	var label *string
	if labelRaw != "" {
		label = &labelRaw
	}

	err = h.q.UpdateUser(c.Request().Context(), db.UpdateUserParams{
		UserID:      int32(userID),
		DisplayName: displayName,
		IconPath:    iconPath,
		IsAdmin:     isAdmin,
		Label:       label,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusSeeOther, basePath+"/admin/users")
}

func (h *Handler) postUserFetchIcon(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
	}
	row, err := h.q.GetUserByID(c.Request().Context(), int32(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	go func() {
		err := account.FetchIcon(context.Background(), h.q, int(row.UserID))
		if err != nil {
			log.Printf("%v", err)
			// The failure is intentionally ignored. Retry manually if needed.
		}
	}()
	return c.Redirect(http.StatusSeeOther, basePath+"/admin/users")
}

func (h *Handler) getGames(c echo.Context) error {
	rows, err := h.q.ListAllGames(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	games := make([]echo.Map, len(rows))
	for i, g := range rows {
		var startedAt string
		if g.StartedAt.Valid {
			startedAt = g.StartedAt.Time.In(jst).Format("2006-01-02T15:04")
		}
		games[i] = echo.Map{
			"GameID":          g.GameID,
			"GameType":        g.GameType,
			"IsPublic":        g.IsPublic,
			"DisplayName":     g.DisplayName,
			"DurationSeconds": g.DurationSeconds,
			"StartedAt":       startedAt,
			"ProblemID":       g.ProblemID,
		}
	}

	return c.Render(http.StatusOK, "games", echo.Map{
		"BasePath": basePath,
		"Title":    "Games",
		"Games":    games,
	})
}

func (h *Handler) getGameEdit(c echo.Context) error {
	gameID, err := strconv.Atoi(c.Param("gameID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
	}
	row, err := h.q.GetGameByID(c.Request().Context(), int32(gameID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var startedAt string
	if row.StartedAt.Valid {
		startedAt = row.StartedAt.Time.In(jst).Format("2006-01-02T15:04")
	}

	mainPlayerRows, err := h.q.ListMainPlayers(c.Request().Context(), []int32{int32(gameID)})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	mainPlayer1 := 0
	if len(mainPlayerRows) > 0 {
		mainPlayer1 = int(mainPlayerRows[0].UserID)
	}
	mainPlayer2 := 0
	if len(mainPlayerRows) > 1 {
		mainPlayer2 = int(mainPlayerRows[1].UserID)
	}

	userRows, err := h.q.ListUsers(c.Request().Context())
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	var users []echo.Map
	for _, r := range userRows {
		users = append(users, echo.Map{
			"UserID":   int(r.UserID),
			"Username": r.Username,
		})
	}

	return c.Render(http.StatusOK, "game_edit", echo.Map{
		"BasePath": basePath,
		"Title":    "Game Edit",
		"Game": echo.Map{
			"GameID":          row.GameID,
			"GameType":        row.GameType,
			"IsPublic":        row.IsPublic,
			"DisplayName":     row.DisplayName,
			"DurationSeconds": row.DurationSeconds,
			"StartedAt":       startedAt,
			"ProblemID":       row.ProblemID,
			"MainPlayer1":     mainPlayer1,
			"MainPlayer2":     mainPlayer2,
		},
		"Users": users,
	})
}

func (h *Handler) postGameEdit(c echo.Context) error {
	gameID, err := strconv.Atoi(c.Param("gameID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
	}

	gameType := c.FormValue("game_type")
	isPublic := (c.FormValue("is_public") != "")
	displayName := c.FormValue("display_name")
	durationSeconds, err := strconv.Atoi(c.FormValue("duration_seconds"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid duration_seconds")
	}
	var problemID int
	{
		problemIDRaw := c.FormValue("problem_id")
		problemIDInt, err := strconv.Atoi(problemIDRaw)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid problem_id")
		}
		problemID = problemIDInt
	}
	var startedAt *time.Time
	{
		startedAtRaw := c.FormValue("started_at")
		if startedAtRaw != "" {
			startedAtJST, err := time.ParseInLocation("2006-01-02T15:04", startedAtRaw, jst)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid started_at")
			}
			startedAtUTC := startedAtJST.UTC()
			startedAt = &startedAtUTC
		}
	}

	var changedStartedAt pgtype.Timestamp
	if startedAt == nil {
		changedStartedAt = pgtype.Timestamp{
			Valid: false,
		}
	} else {
		changedStartedAt = pgtype.Timestamp{
			Time:  *startedAt,
			Valid: true,
		}
	}

	err = h.q.UpdateGame(c.Request().Context(), db.UpdateGameParams{
		GameID:          int32(gameID),
		GameType:        gameType,
		IsPublic:        isPublic,
		DisplayName:     displayName,
		DurationSeconds: int32(durationSeconds),
		StartedAt:       changedStartedAt,
		ProblemID:       int32(problemID),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	mainPlayers := []int{}
	mainPlayer1Raw := c.FormValue("main_player_1")
	if mainPlayer1Raw != "" && mainPlayer1Raw != "0" {
		mainPlayer1, err := strconv.Atoi(mainPlayer1Raw)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid main_player_1")
		}
		mainPlayers = append(mainPlayers, mainPlayer1)
	}
	mainPlayer2Raw := c.FormValue("main_player_2")
	if mainPlayer2Raw != "" && mainPlayer2Raw != "0" {
		mainPlayer2, err := strconv.Atoi(mainPlayer2Raw)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid main_player_2")
		}
		mainPlayers = append(mainPlayers, mainPlayer2)
	}

	err = h.q.RemoveAllMainPlayers(c.Request().Context(), int32(gameID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, userID := range mainPlayers {
		err = h.q.AddMainPlayer(c.Request().Context(), db.AddMainPlayerParams{
			GameID: int32(gameID),
			UserID: int32(userID),
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.Redirect(http.StatusSeeOther, basePath+"/admin/games")
}

func (h *Handler) postGameStart(c echo.Context) error {
	gameID, err := strconv.Atoi(c.Param("gameID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
	}

	startedAt := time.Now().Add(11 * time.Second)

	err = h.q.UpdateGameStartedAt(c.Request().Context(), db.UpdateGameStartedAtParams{
		GameID: int32(gameID),
		StartedAt: pgtype.Timestamp{
			Time:  startedAt,
			Valid: true,
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusSeeOther, basePath+"/admin/games")
}
