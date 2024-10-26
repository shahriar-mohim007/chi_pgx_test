package httpserver

import (
	"fmt"
	"github.com/gofrs/uuid"
	"go_chi_pgx/state"
	"net/http"
	"strconv"
)

type ContactResponse struct {
	ID      string `json:"id"`
	Phone   string `json:"phone"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type ContactsResponse struct {
	TotalCount int               `json:"total_count"`
	Next       string            `json:"next"`
	Previous   string            `json:"previous"`
	Contacts   []ContactResponse `json:"contacts"`
}

func HandlerGetAllContacts(app *state.State) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		userID, _ := GetUserIDFromContext(ctx)
		uuID, err := uuid.FromString(userID)

		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"Context": "Error parsing UUID",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}

		limitParam := req.URL.Query().Get("limit")
		offsetParam := req.URL.Query().Get("offset")

		limit := 10
		offset := 0
		if limitParam != "" {
			limit, err = strconv.Atoi(limitParam)
			if err != nil {
				app.Logger.PrintError(fmt.Errorf("invalid limit value"), map[string]string{
					"context": "pagination",
				})
				_ = BadRequestError.WriteToResponse(w, nil)
				return
			}
		}

		if offsetParam != "" {
			offset, err = strconv.Atoi(offsetParam)
			if err != nil {
				app.Logger.PrintError(fmt.Errorf("invalid offset value"), map[string]string{
					"context": "pagination",
				})
				_ = BadRequestError.WriteToResponse(w, nil)
				return
			}
		}

		contacts, err := app.Repository.GetAllContacts(ctx, uuID, limit, offset)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"Context": "Error fetching contacts",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}

		totalCount, err := app.Repository.GetContactsCount(ctx, uuID)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"Context": "Error fetching contacts count",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}
		scheme := req.Header.Get("X-Forwarded-Proto")
		if scheme == "" {
			scheme = "http"
		}

		// Generate next and previous URLs
		baseURL := scheme + "://" + req.Host + req.URL.Path
		nextOffset := offset + limit
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}

		nextURL := ""
		prevURL := ""

		// Create next URL if more records are available
		if nextOffset < totalCount {
			nextURL = fmt.Sprintf("%s?limit=%d&offset=%d", baseURL, limit, nextOffset)
		}

		// Create previous URL if offset is greater than 0
		if offset > 0 {
			prevURL = fmt.Sprintf("%s?limit=%d&offset=%d", baseURL, limit, prevOffset)
		}

		// Create response
		var contactResponses []ContactResponse
		for _, contact := range contacts {
			contactResponses = append(contactResponses, ContactResponse{
				ID:      contact.ID.String(),
				Phone:   contact.Phone,
				Street:  contact.Street,
				City:    contact.City,
				State:   contact.State,
				ZipCode: contact.ZipCode,
				Country: contact.Country,
			})
		}

		response := ContactsResponse{
			Contacts:   contactResponses,
			TotalCount: totalCount,
			Next:       nextURL,
			Previous:   prevURL,
		}

		_ = ContactRetrieved.WriteToResponse(w, response)
		return

	}
}
