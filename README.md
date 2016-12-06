# miru-syslog

Syslog collector that listens for syslog traffic, parses, and forwards to miru-stumptown.

It is written in golang, housed in a docker container, deployed as a DaemonSet into a Kubernetes cluster.

The syslog protocol is defined via rfc3164 and rfc5424.

## Endpoints

* tcp port 514
* udp port 514

## Environment Variables

* MIRU_SYSLOG_TCP_ADDR_PORT - if empty, do not listen for tcp traffic
* MIRU_SYSLOG_UDP_ADDR_PORT - if empty, do not listen for udp traffic
* MIRU_STUMPTOWN_ADDR_PORT - if empty, do not post to stumptown
* MIRU_STUMPTOWN_INTAKE_URL - default to /miru/stumptown/intake
* CHANNEL_BUFFER_SIZE_PARSE - default to 1024
* CHANNEL_BUFFER_SIZE_POST - default to 1024
* UDP_RECEIVE_BUFFER_SIZE - default to 2*1024*1024

Note, syslog message size _should not_ exceed 1024 bytes per rfc. Though we default to 2mb.

## TODO

### technical

* determine event type (syslog, dates, json)
* parse each event type
* forward set of MiruLogEvent objects via REST POST
* cache map[string]string  event type to remote host

### non-technical

* determine state of existing log aggregators (ELK) research
* ballpark hardware/aws cost for hosting miru-stumptown

### DONE

* determine ballpark cost for existing sumologic usage via mako
* hand off events to parse
* listen to tcp/udp syslog traffic
* provide 12factor environment variables
* determine how to shovel syslog
* $ kubectl logs <mako ms pods> -f | ncat ip:514
* deploy mako MSs/miru-syslog/sample-golang to minikube
* $ journalctl -f | ncat ip:514
* IT in minikube
* get access to ms-integ
* used lemur to generate client cert, kubectl to set config
* used mako servicerepo to determine environment settings
* copy sumologic yaml to miru-syslog k8s spec

## Test Notes

### Run minikube

```
minikube start
minikube ip
kubectl cluster-info
minikube stop
```

### Create kubernetes daemonset

```
kubectl create -f docker/k8s.yml

kubectl get daemonset
kubectl get pods

kubectl logs miru-syslog-xxxxx
kubectl delete daemonset miru-syslog
```

### Run main.go via docker

```
make docker
make run
```

### Test tcp client

```
export MIRU_STUMPTOWN_ADDR_PORT=10.126.5.155:10004
export MIRU_SYSLOG_TCP_ADDR_PORT=`minikube ip`:514
go test -v --run TestTcpClient
```

### Test udp client

```
export MIRU_STUMPTOWN_ADDR_PORT=10.126.5.155:10004
export MIRU_SYSLOG_UDP_ADDR_PORT=`minikube ip`:514
go test -v --run TestUdpClient
```

## References

* https://golang.org/
* https://www.docker.com/
* http://kubernetes.io/
* http://kubernetes.io/docs/admin/daemons/
* https://github.com/kubernetes/minikube
* https://www.ietf.org/rfc/rfc3164.txt
* https://www.ietf.org/rfc/rfc5424.txt
