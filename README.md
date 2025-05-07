# pb stack

This is a system to make it easy to use Taskfile everywhere and to use as a basis for larger things where a system can easily be composed at both compile time as well as runtime. The task system is used to help both and allow anyone to level this via the Task Remote files feature sets.

I work on a lot of systems that are "self sovereign". These are systems where you, the User, control where the data and compute runs, so that you have complete privacy and control. This is mandated on many government systems but it also a right under GDPR laws in the EU and other countries. Its a good thing for everyone, not just governments.

This is a system though where the Users can revoke all their data at anytime. A User can choose to be their own Operator and federate with other Operators, and so automatically retain their own data also. Pick your poison :)

Auth will always support Key based , so that Operators and Users retain fully control. However Authelia Auth can be leveraged if you want a Email, OIDC or Passkey based auth. The Key based system will evolve so that a Browser can be user to authenticate with their own keys on their devices. A Sync system will evolve, so that all your own devices synchronise the keys, bypassing the OS sync system that Google, Apple and Microsoft control and use to sync all your auth keys and passwords. This will be useful for Blockchains also.


## Status 

Its a WIP, but the task files are useful if you need a self running system. 

## docs

This Project uses the top level docs folder.  See [Doc](./doc/README.md) folder for Project Info.

## todo:

OH Shit. Change to .src/reponame. So i can easily search and destroy.

https://taskfile.dev/reference/package#githubcomgo-tasktaskv3taskfile looks like we can call task files from anywhere.
