package main

import (
	"context"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *application) health(ctx context.Context, _ httprouter.Params, _ url.Values) (map[string]string, error) {
	return map[string]string{
		"Status": "OK",
	}, nil
}

func (app *application) docs() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		httpSwagger.WrapHandler(w, r)
	}
}
