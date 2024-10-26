package tests

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

func TestGetContactByIDHandler_Success(t *testing.T) {
	// Create a mock repository
	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Config parsing failed")
	}
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	// Create a sample contact for the mock repository
	contactID, _ := uuid.NewV4()
	mockContact := repository.Contact{
		ID:      contactID,
		Phone:   "123-456-7890",
		Street:  "123 Main St",
		City:    "Sample City",
		State:   "Sample State",
		ZipCode: "12345",
		Country: "Sample Country",
	}

	// Mock repository behavior to simulate successful contact retrieval
	mockRepo.On("GetContactByID", mock.Anything, contactID).Return(mockContact, nil)

	// Prepare the request with the contact ID as a URL parameter
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/contacts/%s", contactID), nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler := httpserver.HandlerGetContactByID(appState)
	handler(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Optionally check the response body if needed
	// (You can implement a specific method to read the response if necessary)

	mockRepo.AssertExpectations(t)
}

func TestGetContactByIDHandler_InvalidID(t *testing.T) {
	// Create a mock repository
	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Config parsing failed")
	}
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	// Prepare a request with an invalid UUID
	req := httptest.NewRequest(http.MethodGet, "/contacts/invalid-uuid", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler := httpserver.HandlerGetContactByID(appState)
	handler(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Invalid contact ID") // Ensure error message is present
}

func TestGetContactByIDHandler_ContactNotFound(t *testing.T) {
	// Create a mock repository
	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Config parsing failed")
	}
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	// Create a sample contact ID
	contactID, _ := uuid.NewV4()

	// Mock repository behavior to simulate "contact not found"
	mockRepo.On("GetContactByID", mock.Anything, contactID).Return(repository.Contact{}, errors.New("no contact found"))

	// Prepare the request with the contact ID as a URL parameter
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/contacts/%s", contactID), nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler := httpserver.HandlerGetContactByID(appState)
	handler(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Contact not found") // Ensure error message is present
}
