package controller_test

import (
	"os"
	"testing"
	"vpa-manager/controller"

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

	env := controller.LoadEnv()

	expected := controller.ControllerEnv{
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

	env := controller.LoadEnv()

	expected := controller.ControllerEnv{
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

	env := controller.LoadEnv()

	assert.Panics(func() { controller.ValidateControllerEnv(env) }, "Expected validate env to panic due to invalid update mode")
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
		env := controller.LoadEnv()
		assert.NotPanics(func() { controller.ValidateControllerEnv(env) }, "Expected validate env to panic due to invalid update mode")
	}
}
