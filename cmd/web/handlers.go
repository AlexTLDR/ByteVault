package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AlexTLDR/ByteVault/internal/models"
	"github.com/julienschmidt/httprouter"
)

// using the templateData struct holding the snippet data in templates.go for rendering multiple pieces of data
// for this to work, in view.html render the struct, instead of a dot {{.Title}} chain the field names like {{.Snippet.Title}}
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	// using the render helper
	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	// TODO: decide in the view.html how to format the time in the snippet footer
	data := app.newTemplateData(r)
	data.Snippet = snippet
	// using the render helper
	app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "create.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// dummy data for test purpose
	title := "test title2"
	content := "test content2\nTesting a new line,\nand this is the 2nd new line\n\n- 2 new lines"
	expires := 7

	// inserting the dummy data
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// user redirtect
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
