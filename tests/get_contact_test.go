package tests

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
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

func TestGetContactByIDHandler(t *testing.T) {

	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		t.Fatalf("Config parsing failed: %v", err)
	}
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	r := chi.NewRouter()
	r.Get("/contacts/{id}", httpserver.HandlerGetContactByID(appState))

	t.Run("Successful Fetch", func(t *testing.T) {
		contactID := uuid.Must(uuid.NewV4())
		mockContact := &repository.ContactWithUserResponse{
			ContactID: contactID,
			Phone:     "123-456-7890",
			Street:    "123 Main St",
			City:      "Sample City",
			State:     "Sample State",
			ZipCode:   "12345",
			Country:   "Sample Country",
			UserName:  "Mohim",
			UserEmail: "mohim@example.com",
		}

		// Mock repository behavior to return ContactWithUserResponse type
		mockRepo.On("GetContactByID", mock.Anything, contactID).Return(mockContact, nil)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/contacts/%s", contactID), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/contacts/invalid-uuid", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Invalid contact ID")
	})

	t.Run("Contact Not Found", func(t *testing.T) {
		contactID := uuid.Must(uuid.NewV4())

		// Mock repository behavior to simulate "contact not found"
		mockRepo.On("GetContactByID", mock.Anything, contactID).Return(repository.Contact{}, errors.New("no contact found"))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/contacts/%s", contactID), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Contact Not found")
	})
}
