package pages

import (
	"fmt"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/marcomaiermm/word-guess/pkg/database"
)

type Page struct {
	Title string
}

type IndexPage struct {
	Page
	Word string
	Rows []Row
}

type Row struct {
	Symbols []string
}

func renderIndexPage(c echo.Context) error {
	sess, _ := session.Get("session", c)

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   604800,
		HttpOnly: true,
	}

	sess.Save(c.Request(), c.Response())

	game, err := database.CreateGame(5)
	if err != nil {
		fmt.Println(err)
		return c.NoContent(500)
	}

	rowSplit := strings.Split(game.Symbols, ",")
	rows := make([]Row, len(rowSplit))

	for i := range rowSplit {
		cols := make([]string, len(game.Word))
		for j := range cols {
			cols[j] = ""
		}
		rows[i] = Row{
			Symbols: cols,
		}
	}

	return c.Render(200, "index.html", IndexPage{
		Page: Page{
			Title: "Play Word Guess!",
		},
		Word: game.Word,
		Rows: rows,
	})
}

func Index(c echo.Context) error {
	return renderIndexPage(c)
}
