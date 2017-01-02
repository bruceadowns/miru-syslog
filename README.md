# miru-syslog

Syslog collector that listens for syslog traffic, parses, and forwards to miru-stumptown.

It is written in golang, housed in a docker container, deployed as a DaemonSet into a Kubernetes cluster.

## Endpoints

* tcp port 514

## Environment Variables

* MIRU_SYSLOG_TCP_ADDR_PORT - required to listen for traffic
* MIRU_STUMPTOWN_ADDR_PORT - required to post to stumptown
* MIRU_STUMPTOWN_INTAKE_URL - default to /miru/stumptown/intake

* CHANNEL_BUFFER_SIZE_PARSE - default to 1024
* CHANNEL_BUFFER_SIZE_MIRU_ACCUM - default to 1024
* CHANNEL_BUFFER_SIZE_MIRU_POST - default to 1024
* CHANNEL_BUFFER_MIRU_ACCUM_BATCH - default to 1000
* CHANNEL_BUFFER_MIRU_ACCUM_DELAY_MS - default to 100

* CHANNEL_BUFFER_SIZE_S3_ACCUM - default to 1000
* CHANNEL_BUFFER_SIZE_S3_POST - default to 1000
* CHANNEL_BUFFER_S3_ACCUM_BATCH_BYTES - default to 10Mb
* CHANNEL_BUFFER_S3_ACCUM_DELAY_MS - default to 1d

* AWS_REGION - default to us-west-2
* AWS_S3_BUCKET_NAME - default to miru-syslog
* AWS_ACCESS_KEY_ID - required to post to S3
* AWS_SECRET_ACCESS_KEY - required to post to S3

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

## References

* https://golang.org/
* https://www.docker.com/
* http://kubernetes.io/
* http://kubernetes.io/docs/admin/daemons/
* https://github.com/kubernetes/minikube/
