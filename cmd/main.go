package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/marcomaiermm/word-guess/pkg/database"
	"github.com/marcomaiermm/word-guess/pkg/pages"
)

type TemplateRenderer struct {
	template *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.template.ExecuteTemplate(w, name, data)
}

// allowOrigin takes the origin as an argument and returns true if the origin
// is allowed or false otherwise.
func allowOrigin(origin string) (bool, error) {
	return regexp.MatchString(`^https:\/\/marcomaier\.dev$`, origin)
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
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: allowOrigin,
		AllowMethods:    []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "cookie:_csrf",
		CookiePath:     "/",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Static("/dist", "dist")

	e.GET("/", pages.Index)
	e.PATCH("/game/:id", pages.Guess)

	e.Logger.Fatal(e.Start(":42069"))
}
