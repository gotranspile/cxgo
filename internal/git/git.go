package git

import (
	"os"
	"os/exec"
)

func Clone(repo, branch, dir string) error {
	args := []string{
		"clone",
		"--depth=1", // we are only interested in the last commit
	}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, repo, dir)
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
