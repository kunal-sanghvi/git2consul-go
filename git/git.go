package git

import (
	"context"

	"gopkg.in/src-d/go-git.v4"
)

type Git interface {
	PlainCloneCtx(ctx context.Context, url string, path string) (*git.Repository, error)
	Checkout(repo *git.Repository, branch string) error
	Fetch(ctx context.Context, repo *git.Repository) error
}
