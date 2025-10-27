package optimize

import (
	"fmt"
	"kubectl-cost/pkg/cost"

	corev1 "k8s.io/api/core/v1"
)

// Recommendation represents a cost optimization recommendation
type Recommendation struct {
	Title       string
	Description string
	Savings     float64
	Priority    string // High, Medium, Low
	Category    string // Rightsizing, Unused, GPU, Spot, etc.
}

// Optimizer generates cost optimization recommendations
type Optimizer struct {
	pricing *cost.Pricing
}

// NewOptimizer creates a new optimizer
func NewOptimizer() *Optimizer {
	return &Optimizer{
		pricing: cost.DefaultPricing(),
	}
}

// Analyze generates optimization recommendations
func (o *Optimizer) Analyze(pods []corev1.Pod, nodes []corev1.Node, costs []cost.PodCost) []Recommendation {
	recommendations := make([]Recommendation, 0)

	// Check for over-provisioned pods
	recommendations = append(recommendations, o.findOverProvisionedPods(pods, costs)...)

	// Check for pods without resource requests
	recommendations = append(recommendations, o.findPodsWithoutRequests(pods)...)

	// Check for unused resources
	recommendations = append(recommendations, o.findUnusedResources(nodes)...)

	// Check for expensive GPU usage
	recommendations = append(recommendations, o.findExpensiveGPUUsage(pods, costs)...)

	// Sort by savings (highest first)
	sortRecommendationsBySavings(recommendations)

	return recommendations
}

func (o *Optimizer) findOverProvisionedPods(pods []corev1.Pod, costs []cost.PodCost) []Recommendation {
	recommendations := make([]Recommendation, 0)

	for i, pod := range pods {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}

		// Check if requests are much larger than typical usage
		// This is a simplified check - in production, you'd use metrics
		for _, container := range pod.Spec.Containers {
			cpuReq := container.Resources.Requests[corev1.ResourceCPU]
			memReq := container.Resources.Requests[corev1.ResourceMemory]

			if cpuReq.IsZero() || memReq.IsZero() {
				continue
			}

			// Assume pods requesting >4 cores or >16GB might be over-provisioned
			if cpuReq.AsApproximateFloat64() > 4.0 || memReq.Value() > 16*1024*1024*1024 {
				if i < len(costs) {
					savings := costs[i].TotalCost * 0.3 // Estimate 30% savings from rightsizing
					recommendations = append(recommendations, Recommendation{
						Title:       "Rightsize over-provisioned pod: " + pod.Name,
						Description: "This pod requests significant resources. Consider rightsizing based on actual usage metrics.",
						Savings:     savings,
						Priority:    "High",
						Category:    "Rightsizing",
					})
				}
			}
		}
	}

	return recommendations
}

func (o *Optimizer) findPodsWithoutRequests(pods []corev1.Pod) []Recommendation {
	recommendations := make([]Recommendation, 0)
	count := 0

	for _, pod := range pods {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}

		hasRequests := false
		for _, container := range pod.Spec.Containers {
			if len(container.Resources.Requests) > 0 {
				hasRequests = true
				break
			}
		}

		if !hasRequests {
			count++
		}
	}

	if count > 0 {
		recommendations = append(recommendations, Recommendation{
			Title:       "Add resource requests to pods without them",
			Description: fmt.Sprintf("%d pods don't have resource requests. This can lead to poor scheduling and cost visibility.", count),
			Savings:     0, // Hard to estimate without knowing actual usage
			Priority:    "Medium",
			Category:    "Best Practice",
		})
	}

	return recommendations
}

func (o *Optimizer) findUnusedResources(nodes []corev1.Node) []Recommendation {
	recommendations := make([]Recommendation, 0)

	// Check for nodes with low allocation
	// In production, you'd compare allocatable vs allocated using metrics
	for _, node := range nodes {
		allocatable := node.Status.Allocatable
		capacity := node.Status.Capacity

		cpuAllocatable := allocatable[corev1.ResourceCPU]
		cpuCapacity := capacity[corev1.ResourceCPU]

		// Simplified check - if allocatable is close to capacity, node might be underutilized
		if cpuAllocatable.AsApproximateFloat64() > cpuCapacity.AsApproximateFloat64()*0.8 {
			// This is a simplified estimation
			nodeCost := o.estimateNodeCost(node)
			savings := nodeCost * 0.5 // Estimate 50% savings if node can be removed/downsized

			recommendations = append(recommendations, Recommendation{
				Title:       "Consider downsizing or removing underutilized node: " + node.Name,
				Description: "This node appears to have low resource allocation. Review if it can be consolidated or removed.",
				Savings:     savings,
				Priority:    "Medium",
				Category:    "Unused",
			})
		}
	}

	return recommendations
}

func (o *Optimizer) findExpensiveGPUUsage(pods []corev1.Pod, costs []cost.PodCost) []Recommendation {
	recommendations := make([]Recommendation, 0)

	for i, pod := range pods {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}

		gpuCount := 0
		for _, container := range pod.Spec.Containers {
			if gpu, ok := container.Resources.Requests["nvidia.com/gpu"]; ok {
				gpuCount += int(gpu.Value())
			}
		}

		if gpuCount > 0 && i < len(costs) {
			// If GPU cost is more than 70% of total, it's a significant expense
			if costs[i].GPUCost > costs[i].TotalCost*0.7 {
				recommendations = append(recommendations, Recommendation{
					Title:       "Review GPU usage for pod: " + pod.Name,
					Description: "This pod uses GPUs which account for most of its cost. Ensure GPU is being utilized efficiently or consider spot instances.",
					Savings:     costs[i].GPUCost * 0.5, // Estimate 50% savings with spot
					Priority:    "High",
					Category:    "GPU",
				})
			}
		}
	}

	return recommendations
}

func (o *Optimizer) estimateNodeCost(node corev1.Node) float64 {
	cpu := node.Status.Capacity[corev1.ResourceCPU]
	mem := node.Status.Capacity[corev1.ResourceMemory]

	cpuCost := o.pricing.CalculateCPUCost(cpu.AsApproximateFloat64())
	memCost := o.pricing.CalculateMemoryCost(mem.Value())

	return cpuCost + memCost
}

func sortRecommendationsBySavings(recommendations []Recommendation) {
	// Simple bubble sort for small arrays
	for i := 0; i < len(recommendations)-1; i++ {
		for j := 0; j < len(recommendations)-i-1; j++ {
			if recommendations[j].Savings < recommendations[j+1].Savings {
				recommendations[j], recommendations[j+1] = recommendations[j+1], recommendations[j]
			}
		}
	}
}
