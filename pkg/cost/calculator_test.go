package cost

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCalculatePodCostsUsesGPURequestBeforeLimit(t *testing.T) {
	pods := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "trainer", Namespace: "ml"},
			Status:     corev1.PodStatus{Phase: corev1.PodRunning},
			Spec: corev1.PodSpec{
				NodeName: "gpu-node",
				Containers: []corev1.Container{
					{
						Name: "worker",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{NVIDIAResourceGPU: resource.MustParse("1")},
							Limits:   corev1.ResourceList{NVIDIAResourceGPU: resource.MustParse("1")},
						},
					},
				},
			},
		},
	}

	results := NewCalculator().CalculatePodCosts(pods, nil)

	if len(results) != 1 {
		t.Fatalf("expected 1 pod cost, got %d", len(results))
	}
	if results[0].GPUCount != 1 {
		t.Fatalf("expected GPU count to avoid double-counting request and limit, got %d", results[0].GPUCount)
	}
}

func TestCalculatePodCostsSkipsNonRunningPods(t *testing.T) {
	pods := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "pending", Namespace: "default"},
			Status:     corev1.PodStatus{Phase: corev1.PodPending},
		},
	}

	results := NewCalculator().CalculatePodCosts(pods, nil)

	if len(results) != 0 {
		t.Fatalf("expected non-running pods to be skipped, got %d results", len(results))
	}
}
