package cmd

import (
	"context"
	"fmt"

	"kubectl-cost/pkg/kubernetes"
	"kubectl-cost/pkg/visualize"

	"github.com/spf13/cobra"
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
  • deployments
  • services
  • all (default)

Examples:
  kubectl cost visualize                     # Visualize all resources
  kubectl cost visualize --type pods         # Show only pods
  kubectl cost visualize -A                  # All namespaces`,
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

	fmt.Printf("📦 Visualizing resources")
	if ns == "" {
		fmt.Printf(" across all namespaces...\n\n")
	} else {
		fmt.Printf(" in namespace: %s...\n\n", ns)
	}

	if resourceType == "all" || resourceType == "nodes" {
		nodes, err := client.GetNodes(ctx)
		if err != nil {
			return fmt.Errorf("failed to get nodes: %w", err)
		}
		visualize.PrintNodeTable(nodes)
		fmt.Println()
	}

	if resourceType == "all" || resourceType == "pods" {
		pods, err := client.GetPods(ctx, ns)
		if err != nil {
			return fmt.Errorf("failed to get pods: %w", err)
		}
		visualize.PrintPodTable(pods)
		fmt.Println()
	}

	return nil
}
