package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"sort"
	"strings"
)

func usage(status int) {
	fmt.Printf(`Usage:
	%s [PKG]
where PKG is the name of a Go package (e.g., github.com/cespare/deplist). If no
package name is given, the current directory is used.
`, os.Args[0])
	os.Exit(status)
}

func three(name string) string {
	parts := strings.Split(name, "/")
	if len(parts) > 3 {
        return strings.Join(parts[:3], "/")
	} else {
        return name
    }
}

func findDeps(soFar map[string]bool, name string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

    pkg, err := build.Import(name, cwd, 0)
	if err != nil {
        if build.IsLocalImport(name) {
            return err
        } else {
            soFar[three(name)] = true
            return nil
        }
	}

	if pkg.Goroot {
		return nil
	}

	soFar[three(name)] = true //pkg.ImportPath)] = true
	for _, imp := range pkg.Imports {
		if !soFar[imp] {
			if err := findDeps(soFar, imp); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	pkg := ""
	switch len(os.Args) {
	case 1:
		pkg = "."
	case 2:
		for _, s := range []string{"-h", "help", "-help", "--help"} {
			if os.Args[1] == s {
				usage(0)
			}
		}
		pkg = os.Args[1]
	default:
		usage(1)
	}

	deps := make(map[string]bool)
	err := findDeps(deps, pkg)
	if err != nil {
		log.Fatalln(err)
	}
	delete(deps, pkg)
	keys := make([]string, 0, len(deps))
	for key, _ := range deps {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, dep := range keys {
        if !build.IsLocalImport(dep) {
			fmt.Println(dep)
        }
	}
}
