package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=72"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	store.User			"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Println("payload:", payload)

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// hash the user password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// token is a code which will be sent to user via email
	token := uuid.New().String()
	fmt.Println("token: ", token)

	// store the user
	if err := app.store.Users.CreateAndInvite(r.Context(), &user, token, app.config.mail.exp); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, store.ErrDuplicateEmail)
		case store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, store.ErrDuplicateUsername)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// mail

	w.WriteHeader(http.StatusCreated)

}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if err := app.store.Users.Activate(r.Context(), token); err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, errors.New("invalid token"))
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
