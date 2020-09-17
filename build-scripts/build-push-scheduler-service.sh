#!/bin/sh
set -e
SVCNAME="scheduler"
echo "Are you sure to deploy ${SVCNAME} service"
read -p "Please confirm [Y/N]? " -n 1 -r
echo # (optional) move to a new line

if [[ $REPLY =~ ^[Yy]$ ]]; then
    PROJECT="message-processing-api"
    DEPLOYNAME=SVCNAME
    TAG="latest"
    REPOSITORY="planxthanee/$PROJECT-$DEPLOYNAME"
    IMAGE="$REPOSITORY:$TAG"
    # HELM="__helm"

    read -p "Do you want to build code in this repo before deploy [Y/N]? " -n 1 -r
    echo # (optional) move to a new line
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cd ./scheduler-service
        docker build -t $IMAGE .
        docker push $IMAGE
    fi
fi