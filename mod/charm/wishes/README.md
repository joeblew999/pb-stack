# wishes

wisshes Is SSH + Extra Steps

Pure GO answer to Infrastructure as Code tools like Ansible

https://github.com/charmbracelet/wish
- see examples to get what we need.

https://github.com/delaneyj/toolbelt/tree/main/wisshes


## Dev

task 
task remote:go:run

.wisshes folder is the config.

## Usage

This is used to SSH to a server, to get the initial binaries and data onto a server.

After that, your agents on the server should take over, so things are automatically updated.

For example with NATS, to get the latest binaries or data in real time.
Or using Cloud FLare R2 is another good way to pull binaries and data, but its not real time.

Process compose works really well with this, because we can restart the processes when the binaries change.
Just do not restart nats though, because NATS is ensuring buffering.

## Clouds

We need to run on Clouds to get going so here are our main ones

Hetzner

https://docs.hetzner.com/robot/dedicated-server/security/ssh/

Fly

https://fly.io/docs/flyctl/ssh/

https://fly.io/docs/flyctl/ssh-console/





