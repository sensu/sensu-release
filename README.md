### Building & Publishing Dockerfiles

Docker images are automatically built and published for every commit. The images
are tagged differently depending on whether or not it was built from a git
commit or a git tag.

#### Git Commits

A git commit with the sha `f412c9cd4181f2b25d4db946914c8096cb25ccbd`will
generate the following Docker tags:

```
sensu/sensu-release:f412c9cd4181f2b25d4db946914c8096cb25ccbd-linux-builder
sensu/sensu-release:f412c9cd4181f2b25d4db946914c8096cb25ccbd-rhel7-builder
sensu/sensu-release:f412c9cd4181f2b25d4db946914c8096cb25ccbd-rhel8-builder
sensu/sensu-release:f412c9cd4181f2b25d4db946914c8096cb25ccbd-packagecloud-pruner
```

#### Git Tags

A git tag with the name `v1.2.3` will generate the following Docker tags:

```
sensu/sensu-release:v1.2.3-linux-builder
sensu/sensu-release:v1.2.3-rhel7-builder
sensu/sensu-release:v1.2.3-rhel8-builder
sensu/sensu-release:v1.2.3-packagecloud-pruner
```
