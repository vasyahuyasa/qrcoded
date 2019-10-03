# qrcoded

__WORK IN PROGRESS__

QR code generator server and cli tool

## Install

```sh
$ go get -u github.com/vasyahuyasa/qrcoded/cmd/qrcoded
```

## Start server

```
Usage of ./qrcoded:
  -d    verbose information
  -h string
        host (default "0.0.0.0")
  -p int
        port (default 80)
```

`$ ./qrcoded -p 80 -d`

## Make query

Open in browser

`http://localhost:80/?q=test&s=1024&r=l`

| Param | Default | Required | Comment |
|-------|---------|----------|---------|
|q      |         | *        | Text for encoding |
|s      | 256     |          | Size of output picture in pixels |
|r      | m       |          | Error correction level, can be l, m, q, h |
