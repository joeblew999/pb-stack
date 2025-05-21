# livekit and redka

```
go install github.com/go-task/task/v3/cmd/task@latest
```

Run these 3:

```sh
task bin


task redka:server

# in new terminal
task server:dev

```



open http://127.0.0.1:7880

```sh

2025-05-20T14:50:37.322+1000	ERROR	livekit	routing/redisrouter.go:221	status update delayed, possible deadlock


```