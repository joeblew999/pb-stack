# vugu


https://github.com/vugu/vugu/


https://github.com/vugu/html
https://github.com/vugu/vgrouter
https://github.com/vugu/vjson
https://github.com/vugu/xxhash

## WASM runtimes

vugu lacks the advanced features we need.

Browser 

https://github.com/nlepage/go-wasm-http-server
users: https://github.com/nlepage/go-wasm-http-server/network/dependents?package_id=UGFja2FnZS01NDU1NzYzMTg0

Cloudlfare

https://github.com/syumai/workers
users: https://github.com/syumai/workers/network/dependents

```sh
bun install wrangler
bunx wrangler init
bunx wrangler dev

```




## users

https://github.com/vugu/vugu/network/dependents


## examples

https://github.com/vugu-examples

https://github.com/vugu-examples/simple

https://github.com/codegod100/lab/tree/main/vugu

## tooling

vugu has gen, fmt, run. Its like golang tooling.

The tooling is highly abstracted with one tool calling another.


## vggen

```sh
task gen:help
task: [gen:help] vugugen -h
Usage of vugugen:
  -r	Run recursively on specified path and subdirectories.
  -s	Merge generated code for a package into a single file.
  -skip-go-mod
    	Do not try to create go.mod as needed
  -skip-main
    	Do not try to create main.go as needed
  -tinygo
    	Generate code intended for compilation under Tinygo


```

## vgfmt

```sh
task fmt:help
task: [fmt:help] vugufmt -h
usage: vugufmt [flags] [path ...]
  -d	display diffs instead of rewriting files
  -i	run goimports instead of gofmt
  -l	list files whose formatting differs from vugufmt's
  -s	simplify code
  -w	write result to (source) file instead of stdout
```

## vgrun

```sh
task run:help

task: [run:help] vgrun -h
Usage of vgrun:
  -1	Run only once and exit after
  -auto-reload-at string
    	Run auto-reload server using this listener.  An empty string will disable it. (default "localhost:8324")
  -bin-dir string
    	Directory of where to place built binary (default "bin")
  -install-tools go install
    	Installs common Vugu tools using go install
  -keep-git
    	With new-from-example causes the .git folder to not be removed after cloning
  -new-from-example string
    	Initialize a new project from example.  Will git clone from github.com/vugu-examples/[value] or if value contains a slash it will be treated as a full URL sent to git clone.  Must be followed by empty or non existent target directory.
  -no-generate go generate
    	Disable go generate
  -v	Verbose output
  -watch-dir string
    	Specifies which directory to watch from (default ".")
  -watch-pattern string
    	Sets the regexp pattern of files to watch (default "\\.vugu$")
```

## Issues

Cant run using go 1.24 due to looking for WASM in old go 1.23 location.

https://github.com/vugu/vugu/issues/362


Cant build with tinygo. Its tired to use docker, but its fails:

https://github.com/vugu/vugu/blob/master/gen/parser-go-pkg.go#L102