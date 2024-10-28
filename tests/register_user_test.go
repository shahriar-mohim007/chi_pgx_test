package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
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
	Message string                   `json:"message"`
	Data    RegistrationResponseData `json:"data"`
}

type RegistrationResponseData struct {
	ActivateToken string `json:"activate_token"`
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
	//
	//// Mock the utils.HashPassword function
	//utils.HashPassword = func(password string) (string, error) {
	//	return "hashedPassword", nil
	//}

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
	})

	//t.Run("Invalid JSON", func(t *testing.T) {
	//	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte("invalid json")))
	//	w := httptest.NewRecorder()
	//
	//	http.HandlerFunc(HandleRegisterUser(appState)).ServeHTTP(w, req)
	//
	//	assert.Equal(t, http.StatusBadRequest, w.Code)
	//})
	//
	//t.Run("User Already Exists", func(t *testing.T) {
	//	mockRepo.On("GetUserByEmail", mock.Anything, validRequest.Email).Return(&repository.User{}, nil)
	//
	//	payload, _ := json.Marshal(validRequest)
	//	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	//	w := httptest.NewRecorder()
	//
	//	http.HandlerFunc(HandleRegisterUser(appState)).ServeHTTP(w, req)
	//
	//	assert.Equal(t, http.StatusConflict, w.Code)
	//})
	//
	//t.Run("Failed to Hash Password", func(t *testing.T) {
	//	utils.HashPassword = func(password string) (string, error) {
	//		return "", fmt.Errorf("hashing error")
	//	}
	//
	//	mockRepo.On("GetUserByEmail", mock.Anything, validRequest.Email).Return(nil, sql.ErrNoRows)
	//
	//	payload, _ := json.Marshal(validRequest)
	//	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	//	w := httptest.NewRecorder()
	//
	//	http.HandlerFunc(HandleRegisterUser(appState)).ServeHTTP(w, req)
	//
	//	assert.Equal(t, http.StatusInternalServerError, w.Code)
	//})
	//
	//t.Run("Failed to Create User", func(t *testing.T) {
	//	mockRepo.On("GetUserByEmail", mock.Anything, validRequest.Email).Return(nil, sql.ErrNoRows)
	//	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*repository.User")).Return(fmt.Errorf("creation error"))
	//
	//	payload, _ := json.Marshal(validRequest)
	//	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	//	w := httptest.NewRecorder()
	//
	//	http.HandlerFunc(HandleRegisterUser(appState)).ServeHTTP(w, req)
	//
	//	assert.Equal(t, http.StatusInternalServerError, w.Code)
	//})

	// Add more tests as needed...
}
