package utils

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/galrub/go/jobSearch/internal/logger"
)

func RenderFragment(filename string, data any, w *http.ResponseWriter) error {
	lp := filepath.Join("web", "views", filename)
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		logger.LOG.Err(err).Msg("Error parsing the index Page")
		return err
	}
	err = tmpl.Execute(*w, data)
	if err != nil {
		logger.LOG.Err(err).Msg("Error parsing the index Page")
		return err
	}
	return nil
}

func RenderFragments(data any, w *http.ResponseWriter, filenames ...string) error {
	var paths []string
	for _, filname := range filenames {
		paths = append(paths, filepath.Join("web", "views", filname))
	}
	tmpl, err := template.ParseFiles(paths...)
	if err != nil {
		logger.LOG.Err(err).Msg("Error parsing the index Page")
		return err
	}
	err = tmpl.Execute(*w, data)
	if err != nil {
		logger.LOG.Err(err).Msg("Error parsing the index Page")
		return err
	}
	return nil
}
