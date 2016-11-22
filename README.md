# miru-syslog

Syslog collector that listens for syslog traffic, parses, and forwards to miru-stumptown.

It is written in golang, housed in a docker container, deployed as a DaemonSet into a kubernetes cluster.

## Endpoints

### Main

* tcp port 514
* udp port 514

## Environment Variables

### Miru

* miru-stumptown host addr:port
* miru-stumptown intake url - default to /miru/stumptown/intake

### Listener

* tcp listen port - default to :514
* udp listen port - default to empty/off
* queue size
* batch size

## TODO

* create miru-syslog/sample-golang in minikube
* IT in minikube
* get access to ms-integ
* copy sumologic yaml k8s spec
* massage to fix miru-syslog
* create miru-syslog in k8s
* provide 12factor environment variables
* listen to tcp/udp syslog traffic
* hand off events to parse
* forward set of MiruLogEvent objects via REST POST
