# docker-compose-daemon
Starting docker containers via docker-compose, redirect docker-compose logs to stderr and stdout, monitoring container state

USAGE:
   docker-compose-daemon [global options] command [command options] [arguments...]

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --configuration, -f 						Docker compose config file: -f /path/to/docker-compose.yml
   --container, -c [--container option --container option]	Full container name: -c demo_app_1 -c demo_db_1 -c demo_web_1
   --timeout, -t "5"						Timeout for container monitoring
   --help, -h							show help
   --version, -v
