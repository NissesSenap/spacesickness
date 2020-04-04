package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo"

	"github.com/NissesSenap/spacesickness/handler"
)

// Define the template registry struct
type TemplateRegistry struct {
	templates *template.Template
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Echo instance
	e := echo.New()

	// Instantiate a template registry and register all html files inside the view folder
	e.Renderer = &TemplateRegistry{
		templates: template.Must(template.ParseGlob("../view/*.html")),
	}

	// Route => handler
	e.GET("/", handler.HomeHandler)

	// Start the Echo server
	e.Logger.Fatal(e.Start(":12345"))
}
