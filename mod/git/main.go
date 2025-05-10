package main

import (
	"flag"
	"os"

	"github.com/go-git/go-git/v6"
	//. "github.com/go-git/go-git/v6/_examples"
)

// Example of how to:
// - Access basic local (i.e. ./.git/config) configuration params
// - Set basic local config params

func main() {
	repoPath := flag.String("path", "", "Directory to initialize the git repository. If empty, a temporary directory is used.")
	flag.Parse()

	var err error
	var workDir string

	if *repoPath == "" {
		workDir, err = os.MkdirTemp("", "go-git-example")
		CheckIfError(err)
		Info("Using temporary directory: %s", workDir)
		// Schedule cleanup only if we created a temporary directory
		defer os.RemoveAll(workDir)
	} else {
		workDir = *repoPath
		Info("Using specified directory: %s", workDir)
		// For a user-specified path, we don't automatically clean it up.
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	cfg, err := repo.Config()
	if err != nil {
		return err
	}

	Info("git init in %s", workDir)
	r, err := git.PlainInit(workDir, false)
	CheckIfError(err)

	// Load the configuration
	cfg, err := r.Config()
	CheckIfError(err)

	Info("worktree is %s", cfg.Core.Worktree)

	/*
		// Set basic local config params
		cfg.Remotes["origin"] = &config.RemoteConfig{
			Name: "origin",
			URLs: []string{"https://github.com/git-fixtures/basic.git"},
		}
	*/

	Info("origin remote: %+v", cfg.Remotes["origin"])

	//cfg.User.Name = "Local name"

	Info("custom.name is %s", cfg.User.Name)

	// In order to save the config file, you need to call SetConfig
	// After calling this go to .git/config and see the custom.name added and the changes to the remote
	//r.Storer.SetConfig(cfg)
}
