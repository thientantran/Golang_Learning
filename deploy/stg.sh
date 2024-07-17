#!/usr/bin/env bash

APP_NAME=food-delivery

docker load -i ${APP_NAME}.tar
docker rm -rf ${APP_NAME}

docker run -d --name ${APP_NAME} \
  --netword my-net \
  -e VIRTUAL_HOST="stg.thientantran.com" \
  -e LETSENCRYPT_HOST="stg.thientantran.com" \
  -e LETSENCRYPT_EMAIL="thientantm@gmail.com" \
  -e DB_URL="food-delivery-stg.cjxjxjxjxjxj.ap-southeast-1.rds.amazonaws.com" \
  -e S3_BUCKET="food-delivery-stg" \
  -e AWS_REGION="ap-southeast-1" \
  -e S3API_KEY="AKIAJXJXJXJXJXJXJXJX" \
  -e S3API_SECRET="adasdadasdadasdadasdadasdadasdadasdadasd" \
  -e SYSTEM_SECRET="tantonten" \
  -p 8080:8080 \
  ${APP_NAME}
