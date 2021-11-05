# jq-sandbox

## Setup
```
$ brew install nats-server
```

Run NATS server in a Jetstream mode. Jetstream is GA in the NATS server since 2.2.0
```
$ nats-server --js
```

Get the supporting NATS go library
```
$ go get github.com/nats-io/nats.go/
```

## Execution
Run n number of workers in different terminals; example run this in 3 terminals
```
$ cd nats-jsmq-sub
$ go run job_worker.go
```

Run the job producer in a separate terminal
```
$ cd nats-jsmq-pub
$ go run job_producer.go
```

## Happy Path:
10 Jobs are produced by the job producer and the workers randomly get jobs and work on them.

## Sad Path: One consumer crashes
10 Jobs are produced by the job producer and kill one of the workers randomly in the middle. All the jobs will get completed by the remaining workers.

