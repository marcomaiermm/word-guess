package main

import (
	"html/template"
	"io"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/marcomaiermm/go-vite-template/pkg/database"
	"github.com/marcomaiermm/go-vite-vemplate/pkg/pages"
)

type TemplateRenderer struct {
	template *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.template.ExecuteTemplate(w, name, data)
}

func main() {
	tmpls, err := template.New("").ParseGlob("public/views/*.html")
	if err != nil {
		log.Fatalf("couldn't initialize templates: %v", err)
	}

	url := os.Getenv("DB_URL")
	if url == "" {
		url = "./database.db"
	}

	err = database.Init(url)
	if err != nil {
		log.Fatalf("couldn't initialize database: %v", err)
	}

	e := echo.New()
	e.Renderer = &TemplateRenderer{
		template: tmpls,
	}

	e.Use(middleware.Logger())
	e.Static("/dist", "dist")

	e.GET("/", pages.Index)

	e.Logger.Fatal(e.Start(":42069"))
}
