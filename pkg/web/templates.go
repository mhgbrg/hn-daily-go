package web

import (
	"fmt"
	templatelib "html/template"

	"github.com/pkg/errors"
)

var cache map[string]*templatelib.Template

// TODO: Read templates when server is loaded instead of doing it lazily.
func GetTemplate(name string) (*templatelib.Template, error) {
	if cache == nil {
		cache = make(map[string]*templatelib.Template)
	}

	if template, ok := cache[name]; ok {
		return template, nil
	}

	filename := fmt.Sprintf("templates/%s.html", name)
	template, err := templatelib.ParseFiles(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse template %s (%s)", name, filename)
	}

	cache[name] = template

	return template, nil
}
