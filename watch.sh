#!/bin/bash
ls *.go cmd/*.go test_data/*.xml |entr -c docker-compose -f docker/docker-compose.yml up &
ls docker/Dockerfile docker/docker-compose.yml Gopkg.toml |entr -c docker-compose -f docker/docker-compose.yml up --build
