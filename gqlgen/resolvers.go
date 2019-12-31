package gqlgen

import (
	"context"

	"github.com/fwojciec/gqlgen-sqlc-example/pg"
)

// Resolver connects individual resolvers with the datalayer.
type Resolver struct {
	Repository pg.Repository
}

// Agent returns an implementation of the AgentResolver interface.
func (r *Resolver) Agent() AgentResolver {
	return &agentResolver{r}
}

// Author returns an implementation of the AuthorResolver interface.
func (r *Resolver) Author() AuthorResolver {
	return &authorResolver{r}
}

// Book returns an implementation of the BookResolver interface.
func (r *Resolver) Book() BookResolver {
	return &bookResolver{r}
}

// Mutation returns an implementation of the MutationResolver interface.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returns an implementation of the QueryResolver interface.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type agentResolver struct{ *Resolver }

func (r *agentResolver) Authors(ctx context.Context, obj *pg.Agent) ([]pg.Author, error) {
	panic("not implemented")
}

type authorResolver struct{ *Resolver }

func (r *authorResolver) Website(ctx context.Context, obj *pg.Author) (*string, error) {
	panic("not implemented")
}
func (r *authorResolver) Agent(ctx context.Context, obj *pg.Author) (*pg.Agent, error) {
	panic("not implemented")
}
func (r *authorResolver) Books(ctx context.Context, obj *pg.Author) ([]pg.Book, error) {
	panic("not implemented")
}

type bookResolver struct{ *Resolver }

func (r *bookResolver) Authors(ctx context.Context, obj *pg.Book) ([]pg.Author, error) {
	panic("not implemented")
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateAgent(ctx context.Context, data AgentInput) (*pg.Agent, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateAgent(ctx context.Context, id int64, data AgentInput) (*pg.Agent, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteAgent(ctx context.Context, id int64) (*pg.Agent, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateAuthor(ctx context.Context, data AuthorInput) (*pg.Author, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateAuthor(ctx context.Context, id int64, data AuthorInput) (*pg.Author, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteAuthor(ctx context.Context, id int64) (*pg.Author, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateBook(ctx context.Context, data BookInput) (*pg.Book, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateBook(ctx context.Context, id int64, data BookInput) (*pg.Book, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteBook(ctx context.Context, id int64) (*pg.Book, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Agent(ctx context.Context, id int64) (*pg.Agent, error) {
	panic("not implemented")
}
func (r *queryResolver) Agents(ctx context.Context) ([]pg.Agent, error) {
	panic("not implemented")
}
func (r *queryResolver) Author(ctx context.Context, id int64) (*pg.Author, error) {
	panic("not implemented")
}
func (r *queryResolver) Authors(ctx context.Context) ([]pg.Author, error) {
	panic("not implemented")
}
func (r *queryResolver) Book(ctx context.Context, id int64) (*pg.Book, error) {
	panic("not implemented")
}
func (r *queryResolver) Books(ctx context.Context) ([]pg.Book, error) {
	panic("not implemented")
}
