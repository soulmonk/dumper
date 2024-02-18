package web

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5/pgtype"
	"io"
	"soulmonk/dumper/pkg/db/ideas"
	"testing"
)

var testIdeas = []ideas.Ideas{
	{
		Title: pgtype.Text{"Title1", true},
		Body:  pgtype.Text{"Body1", true},
	},
	{
		Title: pgtype.Text{"Title2", true},
		Body:  pgtype.Text{"Body2", true},
	},
}

func TestIdeasList(t *testing.T) {

	r, w := io.Pipe()
	go func() {
		ideasList(testIdeas).Render(context.Background(), w)
		w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}

	//if doc.Find(`[data-testid="ideasListIdea"]`).Length() == 0 {
	//	t.Error("expected data-testid attribute to be rendered, but it wasn't")
	//}

	// Expect both posts to be rendered.
	if actualIdeaCount := doc.Find(`[data-testid="ideasListIdea"]`).Length(); actualIdeaCount != len(testIdeas) {
		t.Fatalf("expected %d posts to be rendered, found %d", len(testIdeas), actualIdeaCount)
	}

	// Expect the posts to contain the author name.
	doc.Find(`[data-testid="ideasListIdea"]`).Each(func(index int, sel *goquery.Selection) {
		expectedTitle := testIdeas[index].Title.String
		if actualTitle := sel.Find(`[data-testid="ideasListIdeaTitle"]`).Text(); actualTitle != expectedTitle {
			t.Errorf("expected name %q, got %q", actualTitle, expectedTitle)
		}
		expectedBody := testIdeas[index].Body.String
		if actualBody := sel.Find(`[data-testid="ideasListIdeaBody"]`).Text(); actualBody != expectedBody {
			t.Errorf("expected author %q, got %q", actualBody, expectedBody)
		}
	})
}

func TestIdeasListWithoutTarget(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		IdeasList(testIdeas, false).Render(context.Background(), w)
		w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}

	// Expect the container to be rendered.
	if doc.Find(`[data-testid="container"]`).Length() != 1 {
		t.Error("expected container to be rendered, but it wasn't")
	}
	// Expect the navigation to be rendered.
	if doc.Find(`[data-testid="nav-bar"]`).Length() != 1 {
		t.Error("expected nav to be rendered, but it wasn't")
	}
	// Expect both posts to be rendered.
	if actualIdeaCount := doc.Find(`[data-testid="ideasListIdea"]`).Length(); actualIdeaCount != len(testIdeas) {
		t.Fatalf("expected %d posts to be rendered, found %d", len(testIdeas), actualIdeaCount)
	}
}

func TestIdeasListWithTarget(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		IdeasList(testIdeas, true).Render(context.Background(), w)
		w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}

	// Expect the container not to be rendered.
	if doc.Find(`[data-testid="container"]`).Length() != 0 {
		t.Error("expected data-testid attribute to be rendered, but it wasn't")
	}
	// Expect the navigation not to be rendered.
	if doc.Find(`[data-testid="nav-bar"]`).Length() != 0 {
		t.Error("expected nav to be rendered, but it wasn't")
	}
	// Expect both posts to be rendered.
	if actualIdeaCount := doc.Find(`[data-testid="ideasListIdea"]`).Length(); actualIdeaCount != len(testIdeas) {
		t.Fatalf("expected %d posts to be rendered, found %d", len(testIdeas), actualIdeaCount)
	}
}
