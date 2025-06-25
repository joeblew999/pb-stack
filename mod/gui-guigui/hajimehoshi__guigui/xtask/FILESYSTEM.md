# File System

the File system provided by https://github.com/opencloud-eu/opencloud, which also uses NATS.

https://github.com/opencloud-eu/reva is also used. Not sure why, but their code is way ahead of reva https://github.com/cs3org/reva

I have not decided know if we should add opencloud as a compiled in Module or Run it using process compose.

How to run it without Docker etc ?

https://docs.opencloud.eu/docs/admin/getting-started/other/bare-metal

It does not have any database but seems to use a Posix FS instead, but i do not know if that works on Windows, Darwin or only Linux. Last thing i need is a Fuse dependency.  


It seems to use NATS KV store to help it works. 

We will not use its GUI, but instead use Datastar of course. SO hopefully it has some sort of API that exposes it over HTTP and SEE  ?

We probably need a Taskfile for OpenCloud that we can then include, so that the main Taskfile does not get bloated.












