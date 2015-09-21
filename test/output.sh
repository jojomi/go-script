#!/usr/bin/env sh

echoerr() { echo "$@" 1>&2; }

echoerr "error"
echo "output"
echoerr "wrong"
echo "alright"

exit 0
