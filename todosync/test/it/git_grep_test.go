// +build integration

package it

import (
	"git-todo-asana-sync/todosync/pkg/git/grep"
	"git-todo-asana-sync/todosync/test"
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

func (s *gitGrepTestSuite) Test_GrepLocalRepo() {
	logger, err := zap.NewDevelopment()
	s.Assert().NoError(err, "New logger")

	gitPath := test.AbsoluteProjectDirPath()
	gitGrep := grep.NewGitCmdExec(logger, gitPath)

	out, err := gitGrep.Exec()
	s.Assert().NoError(err, "Exec")
	s.Assert().Contains(out, "todosync/pkg/git/grep/git_grep.go:\tpattern := \"TODO\"")
	s.Assert().Contains(out, "todosync/test/testdata/scala/ScalaFileB.scala:  private val stage = \"dev\" // TODO inject stage param from build, once we move to multi-stage environments")
	s.Assert().Contains(out, "todosync/test/testdata/scala/ScalaFileC.scala:  // TODO: Use Java.URI for parsing?")
}
