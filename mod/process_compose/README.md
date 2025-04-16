# process-compose

https://github.com/F1bonacc1/process-compose

## Dev-time

copy their go.mod in and tidy ...

```sh task ``` , to build.

```sh task go:run ``` , to run.


```sh task base-bin-pack ``` , to package for distribution.

# to distribute
```sh task base-bin-push ```, to push the binary for usage by others.


## Run-time


cd example/local && task 

``` task base-bin-pull  ``` , to update your local binary.

``` task base-bin-run ``` , to run your local binary.


## Usage

For running sets of binaries in a formation.

Replaces K8, Docker and Docker compose.

TODO:

- Same as agent we need for process compose.

## beszel

https://github.com/henrygd/beszel

TODO:

- Add to pocketbase module. We can work out how best to integrate later. 

- The SQLite can be on server or on nats server where the web gui runs.

- Also how best to have a nats agent . Same as agent we need for process compose.

- And then redo the web gui with Datastar . 






