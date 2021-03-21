package grep

import (
	"git-todo-asana-sync/todosync/pkg/cmd"
	"go.uber.org/zap"
	"os/exec"
)

type GitCmdExec struct {
	gitDir string
	logger *zap.Logger
}

const (
	GitDirArg  = "--git-dir"
	GitGrepCmd = "grep"
)

func NewGitCmdExec(logger *zap.Logger, gitDir string) *GitCmdExec {
	return &GitCmdExec{
		gitDir: gitDir,
		logger: logger,
	}
}

func (git *GitCmdExec) Exec() (string, error) {
	pattern := "'TODO'"

	git.logger.Info("Running git grep",
		zap.String(GitDirArg, git.gitDir),
		zap.String("grep-pattern", pattern))

	c := exec.Command("git", GitDirArg, git.gitDir, GitGrepCmd, pattern)

	out, err := cmd.RunCommandWithOutput(c)
	return string(out), err
}
