#!/usr/bin/env sh

test () {
    SCRIPT_NAME=$1
    RES=$(go run . lox_scripts/${SCRIPT_NAME})

    EXP=$2

    if [ "${RES}" = "${EXP}" ]; then
        echo "${SCRIPT_NAME}: passed"
    else
        echo "test failed"
        echo "${SCRIPT_NAME}: expected \"${EXP}\" to be equal to \"${RES}\""
    fi
}

test "hello_world.lox" "Hello, World!"

test "scope_test.lox" "inner a
outer b
global c
outer a
outer b
global c
global a
global b
global c"

test "if_test.lox" "was hello!"
test "and_or.lox" "hi
yes"

test "assign.lox" "init
after
last"
