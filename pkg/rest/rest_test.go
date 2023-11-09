package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"soulmonk/dumper/pkg/db"
	"soulmonk/dumper/pkg/db/ideas"
	"strconv"
	"testing"
	"time"
)

type IdeasQuerierMock struct {
	lastId int64
	doneId int64
}

func defaultIdeasMock() IdeasQuerierMock {
	return IdeasQuerierMock{1, 2}
}

func (q IdeasQuerierMock) ListIdeas(ctx context.Context) ([]ideas.Ideas, error) {
	return []ideas.Ideas{
		{
			ID:    1,
			Title: pgtype.Text{"Title", true},
			Body:  pgtype.Text{"Body", true},
		},
	}, nil
}

func (q IdeasQuerierMock) CreateIdea(ctx context.Context, arg ideas.CreateIdeaParams) (ideas.CreateIdeaRow, error) {
	return ideas.CreateIdeaRow{
		ID:        q.lastId,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}, nil
}

func (q IdeasQuerierMock) DoneIdea(ctx context.Context, id int64) (pgtype.Timestamp, error) {
	if q.doneId == id {
		return pgtype.Timestamp{Time: time.Now(), Valid: false}, fmt.Errorf("already done")
	}
	return pgtype.Timestamp{Time: time.Now(), Valid: true}, nil
}

// Required run migration on test db

func TestPing(t *testing.T) {
	//dao := db.GetDao(context.TODO(), "postgres://cuppa:toor@localhost:5432/cuppa-dumper-test")

	router := setupRouter(defaultIdeasMock())
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{
    "status": "ok"
}`, w.Body.String())
}

func TestCreateIdea(t *testing.T) {
	router := setupRouter(defaultIdeasMock())
	w := httptest.NewRecorder()
	jsonBody := `{"title":"Title","body":"Body"}`
	req, _ := http.NewRequest(http.MethodPost, "/ideas", bytes.NewBuffer([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
	var createResponse ideas.CreateIdeaRow
	if err := json.Unmarshal(w.Body.Bytes(), &createResponse); err != nil {
		t.Errorf("can't decode response, %s", w.Body.String())
		return
	}
	assert.NotEmpty(t, createResponse.ID)
}

func TestDoneIdea(t *testing.T) {
	router := setupRouter(defaultIdeasMock())
	w := httptest.NewRecorder()
	jsonBody := `{"title":"Title","body":"Body"}`
	req, _ := http.NewRequest(http.MethodPost, "/ideas", bytes.NewBuffer([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResponse ideas.CreateIdeaRow
	if err := json.Unmarshal(w.Body.Bytes(), &createResponse); err != nil {
		t.Errorf("can't decode response, %s", w.Body.String())
		return
	}
	w = httptest.NewRecorder()
	url := "/ideas/" + strconv.FormatInt(createResponse.ID, 10) + "/done"
	req, _ = http.NewRequest(http.MethodPost, url, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
}

func TestDoneSameIdeaTwice(t *testing.T) {
	dao := db.GetDao(context.TODO(), "postgres://cuppa:toor@localhost:5432/cuppa-dumper-test")
	router := setupRouter(dao.IdeasQuerier)
	w := httptest.NewRecorder()
	jsonBody := `{"title":"Title","body":"Body"}`
	req, _ := http.NewRequest(http.MethodPost, "/ideas", bytes.NewBuffer([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createResponse ideas.CreateIdeaRow
	if err := json.Unmarshal(w.Body.Bytes(), &createResponse); err != nil {
		t.Errorf("can't decode response, %s", w.Body.String())
		return
	}
	w = httptest.NewRecorder()
	url := "/ideas/" + strconv.FormatInt(createResponse.ID, 10) + "/done"
	req, _ = http.NewRequest(http.MethodPost, url, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, url, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, w.Body.String())
}

func TestDoneSameIdeaTwiceNotFound(t *testing.T) {
	router := setupRouter(IdeasQuerierMock{1, 1})

	w := httptest.NewRecorder()
	url := "/ideas/1/done"
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code, w.Body.String())
}
