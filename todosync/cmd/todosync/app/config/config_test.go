package config

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"testing"
)

var logger, _ = zap.NewDevelopment()

type configTestSuite struct {
	suite.Suite
	tmpGitFile string
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(configTestSuite))
}

func (s *configTestSuite) SetupSuite() {
	dir, err := ioutil.TempDir(os.TempDir(), "git-todo-asana-sync")
	s.Assert().NoError(err, "TempDir")

	s.tmpGitFile = fmt.Sprintf("%s/.git", dir)
	err = ioutil.WriteFile(s.tmpGitFile, []byte{}, 0644)
	s.Assert().NoError(err, "TempFile")
}

func (s *configTestSuite) Test_GitDirIsValidated() {
	s.noPanicOnLoad("GIT_DIR", s.tmpGitFile)

	s.panicOnLoad("GIT_DIR", "")
	s.panicOnLoad("GIT_DIR", "in valid")
	s.panicOnLoad("GIT_DIR", "in valid/.git")
	s.panicOnLoad("GIT_DIR", "rel/ative/bla.git")
}

func (s *configTestSuite) Test_AsanaServerUrlIsValidated() {
	s.noPanicOnLoad("ASANA_SERVER_URL", "https://app.asana.com")

	s.panicOnLoad("ASANA_SERVER_URL", "app.asana.com")
}

func (s *configTestSuite) Test_AsanaAccessTokenIsValidated() {
	s.noPanicOnLoad("ASANA_ACCESS_TOKEN", "1/1199673690536123:43d4feeb2f75d7b554a1f9cb7cca5e7d")

	s.panicOnLoad("ASANA_ACCESS_TOKEN", "")
}

func (s *configTestSuite) noPanicOnLoad(key string, value string) {
	s.Assert().NotPanics(func() {
		s.setValidConfig()

		s.Assert().NoError(os.Setenv(fmt.Sprintf("APP_%s", key), value))
		MustLoadFromEnv(context.Background(), logger)
	})
}

func (s *configTestSuite) panicOnLoad(key string, value string) {
	s.Assert().Panics(func() {
		s.setValidConfig()

		s.Assert().NoError(os.Setenv(fmt.Sprintf("APP_%s", key), value))
		MustLoadFromEnv(context.Background(), logger)
	})
}

func (s *configTestSuite) setValidConfig() {
	s.Assert().NoError(os.Setenv("APP_GIT_DIR", s.tmpGitFile))
	s.Assert().NoError(os.Setenv("APP_ASANA_SERVER_URL", "https://app.asana.com"))
	s.Assert().NoError(os.Setenv("APP_ASANA_ACCESS_TOKEN", "1/1199673690536123:43d4feeb2f75d7b554a1f9cb7cca5e7d"))
}
