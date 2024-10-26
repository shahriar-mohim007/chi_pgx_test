package tests

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go_chi_pgx/cmd/httpserver"
	"go_chi_pgx/mocks"
	"go_chi_pgx/state"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDeleteContactHandler(t *testing.T) {
	// Create a mock repository and state
	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		t.Fatalf("Config parsing failed: %v", err)
	}
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	// Create a new router for handling requests
	r := chi.NewRouter()
	r.Delete("/contacts/{id}", httpserver.HandlerDeleteContactByID(appState))

	t.Run("Invalid Contact ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/contacts/invalid-id", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Invalid contact ID")
	})

	t.Run("Contact Not Found", func(t *testing.T) {
		contactID, _ := uuid.NewV4() // Generate a new UUID
		mockRepo.On("DeleteContactByID", mock.Anything, contactID).Return(sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodDelete, "/contacts/"+contactID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Contact not found")
	})

	t.Run("Deletion Failed", func(t *testing.T) {
		contactID, _ := uuid.NewV4() // Generate a new UUID
		mockRepo.On("DeleteContactByID", mock.Anything, contactID).Return(errors.New("db error"))

		req := httptest.NewRequest(http.MethodDelete, "/contacts/"+contactID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Failed to delete contact")
	})

	t.Run("Successful Deletion", func(t *testing.T) {
		contactID, _ := uuid.NewV4() // Generate a new UUID
		mockRepo.On("DeleteContactByID", mock.Anything, contactID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/contacts/"+contactID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
		assert.Empty(t, w.Body.String()) // Ensure the body is empty
		mockRepo.AssertCalled(t, "DeleteContactByID", mock.Anything, contactID)
	})
}
