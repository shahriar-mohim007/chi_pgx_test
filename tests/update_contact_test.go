package tests

import (
	"bytes"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
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

func TestUpdateContactHandler(t *testing.T) {
	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		t.Fatalf("Config parsing failed: %v", err)
	}
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	r := chi.NewRouter()
	r.Patch("/contacts/{id}", httpserver.HandlerPatchContactByID(appState))

	t.Run("Invalid Contact ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/contacts/invalid-id", nil)
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

		mockRepo.On("GetContactByID", mock.Anything, contactID).Return(nil, sql.ErrNoRows)
		reqBody := bytes.NewBuffer([]byte(`{"name": "Updated Name", "phone": "123-456-7890"}`))
		req := httptest.NewRequest(http.MethodPatch, "/contacts/"+contactID.String(), reqBody)
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

	t.Run("Update Successful", func(t *testing.T) {
		contactID, _ := uuid.NewV4()
		mockContact := &repository.ContactWithUserResponse{
			ContactID: contactID,
			Phone:     "123-456-7890",
			Street:    "123 Main St",
			City:      "Sample City",
			State:     "Sample State",
			ZipCode:   "12345",
			Country:   "Sample Country",
			UserName:  "Mohim",
			UserEmail: "mohim@example.com",
		}
		mockRepo.On("GetContactByID", mock.Anything, contactID).Return(mockContact, nil)
		mockRepo.On("PatchContactByID", mock.Anything, contactID, mock.AnythingOfType("*repository.Contact")).Return(nil)

		reqBody := bytes.NewBuffer([]byte(`{"name": "Updated Name", "phone": "123-456-7890"}`))

		req := httptest.NewRequest(http.MethodPatch, "/contacts/"+contactID.String(), reqBody)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Contacts Updated successfully")
		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})

	t.Run("Update Failed", func(t *testing.T) {
		contactID, _ := uuid.NewV4()
		mockContact := &repository.ContactWithUserResponse{
			ContactID: contactID,
			Phone:     "123-456-7890",
			Street:    "123 Main St",
			City:      "Sample City",
			State:     "Sample State",
			ZipCode:   "12345",
			Country:   "Sample Country",
			UserName:  "Mohim",
			UserEmail: "mohim@example.com",
		}
		mockRepo.On("GetContactByID", mock.Anything, contactID).Return(mockContact, nil)
		mockRepo.On("PatchContactByID", mock.Anything, contactID, mock.AnythingOfType("*repository.Contact")).Return(errors.New("db error"))

		reqBody := bytes.NewBuffer([]byte(`{"name": "Updated Name", "phone": "123-456-7890"}`))

		req := httptest.NewRequest(http.MethodPatch, "/contacts/"+contactID.String(), reqBody)
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

	t.Run("Invalid Payload", func(t *testing.T) {
		contactID, _ := uuid.NewV4()

		body := []byte(`{"invalid":`)
		req := httptest.NewRequest(http.MethodPatch, "/contacts/"+contactID.String(), bytes.NewReader(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "The provided information is invalid. Please recheck and try again.")

		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})
}
