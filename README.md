# golang-http-filter
This is a simple HTTP server written in Go that only allows requests made by wget and curl. Any other requests are blocked, and the client's IP is added to iptables to prevent further access.
