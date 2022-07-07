#!/bin/sh

# nginx options (taken from new docker of nginx: `docker run --rm -it nginx sh -c 'nginx -h'`
# nginx -h
#nginx version: nginx/1.21.5
#Usage: nginx [-?hvVtTq] [-s signal] [-p prefix]
#             [-e filename] [-c filename] [-g directives]
#Options:
#  -?,-h         : this help
#  -v            : show version and exit
#  -V            : show version and configure options then exit
#  -t            : test configuration and exit
#  -T            : test configuration, dump it and exit
#  -q            : suppress non-error messages during configuration testing
#  -s signal     : send signal to a master process: stop, quit, reopen, reload
#  -p prefix     : set prefix path (default: /etc/nginx/)
#  -e filename   : set error log file (default: /var/log/nginx/error.log)
#  -c filename   : set configuration file (default: /etc/nginx/nginx.conf)
#  -g directives : set global directives out of configuration file

# because we are running a command after the nginx we leave the daemon on (by not passing the -g 'daemon off')
# "$@" - to enable passing the parameters (-s reload)
nginx -p /etc/nginx/ -e /var/log-nginx-error.log -c /etc/nginx/nginx.conf "$@"