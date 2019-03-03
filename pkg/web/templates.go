package web

import (
	"fmt"
	templatelib "html/template"

	"github.com/pkg/errors"
)

type Templates struct {
	Digest  *templatelib.Template
	Archive *templatelib.Template
}

func LoadTemplates() (*Templates, error) {
	digestTemplate, err := loadTemplate("digest")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to load digest template")
	}

	archiveTemplate, err := loadTemplate("archive")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to load archive template")
	}

	return &Templates{
		Digest:  digestTemplate,
		Archive: archiveTemplate,
	}, nil
}

func loadTemplate(name string) (*templatelib.Template, error) {
	filename := fmt.Sprintf("templates/%s.html", name)
	template, err := templatelib.ParseFiles(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse template %s (%s)", name, filename)
	}
	return template, nil
}
