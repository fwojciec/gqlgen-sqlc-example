package dataloaders

//go:generate go run github.com/vektah/dataloaden AgentLoader int64 *github.com/fwojciec/gqlgen-sqlc-example/pg.Agent
//go:generate go run github.com/vektah/dataloaden AuthorSliceLoader int64 []github.com/fwojciec/gqlgen-sqlc-example/pg.Author
//go:generate go run github.com/vektah/dataloaden BookSliceLoader int64 []github.com/fwojciec/gqlgen-sqlc-example/pg.Book

import (
	"context"
	"time"

	"github.com/fwojciec/gqlgen-sqlc-example/pg" // update the username
)

type contextKey string

const key = contextKey("dataloaders")

// Loaders holds references to the individual dataloaders.
type Loaders struct {
	// individual loaders will be defined here
	AgentByAuthorID  *AgentLoader
	AuthorsByAgentID *AuthorSliceLoader
	AuthorsByBookID  *AuthorSliceLoader
	BooksByAuthorID  *BookSliceLoader
}

func newLoaders(ctx context.Context, repo pg.Repository) *Loaders {
	return &Loaders{
		// individual loaders will be initialized here
		AgentByAuthorID:  newAgentByAuthorID(ctx, repo),
		AuthorsByAgentID: newAuthorsByAgentID(ctx, repo),
		AuthorsByBookID:  newAuthorsByBookID(ctx, repo),
		BooksByAuthorID:  newBooksByAuthorID(ctx, repo),
	}
}

// Retriever retrieves dataloaders from the request context.
type Retriever interface {
	Retrieve(context.Context) *Loaders
}

type retriever struct {
	key contextKey
}

func (r *retriever) Retrieve(ctx context.Context) *Loaders {
	return ctx.Value(r.key).(*Loaders)
}

// NewRetriever instantiates a new implementation of Retriever.
func NewRetriever() Retriever {
	return &retriever{key: key}
}

func newAgentByAuthorID(ctx context.Context, repo pg.Repository) *AgentLoader {
	return NewAgentLoader(AgentLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(authorIDs []int64) ([]*pg.Agent, []error) {
			// db query
			res, err := repo.ListAgentsByAuthorIDs(ctx, authorIDs)
			if err != nil {
				return nil, []error{err}
			}
			// map
			groupByAuthorID := make(map[int64]*pg.Agent, len(authorIDs))
			for _, r := range res {
				groupByAuthorID[r.AuthorID] = &pg.Agent{
					ID:    r.ID,
					Name:  r.Name,
					Email: r.Email,
				}
			}
			// order
			result := make([]*pg.Agent, len(authorIDs))
			for i, authorID := range authorIDs {
				result[i] = groupByAuthorID[authorID]
			}
			return result, nil
		},
	})
}

func newAuthorsByAgentID(ctx context.Context, repo pg.Repository) *AuthorSliceLoader {
	return NewAuthorSliceLoader(AuthorSliceLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(agentIDs []int64) ([][]pg.Author, []error) {
			// db query
			res, err := repo.ListAuthorsByAgentIDs(ctx, agentIDs)
			if err != nil {
				return nil, []error{err}
			}
			// group
			groupByAgentID := make(map[int64][]pg.Author, len(agentIDs))
			for _, r := range res {
				groupByAgentID[r.AgentID] = append(groupByAgentID[r.AgentID], r)
			}
			// order
			result := make([][]pg.Author, len(agentIDs))
			for i, agentID := range agentIDs {
				result[i] = groupByAgentID[agentID]
			}
			return result, nil
		},
	})
}

func newAuthorsByBookID(ctx context.Context, repo pg.Repository) *AuthorSliceLoader {
	return NewAuthorSliceLoader(AuthorSliceLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(bookIDs []int64) ([][]pg.Author, []error) {
			// db query
			res, err := repo.ListAuthorsByBookIDs(ctx, bookIDs)
			if err != nil {
				return nil, []error{err}
			}
			// group
			groupByBookID := make(map[int64][]pg.Author, len(bookIDs))
			for _, r := range res {
				groupByBookID[r.BookID] = append(groupByBookID[r.BookID], pg.Author{
					ID:      r.ID,
					Name:    r.Name,
					Website: r.Website,
					AgentID: r.AgentID,
				})
			}
			// order
			result := make([][]pg.Author, len(bookIDs))
			for i, bookID := range bookIDs {
				result[i] = groupByBookID[bookID]
			}
			return result, nil
		},
	})
}

func newBooksByAuthorID(ctx context.Context, repo pg.Repository) *BookSliceLoader {
	return NewBookSliceLoader(BookSliceLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(authorIDs []int64) ([][]pg.Book, []error) {
			// db query
			res, err := repo.ListBooksByAuthorIDs(ctx, authorIDs)
			if err != nil {
				return nil, []error{err}
			}
			// group
			groupByAuthorID := make(map[int64][]pg.Book, len(authorIDs))
			for _, r := range res {
				groupByAuthorID[r.AuthorID] = append(groupByAuthorID[r.AuthorID], pg.Book{
					ID:          r.ID,
					Title:       r.Title,
					Description: r.Description,
					Cover:       r.Cover,
				})
			}
			// order
			result := make([][]pg.Book, len(authorIDs))
			for i, authorID := range authorIDs {
				result[i] = groupByAuthorID[authorID]
			}
			return result, nil
		},
	})
}
