package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *Application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.authorizationError(w, r, fmt.Errorf("authorization token is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.authorizationError(w, r, fmt.Errorf("authorization token form is not correct"))
			return
		}

		token := parts[1]
		jwtToken, err := app.Auth.ValidateToken(token)
		if err != nil {
			app.authorizationError(w, r, err)
			return
		}
		claims := jwtToken.Claims.(jwt.MapClaims)
		id, ok := claims["sub"]
		if !ok {
			app.authorizationError(w, r, err)
			return
		}
		if err != nil {
			app.authorizationError(w, r, err)
			return
		}
		ctx := r.Context()
		parsedId, ok := id.(string)
		if !ok {
			app.authorizationError(w, r, err)
		}
		user, err := app.Store.Admin.GetById(r.Context(), parsedId)
		if err != nil {
			app.authorizationError(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, "adminContext", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
