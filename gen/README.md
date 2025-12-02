# `/gen`

Generated pb files using:

```shell
protoc --go_out=.../gen --go_opt=paths=source_relative \
    --go-grpc_out=../gen --go-grpc_opt=paths=source_relative \
    greet/*.proto
```
> [!WARNING]
> You need to be inside of the folder `/proto`

Or just use [build.sh](/build.sh)
