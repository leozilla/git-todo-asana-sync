// +build integration

package it

import (
	"git-todo-asana-sync/todosync/pkg/git/grep"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type gitGrepTestSuite struct {
	suite.Suite
}

func TestGitGrepTestSuite(t *testing.T) {
	suite.Run(t, new(gitGrepTestSuite))
}

func (s *gitGrepTestSuite) SetupSuite() {

}

func (s *gitGrepTestSuite) Test_Grep() {
	logger, err := zap.NewDevelopment()
	s.Assert().NoError(err, "New logger")

	gitGrep := grep.NewGitCmdExec(logger, "/home/david/repos/git-todo-asana-sync/.git")

	out, err := gitGrep.Exec()
	s.Assert().NoError(err, "Exec")
	s.Assert().Equal("bla", out)
}
