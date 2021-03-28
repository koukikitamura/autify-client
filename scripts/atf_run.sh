#!/usr/bin/env bash

DIR=$(cd $(dirname $BASH_SOURCE); pwd)
ROOT_DIR="${DIR}/.."
cd ${ROOT_DIR}

go build -o ./tmp/atf

./tmp/atf run --project-id=$AUTIFY_TEST_PROJECT_ID --plan-id=$AUTIFY_TEST_PLAN_ID

rm ./tmp/atf
