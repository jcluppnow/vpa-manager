package controller

import (
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

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

func CreateInformers(targetNamespaces []string, resourcesToManage ResourcesToManage) {
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

	jobInformer := factory.Batch().V1().Jobs().Informer()
	podInformer := factory.Core().V1().Pods().Informer()

	if resourcesToManage.Cronjobs {
		cronJobInformer := factory.Batch().V1().CronJobs().Informer()
		cronJobInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				cronJob := obj.(*batchv1.CronJob)
				if isTargetNamespace(targetNamespaces, cronJob.Namespace) {
					CreateVPA(*client, "CronJob", cronJob.Name, cronJob.Namespace)
				}
			},
			DeleteFunc: func(obj interface{}) {
				cronJob := obj.(*batchv1.CronJob)
				if isTargetNamespace(targetNamespaces, cronJob.Namespace) {
					DeleteVPA(*client, cronJob.Name, cronJob.Namespace)
				}
			},
		})
	}

	if resourcesToManage.Deployments {
		deploymentInformer := factory.Apps().V1().Deployments().Informer()
		deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				deployment := obj.(*appsv1.Deployment)
				if isTargetNamespace(targetNamespaces, deployment.Namespace) {
					CreateVPA(*client, "Deployment", deployment.Name, deployment.Namespace)
				}
			},
			DeleteFunc: func(obj interface{}) {
				deployment := obj.(*batchv1.CronJob)
				if isTargetNamespace(targetNamespaces, deployment.Namespace) {
					DeleteVPA(*client, deployment.Name, deployment.Namespace)
				}
			},
		})
	}

	if resourcesToManage.Jobs {
		jobInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				job := obj.(*batchv1.Job)
				if len(job.OwnerReferences) == 0 && isTargetNamespace(targetNamespaces, job.Namespace) {
					CreateVPA(*client, "Job", job.Name, job.Namespace)
				}
			},
			DeleteFunc: func(obj interface{}) {
				job := obj.(*batchv1.Job)
				if len(job.OwnerReferences) == 0 && isTargetNamespace(targetNamespaces, job.Namespace) {
					DeleteVPA(*client, job.Name, job.Namespace)
				}
			},
		})
	}

	if resourcesToManage.Pods {
		podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				if len(pod.OwnerReferences) == 0 && isTargetNamespace(targetNamespaces, pod.Namespace) {
					CreateVPA(*client, "Pod", pod.Name, pod.Namespace)
				}
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				if len(pod.OwnerReferences) == 0 && isTargetNamespace(targetNamespaces, pod.Namespace) {
					DeleteVPA(*client, pod.Name, pod.Namespace)
				}
			},
		})
	}

	factory.Start(wait.NeverStop)
	factory.WaitForCacheSync(wait.NeverStop)
}
