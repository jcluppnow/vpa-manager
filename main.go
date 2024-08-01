package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func createVPA(client dynamic.DynamicClient, sourceResourceType string, resourceName string, targetNamespace string) {
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

	// Create VPA template for specified resource
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

func deleteVPA(client dynamic.DynamicClient, resourceName string, targetNamespace string) {
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

func isTargetNamespace(targetNamespaces []string, namespace string) bool {
	if len(targetNamespaces) == 0 {
		return true
	}

	for index := range targetNamespaces {
		if targetNamespaces[index] == namespace {
			return true
		}
	}

	return false
}

func createListeners(targetNamespaces []string) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error creating in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	}

	factory := informers.NewSharedInformerFactory(clientset, time.Minute)
	cronJobInformer := factory.Batch().V1().CronJobs().Informer()
	deploymentInformer := factory.Apps().V1().Deployments().Informer()
	jobInformer := factory.Batch().V1().Jobs().Informer()
	podInformer := factory.Core().V1().Pods().Informer()

	cronJobInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cronJob := obj.(*batchv1.CronJob)
			if isTargetNamespace(targetNamespaces, cronJob.Namespace) {
				createVPA(*client, "CronJob", cronJob.Name, cronJob.Namespace)
			}
		},
		DeleteFunc: func(obj interface{}) {
			cronJob := obj.(*batchv1.CronJob)
			if isTargetNamespace(targetNamespaces, cronJob.Namespace) {
				deleteVPA(*client, cronJob.Name, cronJob.Namespace)
			}
		},
	})

	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deployment := obj.(*appsv1.Deployment)
			if isTargetNamespace(targetNamespaces, deployment.Namespace) {
				createVPA(*client, "Deployment", deployment.Name, deployment.Namespace)
			}
		},
		DeleteFunc: func(obj interface{}) {
			deployment := obj.(*batchv1.CronJob)
			if isTargetNamespace(targetNamespaces, deployment.Namespace) {
				deleteVPA(*client, deployment.Name, deployment.Namespace)
			}
		},
	})

	jobInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			job := obj.(*batchv1.Job)
			if len(job.OwnerReferences) == 0 && isTargetNamespace(targetNamespaces, job.Namespace) {
				createVPA(*client, "Job", job.Name, job.Namespace)
			}
		},
		DeleteFunc: func(obj interface{}) {
			job := obj.(*batchv1.Job)
			if len(job.OwnerReferences) == 0 && isTargetNamespace(targetNamespaces, job.Namespace) {
				deleteVPA(*client, job.Name, job.Namespace)
			}
		},
	})

	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			if len(pod.OwnerReferences) == 0 && isTargetNamespace(targetNamespaces, pod.Namespace) {
				createVPA(*client, "Pod", pod.Name, pod.Namespace)
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			if len(pod.OwnerReferences) == 0 && isTargetNamespace(targetNamespaces, pod.Namespace) {
				deleteVPA(*client, pod.Name, pod.Namespace)
			}
		},
	})

	factory.Start(wait.NeverStop)
	factory.WaitForCacheSync(wait.NeverStop)
}

func main() {
	targetNamespaces := os.Getenv("TARGET_NAMESPACES")
	formattedNamespaces := []string{}

	if targetNamespaces != "" {
		formattedNamespaces = strings.Split(targetNamespaces, ",")
	}

	// Setup event listeners
	createListeners(formattedNamespaces)

	// Wait forever
	select {}
}
