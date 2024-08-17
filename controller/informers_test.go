package controller

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var (
	Env = ControllerEnv{
		EnableCronjobs:    false,
		EnableDeployments: true,
		EnableJobs:        false,
		EnablePods:        false,
		WatchedNamespaces: []string{},
	}

	FakeClientSet = fake.NewSimpleClientset()
)

func TestValidConfigForCreatingInformers(t *testing.T) {
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

	factory := CreateInformers(Env, config, FakeClientSet)

	assert.NotNil(factory, "Expected that a valid Factory instance is returned after creating informers")
}

func TestInvalidConfigForCreatingInformers(t *testing.T) {
	assert := assert.New(t)
	var config *rest.Config

	assert.Panics(func() { CreateInformers(Env, config, FakeClientSet) }, "Code path was expected to panic due to undefined rest config")
}
