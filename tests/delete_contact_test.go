package tests

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
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

func TestDeleteContactHandler(t *testing.T) {

	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		t.Fatalf("Config parsing failed: %v", err)
	}
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	r := chi.NewRouter()
	r.Delete("/contacts/{id}", httpserver.HandlerDeleteContactByID(appState))

	t.Run("Invalid Contact ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/contacts/invalid-id", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Invalid contact ID")
		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})

	t.Run("Contact Not Found", func(t *testing.T) {
		contactID, _ := uuid.NewV4()
		mockRepo.On("DeleteContactByID", mock.Anything, contactID).Return(sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodDelete, "/contacts/"+contactID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Contact Not found")
		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})

	t.Run("Deletion Failed", func(t *testing.T) {
		contactID, _ := uuid.NewV4()
		mockRepo.On("DeleteContactByID", mock.Anything, contactID).Return(errors.New("db error"))

		req := httptest.NewRequest(http.MethodDelete, "/contacts/"+contactID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Internal server error")
		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})

	t.Run("Successful Deletion", func(t *testing.T) {
		contactID, _ := uuid.NewV4()
		mockRepo.On("DeleteContactByID", mock.Anything, contactID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/contacts/"+contactID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
		assert.Empty(t, w.Body.String())
		mockRepo.AssertCalled(t, "DeleteContactByID", mock.Anything, contactID)

		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})
}
