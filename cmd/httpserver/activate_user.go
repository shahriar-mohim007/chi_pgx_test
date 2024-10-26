package httpserver

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go_chi_pgx/state"
	utils "go_chi_pgx/utils"
	"net/http"
)

func HandleActivateUser(app *state.State) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		tokenString := req.URL.Query().Get("token")
		ctx := req.Context()

		if tokenString == "" {
			app.Logger.PrintError(fmt.Errorf("missing token"), map[string]string{
				"context": "missing token",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}

		var claims utils.Claims
		secretKey := app.Config.SecretKey

		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			app.Logger.PrintError(err, map[string]string{
				"context": "invalid token",
			})
			_ = InvalidToken.WriteToResponse(w, nil)
			return
		}

		userID := claims.UserID

		err = app.Repository.ActivateUserByID(ctx, userID)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "failed to activate user",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}
		_ = UserActivated.WriteToResponse(w, nil)
		return

	}
}
