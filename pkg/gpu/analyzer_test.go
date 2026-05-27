package gpu

import (
	"testing"

	"kcavo/pkg/cost"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAnalyzeCalculatesAllocatedGPUsFromScheduledPods(t *testing.T) {
	nodes := []corev1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "gpu-node"},
			Status: corev1.NodeStatus{
				Capacity: corev1.ResourceList{
					cost.NVIDIAResourceGPU: resource.MustParse("4"),
				},
				Allocatable: corev1.ResourceList{
					cost.NVIDIAResourceGPU: resource.MustParse("4"),
				},
			},
		},
	}
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
							Requests: corev1.ResourceList{cost.NVIDIAResourceGPU: resource.MustParse("2")},
						},
					},
				},
			},
		},
	}

	analysis := NewAnalyzer().Analyze(nodes, pods)

	if analysis.AllocatedGPUs != 2 {
		t.Fatalf("expected 2 allocated GPUs, got %d", analysis.AllocatedGPUs)
	}
	if analysis.AvailableGPUs != 2 {
		t.Fatalf("expected 2 available GPUs, got %d", analysis.AvailableGPUs)
	}
	if analysis.UtilizationPct != 50 {
		t.Fatalf("expected 50%% utilization, got %.1f", analysis.UtilizationPct)
	}
}
