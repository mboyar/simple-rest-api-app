# simple-rest-api-app - Rest API application in golang

## Summary

A simple REST API application in golang, which retrieves system bootup duration (kernel+userspace) and responds it to a simple GET request by using systemd-analyze command, and net/http and os/exec packages

The REST API provided by server consists of following api calls:

- Version -> displays the version of the server app on the cli
- Duration -> displays startup duration of the system
## Execution output

```
$ go build -o sysinfo_server
$ ./sysinfo_server &
Server ready, endpoints: /version and /duration 

$ curl http://localhost:8080/
Server ready, endpoints: /version and /duration 

$ curl http://localhost:8080/version
Version: v0.1

$ curl http://localhost:8080/duration
Startup duration of the system: Startup finished in 5.225s (kernel) + 2min 15.289s (userspace) = 2min 20.515s 
graphical.target reached after 2min 15.220s in userspace

$ curl http://localhost:8080/duration.json
{"bootup":{"kernel":5.763,"initrd":0,"userspace":9.894,"graphical.target":9.883},"time-unit":"seconds"}
```
