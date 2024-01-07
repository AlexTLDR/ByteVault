package main

import "net/http"

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
