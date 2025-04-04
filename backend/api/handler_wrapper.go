// Code generated by go generate; DO NOT EDIT.

package api

import (
	"context"
	"errors"
	"strings"

	"github.com/nsfisis/phperkaigi-2025-albatross/backend/auth"
	"github.com/nsfisis/phperkaigi-2025-albatross/backend/db"
)

var _ StrictServerInterface = (*HandlerWrapper)(nil)

type HandlerWrapper struct {
	impl Handler
}

func NewHandler(queries *db.Queries, hub GameHubInterface) *HandlerWrapper {
	return &HandlerWrapper{
		impl: Handler{
			q:   queries,
			hub: hub,
		},
	}
}

func parseJWTClaimsFromAuthorizationHeader(authorization string) (*auth.JWTClaims, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return nil, errors.New("invalid authorization header")
	}
	token := authorization[len(prefix):]
	claims, err := auth.ParseJWT(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (h *HandlerWrapper) GetGame(ctx context.Context, request GetGameRequestObject) (GetGameResponseObject, error) {
	user, err := parseJWTClaimsFromAuthorizationHeader(request.Params.Authorization)
	if err != nil {
		return GetGame401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "Unauthorized",
			},
		}, nil
	}
	return h.impl.GetGame(ctx, request, user)
}

func (h *HandlerWrapper) GetGamePlayLatestState(ctx context.Context, request GetGamePlayLatestStateRequestObject) (GetGamePlayLatestStateResponseObject, error) {
	user, err := parseJWTClaimsFromAuthorizationHeader(request.Params.Authorization)
	if err != nil {
		return GetGamePlayLatestState401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "Unauthorized",
			},
		}, nil
	}
	return h.impl.GetGamePlayLatestState(ctx, request, user)
}

func (h *HandlerWrapper) GetGameWatchLatestStates(ctx context.Context, request GetGameWatchLatestStatesRequestObject) (GetGameWatchLatestStatesResponseObject, error) {
	user, err := parseJWTClaimsFromAuthorizationHeader(request.Params.Authorization)
	if err != nil {
		return GetGameWatchLatestStates401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "Unauthorized",
			},
		}, nil
	}
	return h.impl.GetGameWatchLatestStates(ctx, request, user)
}

func (h *HandlerWrapper) GetGameWatchRanking(ctx context.Context, request GetGameWatchRankingRequestObject) (GetGameWatchRankingResponseObject, error) {
	user, err := parseJWTClaimsFromAuthorizationHeader(request.Params.Authorization)
	if err != nil {
		return GetGameWatchRanking401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "Unauthorized",
			},
		}, nil
	}
	return h.impl.GetGameWatchRanking(ctx, request, user)
}

func (h *HandlerWrapper) GetGames(ctx context.Context, request GetGamesRequestObject) (GetGamesResponseObject, error) {
	user, err := parseJWTClaimsFromAuthorizationHeader(request.Params.Authorization)
	if err != nil {
		return GetGames401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "Unauthorized",
			},
		}, nil
	}
	return h.impl.GetGames(ctx, request, user)
}

func (h *HandlerWrapper) PostGamePlayCode(ctx context.Context, request PostGamePlayCodeRequestObject) (PostGamePlayCodeResponseObject, error) {
	user, err := parseJWTClaimsFromAuthorizationHeader(request.Params.Authorization)
	if err != nil {
		return PostGamePlayCode401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "Unauthorized",
			},
		}, nil
	}
	return h.impl.PostGamePlayCode(ctx, request, user)
}

func (h *HandlerWrapper) PostGamePlaySubmit(ctx context.Context, request PostGamePlaySubmitRequestObject) (PostGamePlaySubmitResponseObject, error) {
	user, err := parseJWTClaimsFromAuthorizationHeader(request.Params.Authorization)
	if err != nil {
		return PostGamePlaySubmit401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "Unauthorized",
			},
		}, nil
	}
	return h.impl.PostGamePlaySubmit(ctx, request, user)
}

func (h *HandlerWrapper) PostLogin(ctx context.Context, request PostLoginRequestObject) (PostLoginResponseObject, error) {
	return h.impl.PostLogin(ctx, request)
}
