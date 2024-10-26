package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

// Sample contact request payload
type ContactRequestPayload struct {
	Phone   string `json:"phone"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

func TestCreateContactHandler_Success(t *testing.T) {
	// Create a mock repository
	cfg, err := state.NewConfig()
	logger := state.New(os.Stdout, state.LevelInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("Config parsing failed")
	}
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)

	if err != nil {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	// Create a sample valid request payload
	requestPayload := ContactRequestPayload{
		Phone:   "1234567890",
		Street:  "123 Test St",
		City:    "Test City",
		State:   "TS",
		ZipCode: "12345",
		Country: "Testland",
	}

	// Create mock user ID and contact ID
	userID := "b7358195-6291-4138-b115-2a046fe848f1"
	// Mock repository behavior to simulate successful contact creation
	mockRepo.On("CreateContact", mock.Anything, mock.Anything).Return(nil)

	// Prepare the request
	body, _ := json.Marshal(requestPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/contacts", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), "userid", userID)
	req = req.WithContext(ctx) // Add user ID to context
	w := httptest.NewRecorder()

	// Call the handler
	handler := httpserver.HandlerCreateContact(appState)
	handler(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestCreateContactHandler_InvalidPayload(t *testing.T) {
	// Create a mock repository
	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Config parsing failed")
	}
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)

	if err != nil {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	// Prepare an invalid JSON payload
	body := []byte(`{"invalid":`)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/contacts", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Call the handler
	handler := httpserver.HandlerCreateContact(appState)
	handler(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	mockRepo.AssertNotCalled(t, "CreateContact")
}
