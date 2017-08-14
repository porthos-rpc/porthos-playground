#!/bin/bash

DOCKERFILE_BUILD="Dockerfile.build"
DOCKERFILE_PUSH="Dockerfile.push"
DIST_DIR="/go/src/github.com/porthos-rpc/porthos-playground/dist/"
BUILD_TAG="porthos/porthos-playground-build"
PUSH_TAG="porthos/porthos-playground"

echo "Building $DOCKERFILE_BUILD"
docker build -t $BUILD_TAG -f $DOCKERFILE_BUILD --no-cache . || exit 1

echo "Running $BUILD_TAG"
docker run -v "$(pwd)/dist:$DIST_DIR" -t $BUILD_TAG || exit 1

echo "Building $DOCKERFILE_PUSH"
docker build -t $PUSH_TAG -f $DOCKERFILE_PUSH --no-cache . || exit 1

echo "Deleting image $BUILD_TAG"
docker rmi $BUILD_TAG -f

echo "Pushing $PUSH_TAG"
docker push $PUSH_TAG
