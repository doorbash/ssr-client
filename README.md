## Build
```
go build
```

## Usage:
```
ssr-client [OPTIONS]
```

**Options:**
```
  -s=         server address
  -p=         server port (default: 8388)
  -b=         local binding address (default: 127.0.0.1)
  -l=         local socks port (default: 1080)
  -r=         local http port (default: 1081)
  -k=         password
  -m=         encryption method (default: aes-256-cfb)
  -o=         obfsplugin (default: http_simple)
      --op=   obfs param
  -O=         protocol (default: origin)
      --Op=   protocol param
  -f=         socks5 proxy address. example: 127.0.0.1:8080
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
docker run --net host --name ssr1 --restart always -d doorbash/ssr-client -s 1.2.3.4 -p 11800 -b 0.0.0.0 -l 8080 -r 8081 -k 1234 -m aes-256-cfb

docker run --net host --name ssr2 --restart always -d doorbash/ssr-client -s 5.6.7.8 -p 11800 -b 0.0.0.0 -l 9090 -r 9091 -k 1234 -m aes-256-cfb -f 127.0.0.1:8080
```