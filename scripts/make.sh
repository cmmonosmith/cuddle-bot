#!/usr/bin/env bash

VERBOSE=0

clean () {
    go clean $( [ "${VERBOSE}" = "1" ] && echo "-x" )
    rm -f bin/${project}
}

compile () {
    go build -o bin/${project} $( [ "${VERBOSE}" = "1" ] && echo "-x -v" )
    chmod +x bin/${project}
}

package () {
    APP_NAME="${project}" docker build -f scripts/Dockerfile .
}

build () {
    compile
    package    
}

#start () {
#    stop
#}

#stop () {

#}

# entry point
cd -- $( dirname -- "${BASH_SOURCE[0]}")/..
project="${PWD##*/}"

# parse args
for arg in "$@"
do
    case "$arg" in
        "--verbose")
            VERBOSE=1
            ;;
        "--debug")
            export DISCORD_BOT_LOG_DEBUG=1
            ;;
        *)
            command="${command} ${arg}"
            ;; 
    esac
    shift
done
$( echo "${command}" )
