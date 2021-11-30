## Build
```
# windows系统中编译，目标为windows
set GOOS=windows GO111MODULE=on GOARCH=amd64 CGO_ENABLED=0
go build -a -buildmode=exe -trimpath -ldflags "-s -w" -o ssr-client.exe
# linux系统中编译，目标为linux
GOOS=linux GO111MODULE=on GOARCH=amd64 CGO_ENABLED=0 go build -a -buildmode=exe -trimpath -ldflags "-s -w" -o ssr-client
# linux系统中编译，目标为arm64
GOOS=linux GO111MODULE=on GOARCH=arm64 CGO_ENABLED=0 go build -a -buildmode=exe -trimpath -ldflags "-s -w" -o ssr-client
# linux系统中编译，目标为ios
GOOS=darwin GO111MODULE=on GOARCH=arm64 CGO_ENABLED=0 go build -a -buildmode=exe -trimpath -ldflags "-s -w" -o ssr-client
# linux系统中编译，编译目标为Mac
GOOS=darwin GO111MODULE=on GOARCH=amd64 CGO_ENABLED=0 go build -a -buildmode=exe -trimpath -ldflags "-s -w" -o ssr-client
```

## Usage:
```
ssr-client [OPTIONS]
```

**Options:**
```
  -s=            server address
  -p=            server port (default: 8388)
  -b=            local binding address (default: 127.0.0.1)
  -l=            local socks port (default: 1080)
  -r=            local http port (default: 1081)
  -k=            password
  -m=            encryption method (default: aes-256-cfb)
  -o=            obfsplugin (default: http_simple)
      --op=      obfs param
  -O=            protocol (default: origin)
      --Op=      protocol param
  -t=            socket timeout in seconds (default: 10)
  -f, --forward= socks5 forward proxy address. example: 127.0.0.1:8080
  -v             verbose mode
```

## Example
```
./ssr-client -s 1.2.3.4 -p 11800 -b 0.0.0.0 -l 1080 -r 1081 -k 1234 -m aes-256-cfb
```

```
./ssr-client -s 5.6.7.8 -p 11800 -b 0.0.0.0 -l 1090 -r 1091 -k 1234 -m aes-256-cfb -f 127.0.0.1:1080
```
**Docker:**
```
docker run --name ssr --restart always -p 1080:1080 -p 1081:1081 -d doorbash/ssr-client -s 1.2.3.4 -p 11800 -b 0.0.0.0 -k 1234 -m aes-256-cfb
```

```
docker run --net host --name ssr2 --restart always -d doorbash/ssr-client -s 5.6.7.8 -p 11800 -b 0.0.0.0 -l 1090 -r 1091 -k 1234 -m aes-256-cfb -f 127.0.0.1:1080
```
