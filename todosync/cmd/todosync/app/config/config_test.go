package config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os"
	"testing"
)

var logger, _ = zap.NewDevelopment()

func Test_ConfIsValidated(t *testing.T) {
	assert.NoError(t, os.Setenv("GIT_REPO_URL", "127:0.0.1"))
	assert.Panics(t, func() {
		MustLoadFromEnv(context.Background(), logger)
	})
}
