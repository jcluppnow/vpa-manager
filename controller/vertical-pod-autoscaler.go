package controller

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func CreateVPA(client dynamic.DynamicClient, sourceResourceType string, resourceName string, targetNamespace string) {
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
					"updateMode": "Off",
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
		log.Println("Skipping creating VPA resource as it already exists for this resource type", sourceResourceType, resourceName, targetNamespace)
	} else {
		_, err := client.Resource(schema.GroupVersionResource{
			Group:    "autoscaling.k8s.io",
			Version:  "v1",
			Resource: "verticalpodautoscalers",
		}).Namespace(targetNamespace).Create(context.TODO(), vpaTemplate, metav1.CreateOptions{})

		if err != nil {
			log.Println("Error creating vpa resource", err, sourceResourceType, resourceName, targetNamespace)
		} else {
			log.Println("Successfully created vpa resource", sourceResourceType, resourceName, targetNamespace)
		}
	}
}

func DeleteVPA(client dynamic.DynamicClient, resourceName string, targetNamespace string) {
	err := client.Resource(schema.GroupVersionResource{
		Group:    "autoscaling.k8s.io",
		Version:  "v1",
		Resource: "verticalpodautoscalers",
	}).Namespace(targetNamespace).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})

	if err != nil {
		log.Println("Error deleting vpa resource", err, resourceName, targetNamespace)
	} else {
		log.Println("Successfully deleted vpa resource", resourceName, targetNamespace)
	}
}
