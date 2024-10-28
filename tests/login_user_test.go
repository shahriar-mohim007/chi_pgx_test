package tests

import (
	"bytes"
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

func TestHandleLogin(t *testing.T) {
	// Initialize dependencies
	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		t.Fatalf("Config parsing failed: %v", err)
	}
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	// Define router and handler
	r := chi.NewRouter()
	r.Post("/login", httpserver.HandleLogin(appState))

	validRequest := httpserver.LoginRequestPayload{
		Email:    "john.doe@example.com",
		Password: "securepassword",
	}

	t.Run("Successful Login", func(t *testing.T) {
		// Set up mock for a successful login
		//mockRepo.On("GetUserByEmail", mock.Anything, validRequest.Email).
		//
		//Return(&repository.User{
		//		Email:    "john.doe@example.com",
		//		Password: "hashedpassword",
		//		IsActive: true,
		//	}, nil)
		//
		//// Mock password check and token generation
		//utils.CheckPasswordHash = func(hashedPassword, password string) bool {
		//	return true // Simulate successful password check
		//}
		//utils.GenerateJWT = func(userID uuid.UUID, scope string, secret string, ttl time.Duration) (string, error) {
		//	return "test-access-token", nil
		//}
		//utils.GenerateRefreshToken = func(userID string, secret string) (string, error) {
		//	return "test-refresh-token", nil
		//}
		//
		//payload, _ := json.Marshal(validRequest)
		//req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
		//w := httptest.NewRecorder()
		//r.ServeHTTP(w, req)
		//
		//assert.Equal(t, http.StatusOK, w.Code)
		//var resp httpserver.LoginResponsePayload
		//_ = json.Unmarshal(w.Body.Bytes(), &resp)
		//assert.Equal(t, "test-access-token", resp.Token)
		//assert.Equal(t, "test-refresh-token", resp.RefreshToken)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		// Set up mock for invalid password scenario
		//mockRepo.On("GetUserByEmail", mock.Anything, validRequest.Email).
		//	Return(&repository.User{
		//		Email:    "john.doe@example.com",
		//		Password: "hashedpassword",
		//		IsActive: true,
		//	}, nil)
		//
		//// Mock a failed password check
		//utils.CheckPasswordHash = func(hashedPassword, password string) bool {
		//	return false // Simulate password mismatch
		//}
		//
		//payload, _ := json.Marshal(validRequest)
		//req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
		//w := httptest.NewRecorder()
		//r.ServeHTTP(w, req)
		//
		//assert.Equal(t, http.StatusUnauthorized, w.Code)
		//assert.Contains(t, w.Body.String(), "Invalid email or password")
	})

	t.Run("Inactive User", func(t *testing.T) {
		// Set up mock for inactive user
		mockRepo.On("GetUserByEmail", mock.Anything, validRequest.Email).
			Return(&repository.User{
				Email:    "john.doe@example.com",
				Password: "hashedpassword",
				IsActive: false,
			}, nil)

		payload, _ := json.Marshal(validRequest)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "User account is not active")
	})
}
