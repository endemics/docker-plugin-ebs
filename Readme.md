# Docker Volume Plugin for EBS support

This project aims to allow the use of EBS-backed Docker volumes by implementing the [docker volume plugin API](https://github.com/docker/docker/blob/master/docs/extend/plugins_volume.md).

It uses the docker/go-plugins-helpers to do so and wouldn't have been possible without David Calavera's Docker volume garbage generator, from wich it inherits its MIT license.

The plugin is meant to be installed on docker hosts in AWS and to have no added dependency. It directly accesses the EC2 API and should get its privileges and credentials using IAM instance profile.

**This is still a WIP at this point and should not be used ;)**
