# asynq-put

`asynq-put` is a command line tool to put a task into an Asynq queue.

```shell
go install github.com/flaboy/asynq-put@latest
```

```shell
asynq-put mytopic '{"key": "value"}' -q myqueue -p 6379
```