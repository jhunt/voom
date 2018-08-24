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

		if command == "sum" {
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

		sum := NewSummary()
		for _, vm := range vms {
			if !vm.On {
				continue
			}

			dir := vm.Tags["director"]
			dep := vm.Tags["deployment"]
			if dir == "bosh-init" {
				dir = dep /* tricky ... */
			}

			if dir == "" {
				fmt.Fprintf(os.Stderr, "vm %s has no director!\n", vm.ID)
				continue
			}

			sum.Breakout(dir).Breakout(dep).Ingest(vm)
		}

		t := table.NewTable("", "cores", "  compute  ", "  memory  ", "   disk   ")
		t.Row(nil, "ALL",
			fmt.Sprintf("% 5d", sum.Cores),
			fmt.Sprintf("% 7.1f GHz", float64(sum.Compute) / 1024.0),
			fmt.Sprintf("% 7.1f GB", float64(sum.Memory) / 1024.0),
			fmt.Sprintf("% 7.1f GB", float64(sum.Disk) / 1024.0))
		t.Row(nil, "---")

		for _, name := range sum.Keys() {
			bosh := sum.Breakout(name)
			t.Row(nil, name,
				fmt.Sprintf("% 5d", bosh.Cores),
				fmt.Sprintf("% 7.1f GHz", float64(bosh.Compute) / 1024.0),
				fmt.Sprintf("% 7.1f GB", float64(bosh.Memory) / 1024.0),
				fmt.Sprintf("% 7.1f GB", float64(bosh.Disk) / 1024.0))

			for _, name := range bosh.Keys() {
				deployment := bosh.Breakout(name)
				t.Row(nil, "   "+name,
					fmt.Sprintf("% 5d", deployment.Cores),
					fmt.Sprintf("% 7.1f    ", float64(deployment.Compute) / 1024.0),
					fmt.Sprintf("% 7.1f   ", float64(deployment.Memory) / 1024.0),
					fmt.Sprintf("% 7.1f   ", float64(deployment.Disk) / 1024.0))
			}
			t.Row(nil, "---")
		}
		t.Output(os.Stdout)
	}
}
