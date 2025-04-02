# Github

! STATUS !

Read the Docs for now to get a feel for the intent. Ask questions if something seems off.

The scaffolds for the task files and tofu are going in.
Then simple examples / playgrounds, so we can work up the code generator.

## Doc

See [Doc](../doc/README.md) folder for Project Info.



## Task

TASK files are used:

1. Locally for dev.

2. In Github Actions for CI and CD, along with TOFO files.

3. In Production for Upgrades, along with Tofu files.


```sh

task --experiments
* GENTLE_FORCE:     on (1)
* REMOTE_TASKFILES: on (1)
* MAP_VARIABLES:    on (2)
* ENV_PRECEDENCE:   on (1)


task --list-all

task: Available tasks for this project:
* build:                                       
* css:                                         
* default:                                     
* deploy:                                      
* idiomorph:                                   
* kill:                                        
* libpub:                                      
* library:                                     
* qtc:                                         
* sdktspub:                                    
* site:                                        build and run site
* support:                                     
* templ:                                       
* test:                                        
* test-all:                                    
* tools:                                       
* version:                                     
* base:default:                                      (aliases: base)
* git:default:                                       (aliases: git)
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


pb-stack % task git
Hello, from base!

pb-stack % task git
Hello, from git!


```
