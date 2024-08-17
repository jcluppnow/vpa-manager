package config_test

import (
	"os"
	"testing"
	"vpa-manager/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnvVariablesWithEmptyString(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("ENABLE_CRONJOBS", "true")
	os.Setenv("ENABLE_DEPLOYMENTS", "false")
	os.Setenv("ENABLE_JOBS", "true")
	os.Setenv("ENABLE_PODS", "false")
	os.Setenv("UPDATE_MODE", "Off")
	os.Setenv("WATCHED_NAMESPACES", "")

	env := config.LoadEnv()

	expected := config.ControllerEnv{
		EnableCronjobs:    true,
		EnableDeployments: false,
		EnableJobs:        true,
		EnablePods:        false,
		UpdateMode:        "Off",
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
	os.Setenv("UPDATE_MODE", "Off")
	os.Setenv("WATCHED_NAMESPACES", "default, kube-system")

	env := config.LoadEnv()

	expected := config.ControllerEnv{
		EnableCronjobs:    true,
		EnableDeployments: false,
		EnableJobs:        true,
		EnablePods:        false,
		UpdateMode:        "Off",
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

func TestValidateControllerEnvPanics(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("ENABLE_CRONJOBS", "true")
	os.Setenv("ENABLE_DEPLOYMENTS", "false")
	os.Setenv("ENABLE_JOBS", "true")
	os.Setenv("ENABLE_PODS", "false")
	os.Setenv("UPDATE_MODE", "invalid-update-mode")
	os.Setenv("WATCHED_NAMESPACES", "default, kube-system")

	env := config.LoadEnv()

	assert.Panics(func() { config.ValidateControllerEnv(env) }, "Expected validate env to panic due to invalid update mode")
}

func TestValidateControllerEnv(t *testing.T) {
	validVPAUpdateModes := []string{"Auto", "Initial", "Recreate", "Off"}

	assert := assert.New(t)

	os.Setenv("ENABLE_CRONJOBS", "true")
	os.Setenv("ENABLE_DEPLOYMENTS", "false")
	os.Setenv("ENABLE_JOBS", "true")
	os.Setenv("ENABLE_PODS", "false")
	os.Setenv("WATCHED_NAMESPACES", "default, kube-system")

	for _, validUpdateMode := range validVPAUpdateModes {
		os.Setenv("UPDATE_MODE", validUpdateMode)
		env := config.LoadEnv()
		assert.NotPanics(func() { config.ValidateControllerEnv(env) }, "Expected validate env to panic due to invalid update mode")
	}
}
