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
  -l=         local port (default: 1080)
  -k=         password
  -m=         encryption method (default: aes-256-cfb)
  -o=         obfsplugin (default: http_simple)
      --op=   obfs param
  -O=         protocol (default: origin)
      --Op=   protocol param
  -t=         socket timeout in seconds (default: 10)
  -f=         socks5 forward proxy address. example: 127.0.0.1:8080
  -v          verbose mode
```

## Example
```
./ssr-client -s 1.2.3.4 -p 11800 -b 0.0.0.0 -l 1080 -k 1234 -m aes-256-cfb
```

```
./ssr-client -s 5.6.7.8 -p 11800 -b 0.0.0.0 -l 1090 -k 1234 -m aes-256-cfb -f 127.0.0.1:1080
```
**Docker:**
```
docker run --name ssr --restart always -p 1080:1080 -d ghcr.io/doorbash/ssr-client -s 1.2.3.4 -p 11800 -b 0.0.0.0 -k 1234 -m aes-256-cfb
```

```
docker run --net host --name ssr2 --restart always -d ghcr.io/doorbash/ssr-client -s 5.6.7.8 -p 11800 -b 0.0.0.0 -l 1090 -k 1234 -m aes-256-cfb -f 127.0.0.1:1080
```