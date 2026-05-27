package cmd

import (
	"context"
	"fmt"

	"kcavo/pkg/kubernetes"
	"kcavo/pkg/visualize"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var (
	resourceType string
)

var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize cluster resources",
	Long: `Visualize Kubernetes resources in a tree or table format.
	
Supported resource types:
  • pods
  • nodes
  • all (default)

Examples:
  kubectl cost visualize                     # Visualize all resources
  kubectl cost visualize --type pods         # Show only pods
  kubectl cost visualize -A                  # All namespaces`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateResourceType(resourceType)
	},
	RunE: runVisualize,
}

func init() {
	rootCmd.AddCommand(visualizeCmd)

	visualizeCmd.Flags().StringVar(&resourceType, "type", "all", "resource type to visualize")
}

func runVisualize(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := kubernetes.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	ns := getNamespace()

	if output == "table" {
		fmt.Print("Visualizing resources")
		if ns == "" {
			fmt.Printf(" across all namespaces...\n\n")
		} else {
			fmt.Printf(" in namespace: %s...\n\n", ns)
		}
	}

	response := make(map[string]interface{})

	if resourceType == "all" || resourceType == "nodes" {
		nodes, err := client.GetNodes(ctx)
		if err != nil {
			return fmt.Errorf("failed to get nodes: %w", err)
		}
		if output == "table" {
			visualize.PrintNodeTable(nodes)
			fmt.Println()
		} else {
			response["nodes"] = summarizeNodes(nodes)
		}
	}

	if resourceType == "all" || resourceType == "pods" {
		pods, err := client.GetPods(ctx, ns)
		if err != nil {
			return fmt.Errorf("failed to get pods: %w", err)
		}
		if output == "table" {
			visualize.PrintPodTable(pods)
			fmt.Println()
		} else {
			response["pods"] = summarizePods(pods)
		}
	}

	switch output {
	case "json":
		return visualize.PrintJSON(response)
	case "yaml":
		return visualize.PrintYAML(response)
	}

	return nil
}

type nodeSummary struct {
	Name   string `json:"name" yaml:"name"`
	Ready  bool   `json:"ready" yaml:"ready"`
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
	Pods   string `json:"pods" yaml:"pods"`
}

type podSummary struct {
	Name          string `json:"name" yaml:"name"`
	Namespace     string `json:"namespace" yaml:"namespace"`
	Status        string `json:"status" yaml:"status"`
	Node          string `json:"node" yaml:"node"`
	CPURequest    string `json:"cpuRequest" yaml:"cpuRequest"`
	MemoryRequest string `json:"memoryRequest" yaml:"memoryRequest"`
}

func summarizeNodes(nodes []corev1.Node) []nodeSummary {
	summaries := make([]nodeSummary, 0, len(nodes))
	for _, node := range nodes {
		ready := false
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady {
				ready = condition.Status == corev1.ConditionTrue
				break
			}
		}

		cpu := node.Status.Capacity[corev1.ResourceCPU]
		mem := node.Status.Capacity[corev1.ResourceMemory]
		podCapacity := node.Status.Capacity[corev1.ResourcePods]
		summaries = append(summaries, nodeSummary{
			Name:   node.Name,
			Ready:  ready,
			CPU:    cpu.String(),
			Memory: fmt.Sprintf("%dGi", mem.Value()/(1024*1024*1024)),
			Pods:   podCapacity.String(),
		})
	}
	return summaries
}

func summarizePods(pods []corev1.Pod) []podSummary {
	summaries := make([]podSummary, 0, len(pods))
	for _, pod := range pods {
		var cpuReq, memReq int64
		for _, container := range pod.Spec.Containers {
			if req, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
				cpuReq += req.MilliValue()
			}
			if req, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
				memReq += req.Value()
			}
		}

		summaries = append(summaries, podSummary{
			Name:          pod.Name,
			Namespace:     pod.Namespace,
			Status:        string(pod.Status.Phase),
			Node:          pod.Spec.NodeName,
			CPURequest:    fmt.Sprintf("%dm", cpuReq),
			MemoryRequest: fmt.Sprintf("%dMi", memReq/(1024*1024)),
		})
	}
	return summaries
}

func validateResourceType(resourceType string) error {
	switch resourceType {
	case "all", "pods", "nodes":
		return nil
	default:
		return fmt.Errorf("unsupported resource type %q; use all, pods, or nodes", resourceType)
	}
}
