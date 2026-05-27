package cmd

import (
	"context"
	"fmt"
	"sort"

	"kcavo/pkg/cost"
	"kcavo/pkg/kubernetes"
	"kcavo/pkg/visualize"

	"github.com/spf13/cobra"
)

var (
	showBreakdown bool
	sortBy        string
	topN          int
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze costs across your Kubernetes cluster",
	Long: `Analyze costs for pods, nodes, and resources in your cluster.
	
Provides detailed cost breakdowns by:
  • Namespace
  • Pod
  • Node
  • Resource type (CPU, Memory, GPU, Storage)
  
Examples:
  kubectl cost analyze                                    # Analyze current namespace
  kubectl cost analyze -A                                 # Analyze all namespaces
  kubectl cost analyze -n production --breakdown         # Show detailed breakdown
  kubectl cost analyze --sort-by cost --top 10          # Top 10 most expensive`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateSortField(sortBy)
	},
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().BoolVar(&showBreakdown, "breakdown", false, "show detailed cost breakdown")
	analyzeCmd.Flags().StringVar(&sortBy, "sort-by", "cost", "sort by: cost, cpu, memory, gpu")
	analyzeCmd.Flags().IntVar(&topN, "top", 0, "show only top N results (0 = all)")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Initialize Kubernetes client
	client, err := kubernetes.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	ns := getNamespace()

	if output == "table" {
		fmt.Print("Analyzing costs")
		if ns == "" {
			fmt.Print(" across all namespaces...\n")
		} else {
			fmt.Printf(" in namespace: %s...\n", ns)
		}
	}

	// Get pods
	pods, err := client.GetPods(ctx, ns)
	if err != nil {
		return fmt.Errorf("failed to get pods: %w", err)
	}

	// Get nodes
	nodes, err := client.GetNodes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}

	// Calculate costs
	calculator := cost.NewCalculatorWithPricing(configuredPricing())
	results := calculator.CalculatePodCosts(pods, nodes)
	sortPodCosts(results, sortBy)

	// Apply filters
	if topN > 0 && len(results) > topN {
		results = results[:topN]
	}

	// Display results
	switch output {
	case "json":
		return visualize.PrintJSON(results)
	case "yaml":
		return visualize.PrintYAML(results)
	default:
		visualize.PrintCostTable(results, showBreakdown)
	}

	// Print summary
	fmt.Println()
	printSummary(results)

	return nil
}

func printSummary(results []cost.PodCost) {
	var totalCost, totalCPU, totalMemory float64
	var totalGPU int

	for _, r := range results {
		totalCost += r.TotalCost
		totalCPU += r.CPUCost
		totalMemory += r.MemoryCost
		totalGPU += r.GPUCount
	}

	fmt.Println("Summary:")
	fmt.Printf("   Total Monthly Cost: $%.2f\n", totalCost)
	fmt.Printf("   Total Pods: %d\n", len(results))
	if totalGPU > 0 {
		fmt.Printf("   Total GPUs: %d\n", totalGPU)
	}
	fmt.Printf("   CPU Cost: $%.2f (%.1f%%)\n", totalCPU, percentage(totalCPU, totalCost))
	fmt.Printf("   Memory Cost: $%.2f (%.1f%%)\n", totalMemory, percentage(totalMemory, totalCost))
}

func sortPodCosts(results []cost.PodCost, field string) {
	sort.SliceStable(results, func(i, j int) bool {
		switch field {
		case "cpu":
			return results[i].CPUCost > results[j].CPUCost
		case "memory":
			return results[i].MemoryCost > results[j].MemoryCost
		case "gpu":
			return results[i].GPUCost > results[j].GPUCost
		default:
			return results[i].TotalCost > results[j].TotalCost
		}
	})
}

func validateSortField(field string) error {
	switch field {
	case "cost", "cpu", "memory", "gpu":
		return nil
	default:
		return fmt.Errorf("unsupported sort field %q; use cost, cpu, memory, or gpu", field)
	}
}

func percentage(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (part / total) * 100
}
