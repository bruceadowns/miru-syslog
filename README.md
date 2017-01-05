# miru-syslog

Syslog collector that listens for syslog traffic, parses, forwards to miru-stumptown and aws s3.

It is written in golang, housed in a docker container, and may be deployed as a DaemonSet into a Kubernetes cluster.

## Endpoints

* tcp port 514

## Environment Variables

### Minimal

* MIRU_SYSLOG_TCP_ADDR_PORT - required to listen for traffic
* MIRU_STUMPTOWN_ADDR_PORT - required if posting to stumptown
* AWS_ACCESS_KEY_ID - required if posting to S3
* AWS_SECRET_ACCESS_KEY - required if posting to S3

### Full Listing

* MIRU_SYSLOG_TCP_ADDR_PORT - required to listen for traffic
* MIRU_STUMPTOWN_ADDR_PORT - required to post to stumptown
* MIRU_STUMPTOWN_INTAKE_URL - default to /miru/stumptown/intake
* CHANNEL_BUFFER_SIZE_PARSE - default to 1k
* CHANNEL_BUFFER_SIZE_MIRU_ACCUM - default to 1k
* CHANNEL_BUFFER_SIZE_MIRU_POST - default to 1k
* CHANNEL_BUFFER_SIZE_S3_ACCUM - default to 1k
* CHANNEL_BUFFER_SIZE_S3_POST - default to 1k
* CHANNEL_BUFFER_MIRU_ACCUM_BATCH - default to 1k
* CHANNEL_BUFFER_MIRU_ACCUM_DELAY_MS - default to 1s
* CHANNEL_BUFFER_MIRU_DELAY_ON_SUCCESS_MS - default to 1/2 s
* CHANNEL_BUFFER_MIRU_DELAY_ON_ERROR_MS - default to 5s
* CHANNEL_BUFFER_S3_ACCUM_BATCH_BYTES - default to 10Mb
* CHANNEL_BUFFER_S3_ACCUM_DELAY_MS - default to 8h
* CHANNEL_BUFFER_S3_DELAY_ON_SUCCESS_MS - default to 1s
* CHANNEL_BUFFER_S3_DELAY_ON_ERROR_MS - default to 5s
* AWS_REGION - default to us-west-2
* AWS_S3_BUCKET_NAME - default to miru-syslog
* AWS_ACCESS_KEY_ID - required to post to S3
* AWS_SECRET_ACCESS_KEY - required to post to S3

## Execution Notes

### Run via golang

```
go get github.com/bruceadowns/miru-syslog
export MIRU_SYSLOG_TCP_ADDR_PORT=:8514
export MIRU_STUMPTOWN_ADDR_PORT=10.126.5.155:10000
go run main.go
```

or

```
go get github.com/bruceadowns/miru-syslog
export MIRU_SYSLOG_TCP_ADDR_PORT=:8514
export MIRU_STUMPTOWN_ADDR_PORT=10.126.5.155:10000
export AWS_ACCESS_KEY_ID=<my ak id>
export AWS_SECRET_ACCESS_KEY=<my secret ak>
go run main.go
```

### Run via docker

```
make docker
make run
```

### Run via minikube

```
minikube start
minikube ip
kubectl cluster-info

kubectl create -f docker/k8s.yml

kubectl get daemonset
kubectl get pods

kubectl logs miru-syslog-xxxxx
kubectl delete daemonset miru-syslog

minikube stop
```

### Test tcp client via minikube

```
export MIRU_STUMPTOWN_ADDR_PORT=10.126.5.155:10000
export MIRU_SYSLOG_TCP_ADDR_PORT=`minikube ip`:514
go test -v --run TestTcpClient
```

## References

* https://golang.org/
* https://www.docker.com/
* http://kubernetes.io/
* http://kubernetes.io/docs/admin/daemons/
* https://github.com/kubernetes/minikube/
