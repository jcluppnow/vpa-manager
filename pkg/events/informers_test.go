package events_test

import (
	"net/http"
	"testing"
	"vpa-manager/pkg/config"
	"vpa-manager/pkg/events"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func TestValidConfigForCreatingInformers(t *testing.T) {
	assert := assert.New(t)

	controllerEnv := config.ControllerEnv{
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
	factory := events.CreateInformers(controllerEnv, config, fakeClientset)

	assert.NotNil(factory, "Expected that a valid Factory instance is returned after creating informers")
}

func TestInvalidConfigForCreatingInformers(t *testing.T) {
	assert := assert.New(t)
	var restConfig *rest.Config
	fakeClientset := fake.NewSimpleClientset()

	controllerEnv := config.ControllerEnv{
		EnableCronjobs:    false,
		EnableDeployments: true,
		EnableJobs:        false,
		EnablePods:        false,
		WatchedNamespaces: []string{},
	}

	assert.Panics(func() { events.CreateInformers(controllerEnv, restConfig, fakeClientset) }, "Code path was expected to panic due to undefined rest config")
}
