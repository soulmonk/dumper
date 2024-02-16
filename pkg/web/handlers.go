package web

import (
	"embed"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
	"math/rand"
	"net/http"
	"soulmonk/dumper/pkg/db/ideas"
	"strings"
)

//go:embed assets/*
var f embed.FS

func setupRouter(querier ideas.Querier) *gin.Engine {
	r := gin.New()
	// enable templ engine for gin
	r.HTMLRender = Default
	// Add the sloggin middleware to all routes.
	// The middleware will log all requests attributes.
	r.Use(sloggin.New(slog.Default()))
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		//testIdeas := gin.H{
		//	"name": "Hello, World!",
		//}
		//sendResponse(c, http.StatusOK, testIdeas, web.Hello)
		c.HTML(http.StatusOK, "index.html", Main("Tes"))
	})

	r.GET("/api/ping", ping)

	r.GET("/ideas", getGetIdeasHandler(querier))
	r.GET("/ideas/random", getRandomIdeasHandler(querier))
	r.POST("/ideas", getCreateIdeaHandler(querier))
	r.POST("/ideas/:id/done", getDoneIdeaHandler(querier))

	r.GET("/favicon.ico", func(c *gin.Context) {
		file, _ := f.ReadFile("assets/favicon.ico")
		c.Data(
			http.StatusOK,
			"image/x-icon",
			file,
		)
	})
	return r
}

func getRandomIdeasHandler(querier ideas.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		ids, err := querier.GetIdsOfActiveIdeas(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result, err := querier.GetIdea(c, ids[rand.Intn(len(ids))])
		c.JSON(http.StatusOK, result)
	}
}

type IdeaId struct {
	ID int64 `uri:"id" binding:"required"`
}

func getDoneIdeaHandler(querier ideas.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		var idea IdeaId
		if err := c.ShouldBindUri(&idea); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := querier.DoneIdea(c, idea.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func sendResponse(c *gin.Context, code int, result map[string]any, component func(data map[string]any) templ.Component) {
	accept := c.GetHeader("Accept")
	if accept == "application/json" {
		c.JSON(code, result)
		return
	} else if accept == "text/html" {
		c.HTML(code, "", component(result))
		return
	}
}

func isHtmlResponse(c *gin.Context) bool {
	accept := c.GetHeader("Accept")
	return strings.Contains(accept, "text/html")
}

func isHasHXTarget(c *gin.Context) bool {
	target := c.GetHeader("HX-Target")
	return target != ""
}

func getGetIdeasHandler(querier ideas.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := querier.ListIdeas(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if isHtmlResponse(c) {
			c.HTML(http.StatusOK, "", IdeasList(result, isHasHXTarget(c)))
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
func getCreateIdeaHandler(querier ideas.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		var idea ideas.CreateIdeaParams
		if err := c.ShouldBindJSON(&idea); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := querier.CreateIdea(c, idea)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

type pingResponse struct {
	Status string `bson:"status" json:"status"`
}

func ping(c *gin.Context) {
	var data = pingResponse{"ok"}
	c.IndentedJSON(http.StatusOK, data)
}
