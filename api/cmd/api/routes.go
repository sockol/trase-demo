package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()

	mux.NotFound = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	mux.GET("/health", handleQuery(app, app.health))
	mux.GET("/docs/*any", app.docs())

	mux.GET("/api/users", handleQuery(app, app.usersGetAll))
	mux.GET("/api/users/:id", handleQuery(app, app.usersGet))
	mux.POST("/api/users", handleMutation(app, app.usersCreate))
	mux.DELETE("/api/users/:id", handleQuery(app, app.usersDelete))
	mux.PUT("/api/users/:id", handleMutation(app, app.usersUpdate))

	mux.GET("/api/posts", handleQuery(app, app.postsGetAll))
	mux.GET("/api/posts/:id", handleQuery(app, app.postsGet))
	mux.POST("/api/posts", handleMutation(app, app.postsCreate))
	mux.DELETE("/api/posts/:id", handleQuery(app, app.postsDelete))
	mux.PUT("/api/posts/:id", handleMutation(app, app.postsUpdate))

	return app.logAccess(app.recoverPanic(mux))
}
