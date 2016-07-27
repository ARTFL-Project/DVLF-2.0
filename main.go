package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

// Template struct
type Template struct {
	templates *template.Template
}

// Render template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.Delims("<<", ">>").ExecuteTemplate(w, name, data)
}

func index(c echo.Context) error {
	indexByte, _ := ioutil.ReadFile("public/index.html")
	index := string(indexByte)
	return c.HTML(http.StatusOK, index)
}

func main() {
	// Echo instance
	e := echo.New()

	e.SetDebug(true)

	e.Static("/static", "public/static")
	e.Static("/app", "public/app")

	// e.File("/favicon.ico", "images/favicon.ico")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// t := &Template{
	// 	templates: template.Must(template.ParseFiles("public/index.html")),
	// }
	// e.SetRenderer(t)

	// Route => handler
	e.GET("/", index)

	// Start server
	e.Run(standard.New(":8080"))
}
