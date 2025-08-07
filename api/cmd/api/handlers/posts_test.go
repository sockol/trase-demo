package handlers

import (
	"api/cmd/api/utils"
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestPostsGetAllTx(t *testing.T) {
	db := utils.TestNewDB(t)

	tests := []struct {
		description   string
		expectedPosts []*Post
	}{
		{
			description: "Get existing fixture posts",
			expectedPosts: []*Post{
				{
					Title:   "title-1",
					Content: "content-1",
					UserId:  db.Fixture.UserId1,
				},
				{
					Title:   "title-1",
					Content: "content-1",
					UserId:  db.Fixture.UserId2,
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
				posts, err := PostsGetAllTx(tx)
				if err != nil {
					return err
				}
				if len(tc.expectedPosts) != len(posts) {
					return fmt.Errorf("Wrong len:%d!=%d", len(tc.expectedPosts), len(posts))
				}

				for _, u := range posts {
					post, err := PostsGetTx(tx, u.Id)
					if err != nil {
						return err
					}
					if post.Title != u.Title {
						return fmt.Errorf("Title mismatch")
					}
					if post.Content != u.Content {
						return fmt.Errorf("Content mismatch")
					}
					if post.UserId != u.UserId {
						return fmt.Errorf("UserId mismatch")
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

func TestPostsGetTx(t *testing.T) {
	db := utils.TestNewDB(t)

	id, _ := uuid.Parse("4a2b9c00-9daf-11ed-93ce-0242ac120001")
	tests := []struct {
		description  string
		postToGet    uuid.UUID
		expectedPost *Post
		expectError  bool
	}{
		{
			description: "Get existing fixture post",
			postToGet:   db.Fixture.PostId1,
			expectedPost: &Post{
				Title:   "title-1",
				Content: "content-1",
				UserId:  db.Fixture.UserId1,
			},
		},
		{
			description: "Get non-existing fixture post",
			postToGet:   id,
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
				post, err := PostsGetTx(tx, tc.postToGet)
				if tc.expectError {
					if post != nil {
						return fmt.Errorf("Should be not found")
					}
				} else {
					u := tc.expectedPost

					if err != nil {
						return err
					}
					if post.Title != u.Title {
						return fmt.Errorf("Title mismatch")
					}
					if post.Content != u.Content {
						return fmt.Errorf("Content mismatch")
					}
					if post.UserId != u.UserId {
						return fmt.Errorf("UserId mismatch")
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

func TestPostsCreateTx(t *testing.T) {
	db := utils.TestNewDB(t)

	id, _ := uuid.Parse("4a2b9c00-9daf-11ed-93ce-0242ac120001")
	tests := []struct {
		description   string
		postsToCreate []*PostInput
		expectedPosts []*Post
		expectError   bool
	}{
		{
			description: "Create 2 posts",
			postsToCreate: []*PostInput{
				{"one", "1", db.Fixture.UserId1},
				{"two", "2", db.Fixture.UserId2},
			},
			expectedPosts: []*Post{
				{
					Title:   "title-1",
					Content: "content-1",
					UserId:  db.Fixture.UserId1,
				},
				{
					Title:   "title-1",
					Content: "content-1",
					UserId:  db.Fixture.UserId2,
				},
				{
					Title:   "one",
					Content: "1",
					UserId:  db.Fixture.UserId1,
				},
				{
					Title:   "two",
					Content: "2",
					UserId:  db.Fixture.UserId2,
				},
			},
		},
		{
			description: "Create post on non existing user, expect fail",
			postsToCreate: []*PostInput{
				{"one", "1", id},
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
				for _, u := range tc.postsToCreate {
					u, err := PostsCreateTx(tx, u)
					if err != nil {
						return err
					}
					post, err := PostsGetTx(tx, u.Id)
					if err != nil {
						return err
					}
					if post.Title != u.Title {
						return fmt.Errorf("Title mismatch")
					}
					if post.Content != u.Content {
						return fmt.Errorf("Content mismatch")
					}
					if post.UserId != u.UserId {
						return fmt.Errorf("UserId mismatch")
					}
				}
				return nil
			})
			if tc.expectError {
				if err == nil {
					t.Error(err)
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestPostsUpdateTx(t *testing.T) {
	db := utils.TestNewDB(t)
	type updateInput struct {
		id uuid.UUID
		PostInput
	}
	id, _ := uuid.Parse("4a2b9c00-9daf-11ed-93ce-0242ac120001")
	tests := []struct {
		description   string
		postsToUpdate []*updateInput
		expectedPosts []*Post
		expectError   bool
	}{
		{
			description: "Update 1 post",
			postsToUpdate: []*updateInput{
				{db.Fixture.PostId2, PostInput{"title-2-updated", "content-2-updated", db.Fixture.UserId2}},
			},
			expectedPosts: []*Post{
				{
					Title:   "title-2-updated",
					Content: "content-2-updated",
					UserId:  db.Fixture.UserId2,
				},
			},
		},
		{
			description: "Update non-existing user on existing post",
			postsToUpdate: []*updateInput{
				{db.Fixture.PostId2, PostInput{"title-2-updated", "content-2-updated", id}},
			},
			expectError: true,
		},
		{
			description: "Update non-existing post",
			postsToUpdate: []*updateInput{
				{id, PostInput{"title-2-updated", "content-2-updated", db.Fixture.UserId2}},
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
				for _, u := range tc.postsToUpdate {
					u, err := PostsUpdateTx(tx, u.id, &u.PostInput)
					if err != nil {
						return err
					}
					if u == nil {
						return fmt.Errorf("Post not found")
					}
					post, err := PostsGetTx(tx, u.Id)
					if err != nil {
						return err
					}
					if post.Title != u.Title {
						return fmt.Errorf("Title mismatch")
					}
					if post.Content != u.Content {
						return fmt.Errorf("Content mismatch")
					}
					if post.UserId != u.UserId {
						return fmt.Errorf("UserId mismatch")
					}
				}
				return nil
			})
			if tc.expectError {
				if err == nil {
					t.Error(err)
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestPostsDeleteTx(t *testing.T) {
	db := utils.TestNewDB(t)

	id, _ := uuid.Parse("4a2b9c00-9daf-11ed-93ce-0242ac120001")
	tests := []struct {
		description   string
		postsToDelete []uuid.UUID
		expectedPosts []*Post
		expectError   bool
	}{
		{
			description: "Delete 1 post",
			postsToDelete: []uuid.UUID{
				db.Fixture.PostId2,
			},
			expectedPosts: []*Post{
				{
					Title:   "title-2",
					Content: "content-2",
					UserId:  db.Fixture.UserId2,
				},
			},
		},
		{
			description: "Delete non-existing post",
			postsToDelete: []uuid.UUID{
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
				for _, id := range tc.postsToDelete {
					u, err := PostsDeleteTx(tx, id)
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
