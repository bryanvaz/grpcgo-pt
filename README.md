This a personal reference repo for gRPC implementation and perf testing in Go.

If you find some stuff useful, feel free to use it.


## Usage
### Server
Uses Air for hot reloading the server. 

```bash
air
```

### Client
```bash
go run client/main.go plow -n 1 -c 1
```

### Rebuild Protobufs and gRPC libs
```bash
make proto
```


