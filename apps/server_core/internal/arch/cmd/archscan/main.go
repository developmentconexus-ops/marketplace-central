// Command archscan runs the architecture detectors over a tree and prints every
// finding as file:line. Exit 0 means zero findings.
package main

import (
	"flag"
	"fmt"
	"os"

	"marketplace-central/apps/server_core/internal/arch"
)

func main() {
	root := flag.String("root", "internal", "directory to scan")
	flag.Parse()

	var findings []arch.Finding
	for _, scan := range []func(string) (arch.Findings, error){
		arch.ScanCrossContextInternal,
		arch.ScanFloatInContracts,
		func(root string) (arch.Findings, error) {
			return arch.ScanVendorTokens(root, arch.VendorTokens)
		},
	} {
		got, err := scan(*root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "archscan: %v\n", err)
			os.Exit(2)
		}
		findings = append(findings, got...)
	}

	for _, f := range findings {
		fmt.Printf("%s:%d: %s: %s\n", f.File, f.Line, f.Rule, f.Detail)
	}
	if len(findings) > 0 {
		fmt.Fprintf(os.Stderr, "archscan: %d finding(s)\n", len(findings))
		os.Exit(1)
	}
}
