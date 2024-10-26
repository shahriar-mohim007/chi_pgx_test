package httpserver

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go_chi_pgx/state"
	utils "go_chi_pgx/utils"
	"net/http"
	"time"
)

type LoginRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponsePayload struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func HandleLogin(app *state.State) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		request := LoginRequestPayload{}
		ctx := req.Context()

		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Invalid JSON",
			})
			_ = ValidDataNotFound.WriteToResponse(w, nil)
			return
		}

		validate := validator.New()

		err = validate.Struct(request)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Invalid payload",
			})
			_ = ValidDataNotFound.WriteToResponse(w, nil)
			return
		}

		user, err := app.Repository.GetUserByEmail(ctx, request.Email)
		if err != nil {
			_ = InvalidEmailPassword.WriteToResponse(w, nil)
			return
		}

		if !utils.CheckPasswordHash(user.Password, request.Password) {
			_ = InvalidEmailPassword.WriteToResponse(w, nil)
			return
		}

		if !user.IsActive {
			_ = UserNotActive.WriteToResponse(w, nil)
			return
		}

		ttl := 2 * time.Hour

		accessToken, err := utils.GenerateJWT(user.ID, utils.ScopeAuthentication, app.Config.SecretKey, ttl)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Error generating access token",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}

		refreshToken, err := utils.GenerateRefreshToken(user.ID.String(), app.Config.SecretKey)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Error generating refresh token",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}

		response := LoginResponsePayload{
			Token:        accessToken,
			RefreshToken: refreshToken,
		}
		_ = loginSuccess.WriteToResponse(w, response)

		return

	}
}
