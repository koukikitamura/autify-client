#!/usr/bin/env bash

DIR=$(cd $(dirname $BASH_SOURCE); pwd)
ROOT_DIR="${DIR}/.."
cd ${ROOT_DIR}

go build -o ./tmp/atf

./tmp/atf result --project-id=$AUTIFY_TEST_PROJECT_ID --result-id=$AUTIFY_TEST_RESULT_ID

rm ./tmp/atf
