package pages

import (
	"github.com/labstack/echo/v4"
)

type Page struct {
	Title     string
	CreatedAt string
	UpdatedAt string
}

type IndexPage struct {
	Page
}

func renderIndexPage(c echo.Context) error {
	return c.Render(200, "index.html", IndexPage{
		Page: Page{
			Title: "Game Merchant",
		},
	})
}

func Index(c echo.Context) error {
	return renderIndexPage(c)
}
