package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/amir-amirov/go-social-media/internal/utils"
)

func (app *application) AuthTokenMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unAuthorizedErrorResponse(w, r, errors.New("authorization header is missing"))
				return
			}

			// verify token and get userID
			user_id, err := utils.VerifyToken(authHeader)
			if err != nil {
				app.unAuthorizedErrorResponse(w, r, err)
				return
			}

			// fetch user
			user, err := app.store.Users.GetByID(r.Context(), user_id)
			if err != nil {
				app.unAuthorizedErrorResponse(w, r, err)
				return
			}

			// save the user in context
			ctx := context.WithValue(r.Context(), userKeyCtx, user)

			// move to the main handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
