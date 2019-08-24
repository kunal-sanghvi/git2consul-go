package git

import (
	"context"
	"git2consul-go/fs"
	"golang.org/x/exp/errors/fmt"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io/ioutil"
	"os"

	cryptoSSH "golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

const (
	isBare = false
)

type opsGit struct {
	lfs  fs.FS
	auth transport.AuthMethod
}

func getSSHKeyAuth(privateSshKeyFile string) transport.AuthMethod {
	var auth transport.AuthMethod
	sshKey, _ := ioutil.ReadFile(privateSshKeyFile)
	signer, _ := cryptoSSH.ParsePrivateKey([]byte(sshKey))
	auth = &ssh.PublicKeys{User: "git", Signer: signer}
	return auth
}

func (ops *opsGit) PlainCloneCtx(ctx context.Context, url string, path string) (*git.Repository, error) {
	repo, err := git.PlainCloneContext(ctx, ops.lfs.GetRootDir()+path, isBare, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Auth:     ops.auth,
	})
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (ops *opsGit) Checkout(repo *git.Repository, branch string) error {
	wTree, err := repo.Worktree()
	if err != nil {
		return err
	}
	if err = wTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
	}); err != nil {
		return err
	}
	return nil
}

func (ops *opsGit) Fetch(ctx context.Context, repo *git.Repository) error {
	if err := repo.FetchContext(ctx, &git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
		Progress: os.Stdout,
		Auth:     ops.auth,
	}); err != nil {
		return err
	}
	return nil
}

func NewGitOps(lfs fs.FS, pathToPemFile string) Git {
	return &opsGit{
		lfs:  lfs,
		auth: getSSHKeyAuth(pathToPemFile),
	}
}
