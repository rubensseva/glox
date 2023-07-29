#!/usr/bin/env sh

RES=$(go run . lox_scripts/if_test.lox)
EXP="was hello!"

if [ "${RES}" = "${EXP}" ]; then
    echo "test passed"
else
    echo "test failed"
    echo "expected \"${EXP}\" to be equal to \"${RES}\""
fi
