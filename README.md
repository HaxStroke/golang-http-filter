# golang-http-filter
This is a simple HTTP server written in Go that only allows requests made by wget and curl. Any other requests are blocked, and the client's IP is added to iptables to prevent further access.

# Features

✅ Serves files from the current directory

✅ Allows only requests with wget or curl User-Agent

✅ Logs all incoming requests

✅ Blocks IPs making invalid requests for 10 minutes

✅ Adds blocked IPs to iptables for enhanced security

# How to Use

Clone the repository and navigate to the project directory:

git clone https://github.com/your-username/golang-http-filter.git
cd golang-http-filter

Compile and run the server:

go run main.go

The server will run on port 8080, serving files from the current directory.

# Usage Example

Allowed access:

curl http://localhost:8080
wget http://localhost:8080

Denied access and blocking:

python -c "import requests; requests.get('http://localhost:8080')"

# Warning

This script modifies iptables rules. Use it with caution and verify the rules before running it in production.
