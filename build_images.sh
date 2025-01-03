#! /bin/bash
cd ./postgres
docker build . -t postgres_db --build-arg POSTGRES_USER=dbAdmin --build-arg POSTGRES_PASSWORD_FILE=/run/secrets/postgres_password --network=host
cd ../../app
docker build . -t go_app --network host