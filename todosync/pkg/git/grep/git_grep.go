package grep

import (
	"git-todo-asana-sync/todosync/pkg/cmd"
	"go.uber.org/zap"
	"os/exec"
	"strings"
)

type GitCmdExec struct {
	gitPath string
	logger  *zap.Logger
}

const (
	GitPathArg = "-C"
	GitGrepCmd = "grep"
)

func NewGitCmdExec(logger *zap.Logger, gitPath string) *GitCmdExec {
	return &GitCmdExec{
		gitPath: gitPath,
		logger:  logger,
	}
}

func (git *GitCmdExec) Exec() ([]string, error) {
	pattern := "TODO"

	git.logger.Info("Running git grep",
		zap.String(GitPathArg, git.gitPath),
		zap.String("grep-pattern", pattern))

	c := exec.Command("git", GitPathArg, git.gitPath, GitGrepCmd, pattern)

	out, err := cmd.RunCommand(c)

	outLines := strings.Split(string(out), "\n")
	return outLines, err
}
