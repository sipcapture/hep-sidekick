package sidekick

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Sidekick manages the operations of attaching to pods.
type Sidekick struct {
	Clientset    kubernetes.Interface
	Selector     string
	HomerAddress string
}

// New creates a new Sidekick instance.
func New(clientset kubernetes.Interface, selector, homerAddress string) *Sidekick {
	return &Sidekick{
		Clientset:    clientset,
		Selector:     selector,
		HomerAddress: homerAddress,
	}
}

// ListPods finds pods based on the configured selector.
func (s *Sidekick) ListPods(ctx context.Context) (*corev1.PodList, error) {
	pods, err := s.Clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		LabelSelector: s.Selector,
	})
	if err != nil {
		return nil, fmt.Errorf("error listing pods with selector '%s': %w", s.Selector, err)
	}
	return pods, nil
}

// createHeplifyPodSpec defines the pod spec for the heplify sidekick.
func (s *Sidekick) createHeplifyPodSpec(targetPod *corev1.Pod) *corev1.Pod {
	privileged := true
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("heplify-%s", targetPod.Name),
			Namespace: targetPod.Namespace,
			Labels: map[string]string{
				"app":                         "heplify",
				"hep-sidekick/target-pod":     targetPod.Name,
				"hep-sidekick/target-namespace": targetPod.Namespace,
			},
		},
		Spec: corev1.PodSpec{
			NodeName:      targetPod.Spec.NodeName,
			HostNetwork:   true, // To sniff traffic from the node
			Containers: []corev1.Container{
				{
					Name:  "heplify",
					Image: "ghcr.io/sipcapture/heplify",
					Args: []string{
						"-i", "any",
						"-t", "pcap",
						"-hs", s.HomerAddress,
					},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &privileged, // Required for packet sniffing
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
}

// AttachToPod creates and starts a heplify pod to sniff a target pod.
func (s *Sidekick) AttachToPod(ctx context.Context, targetPod *corev1.Pod) (*corev1.Pod, error) {
	heplifyPod := s.createHeplifyPodSpec(targetPod)
	
	createdPod, err := s.Clientset.CoreV1().Pods(heplifyPod.Namespace).Create(ctx, heplifyPod, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not create heplify pod: %w", err)
	}

	return createdPod, nil
} 