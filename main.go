package main

import (
	"encoding/json"
	"os"
	"strings"
	"sort"

	fmt "github.com/jhunt/go-ansi"
	"github.com/jhunt/go-cli"
	env "github.com/jhunt/go-envirotron"
	"github.com/jhunt/go-table"
	"github.com/jhunt/voom/client/voom"
)

var opt struct {
	Help    bool `cli:"-h, --help"`
	Version bool `cli:"-v, --version"`

	URL      string `cli:"-u, --url"      env:"VOOM_URL"`
	Username string `cli:"-U, --username" env:"VOOM_USERNAME"`
	Password string `cli:"-P, --password" env:"VOOM_PASSWORD"`

	Dump struct{} `cli:"dump"`
	List struct{} `cli:"ls"`
	Summary struct {} `cli:"sum, summary"`
}

func bail(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s:} @Y{%s}\n", msg, err)
		os.Exit(1)
	}
}

func main() {
	env.Override(&opt)
	command, args, err := cli.Parse(&opt)
	bail("Failed to parse command-line flags", err)

	if opt.Help {
		if command == "" || len(args) > 0 {
			fmt.Printf("USAGE: @G{voom} [options] @C{<command>} [options]\n")
			fmt.Printf("\n")
			printGlobalOptionsHelp()
			fmt.Printf("\n")
			printGlobalCommandsHelp()
			fmt.Printf("\n")
			os.Exit(0)
		}

		if command == "ls" {
			fmt.Printf("USAGE: @G{voom} [options] @C{ls}\n")
			fmt.Printf("\n")
			fmt.Printf("List VMs.\n")
			fmt.Printf("\n")
			printGlobalOptionsHelp()
			fmt.Printf("\n")
			os.Exit(0)
		}

		if command == "dump" {
			fmt.Printf("USAGE: @G{voom} [options] @C{dump}\n")
			fmt.Printf("\n")
			fmt.Printf("Dumps all found VMs to standard output, in JSON.\n")
			fmt.Printf("\n")
			printGlobalOptionsHelp()
			fmt.Printf("\n")
			os.Exit(0)
		}

		if command == "summary" {
			fmt.Printf("USAGE: @G{voom} [options] @C{summary}\n")
			fmt.Printf("\n")
			fmt.Printf("Summarize resource usage, by BOSH director and deployment.\n")
			fmt.Printf("\n")
			printGlobalOptionsHelp()
			fmt.Printf("\n")
			os.Exit(0)
		}
	}

	c, err := voom.Connect(opt.URL, opt.Username, opt.Password)
	bail("Failed to connect to vCenter endpoint", err)

	if command == "ls" {
		vms, err := c.VMs()
		bail("Failed to retrieve list of VMs", err)

		t := table.NewTable("ID", "Status", "IP Address", "Uptime", "CPUs", "Memory")
		for _, vm := range vms {
			pow := "off"
			if vm.On {
				pow = "on"
			}

			tags := []string{}
			for k, v := range vm.Tags {
				tags = append(tags, fmt.Sprintf("  %s = %s\n", k, v))
			}
			sort.Strings(tags)
			t.Row(vm, vm.ID+"\n"+strings.Join(tags, "")+"\n", pow, vm.IP, timeString(vm.Uptime),
				fmt.Sprintf("%d cores\n%dMHz used\n%dMHz demand", vm.CPUs, vm.CPUUsage, vm.CPUDemand),
				fmt.Sprintf("%d MB\n%d MB resv\n%d MB alloc", vm.MemoryUsage, vm.MemoryReserved, vm.MemoryUsed))
		}
		t.Output(os.Stdout)
		os.Exit(0)
	}

	if command == "dump" {
		vms, err := c.VMs()
		bail("Failed to retrieve list of VMs", err)

		var out struct {
			VMs []voom.VM `json:"vms"`
		}
		out.VMs = vms

		b, _ := json.Marshal(out)
		fmt.Printf("%s\n", string(b))
	}

	if command == "sum" {
		vms, err := c.VMs()
		bail("Failed to retrieve list of VMs", err)

		director := make(map[string]voom.VM)
		deployment := make(map[string]voom.VM)

		merge := func(dst voom.VM, src voom.VM) voom.VM {
			dst.MemoryUsed += src.MemoryUsed
			dst.MemoryReserved += src.MemoryReserved
			dst.CPUUsage += src.CPUUsage
			dst.CPUDemand += src.CPUDemand
			dst.MemoryUsage += src.MemoryUsage
			dst.CPUs += src.CPUs
			return dst
		}

		for _, vm := range vms {
			if !vm.On {
				continue
			}

			dir := vm.Tags["director"]
			if dir == "bosh-init" {
				dir = vm.Tags["deployment"] /* tricky ... */
			}
			dep := dir+"/"+vm.Tags["deployment"]

			if dir == "" {
				fmt.Fprintf(os.Stderr, "vm %s has no director!\n", vm.ID)
				continue
			}

			if _, ok := director[dir]; ok {
				director[dir] = merge(director[dir], vm)
			} else {
				director[dir] = vm
			}

			if _, ok := director[dep]; ok {
				deployment[dep] = merge(deployment[dep], vm)
			} else {
				deployment[dep] = vm
			}
		}

		keys := make([]string, 0)
		for k := range director {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, dir := range keys {
			vm := director[dir]
			fmt.Printf("BOSH Director :: @G{%s}\n", dir)
			fmt.Printf("  Memory: %dMB (%dMB resv / %dMB alloc)\n", vm.MemoryUsage, vm.MemoryUsed, vm.MemoryReserved)
			fmt.Printf("  CPU:    %d cores (%dMHz used / %dMHZ demand)\n", vm.CPUs, vm.CPUUsage, vm.CPUDemand)
			fmt.Printf("\n")
		}
		fmt.Printf("\n")

		keys = make([]string, 0)
		for k := range deployment {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, dep := range keys {
			vm := deployment[dep]
			fmt.Printf("BOSH Deployment :: @Y{%s}\n", dep)
			fmt.Printf("  Memory: %dMB (%dMB resv / %dMB alloc)\n", vm.MemoryUsage, vm.MemoryUsed, vm.MemoryReserved)
			fmt.Printf("  CPU:    %d cores (%dMHz used / %dMHZ demand)\n", vm.CPUs, vm.CPUUsage, vm.CPUDemand)
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}
}
