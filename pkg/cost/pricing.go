package cost

// Pricing contains the pricing information for resources
type Pricing struct {
	CPUHourlyCost    float64 // Cost per CPU core per hour
	MemoryGBHourly   float64 // Cost per GB memory per hour
	GPUHourlyCost    float64 // Cost per GPU per hour
	StorageGBMonthly float64 // Cost per GB storage per month
}

// DefaultPricing returns default AWS-like pricing
// Based on typical m5.large pricing (~$0.096/hour)
func DefaultPricing() *Pricing {
	return &Pricing{
		CPUHourlyCost:    0.024, // ~$17.28/month per core
		MemoryGBHourly:   0.003, // ~$2.16/month per GB
		GPUHourlyCost:    0.90,  // ~$648/month per GPU (T4)
		StorageGBMonthly: 0.10,  // ~$0.10/month per GB (EBS gp3)
	}
}

// GCPPricing returns Google Cloud pricing
func GCPPricing() *Pricing {
	return &Pricing{
		CPUHourlyCost:    0.022, // n2-standard pricing
		MemoryGBHourly:   0.003,
		GPUHourlyCost:    0.85, // T4 GPU
		StorageGBMonthly: 0.10,
	}
}

// AzurePricing returns Azure pricing
func AzurePricing() *Pricing {
	return &Pricing{
		CPUHourlyCost:    0.025,
		MemoryGBHourly:   0.003,
		GPUHourlyCost:    0.95, // NC-series
		StorageGBMonthly: 0.12,
	}
}

// CalculateCPUCost calculates monthly cost for CPU cores
func (p *Pricing) CalculateCPUCost(cores float64) float64 {
	hoursPerMonth := 730.0 // Average hours in a month
	return cores * p.CPUHourlyCost * hoursPerMonth
}

// CalculateMemoryCost calculates monthly cost for memory
func (p *Pricing) CalculateMemoryCost(bytes int64) float64 {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	hoursPerMonth := 730.0
	return gb * p.MemoryGBHourly * hoursPerMonth
}

// CalculateGPUCost calculates monthly cost for GPUs
func (p *Pricing) CalculateGPUCost(count int) float64 {
	hoursPerMonth := 730.0
	return float64(count) * p.GPUHourlyCost * hoursPerMonth
}

// CalculateStorageCost calculates monthly cost for storage
func (p *Pricing) CalculateStorageCost(bytes int64) float64 {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	return gb * p.StorageGBMonthly
}
