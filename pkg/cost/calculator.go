package cost

import (
	"sort"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// PodCost represents the cost breakdown for a pod
type PodCost struct {
	Name       string
	Namespace  string
	Node       string
	CPUCost    float64
	MemoryCost float64
	GPUCost    float64
	GPUCount   int
	TotalCost  float64
	CPURequest string
	MemRequest string
	CPULimit   string
	MemLimit   string
}

// Calculator handles cost calculations
type Calculator struct {
	pricing *Pricing
}

// NewCalculator creates a new cost calculator
func NewCalculator() *Calculator {
	return &Calculator{
		pricing: DefaultPricing(),
	}
}

// NewCalculatorWithPricing creates a calculator with custom pricing
func NewCalculatorWithPricing(pricing *Pricing) *Calculator {
	return &Calculator{
		pricing: pricing,
	}
}

// CalculatePodCosts calculates costs for all pods
func (c *Calculator) CalculatePodCosts(pods []corev1.Pod, nodes []corev1.Node) []PodCost {
	results := make([]PodCost, 0, len(pods))

	for _, pod := range pods {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}

		cost := c.calculatePodCost(pod)
		results = append(results, cost)
	}

	// Sort by total cost (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalCost > results[j].TotalCost
	})

	return results
}

// calculatePodCost calculates the cost for a single pod
func (c *Calculator) calculatePodCost(pod corev1.Pod) PodCost {
	var cpuRequest, memRequest, cpuLimit, memLimit resource.Quantity
	gpuCount := 0

	// Sum up all container resources
	for _, container := range pod.Spec.Containers {
		if req, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
			cpuRequest.Add(req)
		}
		if req, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
			memRequest.Add(req)
		}
		if lim, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
			cpuLimit.Add(lim)
		}
		if lim, ok := container.Resources.Limits[corev1.ResourceMemory]; ok {
			memLimit.Add(lim)
		}

		// Check for GPU requests
		if gpu, ok := container.Resources.Requests["nvidia.com/gpu"]; ok {
			gpuCount += int(gpu.Value())
		}
		if gpu, ok := container.Resources.Limits["nvidia.com/gpu"]; ok {
			gpuCount += int(gpu.Value())
		}
	}

	// Calculate costs based on requests (or limits if requests not set)
	cpuToUse := cpuRequest
	if cpuToUse.IsZero() {
		cpuToUse = cpuLimit
	}
	memToUse := memRequest
	if memToUse.IsZero() {
		memToUse = memLimit
	}

	cpuCost := c.pricing.CalculateCPUCost(cpuToUse.AsApproximateFloat64())
	memCost := c.pricing.CalculateMemoryCost(memToUse.Value())
	gpuCost := c.pricing.CalculateGPUCost(gpuCount)

	return PodCost{
		Name:       pod.Name,
		Namespace:  pod.Namespace,
		Node:       pod.Spec.NodeName,
		CPUCost:    cpuCost,
		MemoryCost: memCost,
		GPUCost:    gpuCost,
		GPUCount:   gpuCount,
		TotalCost:  cpuCost + memCost + gpuCost,
		CPURequest: cpuRequest.String(),
		MemRequest: memRequest.String(),
		CPULimit:   cpuLimit.String(),
		MemLimit:   memLimit.String(),
	}
}

// CalculateNodeCost calculates the total cost for a node
func (c *Calculator) CalculateNodeCost(node corev1.Node) float64 {
	cpu := node.Status.Capacity[corev1.ResourceCPU]
	mem := node.Status.Capacity[corev1.ResourceMemory]

	cpuCost := c.pricing.CalculateCPUCost(cpu.AsApproximateFloat64())
	memCost := c.pricing.CalculateMemoryCost(mem.Value())

	// Check for GPUs
	gpuCount := 0
	if gpu, ok := node.Status.Capacity["nvidia.com/gpu"]; ok {
		gpuCount = int(gpu.Value())
	}
	gpuCost := c.pricing.CalculateGPUCost(gpuCount)

	return cpuCost + memCost + gpuCost
}
