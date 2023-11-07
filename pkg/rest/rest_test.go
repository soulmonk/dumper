package rest

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"soulmonk/dumper/pkg/db"
	"testing"
)

// Required run migration on test db

func TestPing(t *testing.T) {
	dao := db.GetDao(context.TODO(), "postgres://cuppa:toor@localhost:5432/cuppa-dumper-test")

	router := setupRouter(dao)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{
    "status": "ok"
}`, w.Body.String())
}

type IdeaCreateParams struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func TestCreateIdea(t *testing.T) {
	dao := db.GetDao(context.TODO(), "postgres://cuppa:toor@localhost:5432/cuppa-dumper-test")
	router := setupRouter(dao)
	w := httptest.NewRecorder()
	jsonBody := `{"title":"Title","body":"Body"}`
	req, _ := http.NewRequest(http.MethodPost, "/ideas", bytes.NewBuffer([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// todo how to compare ID
	// todo how to mock db
	//assert.Equal(t, `{"Title":"Title","Body":"Body"}`, w.Body.String())
}
