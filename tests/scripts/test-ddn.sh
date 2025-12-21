#!/bin/bash

# NOTE: run this script at the root project folder or with make:
# make test-ddn
trap "cd tests/ddn && docker compose down --remove-orphans -v" EXIT


pushd tests/ddn \
    && docker compose up -d --build plugin \
    && popd

export DDN_ENGINE_HOST=http://localhost:3280
export GRAPHQL_SERVER_URL=http://localhost:3280/graphql

make test