// Code generated by go generate; DO NOT EDIT.

package api

import (
	"context"
	"errors"
	"strings"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/auth"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

var _ StrictServerInterface = (*ApiHandlerWrapper)(nil)

type ApiHandlerWrapper struct {
	innerHandler Handler
}

func NewHandler(queries *db.Queries, hubs GameHubsInterface) *ApiHandlerWrapper {
	return &ApiHandlerWrapper{
		innerHandler: Handler{
			q:    queries,
			hubs: hubs,
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

func (h *ApiHandlerWrapper) GetGame(ctx context.Context, request GetGameRequestObject) (GetGameResponseObject, error) {
	user, err := parseJWTClaimsFromAuthorizationHeader(request.Params.Authorization)
	if err != nil {
		return GetGame401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "Unauthorized",
			},
		}, nil
	}
	return h.innerHandler.GetGame(ctx, request, user)
}

func (h *ApiHandlerWrapper) GetGames(ctx context.Context, request GetGamesRequestObject) (GetGamesResponseObject, error) {
	user, err := parseJWTClaimsFromAuthorizationHeader(request.Params.Authorization)
	if err != nil {
		return GetGames401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "Unauthorized",
			},
		}, nil
	}
	return h.innerHandler.GetGames(ctx, request, user)
}

func (h *ApiHandlerWrapper) GetToken(ctx context.Context, request GetTokenRequestObject) (GetTokenResponseObject, error) {
	user, err := parseJWTClaimsFromAuthorizationHeader(request.Params.Authorization)
	if err != nil {
		return GetToken401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "Unauthorized",
			},
		}, nil
	}
	return h.innerHandler.GetToken(ctx, request, user)
}

func (h *ApiHandlerWrapper) PostLogin(ctx context.Context, request PostLoginRequestObject) (PostLoginResponseObject, error) {
	return h.innerHandler.PostLogin(ctx, request)
}
