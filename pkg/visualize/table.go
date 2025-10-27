package visualize

import (
	"encoding/json"
	"fmt"
	"os"

	"kubectl-cost/pkg/cost"
	"kubectl-cost/pkg/gpu"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
)

// PrintCostTable prints costs in a formatted table
func PrintCostTable(costs []cost.PodCost, showBreakdown bool) {
	table := tablewriter.NewWriter(os.Stdout)

	if showBreakdown {
		table.SetHeader([]string{"Pod", "Namespace", "Node", "CPU Cost", "Memory Cost", "GPU Cost", "Total Cost"})
	} else {
		table.SetHeader([]string{"Pod", "Namespace", "Total Cost"})
	}

	table.SetBorder(true)
	table.SetRowLine(false)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(true)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, c := range costs {
		if showBreakdown {
			table.Append([]string{
				c.Name,
				c.Namespace,
				c.Node,
				fmt.Sprintf("$%.2f", c.CPUCost),
				fmt.Sprintf("$%.2f", c.MemoryCost),
				fmt.Sprintf("$%.2f", c.GPUCost),
				fmt.Sprintf("$%.2f", c.TotalCost),
			})
		} else {
			table.Append([]string{
				c.Name,
				c.Namespace,
				fmt.Sprintf("$%.2f/mo", c.TotalCost),
			})
		}
	}

	table.Render()
}

// PrintGPUTable prints GPU analysis in a table
func PrintGPUTable(analysis gpu.Analysis) {
	// Nodes table
	fmt.Println("🖥️  GPU Nodes:")
	if len(analysis.Nodes) == 0 {
		fmt.Println("   No GPU nodes found in cluster")
		return
	}

	nodeTable := tablewriter.NewWriter(os.Stdout)
	nodeTable.SetHeader([]string{"Node", "GPU Type", "Total", "Allocated", "Available", "Utilization"})
	nodeTable.SetBorder(false)
	nodeTable.SetHeaderLine(true)
	nodeTable.SetTablePadding("\t")
	nodeTable.SetNoWhiteSpace(true)

	for _, node := range analysis.Nodes {
		util := 0.0
		if node.TotalGPUs > 0 {
			util = (float64(node.AllocatedGPUs) / float64(node.TotalGPUs)) * 100
		}
		nodeTable.Append([]string{
			node.NodeName,
			node.GPUType,
			fmt.Sprintf("%d", node.TotalGPUs),
			fmt.Sprintf("%d", node.AllocatedGPUs),
			fmt.Sprintf("%d", node.AvailableGPUs),
			fmt.Sprintf("%.1f%%", util),
		})
	}
	nodeTable.Render()

	// Pods table
	fmt.Println()
	fmt.Println("🎮 GPU Pods:")
	if len(analysis.Pods) == 0 {
		fmt.Println("   No pods with GPU requests found")
		return
	}

	podTable := tablewriter.NewWriter(os.Stdout)
	podTable.SetHeader([]string{"Pod", "Namespace", "Node", "GPUs"})
	podTable.SetBorder(false)
	podTable.SetHeaderLine(true)
	podTable.SetTablePadding("\t")
	podTable.SetNoWhiteSpace(true)

	for _, pod := range analysis.Pods {
		podTable.Append([]string{
			pod.PodName,
			pod.Namespace,
			pod.Node,
			fmt.Sprintf("%d", pod.GPUCount),
		})
	}
	podTable.Render()

	// Summary
	fmt.Println()
	fmt.Printf("📊 GPU Summary:\n")
	fmt.Printf("   Total GPUs: %d\n", analysis.TotalGPUs)
	fmt.Printf("   Allocated: %d\n", analysis.AllocatedGPUs)
	fmt.Printf("   Available: %d\n", analysis.AvailableGPUs)
	fmt.Printf("   Utilization: %.1f%%\n", analysis.UtilizationPct)
}

// PrintNodeTable prints nodes in a table
func PrintNodeTable(nodes []corev1.Node) {
	fmt.Println("🖥️  Nodes:")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Status", "CPU", "Memory", "Pods"})
	table.SetBorder(false)
	table.SetHeaderLine(true)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, node := range nodes {
		status := "Ready"
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status != corev1.ConditionTrue {
				status = "Not Ready"
			}
		}

		cpu := node.Status.Capacity[corev1.ResourceCPU]
		mem := node.Status.Capacity[corev1.ResourceMemory]
		pods := node.Status.Capacity[corev1.ResourcePods]

		memGB := mem.Value() / (1024 * 1024 * 1024)

		table.Append([]string{
			node.Name,
			status,
			cpu.String(),
			fmt.Sprintf("%dGi", memGB),
			pods.String(),
		})
	}

	table.Render()
}

// PrintPodTable prints pods in a table
func PrintPodTable(pods []corev1.Pod) {
	fmt.Println("📦 Pods:")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Namespace", "Status", "Node", "CPU Request", "Memory Request"})
	table.SetBorder(false)
	table.SetHeaderLine(true)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, pod := range pods {
		// Sum up resources across containers
		var cpuReq, memReq int64
		for _, container := range pod.Spec.Containers {
			if req, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
				cpuReq += req.MilliValue()
			}
			if req, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
				memReq += req.Value()
			}
		}

		cpuStr := fmt.Sprintf("%dm", cpuReq)
		if cpuReq == 0 {
			cpuStr = "-"
		}

		memStr := fmt.Sprintf("%dMi", memReq/(1024*1024))
		if memReq == 0 {
			memStr = "-"
		}

		table.Append([]string{
			pod.Name,
			pod.Namespace,
			string(pod.Status.Phase),
			pod.Spec.NodeName,
			cpuStr,
			memStr,
		})
	}

	table.Render()
}

// PrintJSON prints data as JSON
func PrintJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// PrintYAML prints data as YAML
func PrintYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	return encoder.Encode(data)
}
