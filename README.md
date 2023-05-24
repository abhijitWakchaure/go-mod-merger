# Go Mod Merger

This is a command line tool which can be used to create a common `go.mod` file along with `imports.go` for a large code base where multiple repositories are combined to have a common vendor directory.

The `imports.go` file will have all the direct imports under a single import statement.

## How to use

Just execute the binary passing path to multiple go.mod file(s)

```bash
$ ./go-mod-merger ~/go/src/github.com/project-flogo/core/go.mod ~/go/src/github.com/project-flogo/flow/go.mod ~/go/src/github.com/project-flogo/legacybridge/go.mod ./go.mod
Starting Go Mod Merger [v1.2.0]...

Using config file: ~/go/src/github.com/abhijitWakchaure/go-mod-merger/go-mod-merger.json
Using replace map:
        github.com/abhijitWakchaure/project-flogo-core
        github.com/abhijitWakchaure/project-flogo-flow
Parsing go.mod at: ~/go/src/github.com/project-flogo/core/go.mod
Parsing go.mod at: ~/go/src/github.com/project-flogo/flow/go.mod
Parsing go.mod at: ~/go/src/github.com/project-flogo/legacybridge/go.mod

Mismatched version for github.com/project-flogo/core
        want  : v1.6.5  mod file: ~/go/src/github.com/project-flogo/legacybridge/go.mod
        want  : v1.6.4  mod file: ~/go/src/github.com/project-flogo/flow/go.mod
        picked: v1.6.5 🔼

Mismatched version for github.com/stretchr/testify
        want  : v1.8.2  mod file: ~/go/src/github.com/project-flogo/legacybridge/go.mod
        want  : v1.4.0  mod file: ~/go/src/github.com/project-flogo/flow/go.mod
        picked: v1.8.2 🔼
Parsing go.mod at: ./go.mod

Mismatched version for golang.org/x/sys
        want  : v0.3.0  mod file: ./go.mod
        want  : v0.0.0-20220715151400-c0bba94af5f8      mod file: ~/go/src/github.com/project-flogo/legacybridge/go.mod
        picked: v0.3.0 🔼

Artifacts generated successfully for module name 'dummy'
```

You can also pass an optioanl package name so that generated artifacts belong to the package you provided.

```bash
$ ./go-mod-merger -p test ~/go/src/github.com/project-flogo/core/go.mod ~/go/src/github.com/project-flogo/flow/go.mod
...
...
...

Artifacts generated successfully for module name 'test'
```

You can also pass an optioanl output directory path if you wish to store generated artifacts somewhere else rather than on default `./output`.

```bash
./go-mod-merger -p test -o ../output ~/go/src/github.com/project-flogo/core/go.mod ~/go/src/github.com/project-flogo/flow/go.mod
```

## Config

Currently the tool supports passing a map of module names to replace via a config file named as `go-mod-merger.json`. This is how the sample config file should look like:

```json
{
  "replace": {
    "github.com/project-flogo/core": "github.com/abhijitWakchaure/project-flogo-core",
    "github.com/project-flogo/flow": "github.com/abhijitWakchaure/project-flogo-flow"
  },
  "ignoreMajorVersionMismatch": ["github.com/project-flogo/core", "github.com/project-flogo/flow"],
  "forceMaster": ["github.com/project-flogo/core", "github.com/project-flogo/flow"]
}
```

This config will ignore the major version mismatch for all the packages in the array `ignoreMajorVersionMismatch` and will just pick the latest major version. It will directly use version as a `master` for packages in the array `forceMaster`. Also it will add the `replace packageA version => packageB version` statements for provided packages in the `replace` map.

## Artifacts

The tool will create these artifacts: `go.mod`, `imports.go`, `depMismatch.json` and `allDeps.json`

- go.mod: This mod file will have the given module name (if not specified default `dummy` will be used), current go version, direct and indirect dependencies.

- imports.go: This will contain import statements for all the required dependencies using blank identifier.

- depMismatch.json: This will contain a tree of all the mismatched dependencies along with their versions and go.mod file path. e.g. The sample depMismatch.json will look like:

```json
{
  "github.com/project-flogo/core": {
    "v1.6.4": ["~/go/src/github.com/project-flogo/flow/go.mod"],
    "v1.6.5": ["~/go/src/github.com/project-flogo/legacybridge/go.mod"]
  },
  "github.com/stretchr/testify": {
    "v1.4.0": [
      "~/go/src/github.com/project-flogo/core/go.mod",
      "~/go/src/github.com/project-flogo/flow/go.mod"
    ],
    "v1.8.2": ["~/go/src/github.com/project-flogo/legacybridge/go.mod"]
  }
}
```

- allDeps.json: This will contain a tree of all the dependencies along with their versions and go.mod file path. e.g. The sample allDeps.json will look like:

```json
{
  "github.com/project-flogo/core": {
    "master": [
      "~/go/src/github.com/project-flogo/flow/go.mod",
      "~/go/src/github.com/project-flogo/legacybridge/go.mod"
    ],
    "v0.9.2": [
      "~/go/src/github.com/project-flogo/grpc/go.mod"
    ]
  }
}
```
