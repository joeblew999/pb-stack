# variables

https://github.com/go-task/task/issues/2253#issuecomment-2868075833

best working system i can get that works on darwin and mac. not tested on linux.

.env overrides work

task env can be overridden, but not task vars

it deterministic - yes.

override from .sh: 
```sh
./run.sh
```

override from inside task: 
```sh
task
task override-01
```

override from outside task. 
HTML_DEEP_PATH fails because is VAR, not ENV.

```sh
HTML_DEEP_PATH=$PWD/over-from-shell HTML_DEEP_NAME=over-from-shell task
```


