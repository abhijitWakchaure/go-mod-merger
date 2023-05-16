package modparser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/abhijitWakchaure/go-mod-merger/gogenerator"
	"golang.org/x/mod/modfile"
)

// depMeta ...
type depMeta struct {
	source, path, version string
	indirect              bool
}

var modReplace = map[string]string{
	"github.com/project-flogo/core": "github.com/abhijitWakchaure/project-flogo-core",
	"github.com/project-flogo/flow": "github.com/abhijitWakchaure/project-flogo-flow",
}

// Parse ...
func Parse(moduleName string, files []string) error {
	deps := make(map[string]depMeta, 0)
	// Create a new modfile.File object
	mod := new(modfile.File)
	if err := mod.AddModuleStmt(moduleName); err != nil {
		return err
	}
	if err := mod.AddGoStmt("1.20"); err != nil {
		return err
	}
	for _, v := range files {
		fmt.Println("Parsing go.mod file from path:", v)
		if _, err := os.Stat(v); err != nil {
			return err
		}
		data, err := ioutil.ReadFile(v)
		if err != nil {
			return err
		}

		// Parse the go.mod file
		mod, err := modfile.Parse(v, data, nil)
		if err != nil {
			return err
		}

		// Print the module's dependencies
		for _, req := range mod.Require {
			dep := depMeta{
				source:   v,
				path:     req.Mod.Path,
				version:  req.Mod.Version,
				indirect: req.Indirect,
			}
			if d, ok := deps[req.Mod.Path]; ok && d.version != dep.version {
				fmt.Printf("Error! Mismatched version for %s\n", req.Mod.Path)
				fmt.Printf("\twant: %s \tmod file: %s\n", dep.version, dep.source)
				fmt.Printf("\twant: %s \tmod file: %s\n", d.version, d.source)
			}
			deps[req.Mod.Path] = dep
		}
	}
	for k, v := range deps {
		mod.AddNewRequire(k, v.version, v.indirect)
		if re, ok := modReplace[k]; ok {
			// Replace a module with a new version
			err := mod.AddReplace(k, v.version, re, "master")
			if err != nil {
				return err
			}
		}
	}

	mod.SetRequireSeparateIndirect(mod.Require)
	mod.SortBlocks()
	mod.Cleanup()
	b, err := mod.Format()
	if err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	outputDir := filepath.Join(pwd, "output")
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}
	outputModFile := filepath.Join(outputDir, "go.mod")
	err = ioutil.WriteFile(outputModFile, b, 0644)
	if err != nil {
		return err
	}
	outputImportsFile := filepath.Join(outputDir, "imports.go")
	imports := gogenerator.Imports{
		PackageName: moduleName,
		ImportsArr:  mod.Require,
	}
	f, err := os.Create(outputImportsFile)
	if err != nil {
		return err
	}
	defer f.Close()
	err = gogenerator.GenerateImportsFile(imports, f)
	if err != nil {
		return err
	}
	fmt.Printf("\nArtifacts generated successfully with module name '%s'\n", moduleName)
	return nil
}
