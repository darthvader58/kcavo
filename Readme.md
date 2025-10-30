# KCAVO - Kubernetes Cluster Analyse Visualize and Optimize
A kubectl or Kubernetes Command Line Tool based plugin that visualizes resources, analyzes costs, and provides optimization recommendations. 

## What it helps you with?
- **Cost Analysis:** Calculate and visualize costs for pods, nodes, and resources
- **GPU Analysis:** Specialized GPU resource tracking and optimization
- **Resource Visualization:** View cluster resources in clean tables
- **Smart Recommendations:** Get actionable cost-saving suggestions

## Install
```bash
make install
kubectl cost --version
```

## 🚀 Quick Start

```bash
# Analyze costs in current namespace
kubectl cost analyze

# Analyze across all namespaces
kubectl cost analyze --all-namespaces

# Visualize cluster resources
kubectl cost visualize

# Analyze GPU usage
kubectl cost gpu

# Get optimization recommendations
kubectl cost optimize
```

## Commands

### `kubectl cost analyze`

Analyze costs for pods and resources in your cluster.

```bash
# Basic usage
kubectl cost analyze

# Across all namespaces
kubectl cost analyze -A

# With proper breakdown
kubectl cost analyze --breakdown

# Show top 10 expensive
kubectl cost analyze --top 10 --sort-by cost

# Different output formats
kubectl cost analyze -o json
kubectl cost analyze -o yaml
```

### `kubectl cost visualize`

Visualize Kubernetes resources in table format.

```bash
# Visualize all resources
kubectl cost visualize

# Specific resource type
kubectl cost visualize --type pods
kubectl cost visualize --type nodes

# All namespaces
kubectl cost visualize -A
```

### `kubectl cost gpu`

Analyze GPU resource allocation and usage.

```bash
# Analyze GPU usage
kubectl cost gpu

# All namespaces
kubectl cost gpu -A
```

### `kubectl cost optimize`

Get cost optimization recommendations.

```bash
# Get recommendations
kubectl cost optimize

# All namespaces
kubectl cost optimize -A
```

## Configuration

Create `~/.kubectl-cost.yaml` to customize pricing or use existing cloud providers' pricing:

```yaml
pricing:
  cpu_hourly: 0.024           # $17.28/month per core
  memory_gb_hourly: 0.003     # $2.16/month per GB
  gpu_hourly: 0.90            # $648/month per GPU
  storage_gb_monthly: 0.10    # $0.10/month per GB

# Or use cloud provider presets
provider: aws  # aws, gcp, or azure
```

## Calculation for Costs

Costs are calculated based on:
- **CPU**: Resource requests (or limits if requests not set)
- **Memory**: Resource requests (or limits if requests not set)
- **GPU**: GPU resource requests
- **Time**: Monthly basis (730 hours/month)

Formula:
```
Monthly Cost = (CPU_cores × CPU_rate + Memory_GB × Memory_rate + GPU_count × GPU_rate) × 730_hours
```

## Acknowledgements
https://github.com/spf13/cobra for CLI.
https://github.com/olekukonko/tablewriter for the tables.

