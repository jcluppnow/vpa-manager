package vpa

import (
	"context"
	"log/slog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func isWatchedNamespace(watchedNamespaces []string, targetNamespace string) bool {
	if len(watchedNamespaces) == 0 {
		return true
	}

	for index := range watchedNamespaces {
		if watchedNamespaces[index] == targetNamespace {
			return true
		}
	}

	return false
}

func CreateVPA(client dynamic.Interface, watchedNamespaces []string, sourceResourceType string, resourceName string, targetNamespace string, updateMode string) {
	if !isWatchedNamespace(watchedNamespaces, targetNamespace) {
		return
	}

	type ApiDetails struct {
		version string
		kind    string
	}

	apiDetails := map[string]ApiDetails{
		"CronJob":    {"batch/v1", "CronJob"},
		"Deployment": {"apps/v1", "Deployment"},
		"Job":        {"batch/v1", "Job"},
		"Pod":        {"v1", "Pod"},
	}[sourceResourceType]

	vpaTemplate := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "autoscaling.k8s.io/v1",
			"kind":       "VerticalPodAutoscaler",
			"metadata": map[string]interface{}{
				"name": resourceName,
			},
			"spec": map[string]interface{}{
				"targetRef": map[string]interface{}{
					"apiVersion": apiDetails.version,
					"kind":       apiDetails.kind,
					"name":       resourceName,
				},
				"updatePolicy": map[string]interface{}{
					"updateMode": updateMode,
				},
			},
		},
	}

	_, getErr := client.Resource(schema.GroupVersionResource{
		Group:    "autoscaling.k8s.io",
		Version:  "v1",
		Resource: "verticalpodautoscalers",
	}).Namespace(targetNamespace).Get(context.TODO(), resourceName, metav1.GetOptions{})

	if getErr == nil {
		slog.Info(
			"Skipping creating VPA resource as it already exists for this resource type",
			"sourceResourceType", sourceResourceType,
			"resourceName", resourceName,
			"targetNamespace", targetNamespace,
		)
	} else {
		_, err := client.Resource(schema.GroupVersionResource{
			Group:    "autoscaling.k8s.io",
			Version:  "v1",
			Resource: "verticalpodautoscalers",
		}).Namespace(targetNamespace).Create(context.TODO(), vpaTemplate, metav1.CreateOptions{})

		if err != nil {
			slog.Warn(
				"Error creating vpa resource",
				"error", err,
				"sourceResourceType", sourceResourceType,
				"resourceName", resourceName,
				"targetNamespace", targetNamespace,
			)
		} else {
			slog.Info(
				"Successfully created vpa resource",
				"sourceResourceType", sourceResourceType,
				"resourceName", resourceName,
				"targetNamespace", targetNamespace,
			)
		}
	}
}

func DeleteVPA(client dynamic.Interface, watchedNamespaces []string, resourceName string, targetNamespace string) {
	if !isWatchedNamespace(watchedNamespaces, targetNamespace) {
		return
	}

	err := client.Resource(schema.GroupVersionResource{
		Group:    "autoscaling.k8s.io",
		Version:  "v1",
		Resource: "verticalpodautoscalers",
	}).Namespace(targetNamespace).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})

	if err != nil {
		slog.Warn(
			"Error deleting vpa resource",
			"error", err,
			"resourceName", resourceName,
			"targetNamespace", targetNamespace,
		)
	} else {
		slog.Info(
			"Successfully deleted vpa resource",
			"resourceName", resourceName,
			"targetNamespace", targetNamespace,
		)
	}
}
