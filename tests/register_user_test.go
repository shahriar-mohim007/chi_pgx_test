package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
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

type RegistrationResponsePayload struct {
	Message string                                 `json:"message"`
	Data    httpserver.RegistrationResponsePayload `json:"data"`
}

func TestHandleRegisterUser(t *testing.T) {
	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		t.Fatalf("Config parsing failed: %v", err)
	}
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)
	r := chi.NewRouter()
	r.Post("/users", httpserver.HandleRegisterUser(appState))

	validRequest := httpserver.RegistrationRequestPayload{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "securepassword",
	}

	t.Run("Successful Registration", func(t *testing.T) {

		mockRepo.On("GetUserByEmail", mock.Anything, validRequest.Email).Return((*repository.User)(nil), sql.ErrNoRows)
		mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*repository.User")).Return(nil)

		payload, _ := json.Marshal(validRequest)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(payload))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response RegistrationResponsePayload
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Data.ActivateToken)

		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("invalid json")))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "The provided information is invalid. Please recheck and try again.")
		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})

	t.Run("User Already Exists", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", mock.Anything, validRequest.Email).Return(&repository.User{}, nil)

		payload, _ := json.Marshal(validRequest)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(payload))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "User Already Exist With this Email")
		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})

	t.Run("Failed to Create User", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", mock.Anything, validRequest.Email).Return((*repository.User)(nil), sql.ErrNoRows)
		mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*repository.User")).Return(fmt.Errorf("creation error"))

		payload, _ := json.Marshal(validRequest)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})
}
