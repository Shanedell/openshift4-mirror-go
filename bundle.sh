#!/usr/bin/env bash

go run main.go bundle \
    --openshift-version 4.13.0-rc.4 \
    --catalogs redhat-operators \
    --platform aws \
    --skip-existing \
    --skip-release \
    --pull-secret $REDHAT_PULL_SECRET
