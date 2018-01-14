#!/usr/bin/sh

REMOTE=$1
git remote add resin $REMOTE
git push resin master
