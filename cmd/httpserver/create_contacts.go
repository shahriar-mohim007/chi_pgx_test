package httpserver

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"go_chi_pgx/repository"
	"go_chi_pgx/state"
	"net/http"
)

type ContactRequestPayload struct {
	Phone   string `json:"phone" validate:"required"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

func HandlerCreateContact(app *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		requestPayload := ContactRequestPayload{}
		ctx := req.Context()

		err := json.NewDecoder(req.Body).Decode(&requestPayload)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Invalid JSON",
			})
			_ = ValidDataNotFound.WriteToResponse(w, nil)
			return
		}

		userID, _ := GetUserIDFromContext(ctx)
		uuID, err := uuid.FromString(userID)

		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Error parsing UUID",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}
		ID, err := uuid.NewV4()
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Error generating UUID",
			})
			_ = InternalError.WriteToResponse(w, nil)
		}

		contact := repository.Contact{
			ID:      ID,
			UserID:  uuID,
			Phone:   requestPayload.Phone,
			Street:  requestPayload.Street,
			City:    requestPayload.City,
			State:   requestPayload.State,
			ZipCode: requestPayload.ZipCode,
			Country: requestPayload.Country,
		}

		if err = app.Repository.CreateContact(ctx, &contact); err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Error creating contact",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}
		_ = ContactCreated.WriteToResponse(w, contact)

		return
	}
}
