package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKeyType string

const userKeyCtx userKeyType = "user"

// GetUser godoc
//
//	@Summary		Fetches user profile
//	@Description	Fetches user profile by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	user := app.getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follow a user by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		200		{string}	string	"User followed"
//	@Failure		400		{object}	error	"Invalid user id"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		500		{object}	error	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {

	follower := app.getUserFromCtx(r) // This is the user who is following (making the request)

	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if userID == follower.ID {
		app.badRequestResponse(w, r, errors.New("you cannot follow yourself"))
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, userID, follower.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	follower := app.getUserFromCtx(r) // This is the user who is unfollowing (making the request)

	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if userID == follower.ID {
		app.badRequestResponse(w, r, errors.New("you cannot unfollow yourself"))
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.UnFollow(ctx, userID, follower.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) getUserFromCtx(r *http.Request) *store.User {
	user := r.Context().Value(userKeyCtx).(*store.User)
	return user
}
