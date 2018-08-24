package main

import (
	fmt "github.com/jhunt/go-ansi"
)

func printGlobalOptionsHelp() {
	fmt.Printf("GLOBAL OPTIONS\n")
	fmt.Printf("\n")
	fmt.Printf("  -h, --help      Print this help screen.\n")
	fmt.Printf("\n")
	fmt.Printf("  -v, --version   Print the version of @G{voom} and exit.\n")
	fmt.Printf("\n")
	fmt.Printf("  -u, --url       @Y{(required)} The URL of the VMWare vCenter.\n")
	fmt.Printf("                  Can be set via the @W{$VOOM_URL} environment variable.\n")
	fmt.Printf("\n")
	fmt.Printf("  -U, --username  @Y{(required)} Your VMWare vCenter Username.\n")
	fmt.Printf("                  Can be set via the @W{$VOOM_USERNAME} environment variable.\n")
	fmt.Printf("\n")
	fmt.Printf("  -P, --password  @Y{(required)} Your VMWare vCenter Password.\n")
	fmt.Printf("                  Can be set via the @W{$VOOM_PASSWORD} environment variable.\n")
}

func printGlobalCommandsHelp() {
	fmt.Printf("COMMANDS\n")
	fmt.Printf("\n")
	fmt.Printf("  ls             List VMs.\n")
	fmt.Printf("  dump           Dump the list of VMs, in JSON.\n")
}
