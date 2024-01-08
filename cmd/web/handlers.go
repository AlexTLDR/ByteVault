package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) // using the notFound() helper
		return
	}
	files := []string{"./ui/html/base.html", "./ui/html/pages/home.html", "./ui/html/partials/nav.html"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err) // using the serverError() helper
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err) // using the serverError() helper
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // using the notFound() helper
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed) // using the clientError() helper
		return
	}

	// dummy data for test purpose
	title := "test title"
	content := "test content\nTesting a new line,\nand this is the 2nd new line\n\n- 2 new lines"
	expires := 7

	// inserting the dummy data
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// user redirtect
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
