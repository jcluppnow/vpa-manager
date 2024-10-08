package controller_test

import (
	"context"
	"net/http"
	"testing"

	"vpa-manager/pkg/config"
	"vpa-manager/pkg/controller"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func TestNewController(t *testing.T) {
	assert := assert.New(t)

	controllerEnv := config.ControllerEnv{
		EnableCronjobs:    false,
		EnableDeployments: true,
		EnableJobs:        false,
		EnablePods:        false,
		UpdateMode:        "Off",
		WatchedNamespaces: []string{},
	}

	config := &rest.Config{
		Host:    "https://localhost",
		APIPath: "/api",
		ContentConfig: rest.ContentConfig{
			GroupVersion:         &v1.SchemeGroupVersion,
			NegotiatedSerializer: scheme.Codecs,
		},
		Transport: &http.Transport{}, // Dummy transport
		QPS:       10,
		Burst:     20,
		UserAgent: rest.DefaultKubernetesUserAgent(),
	}

	fakeClientSet := fake.NewSimpleClientset()
	controller := controller.NewController(controllerEnv, config, fakeClientSet)

	assert.NotNil(controller, "Expected controller to be created correctly")
}

func TestNewControllerWithInvalidEnv(t *testing.T) {
	assert := assert.New(t)

	controllerEnv := config.ControllerEnv{
		EnableCronjobs:    false,
		EnableDeployments: true,
		EnableJobs:        false,
		EnablePods:        false,
		UpdateMode:        "invalid-update-mode",
		WatchedNamespaces: []string{},
	}

	config := &rest.Config{
		Host:    "https://localhost",
		APIPath: "/api",
		ContentConfig: rest.ContentConfig{
			GroupVersion:         &v1.SchemeGroupVersion,
			NegotiatedSerializer: scheme.Codecs,
		},
		Transport: &http.Transport{}, // Dummy transport
		QPS:       10,
		Burst:     20,
		UserAgent: rest.DefaultKubernetesUserAgent(),
	}

	fakeClientSet := fake.NewSimpleClientset()
	assert.Panics(func() { controller.NewController(controllerEnv, config, fakeClientSet) }, "Controller was expected to panic with invalid environment variables")
}

func TestControllerRun(t *testing.T) {
	assert := assert.New(t)

	controllerEnv := config.ControllerEnv{
		EnableCronjobs:    false,
		EnableDeployments: true,
		EnableJobs:        false,
		EnablePods:        false,
		UpdateMode:        "Off",
		WatchedNamespaces: []string{},
	}

	config := &rest.Config{
		Host:    "https://localhost",
		APIPath: "/api",
		ContentConfig: rest.ContentConfig{
			GroupVersion:         &v1.SchemeGroupVersion,
			NegotiatedSerializer: scheme.Codecs,
		},
		Transport: &http.Transport{}, // Dummy transport
		QPS:       10,
		Burst:     20,
		UserAgent: rest.DefaultKubernetesUserAgent(),
	}

	fakeClientSet := fake.NewSimpleClientset()
	controller := controller.NewController(controllerEnv, config, fakeClientSet)
	ctx, cancel := context.WithCancel(context.Background())

	assert.NotPanics(
		func() {
			err := controller.Run(ctx)
			cancel()
			assert.Nil(err, "Expected no errors to be thrown while running controller")
		},
		"Controller run was expected to run without panicking",
	)
}
