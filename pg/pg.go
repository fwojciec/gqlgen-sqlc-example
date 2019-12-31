package pg

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // required
)

// Repository is the application's data layer functionality.
type Repository interface {
	// agent queries
	CreateAgent(ctx context.Context, arg CreateAgentParams) (Agent, error)
	DeleteAgent(ctx context.Context, id int64) (Agent, error)
	GetAgent(ctx context.Context, id int64) (Agent, error)
	ListAgents(ctx context.Context) ([]Agent, error)
	UpdateAgent(ctx context.Context, arg UpdateAgentParams) (Agent, error)

	// author queries
	CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error)
	DeleteAuthor(ctx context.Context, id int64) (Author, error)
	GetAuthor(ctx context.Context, id int64) (Author, error)
	ListAuthors(ctx context.Context) ([]Author, error)
	UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) (Author, error)
	ListAuthorsByAgentID(ctx context.Context, agentID int64) ([]Author, error)
	ListAuthorsByBookID(ctx context.Context, bookID int64) ([]Author, error)

	// book queries
	CreateBook(ctx context.Context, bookArg CreateBookParams, authorIDs []int64) (*Book, error)
	UpdateBook(ctx context.Context, bookArg UpdateBookParams, authorIDs []int64) (*Book, error)
	DeleteBook(ctx context.Context, id int64) (Book, error)
	GetBook(ctx context.Context, id int64) (Book, error)
	ListBooks(ctx context.Context) ([]Book, error)
	ListBooksByAuthorID(ctx context.Context, authorID int64) ([]Book, error)
}

type repoSvc struct {
	*Queries
	db *sql.DB
}

func (r *repoSvc) withTx(ctx context.Context, txFn func(*Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = txFn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			err = fmt.Errorf("tx failed: %v, unable to rollback: %v", err, rbErr)
		}
	} else {
		err = tx.Commit()
	}
	return err
}

func (r *repoSvc) CreateBook(ctx context.Context, bookArg CreateBookParams, authorIDs []int64) (*Book, error) {
	book := new(Book)
	err := r.withTx(ctx, func(q *Queries) error {
		res, err := q.CreateBook(ctx, bookArg)
		if err != nil {
			return err
		}
		for _, authorID := range authorIDs {
			if err := q.SetBookAuthor(ctx, SetBookAuthorParams{
				BookID:   res.ID,
				AuthorID: authorID,
			}); err != nil {
				return err
			}
		}
		book = &res
		return nil
	})
	return book, err
}

func (r *repoSvc) UpdateBook(ctx context.Context, bookArg UpdateBookParams, authorIDs []int64) (*Book, error) {
	book := new(Book)
	err := r.withTx(ctx, func(q *Queries) error {
		res, err := q.UpdateBook(ctx, bookArg)
		if err != nil {
			return err
		}
		if err = q.UnsetBookAuthors(ctx, res.ID); err != nil {
			return err
		}
		for _, authorID := range authorIDs {
			if err := q.SetBookAuthor(ctx, SetBookAuthorParams{
				BookID:   res.ID,
				AuthorID: authorID,
			}); err != nil {
				return err
			}
		}
		book = &res
		return nil
	})
	return book, err
}

// NewRepository returns an implementation of the Repository interface.
func NewRepository(db *sql.DB) Repository {
	return &repoSvc{
		Queries: New(db),
		db:      db,
	}
}

// Open opens a database specified by the data source name.
// Format: host=foo port=5432 user=bar password=baz dbname=qux sslmode=disable"
func Open(dataSourceName string) (*sql.DB, error) {
	return sql.Open("postgres", dataSourceName)
}

// StringPtrToNullString converts *string to sql.NullString.
func StringPtrToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}
