#!/bin/bash
version=latest
docker build . -t jonas27test/goserver:$version
# docker run -p 8080:8080 -p 4443:4443 -v $HOME/repos/goserver/cert:/cert -v $HOME/repos/goserver/static:/static jonas27test/goserver:v0.0.1 
docker push jonas27test/goserver:$version