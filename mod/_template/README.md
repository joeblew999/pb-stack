# _template

There are 3 Archetypes so far:

1. Remote-Source pulls remote source code from a repo, and packages it.

2. Binary pulls binaries from a repo and packages it.

3. Source has source code locally.



## Archetypes Traits 

Root needs to have:

**taskfile.yml** thats includes **name.taskfile.yaml**

Examples needs to have **local**, **none**, **remote** setups to simulate the Software Flows.

## Runtimes

local is when your using local task files.

none is when there is nothing in the foolder. 

remote is when your pulling the tasks remotely, as your running remotely from some other place.
