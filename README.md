# Hikari-go
The Golang version of Hikari

### About
Hikari, A network proxy tool which helps you break through network restrictions.

How it works:

```
Local
<---[SOCKS5 or HTTP protocol over local machine]--->
Hikari client
<---[Hikari protocol(encrypted) over Internet]--->
Hikari server
<---[Local data]--->
Target
```

### Usage
Client side:
> ./hikari-client-xx-xxx client.json

Server side:
> ./hikari-server-xx-xxx server.json

### Sample configuration
Client side:

```
{
  "socksAddress": "localhost", // local SOCKS5 proxy server address
  "socksPort": 1180, // local SOCKS5 proxy server port, using '0' to turn off SOCKS5 proxy server
  "httpAddress": "localhost", // local HTTP proxy server address
  "httpPort": 1190, // local HTTP proxy server port, using '0' to turn off HTTP proxy server
  "serverAddress": "localhost", // remote Hikari server address
  "serverPort": 9670, // remote Hikari server port
  "privateKey": "hikari", // authentication key, must be same with server side
  "secret": "hikari-secret" // encryption key, must be same with server side
}
```

Server side:

```
{
  "listenAddress": "localhost", // Hikari server address
  "listenPort": 9670, // associated port
  "privateKeyList": [ // authentication key list
    "hikari"
  ],
  "secret": "hikari-secret" // encryption key
}
```
