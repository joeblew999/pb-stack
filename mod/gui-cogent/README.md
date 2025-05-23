# cogent

HIghly geared to Science, which is what EU and US need for collab science.

Does 2d and 3d on web, mobile and desktop.

https://github.com/cogentcore/core/

https://github.com/cogentcore/core/releases/tag/v0.3.11

https://github.com/cogentcore/core/network/dependents

---

https://github.com/cogentcore/webgpu

---

https://github.com/cogentcore/lab
- has cool GPU Render farming in example called baremetal and simmer, which looks really useful for Science teams.

# taskfile

Best to let it find the source ? COGENT_SRC_PATH is used for now, but we can use USER_WORKING_DIR ?

dev.env in any empty repo works. Is decent solution for now.

Or pass args:

```sh
# get status of git already cloned repo ( assuming we are in the repo.)

task ENV_GIT_REPO=https://github.com/devilcove/bbolteditor ENV_GIT_REPO_NAME=bbolteditor ENV_GIT_REPO_VERSION=master COGENT_ENV_SRC_PREFIX=. src:status


# creates a new git repo
rm -rf joeblew999__test01

rm -rf joeblew999__test01 && mkdir -p joeblew999__test01 && cd joeblew999__test01 && task ENV_GIT_REPO=git@github.com-joeblew999:joeblew999/test01 ENV_GIT_REPO_NAME=test01 ENV_GIT_REPO_VERSION=main git:create

# create a new ssh 
rm -rf joeblew999__test02

rm -rf joeblew999__test02 && mkdir -p joeblew999__test02 && cd joeblew999__test02 && task ENV_GIT_REPO=git@github.com-joeblew999:joeblew999/test01 ENV_GIT_REPO_NAME=test01 ENV_GIT_REPO_VERSION=main ssh:create



```



## archi

Content is reused. for creating content-focused sites consisting of Markdown, HTML, and Cogent Core

So can we reuse the markdown in the Cogent GUI too ?  It can render SVG, so Deck SVG should in theory work.
Editing of the SVG, needs to know the DeckXML line that matches the SVG, so that we can work out the changes to the DeckXML.
The Editing ability of Cogent should come in handy to do this.

## cool examples

https://github.com/devilcove/bbolteditor



