Forward2Docker
===
Utility to auto forward a port from localhost into ports on Docker containers running in a boot2docker or Docker Machine VM.

How it works?
---------
When it started it will listen for Docker events (start and die) and reconfigure port forwarding rules for your VirtualBox VM with Docker.

Why?
---------
Currently even with Boot2Docker you wouldn't get all Docker experience on your OS X host because
your Docker daemon will run in VM, not on localhost, which means that you will have to use
Docker VM's IP address instead of "localhost". It causes some fragmentation between native
and non-native Docker users. But we can solve it with "Port forwarding" feature in VirtualBox.

Install
---------
Tool is available in two options:

  1. Binaries. You can download the latest binaries from here: https://github.com/bsideup/forward2docker/releases/
  1. Go distribution. Install by running 'go get github.com/bsideup/forward2docker'
  1. Build it yourself. See "Contributing" section

Usage
---------

  1. Run some Docker container: `$ docker run --name f2dtest -d -p 8000:80 nginx`
  1. Open terminal and run forward2docker: `$ forward2docker` (NOTE: it runs in foreground, do not kill it, otherwise mappings will not be updated)
  1. Ensure that port is mapped to your host: `$ curl http://localhost:8000`
  1. Kill your container: `$ docker kill f2dtest`
  1. Ensure that port is unmapped (you should see 'Connection refused'): `$ curl http://localhost:8000`
  1. Run few more containers and verify that you can access them on localhost

Configuration
---------
no configuration required, but you can pass `--run-once` flag to prevent forward2docker to listen for events and quit right after the first port assignment.

Contributing
---------
GNU Make is used as a build tool. Following commands are available:
  - `make bootstrap` - you should call it (once) before you start. Will download all dependencies
  - `make build` - will run `go vet`, `go fmt` and compile binary for current platform
  - `make build_all` - will compile binaries for every supported platform. All binaries will be saved to `./bin/` folder