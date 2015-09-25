# docker-compose-daemon
Starting docker containers via docker-compose, redirect docker-compose logs to stderr and stdout, monitoring containers state

## Compiling

This section assumes you're familiar with the Go language.

Use <code>go get</code> to get the source local:

```bash
$ go get github.com/maxposter/docker-compose-daemon
```

Change to the directory, e.g.:

```bash
$ cd $GOPATH/src/github.com/maxposter/docker-compose-daemon
```

Get the dependencies:

```bash
$ go get ./...
```

Then build and/or install:

```bash
$ go build
$ go install
```

### USAGE
```
USAGE:
   docker-compose-daemon --configuration /path/to/docker-compose.yml --container name [--container name, ...]

OPTIONS:
   --configuration, -f 						Docker compose config file: -f /path/to/docker-compose.yml
   --container, -c [--container option --container option]	Full container name: -c demo_app_1 -c demo_db_1 -c demo_web_1
   --timeout, -t "5"						Timeout for container monitoring
   --help, -h							show help
   --version, -v
```
