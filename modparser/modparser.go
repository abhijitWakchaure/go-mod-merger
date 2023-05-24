package modparser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/abhijitWakchaure/go-mod-merger/config"
	"github.com/abhijitWakchaure/go-mod-merger/gogenerator"
	"github.com/abhijitWakchaure/go-mod-merger/semvar"
	"golang.org/x/mod/modfile"
)

// depMeta ...
type depMeta struct {
	source, path, version string
	indirect              bool
}

var depMismatchTree, allDepsTree map[string]interface{}

var modReplace = map[string]string{}

// Parse ...
func Parse(moduleName, outputDir string, files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("no go.mod file(s) provided")
	}
	c := config.Read()
	if c != nil && len(c.Replace) > 0 {
		modReplace = c.Replace
		fmt.Printf("Using replace map:\n")
		for _, v := range modReplace {
			fmt.Printf("\t%s\n", v)
		}
	}
	deps := make(map[string]*depMeta, 0)
	// Create a new modfile.File object
	mod := new(modfile.File)
	if err := mod.AddModuleStmt(moduleName); err != nil {
		return err
	}
	if err := mod.AddGoStmt(goVersion()); err != nil {
		return err
	}
	var versionMiss bool
	depMismatchTree = make(map[string]interface{})
	allDepsTree = make(map[string]interface{})
	for _, v := range files {
		if filepath.Base(v) != "go.mod" {
			return fmt.Errorf("invalid go.mod file path: %s", v)
		}
		fmt.Printf("\nParsing go.mod at: %s", v)
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
			dep := &depMeta{
				source:   v,
				path:     req.Mod.Path,
				version:  req.Mod.Version,
				indirect: req.Indirect,
			}
			// check if master is forced for the package
			if isMasterForced(dep.path) {
				fmt.Printf("\n\t🚩 Overriding version with 'master' for module [%s]", dep.path)
				dep.version = "master"
			}
			if d, ok := deps[req.Mod.Path]; ok && d.version != dep.version {
				fmt.Printf("\n\tMismatched version for [%s]\n", req.Mod.Path)
				fmt.Printf("\t\twant  : %s \tmod file: %s\n", dep.version, dep.source)
				fmt.Printf("\t\twant  : %s \tmod file: %s\n", d.version, d.source)
				latest, err := semvar.Compare(req.Mod.Path, d.version, dep.version)
				if err != nil {
					fmt.Printf("❌ Error! %s\n", err.Error())
					versionMiss = true
				} else {
					fmt.Printf("\t\tpicked: %s 🔼\n", latest)
					dep.version = latest
					deps[req.Mod.Path].version = latest
				}
			}
			deps[req.Mod.Path] = dep
			addDepTree(req.Mod.Path, dep)
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
	_, err = os.Stat(outputDir)
	if err != nil {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return err
		}
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
	filterDepTree()
	writeJSON(filepath.Join(outputDir, "depMismatch.json"), depMismatchTree)
	writeJSON(filepath.Join(outputDir, "allDeps.json"), allDepsTree)
	if versionMiss {
		fmt.Printf("\nArtifacts generated with error(s) for module name '%s'\n", moduleName)
	} else {
		fmt.Printf("\nArtifacts generated successfully for module name '%s'\n", moduleName)
	}
	return nil
}

func addDepTree(modPath string, dep *depMeta) {
	v, ok := depMismatchTree[modPath]
	if !ok {
		depMismatchTree[modPath] = map[string][]string{
			dep.version: {dep.source},
		}
		return
	}
	versionList := v.(map[string][]string)
	sources, ok := versionList[dep.version]
	if ok {
		sources = append(sources, dep.source)
		versionList[dep.version] = sources
	} else {
		versionList[dep.version] = []string{dep.source}
	}
}

func filterDepTree() {
	for k, v := range depMismatchTree {
		allDepsTree[k] = v
		versionList := v.(map[string][]string)
		if len(versionList) == 1 {
			delete(depMismatchTree, k)
		}
	}
}

func goVersion() string {
	v := runtime.Version()
	v = strings.TrimLeft(v, "go")
	vArr := strings.Split(v, ".")
	return fmt.Sprintf("%s.%s", vArr[0], vArr[1])
}

func writeJSON(filePath string, data any) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func isMasterForced(packageName string) bool {
	for _, v := range config.Read().ForceMaster {
		if packageName == v {
			return true
		}
	}
	return false
}
