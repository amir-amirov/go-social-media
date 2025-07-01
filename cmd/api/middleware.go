package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/amir-amirov/go-social-media/internal/utils"
)

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := app.getUserFromCtx(r)
		post := app.getPostFromCtx(r)

		// check is this user the owner of the post
		if user.ID == post.UserID {
			next.ServeHTTP(w, r)
			return
		}

		// role precedence check (for admin/moderator)
		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			switch err.Error() {
			case "not allowed role":
				app.badRequestResponse(w, r, fmt.Errorf("not allowed role"))
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		if !allowed {
			app.forbiddenResponse(w, r, store.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	userRole, err := app.store.Roles.GetByID(ctx, user.RoleID)
	if err != nil {
		return false, err
	}

	return role.Level <= userRole.Level, nil
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
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

// func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			authHeader := r.Header.Get("Authorization")
// 			if authHeader == "" {
// 				app.unAuthorizedErrorResponse(w, r, errors.New("authorization header is missing"))
// 				return
// 			}

// 		}
// 	}
// }
