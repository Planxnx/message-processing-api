#!/bin/sh
set -e

echo "Are you sure to deploy below services \n ${SVC_LIST// /\n }"
read -p "Please confirm [Y/N]? " -n 1 -r
echo # (optional) move to a new line

if [[ $REPLY =~ ^[Yy]$ ]]; then
    PROJECT="message-processing-api"
    DEPLOYNAME="apigateway"
    TAG="test"
    REPOSITORY="planxthanee/$PROJECT-$DEPLOYNAME"
    IMAGE="$REPOSITORY:$TAG"
    # HELM="__helm"

    read -p "Do you want to build code in this repo before deploy [Y/N]? " -n 1 -r
    echo # (optional) move to a new line
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cd ./producer-service
        docker build -t $IMAGE .
        docker push $IMAGE
    fi
fi