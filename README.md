### Building & Publishing Dockerfiles

#### Logging in to registry.redhat.io

**NOTE:** This may not be necessary.

1. Go to https://access.redhat.com/terms-based-registry/ and create an account.
2. Login and create a service account.
3. After receiving a token, run:

``` shell
docker login registry.redhat.io
```

#### Building images for testing

``` shell
# use one of the directory names in the dockerfiles directory (e.g. linux-builder)
export IMAGE_NAME="image_name"
export DOCKERHUB_USERNAME="your_dockerhub_username_here"
export TAG_NAME="${IMAGE_NAME}-${DOCKERHUB_USERNAME}"
docker build . -f dockerfiles/linux-publisher/Dockerfile -t sensu/sensu-release:${TAG_NAME}
docker push sensu/sensu-release:${TAG_NAME}
```

#### Building images for release

``` shell
# use one of the directory names in the dockerfiles directory (e.g. linux-builder)
export IMAGE_NAME="image_name"
docker build . -f dockerfiles/linux-publisher/Dockerfile -t sensu/sensu-release:${IMAGE_NAME}
docker push sensu/sensu-release:${IMAGE_NAME}
```
