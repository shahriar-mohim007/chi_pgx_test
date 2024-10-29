package tests

import (
	"bytes"
	"context"
	"encoding/json"
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

func Test_CreateContactHandler_Success(t *testing.T) {

	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		logger.PrintError(err, map[string]string{
			"context": "Error loading env value",
		})
	}

	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	requestPayload := httpserver.ContactRequestPayload{
		Phone:   "1234567890",
		Street:  "123 Test St",
		City:    "Test City",
		State:   "TS",
		ZipCode: "12345",
		Country: "Testland",
	}

	userID := "b7358195-6291-4138-b115-2a046fe848f1"

	mockRepo.On("CreateContact", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(requestPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/contacts", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), "userid", userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler := httpserver.HandlerCreateContact(appState)
	handler(w, req)

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
}

func Test_CreateContactHandler_InvalidPayload(t *testing.T) {

	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		logger.PrintError(err, map[string]string{
			"context": "Error loading env value",
		})
	}

	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	body := []byte(`{"invalid":`)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/contacts", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := httpserver.HandlerCreateContact(appState)
	handler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	mockRepo.AssertNotCalled(t, "CreateContact")
}
