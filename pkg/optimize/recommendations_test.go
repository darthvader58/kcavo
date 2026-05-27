package optimize

import (
	"testing"

	"kcavo/pkg/cost"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAnalyzeMatchesCostsByPodKey(t *testing.T) {
	pods := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "small", Namespace: "default"},
			Status:     corev1.PodStatus{Phase: corev1.PodRunning},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "app",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("500m"),
								corev1.ResourceMemory: resource.MustParse("512Mi"),
							},
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "large", Namespace: "default"},
			Status:     corev1.PodStatus{Phase: corev1.PodRunning},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "app",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("8"),
								corev1.ResourceMemory: resource.MustParse("32Gi"),
							},
						},
					},
				},
			},
		},
	}
	costs := []cost.PodCost{
		{Name: "large", Namespace: "default", TotalCost: 300},
		{Name: "small", Namespace: "default", TotalCost: 10},
	}

	recommendations := NewOptimizer().Analyze(pods, nil, costs)

	for _, rec := range recommendations {
		if rec.Category == "Rightsizing" {
			if rec.Savings != 90 {
				t.Fatalf("expected rightsizing savings to use large pod cost, got %.2f", rec.Savings)
			}
			return
		}
	}
	t.Fatal("expected rightsizing recommendation")
}
