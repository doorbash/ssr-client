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
      --dns=  custom dns (default: 8.8.8.8:53)
```

## Example
```
./ssr-client -s 1.2.3.4 -p 11800 -b 0.0.0.0 -l 1080 -k 1234 -m aes-256-cfb --dns=8.8.4.4:53
```

**Docker:**
```
docker run --name ssr --restart always -d -p 1080:1080 doorbash/ssr-client -s 1.2.3.4 -p 11800 -b 0.0.0.0 -l 1080 -k 1234 -m aes-256-cfb --dns=8.8.4.4:53
```
