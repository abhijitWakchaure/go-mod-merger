# Go Mod Merger

This is a command line tool which can be used to create a common `go.mod` file along with `imports.go` for a large code base where multiple repositories are combined to have a common vendor directory.

The `imports.go` file will have all the direct imports under a single import statement.

## How to use

Just execute the binary passing path to multiple go.mod file(s)

```bash
$ ./go-mod-merger ~/go/src/github.com/project-flogo/core/go.mod ~/go/src/github.com/project-flogo/flow/go.mod
Parsing go.mod file from path: /home/abhijit/dev/godev/src/github.com/project-flogo/core/go.mod
Parsing go.mod file from path: /home/abhijit/dev/godev/src/github.com/project-flogo/flow/go.mod
Error! Mismatched version for github.com/project-flogo/core
        want: v1.6.5    mod file: /home/abhijit/dev/godev/src/github.com/project-flogo/legacybridge/go.mod
        want: v1.6.4    mod file: /home/abhijit/dev/godev/src/github.com/project-flogo/flow/go.mod
Error! Mismatched version for github.com/stretchr/testify
        want: v1.8.2    mod file: /home/abhijit/dev/godev/src/github.com/project-flogo/legacybridge/go.mod
        want: v1.4.0    mod file: /home/abhijit/dev/godev/src/github.com/project-flogo/flow/go.mod

Artifacts generated successfully with module name 'dummy'
```

You can also pass a optioanl package name so that generated artifacts belong to the package you provided.

```bash
$ ./go-mod-merger -p test ~/go/src/github.com/project-flogo/core/go.mod ~/go/src/github.com/project-flogo/flow/go.mod
...
...
...

Artifacts generated successfully with module name 'test'
```
