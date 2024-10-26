package tests

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go_chi_pgx/cmd/httpserver"
	"go_chi_pgx/mocks"
	"go_chi_pgx/repository"
	"go_chi_pgx/state"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetAllContactsHandler(t *testing.T) {
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
	r.Get("/contacts", httpserver.HandlerGetAllContacts(appState))

	// Mock user ID for context
	userID := uuid.Must(uuid.NewV4()).String()

	t.Run("Successful Fetch", func(t *testing.T) {
		// Mocking repository behavior
		contactID, _ := uuid.NewV4()
		mockRepo.On("GetAllContacts", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return([]repository.Contact{
			{
				ID:      contactID,
				Phone:   "123-456-7890",
				Street:  "123 Main St",
				City:    "Sample City",
				State:   "Sample State",
				ZipCode: "12345",
				Country: "Sample Country",
			},
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userid", userID)) // Set the user ID in the context
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Check response content here as needed
	})

	t.Run("Error Parsing UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userid", "invalid-uuid")) // Invalid UUID
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Internal server error")
	})

	t.Run("No Contacts Found", func(t *testing.T) {
		mockRepo.On("GetAllContacts", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return([]repository.Contact{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userid", userID))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		// Check the response body for empty contacts if needed
	})
}
