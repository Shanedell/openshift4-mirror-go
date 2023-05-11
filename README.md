# openshift4-mirror-go

This is a Golang implementation of [openshift4-mirror](https://repo1.dso.mil/platform-one/distros/red-hat/ocp4/openshift4-mirror). The purpose behind this was to make distribution easier.

Executable files can be downloaded from the [releases page](https://github.com/Shanedell/openshift4-mirror-go/releases). Also a docker image ghcr.io/shanedell/openshift4-mirror-go.

## Options

- `build` - creates local docker image `localhost/openshift4-mirror:latest`. This image included the needed packages installed, code and prebuilt binary `openshift4_mirror`.
- `shell` - opens an interactive shell to local docker image. If the local docker image does not exist it will be created.
- `bundle` - downloads/bundles different parts of OpenShift content, based on args given.
- `prune` - prunes index images, currently only support redhat-operators, so you don't have to use/download all operators.

**NOTES**:

- Currently `prune` only supports pruning the redhat-operators index image.
- Currently `prune` only seems to work on Linux with Podman (not tested with Linux and docker).

## Usage

**Root Usage:**

```bash
openshift4_mirror - CLI for mirroring OpenShift 4 content.

Usage:
  openshift4_mirror [flags]
  openshift4_mirror [command]

Available Commands:
  build       build the container image
  bundle      bundle the OpenShift content
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  prune       prune the Red Hat Operator index image
  shell       open a shell in the container environment

Flags:
      --bundle-dir string             directory to save downloaded content
  -c, --containerRuntime string       container runtime. supported options [docker, podman]. if not specified, code looks for both and uses whichever is found first.
      --containerRuntimePath string   full to container runtime. needed if executable not in /usr/bin or /usr/local/bin
  -h, --help                          help for openshift4_mirror
  -v, --openshift-version string      the OpenShift version (e.g. 4.5.11)
      --pre-release                   pre-release version of OpenShift
      --pull-secret string            the content of your pull secret (can be found at https://cloud.redhat.com/openshift/install/pull-secret)
  -r, --target-registry string        target registry to tag the image with (default "example.registry.com")

Use "openshift4_mirror [command] --help" for more information about a command.
```

**Build Usage:**

```bash

Usage:
  openshift4_mirror build [flags]

Flags:
  -h, --help   help for build

Global Flags:
      --bundle-dir string             directory to save downloaded content
  -c, --containerRuntime string       container runtime. supported options [docker, podman]. if not specified, code looks for both and uses whichever is found first.
      --containerRuntimePath string   full to container runtime. needed if executable not in /usr/bin or /usr/local/bin
  -v, --openshift-version string      the OpenShift version (e.g. 4.5.11)
      --pre-release                   pre-release version of OpenShift
      --pull-secret string            the content of your pull secret (can be found at https://cloud.redhat.com/openshift/install/pull-secret)
  -r, --target-registry string        target registry to tag the image with (default "example.registry.com")
```

**Shell Usage:**

```bash
open a shell in the container environment

Usage:
  openshift4_mirror shell [flags]

Flags:
  -h, --help   help for shell

Global Flags:
      --bundle-dir string             directory to save downloaded content
  -c, --containerRuntime string       container runtime. supported options [docker, podman]. if not specified, code looks for both and uses whichever is found first.
      --containerRuntimePath string   full to container runtime. needed if executable not in /usr/bin or /usr/local/bin
  -v, --openshift-version string      the OpenShift version (e.g. 4.5.11)
      --pre-release                   pre-release version of OpenShift
      --pull-secret string            the content of your pull secret (can be found at https://cloud.redhat.com/openshift/install/pull-secret)
  -r, --target-registry string        target registry to tag the image with (default "example.registry.com")
```

**Bundle Usage:**

```bash
bundle the OpenShift content

Usage:
  openshift4_mirror bundle [flags]

Flags:
      --catalog-version string   version of images to use for catalogs
      --catalogs strings         the catalog(s) content to download. catalogs: [redhat-operators, certified-operators, redhat-marketplace, community-operators]. defaults to all
  -h, --help                     help for bundle
      --platform string          target platform for install. platforms: [aws, azure, gcp, metal, openstack, vmware]
      --skip-catalogs            skip downloading of catalog content
      --skip-existing            skip downloading content that already exists on disk (default true)
      --skip-release             skip downloading of release content
      --skip-rhcos               skip downloading of RHCOS image

Global Flags:
      --bundle-dir string             directory to save downloaded content
  -c, --containerRuntime string       container runtime. supported options [docker, podman]. if not specified, code looks for both and uses whichever is found first.
      --containerRuntimePath string   full to container runtime. needed if executable not in /usr/bin or /usr/local/bin
  -v, --openshift-version string      the OpenShift version (e.g. 4.5.11)
      --pre-release                   pre-release version of OpenShift
      --pull-secret string            the content of your pull secret (can be found at https://cloud.redhat.com/openshift/install/pull-secret)
  -r, --target-registry string        target registry to tag the image with (default "example.registry.com")
```

**Prune Usage:**

```bash
prune the Red Hat Operator index image

Usage:
  openshift4_mirror prune [flags]

Flags:
  -h, --help                help for prune
  -n, --namespace string    full path of tagged image (default "redhat")
  -o, --operators strings   the operator(s) desired. Rest are pruned out
  -t, --tag string          version tag for image (default "latest")

Global Flags:
      --bundle-dir string             directory to save downloaded content
  -c, --containerRuntime string       container runtime. supported options [docker, podman]. if not specified, code looks for both and uses whichever is found first.
      --containerRuntimePath string   full to container runtime. needed if executable not in /usr/bin or /usr/local/bin
  -v, --openshift-version string      the OpenShift version (e.g. 4.5.11)
      --pre-release                   pre-release version of OpenShift
      --pull-secret string            the content of your pull secret (can be found at https://cloud.redhat.com/openshift/install/pull-secret)
  -r, --target-registry string        target registry to tag the image with (default "example.registry.com")
```
