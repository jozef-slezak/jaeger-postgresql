This is the repository that contains PostgreSQL Storage gRPC plugin for Jaeger.

> IMPORTANT: This plugin is still under development. We are using it internally
> already but the way we store data in PostgreSQL can change based on what do we
> learn about the data structure!

The Jaeger community made a big work supporting external gRPC plugin to manage
integration with external backend without merge them as part of the Jaeger
Tracer. You can see issue [#1461](https://github.com/jaegertracing/jaeger/pull/1461).

The plugin uses `go/mod` to manage its dependencies.

## License

The PostgreSQL Storage gRPC Plugin for Jaeger is an [MIT licensed](LICENSE) open source project.

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
