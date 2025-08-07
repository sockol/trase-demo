package handlers

import (
	"api/cmd/api/utils"
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestUsersGetAllTx(t *testing.T) {
	db := utils.TestNewDB(t)

	tests := []struct {
		description   string
		expectedUsers []*User
	}{
		{
			description: "Get existing fixture users",
			expectedUsers: []*User{
				{
					Name:  "user-1",
					Email: "email-1",
				},
				{
					Name:  "user-2",
					Email: "email-2",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {

			_, err := db.Open()
			if err != nil {
				t.Error(err)
				return
			}

			ctx := context.Background()
			err = db.BeginTx(ctx, nil, func(tx *sql.Tx) error {
				users, err := UsersGetAllTx(tx)
				if err != nil {
					return err
				}
				if len(tc.expectedUsers) != len(users) {
					return fmt.Errorf("Wrong len:%d!=%d", len(tc.expectedUsers), len(users))
				}

				for _, u := range users {
					user, err := UsersGetTx(tx, u.Id)
					if err != nil {
						return err
					}
					if user.Name != u.Name {
						return fmt.Errorf("Name mismatch")
					}
					if user.Email != u.Email {
						return fmt.Errorf("Email mismatch")
					}
				}
				return nil
			})
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUsersGetTx(t *testing.T) {
	db := utils.TestNewDB(t)
	id, _ := uuid.Parse("4a2b9c00-9daf-11ed-93ce-0242ac120001")
	tests := []struct {
		description  string
		userToGet    uuid.UUID
		expectedUser *User
		expectError  bool
	}{
		{
			description: "Get existing fixture user",
			userToGet:   db.Fixture.UserId1,
			expectedUser: &User{
				Name:  "user-1",
				Email: "email-1",
			},
		},
		{
			description: "Get non-existing fixture user",
			userToGet:   id,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {

			_, err := db.Open()
			if err != nil {
				t.Error(err)
				return
			}

			ctx := context.Background()
			err = db.BeginTx(ctx, nil, func(tx *sql.Tx) error {
				user, err := UsersGetTx(tx, tc.userToGet)
				if tc.expectError {
					if user != nil {
						return fmt.Errorf("Should be not found")
					}
				} else {
					u := tc.expectedUser

					if err != nil {
						return err
					}
					if user.Name != u.Name {
						return fmt.Errorf("Name mismatch")
					}
					if user.Email != u.Email {
						return fmt.Errorf("Email mismatch")
					}
				}
				return nil
			})
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUsersCreateTx(t *testing.T) {
	db := utils.TestNewDB(t)

	tests := []struct {
		description   string
		usersToCreate []*UserInput
		expectedUsers []*User
	}{
		{
			description: "Create 2 users",
			usersToCreate: []*UserInput{
				{"one", "1"},
				{"two", "2"},
			},
			expectedUsers: []*User{
				{
					Name:  "user-1",
					Email: "email-1",
				},
				{
					Name:  "user-2",
					Email: "email-2",
				},
				{
					Name:  "one",
					Email: "1",
				},
				{
					Name:  "two",
					Email: "2",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {

			_, err := db.Open()
			if err != nil {
				t.Error(err)
				return
			}

			ctx := context.Background()
			err = db.BeginTx(ctx, nil, func(tx *sql.Tx) error {
				for _, u := range tc.usersToCreate {
					u, err := UsersCreateTx(tx, u)
					if err != nil {
						return err
					}
					user, err := UsersGetTx(tx, u.Id)
					if err != nil {
						return err
					}
					if user.Name != u.Name {
						return fmt.Errorf("Name mismatch")
					}
					if user.Email != u.Email {
						return fmt.Errorf("Email mismatch")
					}
				}
				return nil
			})
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUsersUpdateTx(t *testing.T) {
	db := utils.TestNewDB(t)
	type updateInput struct {
		id uuid.UUID
		UserInput
	}
	id, _ := uuid.Parse("4a2b9c00-9daf-11ed-93ce-0242ac120001")
	tests := []struct {
		description   string
		usersToUpdate []*updateInput
		expectedUsers []*User
		expectError   bool
	}{
		{
			description: "Update 1 user",
			usersToUpdate: []*updateInput{
				{db.Fixture.UserId2, UserInput{"user-2-updated", "user-2-updated"}},
			},
			expectedUsers: []*User{
				{
					Name:  "user-1",
					Email: "email-1",
				},
				{
					Name:  "user-2",
					Email: "email-2-updated",
				},
			},
		},
		{
			description: "Update non-existing user",
			usersToUpdate: []*updateInput{
				{id, UserInput{"user-2-updated", "user-2-updated"}},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {

			_, err := db.Open()
			if err != nil {
				t.Error(err)
				return
			}

			ctx := context.Background()
			err = db.BeginTx(ctx, nil, func(tx *sql.Tx) error {
				for _, u := range tc.usersToUpdate {
					u, err := UsersUpdateTx(tx, u.id, &u.UserInput)
					if err != nil {
						return err
					}
					if tc.expectError {
						if u != nil {
							return fmt.Errorf("Should be not found")
						}
					} else {
						user, err := UsersGetTx(tx, u.Id)
						if err != nil {
							return err
						}
						if user.Name != u.Name {
							return fmt.Errorf("Name mismatch")
						}
						if user.Email != u.Email {
							return fmt.Errorf("Email mismatch")
						}
					}
				}
				return nil
			})
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUsersDeleteTx(t *testing.T) {
	db := utils.TestNewDB(t)

	id, _ := uuid.Parse("4a2b9c00-9daf-11ed-93ce-0242ac120001")
	tests := []struct {
		description   string
		usersToDelete []uuid.UUID
		expectedUsers []*User
		expectError   bool
	}{
		{
			description: "Delete 1 user",
			usersToDelete: []uuid.UUID{
				db.Fixture.UserId2,
			},
			expectedUsers: []*User{
				{
					Name:  "user-1",
					Email: "email-1",
				},
			},
		},
		{
			description: "Delete non-existing user",
			usersToDelete: []uuid.UUID{
				id,
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {

			_, err := db.Open()
			if err != nil {
				t.Error(err)
				return
			}

			ctx := context.Background()
			err = db.BeginTx(ctx, nil, func(tx *sql.Tx) error {
				for _, id := range tc.usersToDelete {
					u, err := UsersDeleteTx(tx, id)
					if tc.expectError {
						if u != nil {
							return fmt.Errorf("Should be not found")
						}
					} else {
						if err != nil {
							return err
						}
					}
				}
				return nil
			})
			if err != nil {
				t.Error(err)
			}
		})
	}
}
