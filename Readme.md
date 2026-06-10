# KCAVO

KCAVO is a `kubectl` plugin for Kubernetes cost visibility. It estimates monthly pod costs from resource requests, reports GPU allocation, and provides optimization recommendations from cluster state.

KCAVO works with only `kubectl` access — no Prometheus, no Helm, no in-cluster agents. 
Just point it at a cluster and get cost estimates in seconds.

This plugin is an estimator, not a billing system. Estimates are based on Kubernetes resource 
requests (falling back to limits when requests are unset), and are not a substitute for 
cloud billing exports or metrics-backed tools like Kubecost/OpenCost.

## Features

- Cost estimates for running pods based on CPU, memory, and GPU requests
- Cluster resource tables for pods and nodes
- GPU allocation summaries by node and pod
- Optimization checks for missing requests, large requests, underused nodes, and GPU-heavy workloads
- Table, JSON, and YAML output
- Configurable pricing with provider presets

## Installation

```bash
make install
kubectl cost --version
```

`make install` builds and installs the plugin as `kubectl-cost` in `~/.local/bin`. Make sure `~/.local/bin` is on your `PATH`.

## Quick Start

```bash
# Analyze costs in the current namespace
kubectl cost analyze

# Analyze all namespaces
kubectl cost analyze --all-namespaces

# Show detailed cost columns
kubectl cost analyze --breakdown

# Show the top 10 pods by estimated monthly cost
kubectl cost analyze --top 10 --sort-by cost

# Visualize pods and nodes
kubectl cost visualize

# Analyze GPU allocation
kubectl cost gpu

# Get optimization recommendations
kubectl cost optimize
```

## Commands

### `kubectl cost analyze`

Estimates monthly cost for running pods.

```bash
kubectl cost analyze
kubectl cost analyze -A
kubectl cost analyze -n production --breakdown
kubectl cost analyze --sort-by cpu --top 10
kubectl cost analyze -o json
kubectl cost analyze -o yaml
```

Supported sort fields: `cost`, `cpu`, `memory`, `gpu`.

### `kubectl cost visualize`

Prints cluster resource tables.

```bash
kubectl cost visualize
kubectl cost visualize --type pods
kubectl cost visualize --type nodes
kubectl cost visualize -A
```

### `kubectl cost gpu`

Shows GPU capacity, scheduled GPU requests, available GPUs, and GPU pods.

```bash
kubectl cost gpu
kubectl cost gpu -A
```

### `kubectl cost optimize`

Generates lightweight recommendations from Kubernetes resource specs.

```bash
kubectl cost optimize
kubectl cost optimize -A
```

## Configuration

Create `~/.kubectl-cost.yaml` to customize pricing:

```yaml
pricing:
  cpu_hourly: 0.024
  memory_gb_hourly: 0.003
  gpu_hourly: 0.90
  storage_gb_monthly: 0.10
```

You can also use a provider preset:

```yaml
provider: gcp # aws, gcp, or azure
```

The default pricing is AWS-like. For production reporting, replace these values with your actual blended or negotiated rates.

## Cost Model

KCAVO calculates estimated monthly cost as:

```text
Monthly cost = (CPU cores * CPU hourly rate + memory GiB * memory hourly rate + GPUs * GPU hourly rate) * 730
```

Resource selection rules:

- CPU and memory use requests first, then limits if requests are unset.
- GPU uses `nvidia.com/gpu` requests first, then limits if requests are unset.
- Only running pods are included in pod cost estimates.

## Development

```bash
make fmt
make test
make build
```

## Acknowledgements

- [Cobra](https://github.com/spf13/cobra) for the CLI framework
- [tablewriter](https://github.com/olekukonko/tablewriter) for terminal tables

Made with <3 by [Shashwat Raj](shashwatraj.com)
