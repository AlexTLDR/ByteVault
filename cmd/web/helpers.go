package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
)

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

// The serverError helper logs an entry at the Error level (including the request
// method and URI as attributes) and then sends a standard 500 Internal Server Error
// response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper transmits a designated status code and its corresponding description
// to the user. Later in the book, we will utilize this to send responses such as 400 "Bad
// Request" when there's an issue with the user's sent request.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// To maintain consistency, I will also create a notFound helper. This is essentially a
// convenient wrapper around clientError, designed to dispatch a 404 Not Found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Get the right template set from the cache using the page name (e.g., 'home.tmpl').
	// If there's no entry in the cache with the given name, create a new error and use the serverError()
	// method we created before, then return.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("template %s not found", page)
		app.serverError(w, r, err)
		return
	}
	// using a buffer to be able to render errors to the user
	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}
