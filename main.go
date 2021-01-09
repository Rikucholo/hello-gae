package main

import (
	"net/http"
	"io"
	"os"
	"fmt"
	"html/template"
	"github.com/labstack/echo"
	"github.com/ymotongpoo/datemaki"
)
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Renderer = &Template{
        templates: template.Must(template.ParseGlob("front/*.html")),
	}
	
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK,"index",nil)
	})
	e.GET("/date", parseDate)
	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
        e.Logger.Printf("Defaulting to port %s", port)
    }
    e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))}

func parseDate(c echo.Context) error {
	date, err := datemaki.Parse(c.QueryParam("value"))
	if err != nil {
		return c.JSON(http.StatusOK, "日付ではありません")
	}
	return c.JSON(http.StatusOK, date)
}
