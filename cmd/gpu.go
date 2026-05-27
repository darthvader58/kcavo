package cmd

import (
	"context"
	"fmt"

	"kcavo/pkg/gpu"
	"kcavo/pkg/kubernetes"
	"kcavo/pkg/visualize"

	"github.com/spf13/cobra"
)

var gpuCmd = &cobra.Command{
	Use:   "gpu",
	Short: "Analyze GPU resource usage and scheduling",
	Long: `Analyze GPU allocation, utilization, and scheduling across your cluster.
	
Provides insights on:
  • GPU allocation per pod/node
  • Underutilized GPUs
  • GPU scheduling recommendations
  • Cost per GPU

Examples:
  kubectl cost gpu                    # Analyze GPU usage
  kubectl cost gpu -A                 # All namespaces`,
	RunE: runGPU,
}

func init() {
	rootCmd.AddCommand(gpuCmd)
}

func runGPU(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := kubernetes.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	ns := getNamespace()

	if output == "table" {
		fmt.Println("Analyzing GPU resources...")
		fmt.Println()
	}

	// Get nodes with GPUs
	nodes, err := client.GetNodes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}

	// Get pods
	pods, err := client.GetPods(ctx, ns)
	if err != nil {
		return fmt.Errorf("failed to get pods: %w", err)
	}

	// Analyze GPU usage
	analyzer := gpu.NewAnalyzer()
	analysis := analyzer.Analyze(nodes, pods)

	switch output {
	case "json":
		return visualize.PrintJSON(analysis)
	case "yaml":
		return visualize.PrintYAML(analysis)
	}

	// Display results
	visualize.PrintGPUTable(analysis)

	// Print recommendations
	fmt.Println()
	fmt.Println("Recommendations:")
	for i, rec := range analysis.Recommendations {
		fmt.Printf("   %d. %s\n", i+1, rec)
	}

	if len(analysis.Recommendations) == 0 {
		fmt.Println("   No GPU optimization recommendations at this time.")
	}

	return nil
}
