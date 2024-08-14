package controller

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnvVariablesWithEmptyString(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("ENABLE_CRONJOBS", "true")
	os.Setenv("ENABLE_DEPLOYMENTS", "false")
	os.Setenv("ENABLE_JOBS", "true")
	os.Setenv("ENABLE_PODS", "false")
	os.Setenv("WATCHED_NAMESPACES", "")

	env := LoadEnv()

	expected := ControllerEnv{
		EnableCronjobs:    true,
		EnableDeployments: false,
		EnableJobs:        true,
		EnablePods:        false,
		WatchedNamespaces: []string{},
	}

	assert.Equal(
		env,
		expected,
		"Environment variables did not match the expected values.\nActual: %+v\nExpected: %+v",
		env,
		expected,
	)
}

func TestLoadEnvVariablesWithNamespacesDefined(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("ENABLE_CRONJOBS", "true")
	os.Setenv("ENABLE_DEPLOYMENTS", "false")
	os.Setenv("ENABLE_JOBS", "true")
	os.Setenv("ENABLE_PODS", "false")
	os.Setenv("WATCHED_NAMESPACES", "default, kube-system")

	env := LoadEnv()

	expected := ControllerEnv{
		EnableCronjobs:    true,
		EnableDeployments: false,
		EnableJobs:        true,
		EnablePods:        false,
		WatchedNamespaces: []string{"default", "kube-system"},
	}

	assert.Equal(
		env,
		expected,
		"Environment variables did not match the expected values.\nActual: %+v\nExpected: %+v",
		env,
		expected,
	)
}
