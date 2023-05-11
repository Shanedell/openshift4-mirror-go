#!/usr/bin/env bash

# duplicate of https://repo1.dso.mil/platform-one/distros/red-hat/ocp4/openshift4-mirror/-/raw/main/entrypoint.sh

set -e

if [[ $# == 0 ]]; then
  if [[ -t 0 ]]; then
    echo
    echo "Starting shell..."
    echo

    exec "bash"
  else
    echo "An interactive shell was not detected."
    echo
    echo "By default, this container starts a bash shell, be sure you are passing '-it' to your run command."

    exit 1
  fi
else
  exec "$@"
fi
