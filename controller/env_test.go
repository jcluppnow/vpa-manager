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
