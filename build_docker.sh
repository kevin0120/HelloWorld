#!/bin/bash

version="0.2.55"

docker_repo="kevin/hello1"

docker build -t ${docker_repo}:${version} -t ${docker_repo}:latest .

docker push ${docker_repo}:${version}
docker push ${docker_repo}:latest
