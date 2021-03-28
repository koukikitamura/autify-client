#!/usr/bin/env bash

DIR=$(cd $(dirname $BASH_SOURCE); pwd)
ROOT_DIR="${DIR}/.."
cd ${ROOT_DIR}

go build -o ./tmp/atf

./tmp/atf scenario --project-id=$AUTIFY_TEST_PROJECT_ID --scenario-id=$AUTIFY_TEST_SCENRIO_ID

rm ./tmp/atf
