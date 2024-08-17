package controller

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func newControllerEnv() ControllerEnv {
	return ControllerEnv{
		EnableCronjobs:    false,
		EnableDeployments: true,
		EnableJobs:        false,
		EnablePods:        false,
		WatchedNamespaces: []string{},
	}
}

func TestNewController(t *testing.T) {
	assert := assert.New(t)

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
	controller := NewController(newControllerEnv(), config, fakeClientSet)

	assert.NotNil(controller, "Expected controller to be created correctly")
}

func TestControllerRun(t *testing.T) {
	assert := assert.New(t)
	fakeClientSet := fake.NewSimpleClientset()

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

	controller := NewController(newControllerEnv(), config, fakeClientSet)
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
