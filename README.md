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

### technical

* provide 12factor environment variables
* listen to tcp/udp syslog traffic
* hand off events to parse
* forward set of MiruLogEvent objects via REST POST

### non-technical

* determine state of existing log aggregators (ELK) research
* determine ballpark cost for existing sumologic usage via mako
* ballpark hardware/aws cost for hosting miru-stumptown

### DONE

* determine how to shovel syslog
* $ kubectl logs <mako ms pods> -f | ncat ip:514
* deploy mako ms to minikube
* $ journalctl -f | ncat ip:514

* create miru-syslog/sample-golang in minikube
* IT in minikube
* get access to ms-integ
* used lemur to generate client cert, kubectl to set config
* used mako servicerepo to determine environment settings
* copy sumologic yaml to miru-syslog k8s spec

## Test Notes

```
make docker
make run

export MIRU_STUMPTOWN_HOST_PORT=10.126.5.155:10004
export MIRU_SYSLOG_HOST_PORT=`docker-machine ip`:514
go test -v --run Test.*Client
```

```
minikube start
minikube ip

kubectl create -f k8s.yml
kubectl get pod
kubectl get daemonset

export MIRU_STUMPTOWN_HOST_PORT=10.126.5.155:10004
export MIRU_SYSLOG_HOST_PORT=`minikube ip`:514
go test -v --run Test.*Client

kubectl logs miru-syslog-xxxxx
kubectl delete daemonset miru-syslog
```
