### Building & Publishing Dockerfiles

``` shell
for i in `ls dockerfiles`; do docker build . -f dockerfiles/$i/Dockerfile -t sensu/sensu-release:$i; docker push sensu/sensu-release:$i; done
```
