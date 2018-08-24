package voom

type VM struct {
	ID             string            `json:"id"`
	Uptime         int32             `json:"uptime"`
	Type           string            `json:"type"`
	IP             string            `json:"ip"`
	On             bool              `json:"on"`
	MemoryUsed     int32             `json:"mem_alloc"`
	MemoryReserved int32             `json:"mem_reserved"`
	CPUUsage       int32             `json:"cpu_usage"`
	CPUDemand      int32             `json:"cpu_demand"`
	MemoryUsage    int32             `json:"mem_usage"`
	CPUs           int32             `json:"cpus"`
	Tags           map[string]string `json:"tags"`
}
