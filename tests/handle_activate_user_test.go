package tests

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go_chi_pgx/cmd/httpserver"
	"go_chi_pgx/mocks"
	"go_chi_pgx/state"
	utils "go_chi_pgx/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandleActivateUser_Success(t *testing.T) {
	// Create a mock repository and state
	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Config parsing failed")
	}
	cfg.SecretKey = "secret"
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	// Create a valid JWT token for testing
	userID := "b7358195-6291-4138-b115-2a046fe848f1"
	claims := utils.Claims{UserID: uuid.FromStringOrNil(userID)} // Convert to UUID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.SecretKey))

	// Mock the repository behavior for a successful user activation
	mockRepo.On("ActivateUserByID", mock.Anything, claims.UserID).Return(nil)

	// Create the request
	req := httptest.NewRequest(http.MethodGet, "/activate?token="+tokenString, nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler := httpserver.HandleActivateUser(appState)
	handler(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify that the mock repository method was called with the correct parameters
	mockRepo.AssertCalled(t, "ActivateUserByID", mock.Anything, claims.UserID)
}

func TestHandleActivateUser_InvalidToken(t *testing.T) {
	// Create a mock repository and state
	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Config parsing failed")
	}
	cfg.SecretKey = "secret"
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	// Create an invalid token (incorrect signature)
	userID := "b7358195-6291-4138-b115-2a046fe848f1"
	claims := utils.Claims{UserID: uuid.FromStringOrNil(userID)} // Convert to UUID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	invalidTokenString, _ := token.SignedString([]byte("invalid_secret"))

	// Create the request
	req := httptest.NewRequest(http.MethodGet, "/activate?token="+invalidTokenString, nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler := httpserver.HandleActivateUser(appState)
	handler(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Invalid token")

	// Verify that the repository method was not called
	mockRepo.AssertNotCalled(t, "ActivateUserByID")
}

func TestHandleActivateUser_ActivationError(t *testing.T) {
	// Initialize the necessary components
	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Config parsing failed")
	}
	cfg.SecretKey = "secret"
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	// Create a valid JWT token for testing
	userID := "b7358195-6291-4138-b115-2a046fe848f1"
	claims := utils.Claims{UserID: uuid.FromStringOrNil(userID)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.SecretKey))

	// Mock the repository behavior to return an error
	mockRepo.On("ActivateUserByID", mock.Anything, claims.UserID).Return(errors.New("db error"))

	// Create the request
	req := httptest.NewRequest(http.MethodGet, "/activate?token="+tokenString, nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler := httpserver.HandleActivateUser(appState)
	handler(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Internal server error")

	// Verify that the mock repository method was called
	mockRepo.AssertCalled(t, "ActivateUserByID", mock.Anything, claims.UserID)
}
