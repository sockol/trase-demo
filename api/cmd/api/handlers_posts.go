package main

import (
	"api/cmd/api/handlers"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// postsGetAll godoc
// @Summary      Get all posts
// @Description  Returns a list of all posts
// @Tags         posts
// @Produce      json
// @Success      200  {array}  handlers.Post
// @Failure      500  {object}  error
// @Router       /api/posts [get]
func (app *application) postsGetAll(ctx context.Context, _ httprouter.Params, _ url.Values) ([]*handlers.Post, error) {
	posts := []*handlers.Post{}
	err := app.db.BeginTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		u, err := handlers.PostsGetAllTx(tx)
		if err != nil {
			return err
		}
		posts = u
		return nil
	})
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// postsGet godoc
// @Summary      Get post by ID
// @Description  Returns a single post by UUID
// @Tags         posts
// @Param        id   path      string  true  "Post ID"
// @Produce      json
// @Success      200  {object}  handlers.Post
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /api/posts/{id} [get]
func (app *application) postsGet(ctx context.Context, params httprouter.Params, _ url.Values) (*handlers.Post, error) {
	var post *handlers.Post
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusBadRequest, err)
	}

	err = app.db.BeginTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		u, err := handlers.PostsGetTx(tx, id)
		if err != nil {
			return err
		}
		post = u
		return nil
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}

// postsCreate godoc
// @Summary      Create post
// @Description  Creates a new post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        post  body      handlers.PostInput  true  "Post Input"
// @Success      201   {object}  handlers.Post
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500   {object}  error
// @Router       /api/posts [post]
func (app *application) postsCreate(ctx context.Context, params httprouter.Params, body []byte) (*handlers.Post, error) {
	var input *handlers.PostInput
	err := json.Unmarshal(body, &input)
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusBadRequest, err)
	}

	var post *handlers.Post
	err = app.db.BeginTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		u, err := handlers.UsersGetTx(tx, input.UserId)
		if err != nil {
			return err
		}
		if u == nil {
			return handlers.NewHTTPError(http.StatusBadRequest, fmt.Errorf("user does not exist"))
		}

		p, err := handlers.PostsCreateTx(tx, input)
		if err != nil {
			return err
		}
		post = p
		return nil
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}

// postsUpdate godoc
// @Summary      Update post
// @Description  Updates an existing post by ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "Post ID"
// @Param        post  body      handlers.PostInput   true  "Updated Post"
// @Success      200   {object}  handlers.Post
// @Failure      400   {object}  error
// @Failure      404   {object}  error
// @Failure      500   {object}  error
// @Router       /api/posts/{id} [put]
func (app *application) postsUpdate(ctx context.Context, params httprouter.Params, body []byte) (*handlers.Post, error) {
	var input *handlers.PostInput
	err := json.Unmarshal(body, &input)
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusBadRequest, err)
	}

	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusBadRequest, err)
	}

	var post *handlers.Post
	err = app.db.BeginTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		u, err := handlers.UsersGetTx(tx, input.UserId)
		if err != nil {
			return err
		}
		if u == nil {
			return handlers.NewHTTPError(http.StatusBadRequest, fmt.Errorf("user does not exist"))
		}

		p, err := handlers.PostsUpdateTx(tx, id, input)
		if err != nil {
			return err
		}
		if p == nil {
			return handlers.NewHTTPError(http.StatusNotFound, fmt.Errorf("post does not exist"))
		}
		post = p
		return nil
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}

// postsDelete godoc
// @Summary      Delete post
// @Description  Deletes a post by ID
// @Tags         posts
// @Produce      json
// @Param        id    path      string  true  "Post ID"
// @Success      200   {object}  handlers.Post
// @Failure      404   {object}  error
// @Failure      500   {object}  error
// @Router       /api/posts/{id} [delete]
func (app *application) postsDelete(ctx context.Context, params httprouter.Params, _ url.Values) (*handlers.Post, error) {
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusNotFound, fmt.Errorf("post does not exist"))
	}

	var post *handlers.Post
	err = app.db.BeginTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		p, err := handlers.PostsDeleteTx(tx, id)
		if err != nil {
			return err
		}
		if p == nil {
			return handlers.NewHTTPError(http.StatusNotFound, fmt.Errorf("post does not exist"))
		}
		post = p
		return nil
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}
