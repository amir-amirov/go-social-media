package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKeyType string

const userKeyCtx userKeyType = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	user := app.getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {

	followUser := app.getUserFromCtx(r)
	var userID int64 = 5
	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, followUser.ID, userID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	unfollowUser := app.getUserFromCtx(r)
	var userID int64 = 5
	ctx := r.Context()

	if err := app.store.Followers.UnFollow(ctx, unfollowUser.ID, userID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		user_id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid user id"))
			return
		}

		user, err := app.store.Users.GetByID(r.Context(), user_id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), userKeyCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUserFromCtx(r *http.Request) *store.User {
	user := r.Context().Value(userKeyCtx).(*store.User)
	return user
}
