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

// usersGetAll godoc
// @Summary      Get all users
// @Description  Returns a list of all users
// @Tags         users
// @Produce      json
// @Success      200  {array}  handlers.User
// @Failure      500  {object}  error
// @Router       /api/users [get]
func (app *application) usersGetAll(ctx context.Context, _ httprouter.Params, _ url.Values) ([]*handlers.User, error) {
	users := []*handlers.User{}
	err := app.db.BeginTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		u, err := handlers.UsersGetAllTx(tx)
		if err != nil {
			return err
		}
		users = u
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

// usersGet godoc
// @Summary      Get user by ID
// @Description  Returns a single user by UUID
// @Tags         users
// @Param        id   path      string  true  "User ID"
// @Produce      json
// @Success      200  {object}  handlers.User
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /api/users/{id} [get]
func (app *application) usersGet(ctx context.Context, params httprouter.Params, _ url.Values) (*handlers.User, error) {
	var user *handlers.User
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusBadRequest, err)
	}

	err = app.db.BeginTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		u, err := handlers.UsersGetTx(tx, id)
		if err != nil {
			return err
		}
		if u == nil {
			return handlers.NewHTTPError(http.StatusNotFound, fmt.Errorf("user does not exist"))
		}
		user = u
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// usersCreate godoc
// @Summary      Create user
// @Description  Creates a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      handlers.UserInput  true  "User Input"
// @Success      201   {object}  handlers.User
// @Failure      404  {object}  error
// @Failure      500   {object}  error
// @Router       /api/users [post]
func (app *application) usersCreate(ctx context.Context, params httprouter.Params, body []byte) (*handlers.User, error) {
	var input *handlers.UserInput
	err := json.Unmarshal(body, &input)
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusBadRequest, err)
	}

	var user *handlers.User
	err = app.db.BeginTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		u, err := handlers.UsersCreateTx(tx, input)
		if err != nil {
			return err
		}
		user = u
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// usersUpdate godoc
// @Summary      Update user
// @Description  Updates an existing user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "User ID"
// @Param        user  body      handlers.UserInput   true  "Updated User"
// @Success      200   {object}  handlers.User
// @Failure      404   {object}  error
// @Failure      500   {object}  error
// @Router       /api/users/{id} [put]
func (app *application) usersUpdate(ctx context.Context, params httprouter.Params, body []byte) (*handlers.User, error) {
	var input *handlers.UserInput
	err := json.Unmarshal(body, &input)
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusBadRequest, err)
	}

	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusBadRequest, err)
	}

	var user *handlers.User
	err = app.db.BeginTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		u, err := handlers.UsersUpdateTx(tx, id, input)
		if err != nil {
			return err
		}
		if u == nil {
			return handlers.NewHTTPError(http.StatusNotFound, fmt.Errorf("user does not exist"))
		}
		user = u
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// usersDelete godoc
// @Summary      Delete user
// @Description  Deletes a user by ID
// @Tags         users
// @Produce      json
// @Param        id    path      string  true  "User ID"
// @Success      200   {object}  handlers.User
// @Failure      404   {object}  error
// @Failure      500   {object}  error
// @Router       /api/users/{id} [delete]
func (app *application) usersDelete(ctx context.Context, params httprouter.Params, _ url.Values) (*handlers.User, error) {
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusBadRequest, err)
	}

	var user *handlers.User
	err = app.db.BeginTx(ctx, &sql.TxOptions{}, func(tx *sql.Tx) error {
		u, err := handlers.UsersDeleteTx(tx, id)
		if err != nil {
			return err
		}
		if u == nil {
			return handlers.NewHTTPError(http.StatusNotFound, fmt.Errorf("user does not exist"))
		}
		user = u
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
