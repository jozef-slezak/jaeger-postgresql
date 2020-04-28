This is the repository that contains PostgreSQL [Storage gRPC plugin for Jaeger](https://github.com/jaegertracing/jaeger/pull/1461).

> IMPORTANT: This plugin is still under development. We are using it internally
> already but the way we store data in PostgreSQL can change based on what do we
> learn about the data structure!

## Compile
In order to compile the plugin from source code you can use `go build`:

```
CGO_ENABLED=0 go build ./cmd/jaeger-pg-store/
```

## Start
In order to start plugin just tell jaeger the path to a config compiled plugin (password can be passed also as ENV: DB_PASSWORD).

```
jaeger-all-in-one --grpc-storage-plugin.binary=./jaeger-pg-store --grpc-storage-plugin.configuration-file=./config-example.yaml
```

## Tables
This plugins create tables if they not exist:

* spans
* span_logs
* span_refs
* operations
* services

## License

The PostgreSQL Storage gRPC Plugin for Jaeger is an [MIT licensed](LICENSE) open source project.
