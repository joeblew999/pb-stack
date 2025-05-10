# variables

https://github.com/go-task/task/issues/2253#issuecomment-2868075833

env can be overridden, but not var.

it deterministic - yes it seems so.

override from .sh: 
```sh
./run.sh
```

override from inside task: 
```sh
task override-01
```


