package renderer

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"

	"github.com/labstack/echo/v4"
)

type Renderer interface {
	// Render will write the template to given the given writer.
	// Will return an error if we can't write to the writer or we can
	// can't find the template.
	Render(w io.Writer, templateName string, data interface{}) error
}

type ErrTemplateNotFound struct {
	Err  error
	Name string
}

func (e *ErrTemplateNotFound) Error() string {
	return fmt.Sprintf("unable to find the template '%s', err: %v", e.Name, e.Err)
}

type ErrFailedToWrite struct {
	Err error
}

func (e *ErrFailedToWrite) Error() string {
	return fmt.Sprintf("failed to write template, err: %v", e.Err)
}

type HTMLRender struct {
	fs        fs.FS
	pattern   string
	reload    bool
	templates *template.Template
}

func New(fs fs.FS, pattern string, reload bool) (*HTMLRender, error) {
	return &HTMLRender{
		templates: template.Must(template.ParseFS(fs, pattern)),
		fs:        fs,
		pattern:   pattern,
		reload:    reload,
	}, nil
}

func (h *HTMLRender) Render(w io.Writer, templateName string, data interface{}, ctx echo.Context) error {
	if h.reload {
		h.templates = template.Must(template.ParseFS(h.fs, h.pattern))
	}

	return h.templates.ExecuteTemplate(w, templateName, data)
}
