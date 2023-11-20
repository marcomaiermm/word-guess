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

type Symbol struct {
	Letter  string
	Correct bool
}

type GuessRow struct {
	Symbol []Symbol
}

type IndexPage struct {
	Page
	UUID string
	Word string
	Rows []GuessRow
	Won  bool
	Lost bool
}

func createGuess(word string, symbols []string, _ string) *[]GuessRow {
	rows := make([]GuessRow, len(symbols))

	for i := range rows {
		cols := make([]Symbol, len(word))
		for j := range cols {
			char := ""
			correct := false
			if len(symbols[i]) > j {
				char = string(symbols[i][j])
				correct = strings.EqualFold(strings.TrimSpace(char), strings.TrimSpace(string(word[j])))
			}
			symbol := Symbol{
				Letter:  char,
				Correct: correct,
			}
			cols[j] = symbol
		}
		rows[i] = GuessRow{
			Symbol: cols,
		}
	}

	return &rows
}

func setupSession(c echo.Context, path string) *sessions.Session {
	sess, _ := session.Get("session", c)

	sess.Options = &sessions.Options{
		Path:     path,
		MaxAge:   604800,
		HttpOnly: true,
	}

	sess.Save(c.Request(), c.Response())

	return sess
}

func Index(c echo.Context) error {
	setupSession(c, "/")

	game, err := database.CreateGame(6)
	if err != nil {
		fmt.Println(err)
		return c.NoContent(500)
	}

	symbols := strings.Split(game.Symbols, ",")[:game.Rows]

	rows := createGuess(game.Word, symbols, "")

	return c.Render(200, "index.html", IndexPage{
		Page: Page{
			Title: "Play Word Guess!",
		},
		UUID: game.UUID,
		Word: game.Word,
		Rows: *rows,
		Won:  false,
		Lost: false,
	})
}

func NewGameFragment(c echo.Context) error {
	game, err := database.CreateGame(6)
	if err != nil {
		fmt.Println(err)
		return c.NoContent(500)
	}

	symbols := strings.Split(game.Symbols, ",")[:game.Rows]

	rows := createGuess(game.Word, symbols, "")

	return c.Render(200, "content.html", GuessResponse{
		UUID: game.UUID,
		Rows: *rows,
		Won:  false,
		Lost: false,
		Word: game.Word,
	})
}

type GuessResponse struct {
	UUID string
	Word string
	Rows []GuessRow
	Won  bool
	Lost bool
}

func Guess(c echo.Context) error {
	setupSession(c, "/")

	// uuid from /game/:uuid
	uuid := c.Param("id")
	// guess from payload
	guess := c.FormValue("guess")

	if uuid == "" {
		return c.JSON(400, "uuid is required")
	}

	if guess == "" {
		return c.JSON(400, "guess is required")
	}

	updatedGame, err := database.UpdateGame(database.UpdateGameParams{
		UUID:    uuid,
		Symbols: guess,
	})
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, err)
	}

	won := false
	lost := false
	symbols := strings.Split(updatedGame.Symbols, ",")[:updatedGame.Rows]

	var guesses int
	for _, symbol := range symbols {
		if symbol != "" {
			guesses++
		}
	}

	if strings.EqualFold(updatedGame.Word, strings.TrimSpace(guess)) {
		won = true
	} else if guesses == updatedGame.Rows {
		lost = true
	}

	rows := createGuess(updatedGame.Word, symbols, guess)

	return c.Render(200, "content.html", GuessResponse{
		UUID: updatedGame.UUID,
		Rows: *rows,
		Won:  won,
		Lost: lost,
		Word: updatedGame.Word,
	})
}
