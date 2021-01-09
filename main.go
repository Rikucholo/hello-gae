package main

import (
	"net/http"
	"io"
	"os"
	"fmt"
	"html/template"
	"context"
	"log"
	"cloud.google.com/go/logging"
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
	ctx := context.Background()
	projectID := "sample-301100"
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
			log.Fatalf("Failed to create client: %v", err)
	}
	logName := "my-log"
	logger := client.Logger(logName)
	text := "Hello, world!"
	logger.Log(logging.Entry{Payload: text})
	if err := client.Close(); err != nil {
			log.Fatalf("Failed to close client: %v", err)
	}
	fmt.Printf("Logged: %v\n", text)
	
	e := echo.New()
	e.Renderer = &Template{
        templates: template.Must(template.ParseGlob("front/*.html")),
	}
	
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK,"index",nil)
	})
	l := e.GET("/date", parseDate)
	e.Logger.Printf("結果： %s", l)
	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
        e.Logger.Printf("Defaulting to port %s", port)
    }
    e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))}

func parseDate(c echo.Context) error {
	date, err := datemaki.Parse(c.QueryParam("value"))
	if err != nil {
		fmt.Printf("REQUESTされた結果は: %v\n", "日付ではありません")
		return c.JSON(http.StatusOK, "日付ではありません")
	}
	fmt.Printf("REQUESTされた結果は: %v\n",date)
	return c.JSON(http.StatusOK, date)
}
