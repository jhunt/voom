package voom

type VM struct {
	ID              string            `json:"id"`
	Uptime          int32             `json:"uptime"`
	Type            string            `json:"type"`
	IP              string            `json:"ip"`
	On              bool              `json:"on"`
	MemoryAllocated int32             `json:"mem_allocated"`
	MemoryReserved  int32             `json:"mem_reserved"`
	MemoryUsed      int32             `json:"mem_used"`
	CPUUsage        int32             `json:"cpu_usage"`
	CPUDemand       int32             `json:"cpu_demand"`
	CPUs            int32             `json:"cpus"`
	DiskAllocated   int64             `json:"disk_allocated"`
	DiskUsed        int64             `json:"disk_used"`
	DiskFree        int64             `json:"disk_free"`
	Tags            map[string]string `json:"tags"`
}
