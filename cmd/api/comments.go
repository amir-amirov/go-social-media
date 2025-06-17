package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/go-chi/chi/v5"
)

type commentKeyType string

const commentKeyCtx commentKeyType = "comment"

type commentPayload struct {
	Content string `json:"content" validate:"required,max=200"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload commentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := getPostFromCtx(r)
	ctx := r.Context()

	var userID int64 = 65

	comment := &store.Comment{
		Content: payload.Content,
		PostID:  post.ID,
		UserID:  userID,
	}

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	comment := getCommentFromCtx(r)

	user, err := app.store.Users.GetByID(r.Context(), comment.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	comment.User = *user

	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	comment := getCommentFromCtx(r)

	var payload commentPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment.Content = payload.Content

	if err := app.store.Comments.Update(r.Context(), comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	comment := getCommentFromCtx(r)

	if err := app.store.Comments.Delete(r.Context(), comment.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// program will reach end of the func and send response
	// can do nothing, response will be sent as 200
}

func (app *application) commentsContextMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "commentID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid comment id"))
			return
		}

		ctx := r.Context()

		comment, err := app.store.Comments.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, commentKeyCtx, comment)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getCommentFromCtx(r *http.Request) *store.Comment {
	comment := r.Context().Value(commentKeyCtx).(*store.Comment)
	return comment
}
