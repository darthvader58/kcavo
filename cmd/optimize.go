package cmd

import (
	"context"
	"fmt"

	"kcavo/pkg/cost"
	"kcavo/pkg/kubernetes"
	"kcavo/pkg/optimize"
	viz "kcavo/pkg/visualize"

	"github.com/spf13/cobra"
)

var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Get cost optimization recommendations",
	Long: `Analyze your cluster and get actionable recommendations to reduce costs.
	
Recommendations include:
  • Rightsizing pods (over-provisioned resources)
  • Unused resources
  • GPU optimization
  • Spot instance opportunities
  • Resource quotas

Examples:
  kubectl cost optimize               # Get recommendations
  kubectl cost optimize -A            # Cluster-wide analysis`,
	RunE: runOptimize,
}

func init() {
	rootCmd.AddCommand(optimizeCmd)
}

func runOptimize(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := kubernetes.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	ns := getNamespace()

	if output == "table" {
		fmt.Println("Analyzing cluster for cost optimization opportunities...")
		fmt.Println()
	}

	// Get resources
	pods, err := client.GetPods(ctx, ns)
	if err != nil {
		return fmt.Errorf("failed to get pods: %w", err)
	}

	nodes, err := client.GetNodes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}

	// Calculate current costs
	pricing := configuredPricing()
	calculator := cost.NewCalculatorWithPricing(pricing)
	costs := calculator.CalculatePodCosts(pods, nodes)

	// Get optimization recommendations
	optimizer := optimize.NewOptimizerWithPricing(pricing)
	recommendations := optimizer.Analyze(pods, nodes, costs)

	switch output {
	case "json":
		return viz.PrintJSON(recommendations)
	case "yaml":
		return viz.PrintYAML(recommendations)
	}

	// Display recommendations
	fmt.Println("Optimization Recommendations:")
	fmt.Println()

	totalSavings := 0.0
	for i, rec := range recommendations {
		fmt.Printf("   %d. %s\n", i+1, rec.Title)
		fmt.Printf("      %s\n", rec.Description)
		fmt.Printf("      Potential savings: $%.2f/month\n", rec.Savings)
		fmt.Printf("      Priority: %s\n", rec.Priority)
		fmt.Println()
		totalSavings += rec.Savings
	}

	if len(recommendations) == 0 {
		fmt.Println("   No optimization opportunities found.")
	} else {
		fmt.Printf("Total Potential Savings: $%.2f/month (%.1f%% reduction)\n",
			totalSavings, calculateSavingsPercentage(costs, totalSavings))
	}

	return nil
}

func calculateSavingsPercentage(costs []cost.PodCost, savings float64) float64 {
	var totalCost float64
	for _, c := range costs {
		totalCost += c.TotalCost
	}
	if totalCost == 0 {
		return 0
	}
	return (savings / totalCost) * 100
}
