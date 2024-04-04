# protoc-gen-authorize


###  Install
```
go install github.com/gh1st/protoc-gen-authorize
```

### buf.gen.yaml
```yaml
version: v1
plugins:
  - name: go
    out: gen
    opt: paths=source_relative
  - name: go-grpc
    out: gen
    opt: paths=source_relative
  - name: authorize
    out: gen
    opt: paths=source_relative
```

### Example
See: [example.proto](proto/example/v1/example.proto)