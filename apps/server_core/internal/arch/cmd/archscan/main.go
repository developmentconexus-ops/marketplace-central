// Command archscan runs the architecture detectors over one or more trees and
// prints every finding as file:line. Exit 0 means zero findings over a non-empty
// file set; exit 1 is findings; exit 2 is the scanner failing to run, which
// includes having nothing to scan. That last case is deliberate: a scanner that
// parsed zero files and a clean tree used to print the same nothing, and the
// default root quietly never covered cmd/ or tests/.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"marketplace-central/apps/server_core/internal/arch"
)

func main() {
	roots := flag.String("root", "cmd,internal,tests", "comma-separated directories to scan")
	flag.Parse()

	scanned := 0
	var findings []arch.Finding
	for _, root := range strings.Split(*roots, ",") {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}
		n, err := arch.CountFiles(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "archscan: count %s: %v\n", root, err)
			os.Exit(2)
		}
		scanned += n
		for _, scan := range []func(string) (arch.Findings, error){
			arch.ScanCrossContextInternal,
			arch.ScanFloatInContracts,
			func(root string) (arch.Findings, error) {
				return arch.ScanVendorTokens(root, arch.VendorTokens)
			},
			arch.ScanFactValueDiscard,
		} {
			got, err := scan(root)
			if err != nil {
				fmt.Fprintf(os.Stderr, "archscan: %v\n", err)
				os.Exit(2)
			}
			findings = append(findings, got...)
		}
	}

	for _, f := range findings {
		fmt.Printf("%s:%d: %s: %s\n", f.File, f.Line, f.Rule, f.Detail)
	}
	fmt.Printf("archscan: scanned=%d findings=%d\n", scanned, len(findings))
	if scanned == 0 {
		fmt.Fprintln(os.Stderr, "archscan: no Go file scanned. A verdict over the empty set is not a verdict.")
		os.Exit(2)
	}
	if len(findings) > 0 {
		os.Exit(1)
	}
}
