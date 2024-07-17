#!/usr/bin/env bash

APP_NAME=food-delivery
DEPLOY_CONNECT=ec2-user@10.10.10.10

echo "Downloading packages..."
go mod download
echo "Compiling..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app

echo "Docker building..."
docker build -t ${APP_NAME} -f ./Dockerfile .
echo "Docker saving..."
docker save -o ${APP_NAME}.tar ${APP_NAME}

echo "Deploying..."
scp -o StrictHostKeyChecking=no ./${APP_NAME}.tar ${DEPLOY_CONNECT}:~
ssh -o StrictHostKeyChecking=no ${DEPLOY_CONNECT} 'bash -s' < ./deploy/stg.sh
#
echo "Cleaning..."
rm -f ./${APP_NAME}.tar

echo "Done!"