package main

import (
	"sort"

	"github.com/jhunt/voom/client/voom"
)

type Summary struct {
	VMs             int
	Cores           int
	MemoryAllocated int
	MemoryUsed      int
	Compute         int
	DiskAllocated   int
	DiskUsed        int
	DiskFree        int

	Parent *Summary
	Sub    map[string]*Summary
}

func NewSummary() *Summary {
	return &Summary{
		Sub: make(map[string]*Summary),
	}
}

func (s *Summary) Ingest(vm voom.VM) {
	s.VMs = s.VMs + 1
	s.Cores = s.Cores + int(vm.CPUs)
	s.MemoryAllocated = s.MemoryAllocated + int(vm.MemoryAllocated)
	s.MemoryUsed = s.MemoryUsed + int(vm.MemoryUsed)
	s.Compute = s.Compute + int(vm.CPUUsage)
	s.DiskAllocated = s.DiskAllocated + int(vm.DiskAllocated)
	s.DiskFree = s.DiskFree + int(vm.DiskFree)
	s.DiskUsed = s.DiskUsed + int(vm.DiskUsed)

	if s.Parent != nil {
		s.Parent.Ingest(vm)
	}
}

func (s *Summary) Breakout(key string) *Summary {
	if _, ok := s.Sub[key]; !ok {
		s.Sub[key] = NewSummary()
		s.Sub[key].Parent = s
	}
	return s.Sub[key]
}

func (s *Summary) Keys() []string {
	l := make([]string, 0)
	for k := range s.Sub {
		l = append(l, k)
	}
	sort.Strings(l)
	return l
}
