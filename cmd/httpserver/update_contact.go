package httpserver

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"go_chi_pgx/repository"
	"go_chi_pgx/state"
	"net/http"
)

func HandlerPatchContactByID(app *state.State) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		contactID := chi.URLParam(req, "id")
		uuidContactID, err := uuid.FromString(contactID)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Error parsing contact ID",
			})
			_ = InvalidId.WriteToResponse(w, nil)
			return
		}
		ctx := req.Context()

		requestPayload := ContactRequestPayload{}
		err = json.NewDecoder(req.Body).Decode(&requestPayload)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Invalid JSON",
			})
			_ = ValidDataNotFound.WriteToResponse(w, nil)
			return
		}

		contact, err := app.Repository.GetContactByID(ctx, uuidContactID)
		if err != nil {
			_ = NotFound.WriteToResponse(w, nil)
			return
		}

		if requestPayload.Phone != "" {
			contact.Phone = requestPayload.Phone
		}
		if requestPayload.Street != "" {
			contact.Street = requestPayload.Street
		}
		if requestPayload.City != "" {
			contact.City = requestPayload.City
		}
		if requestPayload.State != "" {
			contact.State = requestPayload.State
		}
		if requestPayload.ZipCode != "" {
			contact.ZipCode = requestPayload.ZipCode
		}
		if requestPayload.Country != "" {
			contact.Country = requestPayload.Country
		}

		updatedContact := repository.Contact{
			Phone:   contact.Phone,
			Street:  contact.Street,
			City:    contact.City,
			State:   contact.State,
			ZipCode: contact.ZipCode,
			Country: contact.Country,
		}

		err = app.Repository.PatchContactByID(ctx, uuidContactID, &updatedContact)
		if err != nil {
			_ = InternalError.WriteToResponse(w, err)
			return
		}
		response := ContactResponse{
			ID:      contactID,
			Phone:   contact.Phone,
			Street:  contact.Street,
			City:    contact.City,
			State:   contact.State,
			ZipCode: contact.ZipCode,
			Country: contact.Country,
		}

		_ = ContactUpdated.WriteToResponse(w, response)
		return
	}
}
