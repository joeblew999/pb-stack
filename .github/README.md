# Github

! STATUS !

Read the Docs for now to get a feel for the intent. Ask questions if something seems off.

The scaffolds for the [TASK](../mod/task/README.md) files and [TOFU](../mod/tofu/README.md) are going in currently. We need this to rapidly develop. 

Next will be the [GitHub Actions for CI and CD](../.github/workflows/README.md) round-tripping to your own Desktop.

Then simple examples / playgrounds, so we can work up the code generator against real projects. Its the only way for others to understand how to develop this.

## Documentation

See [Doc](../doc/README.md) folder for Project Info.



## Make

This is just a helper, until everything is converted to Task files.

We need task and golang installed and this does that check for now.

## Task

https://taskfile.dev/reference/

TASK files are used:

1. Locally for dev.

2. In Github Actions for CI and CD, along with TOFU files.

3. In Production for Upgrades, along with Tofu files.


```sh

task --experiments
* GENTLE_FORCE:     on (1)
* REMOTE_TASKFILES: on (1)
* MAP_VARIABLES:    on (2)
* ENV_PRECEDENCE:   on (1)

task --list-all --verbose
task: [https://raw.githubusercontent.com/saydulaev/taskfile/v1.4.3/Taskfile.yml] Fetched remote copy
task: [https://raw.githubusercontent.com/saydulaev/taskfile/v1.4.3/docker-compose/Taskfile.yml] Fetched remote copy
task: [https://raw.githubusercontent.com/saydulaev/taskfile/v1.4.3/terraform/Taskfile.yml] Fetched remote copy
task: [https://raw.githubusercontent.com/saydulaev/taskfile/v1.4.3/security/Taskfile.yml] Fetched remote copy
task: [https://raw.githubusercontent.com/saydulaev/taskfile/v1.4.3/security/sast/Taskfile.yml] Fetched remote copy
task: [https://raw.githubusercontent.com/saydulaev/taskfile/v1.4.3/security/sast/Checkov.yml] Fetched remote copy
task: [https://raw.githubusercontent.com/saydulaev/taskfile/v1.4.3/security/sast/Grype.yml] Fetched remote copy
task: [https://raw.githubusercontent.com/saydulaev/taskfile/v1.4.3/security/sast/Trivy.yml] Fetched remote copy
task: Available tasks for this project:
* base:dep:                                    base dep, installs shell level components.
* base:print:                                  base print
* caddy:dep:                                   caddy install
* caddy:print:                                 caddy print
* cloudflare:print:                            cloudflare print
* compose:print:                               
* datastar:build:                              
* datastar:css:                                
* datastar:deploy:                             
* datastar:idiomorph:                          
* datastar:kill:                               
* datastar:libpub:                             
* datastar:library:                            
* datastar:print:                              
* datastar:qtc:                                
* datastar:sdktspub:                           
* datastar:site:                               build and run site
* datastar:support:                            
* datastar:templ:                              
* datastar:test:                               
* datastar:test-all:                           
* datastar:tools:                              
* datastar:version:                            
* git:print:                                   git print
* go:print:                                    go print
* remote:docker-compose:down:                  Stop and remove containers, networks
* remote:docker-compose:up:                    Create and start containers
* remote:security:sast:checkov:scanner:        Infrastructure as code static analysis
* remote:security:sast:grype:scanner:          A vulnerability scanner for container images, filesystems, and SBOMs
* remote:security:sast:trivy:aws:              [EXPERIMENTAL] Scan AWS account
* remote:security:sast:trivy:config:           Scan config files for misconfigurations
* remote:security:sast:trivy:filesystem:       Scan local filesystem
* remote:security:sast:trivy:image:            Scan a container image
* remote:security:sast:trivy:kubernetes:       [EXPERIMENTAL] Scan kubernetes cluster
* remote:security:sast:trivy:repository:       Scan a repository
* remote:security:sast:trivy:rootfs:           Scan rootfs
* remote:security:sast:trivy:sbom:             Scan SBOM for vulnerabilities and licenses
* remote:security:sast:trivy:vm:               [EXPERIMENTAL] Scan a virtual machine image
* remote:terraform:apply:                      terraform apply -auto-approve
* remote:terraform:destroy:                    terraform destroy
* remote:terraform:doc:                        terraform-docs markdown table
* remote:terraform:fmt:                        terraform fmt
* remote:terraform:init:                       terraform init
* remote:terraform:plan:                       terraform plan
* remote:terraform:terrascan:                  Terrascan static code analyzer
* remote:terraform:tflint:                     tflint
* remote:terraform:validate:                   terraform validate
* todo:print:                                  todo print


task base:print --silent
BASE_SHELL_OS_NAME: darwin
BASE_SHELL_OS_ARCH: arm64
BASE_OS_NAME: darwin
BASE_OS_ARCH: arm64



task git:print
task: [git:print] echo "GIT_NAME:"       /opt/homebrew/opt/git/libexec/git-core/git
GIT_NAME: /opt/homebrew/opt/git/libexec/git-core/git
task: [git:print] echo "GIT_VERSION:"    git version 2.49.0




```
