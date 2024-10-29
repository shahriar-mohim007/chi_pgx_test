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

type LoginResponse struct {
	Message string                          `json:"message"`
	Data    httpserver.LoginResponsePayload `json:"data"`
}

func TestHandleLogin(t *testing.T) {

	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		t.Fatalf("Config parsing failed: %v", err)
	}
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	r := chi.NewRouter()
	r.Post("/login", httpserver.HandleLogin(appState))

	validRequest := httpserver.LoginRequestPayload{
		Email:    "john.doe@example.com",
		Password: "securepassword",
	}

	tests := []struct {
		name             string
		mockUser         *repository.User
		mockError        error
		expectedCode     int
		expectedMessage  string
		expectValidToken bool
	}{
		{
			name: "Successful Login",
			mockUser: &repository.User{
				Email:    "john.doe@example.com",
				Password: "$2a$10$OflXl1si7Vo2ZAEDI6jWDulW17Nq/8C9mME1gZ6w19lo1Ix3j5f4K",
				IsActive: true,
			},
			expectedCode:     http.StatusOK,
			expectValidToken: true,
		},
		{
			name: "Invalid Password",
			mockUser: &repository.User{
				Email:    "john.doe@example.com",
				Password: "$hashedpasgrgweqw1d2eghhthtrhtrswordbdheththtr",
				IsActive: true,
			},
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: "Invalid email or password",
		},
		{
			name: "Inactive User",
			mockUser: &repository.User{
				Email:    "john.doe@example.com",
				Password: "$2a$10$OflXl1si7Vo2ZAEDI6jWDulW17Nq/8C9mME1gZ6w19lo1Ix3j5f4K",
				IsActive: false,
			},
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: "User not active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.On("GetUserByEmail", mock.Anything, validRequest.Email).Return(tt.mockUser, tt.mockError)

			payload, _ := json.Marshal(validRequest)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK && tt.expectValidToken {
				var resp LoginResponse
				_ = json.Unmarshal(w.Body.Bytes(), &resp)

				jwtPattern := `^[A-Za-z0-9-_]+?\.[A-Za-z0-9-_]+?\.[A-Za-z0-9-_]+$`
				assert.Regexp(t, jwtPattern, resp.Data.Token, "Token should match JWT format")
				assert.Regexp(t, jwtPattern, resp.Data.RefreshToken, "RefreshToken should match JWT format")
			} else {

				assert.Contains(t, w.Body.String(), tt.expectedMessage)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
