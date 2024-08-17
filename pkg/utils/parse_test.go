package utils_test

import (
	"os"
	"testing"
	"vpa-manager/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestParseBoolWithValidParam(t *testing.T) {
	assert := assert.New(t)
	const envVar = "TEST_ENV_VAR"

	varValues := []string{"TRUE", "true", "True"}

	for _, envVarValue := range varValues {
		os.Setenv(envVar, envVarValue)
		parsedValue := utils.ParseBoolFromEnv(envVar)
		assert.Equal(parsedValue, true, "Parse Bool from Env failed for value: %s", envVarValue)
	}
}

func TestParseBoolWithInvalidValues(t *testing.T) {
	assert := assert.New(t)
	const envVar = "TEST_ENV_VAR"
	varValues := []string{"", "invalid_bool"}

	for _, envVarValue := range varValues {
		os.Setenv(envVar, envVarValue)
		assert.Panics(func() { utils.ParseBoolFromEnv(envVar) }, "Code path was expected to panic")
	}
}

func TestParseBoolWithEmptyEnvVariable(t *testing.T) {
	assert := assert.New(t)
	const envVar = "TEST_ENV_VAR"

	assert.Panics(func() { utils.ParseBoolFromEnv(envVar) }, "Code path was expected to panic")
}
