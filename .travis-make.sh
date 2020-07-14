#!/bin/sh

cd $(dirname $0)
echo "$0: Running $@ ... "
date
# work around Travis CI heuristics
( while sleep 60; do echo -n '.' ; done) &
trap "echo; date; kill $!" EXIT
"$@" >> .travis-make.log 2>&1
