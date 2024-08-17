package controller_test

import (
	"net/http"
	"testing"
	"vpa-manager/pkg/controller"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func TestValidConfigForCreatingInformers(t *testing.T) {
	assert := assert.New(t)

	controllerEnv := controller.ControllerEnv{
		EnableCronjobs:    false,
		EnableDeployments: true,
		EnableJobs:        false,
		EnablePods:        false,
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

	fakeClientset := fake.NewSimpleClientset()
	factory := controller.CreateInformers(controllerEnv, config, fakeClientset)

	assert.NotNil(factory, "Expected that a valid Factory instance is returned after creating informers")
}

func TestInvalidConfigForCreatingInformers(t *testing.T) {
	assert := assert.New(t)
	var config *rest.Config
	fakeClientset := fake.NewSimpleClientset()

	controllerEnv := controller.ControllerEnv{
		EnableCronjobs:    false,
		EnableDeployments: true,
		EnableJobs:        false,
		EnablePods:        false,
		WatchedNamespaces: []string{},
	}

	assert.Panics(func() { controller.CreateInformers(controllerEnv, config, fakeClientset) }, "Code path was expected to panic due to undefined rest config")
}
