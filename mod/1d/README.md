# 1d

Text and trans

https://github.com/EliCDavis/notes looks like a nice generator to work off.

This creates the data as JSON and gens to markdown.

Lets see how far a Taskfile overlay gets us.

```sh

task tools

mkdir -p $PWD/test
task project:new -- --path $PWD/test "project01"

```
