package controller

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
)

type ApiDetails struct {
	Version string
	Kind    string
}

var (
	Scheme     = runtime.NewScheme()
	FakeClient = fake.NewSimpleDynamicClient(Scheme)

	ApiDetailsMap = map[string]ApiDetails{
		"CronJob":    {"batch/v1", "CronJob"},
		"Deployment": {"apps/v1", "Deployment"},
		"Job":        {"batch/v1", "Job"},
		"Pod":        {"v1", "Pod"},
	}
)

func TestCreateVPAForUnwatchedNamespace(t *testing.T) {
	const resourceName = "nginx"
	const resourceType = "Pod"
	const targetNamespace = "kube-system"
	var watchedNamespaces = []string{"default"}

	assert := assert.New(t)

	CreateVPA(FakeClient, watchedNamespaces, resourceType, resourceName, targetNamespace)

	_, err := FakeClient.Resource(schema.GroupVersionResource{
		Group:    "autoscaling.k8s.io",
		Version:  "v1",
		Resource: "verticalpodautoscalers",
	}).Namespace(targetNamespace).Get(context.TODO(), resourceName, metav1.GetOptions{})

	assert.NotNil(
		err,
		"Expected that VPA resource '%s' of type '%s' should not exist in namespace '%s' but it was found.",
		resourceName,
		resourceType,
		targetNamespace,
	)
}

func TestCreateVPAForWatchedNamespace(t *testing.T) {
	const targetNamespace = "kube-system"
	var watchedNamespaces = []string{"kube-system"}

	assert := assert.New(t)

	var index = 0
	for _, api := range ApiDetailsMap {
		var resourceName = "nginx-" + strconv.Itoa(index)

		CreateVPA(FakeClient, watchedNamespaces, api.Kind, resourceName, targetNamespace)

		fetchedResource, err := FakeClient.Resource(schema.GroupVersionResource{
			Group:    "autoscaling.k8s.io",
			Version:  "v1",
			Resource: "verticalpodautoscalers",
		}).Namespace(targetNamespace).Get(context.TODO(), resourceName, metav1.GetOptions{})

		expectedResource := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "autoscaling.k8s.io/v1",
				"kind":       "VerticalPodAutoscaler",
				"metadata": map[string]interface{}{
					"name":      resourceName,
					"namespace": targetNamespace,
				},
				"spec": map[string]interface{}{
					"targetRef": map[string]interface{}{
						"apiVersion": api.Version,
						"kind":       api.Kind,
						"name":       resourceName,
					},
					"updatePolicy": map[string]interface{}{
						"updateMode": "Off",
					},
				},
			},
		}

		assert.Equal(expectedResource, fetchedResource)
		assert.Nil(err, "Error returned when creating VPA")

		index++
	}
}

func TestDeleteVPAForUnwatchedNamespace(t *testing.T) {
	const resourceName = "nginx"
	const targetNamespace = "kube-system"
	var watchedNamespaces = []string{"default"}

	assert := assert.New(t)

	CreateVPA(FakeClient, []string{"kube-system"}, "Pod", resourceName, targetNamespace)
	DeleteVPA(FakeClient, watchedNamespaces, resourceName, targetNamespace)

	fetchedResource, err := FakeClient.Resource(schema.GroupVersionResource{
		Group:    "autoscaling.k8s.io",
		Version:  "v1",
		Resource: "verticalpodautoscalers",
	}).Namespace(targetNamespace).Get(context.TODO(), resourceName, metav1.GetOptions{})

	assert.Nil(err)
	assert.NotNil(
		fetchedResource,
		"Expected that VPA resource '%s' should exist in namespace '%s' but it was not found.",
		resourceName,
		targetNamespace,
	)
}

func TestDeleteVPAForWatchedNamespace(t *testing.T) {
	const targetNamespace = "kube-system"
	var watchedNamespaces = []string{"kube-system"}

	assert := assert.New(t)

	var index = 0
	for _, api := range ApiDetailsMap {
		var resourceName = "nginx-" + strconv.Itoa(index)

		CreateVPA(FakeClient, watchedNamespaces, api.Kind, resourceName, targetNamespace)

		_, err := FakeClient.Resource(schema.GroupVersionResource{
			Group:    "autoscaling.k8s.io",
			Version:  "v1",
			Resource: "verticalpodautoscalers",
		}).Namespace(targetNamespace).Get(context.TODO(), resourceName, metav1.GetOptions{})

		assert.Nil(err, "Expected VPA to be created")

		DeleteVPA(FakeClient, watchedNamespaces, resourceName, targetNamespace)

		fetchedResource, err := FakeClient.Resource(schema.GroupVersionResource{
			Group:    "autoscaling.k8s.io",
			Version:  "v1",
			Resource: "verticalpodautoscalers",
		}).Namespace(targetNamespace).Get(context.TODO(), resourceName, metav1.GetOptions{})

		assert.NotNil(err)
		assert.Nil(
			fetchedResource,
			"Expected that VPA resource '%s' should not exist in namespace '%s' but it was found.",
			resourceName,
			targetNamespace,
		)

		index++
	}
}
