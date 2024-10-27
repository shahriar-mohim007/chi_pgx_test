package tests

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go_chi_pgx/cmd/httpserver"
	"go_chi_pgx/mocks"
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
	})
}
