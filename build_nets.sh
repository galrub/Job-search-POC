#! /bin/bash
docker network create --driver overlay jobsearch_backend
docker network create --driver overlay jobsearch_frontend --opt com.docker.network.bridge.host_binding_ipv4=127.0.0.1