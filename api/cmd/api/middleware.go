package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"api/cmd/api/handlers"
	"api/internal/response"

	"github.com/julienschmidt/httprouter"
	"github.com/tomasen/realip"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) logAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := response.NewMetricsResponseWriter(w)
		next.ServeHTTP(mw, r)

		var (
			ip     = realip.FromRequest(r)
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)

		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)
		responseAttrs := slog.Group("repsonse", "status", mw.StatusCode, "size", mw.BytesCount)

		app.logger.Info("access", userAttrs, requestAttrs, responseAttrs)
	})
}

func handleQuery[T any](app *application, handler func(context.Context, httprouter.Params, url.Values) (T, error)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()

		done := make(chan struct{})
		var result T
		var err error

		q := r.URL.Query()

		go func() {
			result, err = handler(ctx, p, q)
			close(done)
		}()

		select {
		case <-ctx.Done():
			return

		case <-done:
			if err != nil {
				if httpErr, ok := err.(*handlers.HTTPError); ok {
					app.errorMessage(w, r, httpErr.Code, httpErr.Message.Error(), nil)
				} else {
					app.serverError(w, r, err)
				}
				return
			}

			err = response.JSON(w, http.StatusOK, result)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}
	}
}

func handleMutation[T any](app *application, handler func(context.Context, httprouter.Params, []byte) (*T, error)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()

		done := make(chan struct{})
		var result *T
		var err error

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		go func() {
			result, err = handler(ctx, p, body)
			close(done)
		}()

		select {
		case <-ctx.Done():
			return

		case <-done:
			if err != nil {
				if httpErr, ok := err.(*handlers.HTTPError); ok {
					app.errorMessage(w, r, httpErr.Code, httpErr.Message.Error(), nil)
				} else {
					app.serverError(w, r, err)
				}
				return
			}
			if result == nil {
				app.notFound(w, r)
				return
			}
			err = response.JSON(w, http.StatusOK, result)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}
	}
}
