package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{ //return template data(for rendering templates)
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status) // handle client errors
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound) // wrap for clientErr
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page] // check if the template exists
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page) // throw error if the template does not exist
		app.serverError(w, err)
		return
	}
	buf := new(bytes.Buffer) // create a new buffer for the template
	err := ts.ExecuteTemplate(buf, "base", data) // execute the template 
	if err != nil {
		app.serverError(w, err) // throw server error if the template cannot be executed
		return
	}
	w.WriteHeader(status) // write the status code
	buf.WriteTo(w) // write the buffer to the response writer
}

func (app *application) decodePostForm(r *http.Request, dst any) error { // decoder method

	err := r.ParseForm() // parsing post form
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm) // decoding post form
	if err != nil {

		var invalidDecoderError *form.InvalidDecoderError // checking for invalid decoder error

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool) // checking if the user is authenticated
	if !ok {
		return false
	}

	return isAuthenticated
}

func (app *application) serverError(w http.ResponseWriter, err error) { // server error handler
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()) // format stack trace response string
	app.errorLog.Output(2, trace)

	if app.debug {
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
