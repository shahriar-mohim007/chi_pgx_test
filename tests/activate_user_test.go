package tests

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
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

	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		logger.PrintError(err, map[string]string{
			"context": "Error loading env value",
		})
	}

	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	userID := "b7358195-6291-4138-b115-2a046fe848f1"
	claims := utils.Claims{UserID: uuid.FromStringOrNil(userID)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.SecretKey))

	mockRepo.On("ActivateUserByID", mock.Anything, claims.UserID).Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/activate?token="+tokenString, nil)
	w := httptest.NewRecorder()

	handler := httpserver.HandleActivateUser(appState)
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	mockRepo.AssertCalled(t, "ActivateUserByID", mock.Anything, claims.UserID)
}

func Test_HandleActivateUser_InvalidToken(t *testing.T) {

	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		logger.PrintError(err, map[string]string{
			"context": "Error loading env value",
		})
	}

	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	userID := "b7358195-6291-4138-b115-2a046fe848f1"
	claims := utils.Claims{UserID: uuid.FromStringOrNil(userID)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	invalidTokenString, _ := token.SignedString([]byte("invalid_secret"))

	req := httptest.NewRequest(http.MethodGet, "/activate?token="+invalidTokenString, nil)
	w := httptest.NewRecorder()

	handler := httpserver.HandleActivateUser(appState)
	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Invalid token")

	mockRepo.AssertNotCalled(t, "ActivateUserByID")
}

func Test_HandleActivateUser_ActivationError(t *testing.T) {

	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		logger.PrintError(err, map[string]string{
			"context": "Error loading env value",
		})
	}

	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	userID := "b7358195-6291-4138-b115-2a046fe848f1"
	claims := utils.Claims{UserID: uuid.FromStringOrNil(userID)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.SecretKey))

	mockRepo.On("ActivateUserByID", mock.Anything, claims.UserID).Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/activate?token="+tokenString, nil)
	w := httptest.NewRecorder()

	handler := httpserver.HandleActivateUser(appState)
	handler(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Internal server error")

	mockRepo.AssertCalled(t, "ActivateUserByID", mock.Anything, claims.UserID)
}
