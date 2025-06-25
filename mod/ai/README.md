# AI

Golang based so it can be used at Compile time and runtime.

At Runtime, its essentially using the data to generate new data. 

- the AST files that describe the system.
- the benthos config to map data.
- Perhaps even golang that is compiled ( using compiler as a service ) into WASM, and then loading that WASM into the running system.

CLI and GUI 

## gui

https://github.com/CoreyCole/datastarui

https://github.com/henrygd/beszel 


## Engines

For building own AI systems, instead of running theirs.

https://github.com/gomlx/gomlx 

- by https://github.com/janpfeifer who is a google researcher.

- can run anywhere. In Browser as WASM. On Bare metal with no OS.

- example as wasm in a Browser: https://github.com/janpfeifer/hiveGo

- uses https://github.com/janpfeifer/gonb for Notebook.

- uses github.com/gowebapi/webapi for driving any Web GUI from Golang.

## Runtimes

https://github.com/ergo-services is erlang for golang. 

docs: https://docs.ergo.services

Its good for AI because Actors and workflows are able to be spawns, restarted using the old Supervisions trees approach.

https://github.com/halturin is the main person.

https://github.com/ergo-services/examples 



## Frameworks

## https://github.com/sst/opencode

- uses github.com/mark3labs/mcp-go


## https://github.com/charmbracelet/mods

- uses github.com/mark3labs/mcp-go


## https://github.com/danielmiessler/fabric

- custom mcp 
- web gui using svelte


## https://github.com/ollama/ollama

- so easy and used by everyone to riun local ai

## https://github.com/XiaoConstantine/maestro

Maestro evaluates code across four key dimensions:

Code Defects: Error handling, logic flaws, resource management issues.
Security Vulnerabilities: SQL injection, cross-site scripting, insecure data handling.
Maintainability and Readability: Code organization, documentation, naming conventions.
Performance Issues: Inefficient algorithms, suboptimal data structures, excessive operations.

- uses https://github.com/XiaoConstantine/dspy-go

---

https://github.com/localrivet

## https://github.com/localrivet

reddit: https://www.reddit.com/r/golang/comments/1kx4gew/built_a_go_mcp_server_that_let_claude_generate_a/

medium: https://medium.com/@alma.tuck


he is using task-master to help him make the golang.


https://github.com/localrivet/gomcp

https://github.com/localrivet/projectmemory
- uses gomcp
- sqlite for embedding

https://github.com/localrivet/gocreate
- uses gomcp
- has tools ( https://github.com/localrivet/gocreate/tree/main/tools ) that do the classic things you need to do as a dev.


