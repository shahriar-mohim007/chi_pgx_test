package httpserver

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"go_chi_pgx/state"
	"net/http"
)

func HandlerDeleteContactByID(app *state.State) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		id := chi.URLParam(req, "id")
		contactID, err := uuid.FromString(id)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Error parsing contact ID",
			})
			_ = InvalidId.WriteToResponse(w, nil)
			return
		}

		ctx := req.Context()
		err = app.Repository.DeleteContactByID(ctx, contactID)

		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				app.Logger.PrintError(err, map[string]string{
					"context": "Contact not found",
				})
				_ = NotFound.WriteToResponse(w, nil)

			} else {
				app.Logger.PrintError(err, map[string]string{
					"context": "Error deleting contact",
				})
				_ = InternalError.WriteToResponse(w, nil)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}
}
