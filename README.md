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

* get access to ms-integ
* copy sumologic yaml to miru-syslog k8s spec
* provide 12factor environment variables
* listen to tcp/udp syslog traffic
* hand off events to parse
* forward set of MiruLogEvent objects via REST POST

## DONE

* create miru-syslog/sample-golang in minikube
* IT in minikube

```
minikube start
minikube ip

MIRU_STUMPTOWN_HOST_PORT=10.126.5.155:10004 MIRU_SYSLOG_HOST_PORT=`minikube ip`:514 go test -v --run TestTcpClient
MIRU_STUMPTOWN_HOST_PORT=10.126.5.155:10004 MIRU_SYSLOG_HOST_PORT=`minikube ip`:514 go test -v --run TestUdpClient

kubectl create -f k8s.yml
kubectl get pod
kubectl get daemonset
kubectl logs miru-syslog-xxxxx
kubectl delete daemonset miru-syslog
```
