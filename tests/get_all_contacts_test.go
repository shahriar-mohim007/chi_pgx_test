package tests

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
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

type ContactsResponse struct {
	Message string       `json:"message"`
	Data    ContactsData `json:"data"`
}

type ContactsData struct {
	TotalCount int               `json:"total_count"`
	Next       string            `json:"next"`
	Previous   string            `json:"previous"`
	Contacts   []ContactResponse `json:"contacts"`
}

type ContactResponse struct {
	ID      string `json:"id"`
	Phone   string `json:"phone"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

func Test_GetAllContactsHandler(t *testing.T) {

	logger := state.New(os.Stdout, state.LevelInfo)
	cfg, err := state.NewConfig()
	if err != nil {
		t.Fatalf("Config parsing failed: %v", err)
	}
	mockRepo := new(mocks.MockRepository)
	appState := state.NewState(cfg, mockRepo, logger)

	r := chi.NewRouter()
	r.Get("/contacts", httpserver.HandlerGetAllContacts(appState))

	userID := uuid.Must(uuid.NewV4()).String()

	t.Run("Successful Fetch", func(t *testing.T) {

		contactID, _ := uuid.NewV4()
		contacts := []repository.Contact{
			{
				ID:      contactID,
				Phone:   "123-456-7890",
				Street:  "123 Main St",
				City:    "Sample City",
				State:   "Sample State",
				ZipCode: "12345",
				Country: "Sample Country",
			},
		}
		totalCount := 1

		mockRepo.On("GetAllContacts", mock.Anything, mock.AnythingOfType("uuid.UUID"), 10, 0).Return(contacts, nil)
		mockRepo.On("GetContactsCount", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(totalCount, nil)

		req := httptest.NewRequest(http.MethodGet, "/contacts?limit=10&offset=0", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userid", userID))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response ContactsResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, totalCount, response.Data.TotalCount)
		assert.Len(t, response.Data.Contacts, len(contacts))

		assert.Equal(t, contactID.String(), response.Data.Contacts[0].ID)
		assert.Equal(t, "123-456-7890", response.Data.Contacts[0].Phone)
		assert.Equal(t, "123 Main St", response.Data.Contacts[0].Street)
		assert.Equal(t, "Sample City", response.Data.Contacts[0].City)
		assert.Equal(t, "Sample State", response.Data.Contacts[0].State)
		assert.Equal(t, "12345", response.Data.Contacts[0].ZipCode)
		assert.Equal(t, "Sample Country", response.Data.Contacts[0].Country)

		assert.Equal(t, "", response.Data.Next)     // No next URL as there is only one contact
		assert.Equal(t, "", response.Data.Previous) // No previous URL since offset is 0

		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})

	t.Run("Error Parsing UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userid", "invalid-uuid"))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Invalid user ID")
		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})

	t.Run("No Contacts Found", func(t *testing.T) {
		mockRepo.On("GetAllContacts", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return([]repository.Contact{}, nil)
		mockRepo.On("GetContactsCount", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(0, nil)

		req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userid", userID))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		mockRepo.AssertExpectations(t)
		t.Cleanup(func() {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
		})
	})
}
