# miru-syslog

Syslog collector that listens for syslog traffic, parses, and forwards to miru-stumptown.

It is written in golang, housed in a docker container, deployed via mako, into a kubernetes cluster.

## Endpoints

### Main

* tcp port 514
* udp port 514

## Admin

* /admin root html
* /ping 200 status
* /healthcheck json
* /healthcheck?pretty_print=true
* /metrics json
* /metrics?pretty_print=true

## Environment Variables

### Miru

* miru-stumptown host address
* miru-stumptown host port
* miru-stumptown intake url

### Listener

* tcp listen port
* udp listen port
* queue size
* batch size

### Mako

* MAKO_SERVICE_ID
* MAKO_ENVIRONMENT
* MAKO_PIPELINE
* MAKO_VERSION
* MAKO_STATSD_HOST
* MAKO_STATSD_PORT

## TODO

* provide variables via environment in true 12factor fashion
* listen on tcp/udp port 514 for syslog traffic
* hand off traffic to parse
* forward set of MiruLogEvent objects via REST POST
