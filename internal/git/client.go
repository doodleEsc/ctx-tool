package git

import (
	"fmt"
	"os"

	"github.com/doodleEsc/ctx-tool/internal/i18n"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Client struct {
	repoURL string
	branch  string
}

func NewClient(repoURL, branch string) *Client {
	return &Client{
		repoURL: repoURL,
		branch:  branch,
	}
}

// CloneToTemp clones the repository to a temporary directory
func (c *Client) CloneToTemp() (string, error) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "ctx-tool-*")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}

	fmt.Printf("%s\n", i18n.Tf(i18n.MsgCloningRepository, map[string]interface{}{"Repo": c.repoURL, "Branch": c.branch}))

	// Clone with progress output
	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:           c.repoURL,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", c.branch)),
		SingleBranch:  true,
		Depth:         1, // Shallow clone for speed
		Progress:      os.Stdout,
	})

	if err != nil {
		// Clean up temp directory on error
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("clone repository: %w", err)
	}

	fmt.Printf("%s\n", i18n.Tf(i18n.MsgRepositoryCloned, map[string]interface{}{"Path": tempDir}))
	return tempDir, nil
}

// CloneToDirectory clones the repository to a specific directory
func (c *Client) CloneToDirectory(targetDir string) error {
	fmt.Printf("%s\n", i18n.Tf(i18n.MsgCloningRepository, map[string]interface{}{"Repo": c.repoURL, "Branch": c.branch}))

	// Ensure directory doesn't exist or is empty
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		// Directory exists, check if it's empty
		entries, err := os.ReadDir(targetDir)
		if err != nil {
			return fmt.Errorf("read directory: %w", err)
		}
		if len(entries) > 0 {
			return fmt.Errorf("target directory %s is not empty", targetDir)
		}
	}

	// Clone with progress output
	_, err := git.PlainClone(targetDir, false, &git.CloneOptions{
		URL:           c.repoURL,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", c.branch)),
		SingleBranch:  true,
		Depth:         1, // Shallow clone for speed
		Progress:      os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("clone repository: %w", err)
	}

	fmt.Printf("%s\n", i18n.T(i18n.MsgCloneSuccess))
	return nil
}
