package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sipcapture/hep-sidekick/pkg/sidekick"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	selector := flag.String("selector", "hep-sidekick/enabled=true", "Label selector to find pods to attach to.")
	homerAddress := flag.String("homer-address", "127.0.0.1:9060", "Address of the HOMER server.")
	flag.Parse()

	log.Printf("Using selector: %s", *selector)
	log.Printf("Using HOMER address: %s", *homerAddress)

	config, err := getConfig()
	if err != nil {
		log.Fatalf("Error getting Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes clientset: %v", err)
	}

	sk := sidekick.New(clientset, *selector, *homerAddress)
	pods, err := sk.ListPods(context.Background())
	if err != nil {
		log.Fatalf("Error listing pods: %v", err)
	}

	log.Printf("Found %d pods with selector '%s'", len(pods.Items), *selector)
	for i := range pods.Items {
		pod := &pods.Items[i]
		log.Printf("-> Attaching to pod %s/%s on node %s", pod.Namespace, pod.Name, pod.Spec.NodeName)
		createdPod, err := sk.AttachToPod(context.Background(), pod)
		if err != nil {
			log.Printf("   Error attaching to pod %s/%s: %v", pod.Namespace, pod.Name, err)
			continue
		}
		log.Printf("   Successfully created heplify pod %s/%s", createdPod.Namespace, createdPod.Name)
	}
}

func getConfig() (*rest.Config, error) {
	// First, try to get the in-cluster configuration
	config, err := rest.InClusterConfig()
	if err == nil {
		log.Println("Using in-cluster Kubernetes config")
		return config, nil
	}
	log.Printf("Not in a cluster, attempting to use local kubeconfig: %v", err)

	// If that fails, try the local kubeconfig from the home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get user home directory: %w", err)
	}
	kubeconfigPath := filepath.Join(home, ".kube", "config")

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("error creating client config from flags: %w", err)
	}

	log.Println("Using local kubeconfig from", kubeconfigPath)
	return config, nil
} 