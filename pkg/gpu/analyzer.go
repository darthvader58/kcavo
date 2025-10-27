package gpu

import (
	corev1 "k8s.io/api/core/v1"
)

// Analysis contains GPU usage analysis
type Analysis struct {
	Nodes           []NodeGPU
	Pods            []PodGPU
	TotalGPUs       int
	AllocatedGPUs   int
	AvailableGPUs   int
	UtilizationPct  float64
	Recommendations []string
}

// NodeGPU represents GPU information for a node
type NodeGPU struct {
	NodeName      string
	TotalGPUs     int
	AllocatedGPUs int
	AvailableGPUs int
	GPUType       string
}

// PodGPU represents GPU usage for a pod
type PodGPU struct {
	PodName   string
	Namespace string
	Node      string
	GPUCount  int
}

// Analyzer analyzes GPU resources
type Analyzer struct{}

// NewAnalyzer creates a new GPU analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// Analyze performs GPU analysis on nodes and pods
func (a *Analyzer) Analyze(nodes []corev1.Node, pods []corev1.Pod) Analysis {
	analysis := Analysis{
		Nodes:           make([]NodeGPU, 0),
		Pods:            make([]PodGPU, 0),
		Recommendations: make([]string, 0),
	}

	// Analyze nodes
	for _, node := range nodes {
		nodeGPU := a.analyzeNode(node)
		if nodeGPU.TotalGPUs > 0 {
			analysis.Nodes = append(analysis.Nodes, nodeGPU)
			analysis.TotalGPUs += nodeGPU.TotalGPUs
			analysis.AllocatedGPUs += nodeGPU.AllocatedGPUs
		}
	}

	// Analyze pods
	for _, pod := range pods {
		podGPU := a.analyzePod(pod)
		if podGPU.GPUCount > 0 {
			analysis.Pods = append(analysis.Pods, podGPU)
		}
	}

	analysis.AvailableGPUs = analysis.TotalGPUs - analysis.AllocatedGPUs
	if analysis.TotalGPUs > 0 {
		analysis.UtilizationPct = (float64(analysis.AllocatedGPUs) / float64(analysis.TotalGPUs)) * 100
	}

	// Generate recommendations
	analysis.Recommendations = a.generateRecommendations(analysis)

	return analysis
}

func (a *Analyzer) analyzeNode(node corev1.Node) NodeGPU {
	nodeGPU := NodeGPU{
		NodeName: node.Name,
		GPUType:  "Unknown",
	}

	// Get total GPUs from capacity
	if gpu, ok := node.Status.Capacity["nvidia.com/gpu"]; ok {
		nodeGPU.TotalGPUs = int(gpu.Value())
	}

	// Get allocated GPUs from allocatable (capacity - allocated = available)
	if gpu, ok := node.Status.Allocatable["nvidia.com/gpu"]; ok {
		nodeGPU.AvailableGPUs = int(gpu.Value())
		nodeGPU.AllocatedGPUs = nodeGPU.TotalGPUs - nodeGPU.AvailableGPUs
	}

	// Try to detect GPU type from node labels
	if gpuType, ok := node.Labels["nvidia.com/gpu.product"]; ok {
		nodeGPU.GPUType = gpuType
	} else if gpuType, ok := node.Labels["accelerator"]; ok {
		nodeGPU.GPUType = gpuType
	}

	return nodeGPU
}

func (a *Analyzer) analyzePod(pod corev1.Pod) PodGPU {
	podGPU := PodGPU{
		PodName:   pod.Name,
		Namespace: pod.Namespace,
		Node:      pod.Spec.NodeName,
	}

	// Count GPUs across all containers
	for _, container := range pod.Spec.Containers {
		if gpu, ok := container.Resources.Requests["nvidia.com/gpu"]; ok {
			podGPU.GPUCount += int(gpu.Value())
		}
		if gpu, ok := container.Resources.Limits["nvidia.com/gpu"]; ok {
			podGPU.GPUCount += int(gpu.Value())
		}
	}

	return podGPU
}

func (a *Analyzer) generateRecommendations(analysis Analysis) []string {
	recommendations := make([]string, 0)

	// Low utilization
	if analysis.TotalGPUs > 0 && analysis.UtilizationPct < 50 {
		recommendations = append(recommendations,
			"GPU utilization is below 50%. Consider scaling down GPU nodes or consolidating workloads.")
	}

	// High utilization
	if analysis.UtilizationPct > 85 {
		recommendations = append(recommendations,
			"GPU utilization is above 85%. Consider adding more GPU nodes to prevent scheduling issues.")
	}

	// Fragmented GPUs
	fragmentedNodes := 0
	for _, node := range analysis.Nodes {
		if node.AllocatedGPUs > 0 && node.AvailableGPUs > 0 {
			fragmentedNodes++
		}
	}
	if fragmentedNodes > len(analysis.Nodes)/2 {
		recommendations = append(recommendations,
			"Many nodes have partially allocated GPUs. Consider using node affinity to pack GPU workloads efficiently.")
	}

	// No GPUs but could use them
	if analysis.TotalGPUs == 0 {
		recommendations = append(recommendations,
			"No GPU resources detected. If you have ML/AI workloads, consider adding GPU nodes for better performance.")
	}

	// Single GPU pods that could share
	singleGPUPods := 0
	for _, pod := range analysis.Pods {
		if pod.GPUCount == 1 {
			singleGPUPods++
		}
	}
	if singleGPUPods > 2 {
		recommendations = append(recommendations,
			"Multiple pods requesting single GPUs. Consider MIG (Multi-Instance GPU) or time-slicing for better utilization.")
	}

	return recommendations
}
