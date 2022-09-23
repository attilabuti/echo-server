# Echo Server

## Installation

Just clone this repo and run it:

```
git clone https://github.com/attilabuti/echo-server.git
cd echo-server
go run .
```

## CLI arguments

```shell
--host host             Server host
--enable-http           Enable HTTP server (default: false)
--port-http port        HTTP port (default: 80)
--enable-https          Enable HTTPS server (default: false)
--port-https port       HTTPS port (default: 443)
--crt-file file         Location of the SSL certificate file
--key-file file         Location of the RSA private key file
--content value         Response body (default: "ok")
--content-file file     Response body from file
--content-type value    Content-Type header (default: "text/plain; charset=UTF-8")
--enable-tcp            Enable TCP echo server (default: false)
--port-tcp port         TCP echo port (default: random)
--enable-udp            Enable UDP echo server (default: false)
--port-udp port         UDP echo port (default: random)
--enable-log            Enable file logging (default: false)
--log-dir value         Location of the log directory (default: "log")
--log-requests          Log HTTP(S) requests (default: true)
--log-connections       Log TCP connections (default: true)
--log-packets           Log TCP/UDP echo packets (default: true)
--config file, -c file  Location of the configuration file in .yml format
--quiet, -q             Activate quiet mode (default: false)
--help, -h              Print this help text and exit
--version, -v           Print program version and exit
```

## Configuration file

| Property | Type | Default | Description |
|:---|:---|:---|:---|
| `host` | `string` | | Server host |
| `enable-http` | `bool` | `false` | Enable HTTP server |
| `port-http` | `int` | `80` | HTTP port |
| `enable-https` | `bool` | `false` | Enable HTTPS server |
| `port-https` | `int` | `443` | HTTPS port |
| `crt-file` | `string` | | Location of the SSL certificate file |
| `key-file` | `string` | | Location of the RSA private key file |
| `content` | `string` | `ok` | Response body |
| `content-file` | `string` | | Response body from file |
| `content-type` | `string` | `text/plain; charset=UTF-8` | Content-Type header |
| `enable-tcp` | `bool` | `false` | Enable TCP echo server |
| `port-tcp` | `int` | `0` | TCP echo port |
| `enable-udp` | `bool` | `false` | Enable UDP echo server |
| `port-udp` | `int` | `0` | UDP echo port |
| `enable-log ` | `bool` | `false` | Enable file logging |
| `log-dir` | `string` | `log` | Location of the log directory |
| `log-requests` | `bool` | `true` | Log HTTP(S) requests |
| `log-connections` | `bool` | `true` | Log TCP connections |
| `log-packets` | `bool` | `true` | Log TCP/UDP echo packets |
| `quiet` | `bool` | `false` | Activate quiet mode |

## Issues

Submit the [issues](https://github.com/attilabuti/echo-server/issues) if you find any bug or have any suggestion.

## Contribution

Fork the [repo](https://github.com/attilabuti/echo-server) and submit pull requests.

## License

This project is licensed under the [MIT License](https://github.com/attilabuti/echo-server/blob/main/LICENSE).