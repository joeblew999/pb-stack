# Task

Each Task file must both install the binaries and run them and copy itself ( the task file )

The binaries and task file itself are then placed into the .dep folder.

The .dep folder is path, this automatically exposing the binaries, without polluting any part of the OS.

When we pack to the .pack dot folder, the binaries are deployed to the NATS Object Store, so that Task can pull them from there over at the Deployment side of things.  


TODO

- TASK files for each.

## task

https://github.com/go-task/task

## task TUI 

https://github.com/aleksandersh/task-tui

A TUI for task

## task Web

https://github.com/titpetric/task-ui

A Web GUI for task when its running inside docker.

## examples and references

https://github.com/saydulaev/taskfile 

## livekit task integration

Cli uses task...

https://github.com/livekit/livekit-cli is hardcoded to read projects templates from https://github.com/livekit-examples/index which has a taskfile and template file.

---

Server

https://github.com/livekit/livekit

https://github.com/nalgeon/redka to replace redis server.

opts.DontListen=true














