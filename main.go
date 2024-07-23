package main

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func createPodVPA(client dynamic.DynamicClient, podName string, targetNamespace string) {
	// Create VPA template for a pod
	pod := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "autoscaling.k8s.io/v1",
			"kind":       "VerticalPodAutoscaler",
			"metadata": map[string]interface{}{
				"name": podName,
			},
			"spec": map[string]interface{}{
				"targetRef": map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"name":       podName,
				},
				"updatePolicy": map[string]interface{}{
					"updateMode": "Off",
				},
			},
		},
	}

	_, err := client.Resource(schema.GroupVersionResource{
		Group:    "autoscaling.k8s.io",
		Version:  "v1",
		Resource: "verticalpodautoscalers",
	}).Namespace(targetNamespace).Create(context.TODO(), pod, metav1.CreateOptions{})

	if err != nil {
		log.Println("Error creating vpa resource", err)
	} else {
		log.Println("Successfully created vpa resource")
	}
}

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error creating in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// Create a pod informer
	informer := cache.NewSharedInformer(
		cache.NewListWatchFromClient(
			clientset.CoreV1().RESTClient(),
			"pods",
			metav1.NamespaceAll,
			fields.Everything(),
		),
		&v1.Pod{},
		0, // No resync period
	)

	// Create dynamic client to deal with VPA CRD
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	} else {
		log.Println("Dynamic client created")
	}

	// Register event handlers
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			fmt.Printf("Pod created: %s\n", pod.Name)
			createPodVPA(*client, pod.Name, pod.Namespace)
		},
		// Optionally handle update and delete events
		UpdateFunc: func(oldObj, newObj interface{}) {
			// Handle pod update
		},
		DeleteFunc: func(obj interface{}) {
			// Handle pod deletion
		},
	})

	// Start the informer
	stopCh := make(chan struct{})
	defer close(stopCh)
	go informer.Run(stopCh)

	// Wait forever
	select {}
}
