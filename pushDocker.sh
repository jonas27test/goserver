#!/bin/bash
version="v1.0.1"
docker build . -t jonas27test/goserver:$version
# docker run -p 8080:8080 -v $HOME/repos/goserver/static:/static jonas27test/goserver:v1.0.1
docker push jonas27test/goserver:$version