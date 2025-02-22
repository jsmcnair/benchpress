# Benchpress

Simple utility to send HTTP requests to a server and measure the response times. Requests can be spread across multiple clients. 

Each client sends _n_ requests. Response statuses are counted and reported at the end, along with the number and percentage of total and successful requests. 

Requests can either be sent at a fixed rate per client (RPS), a defined sleep time between requests. In either case, if the duration is shorter than the time taken to make the request, then the request will be made immediately after the previous one completes.

## Installation

```shell
git clone https://github.com/jsmcnair/benchpress && \
  cd benchpress && \
  go build -o bp main.go && \
  sudo install bp /usr/local/bin && \
  rm -rf benchpress
```

## Usage

```shell
-c int
  	Number of clients to create. Defaults to 1. (default 1)
-n int
  	Number of requests to make per client. Defaults to 1. (default 1)
-r int
  	Requests per second to attempt to make per client. Defaults to 0.
-s duration
  	Time to sleep between requests. Defaults to 1 millisecond. (default 1ms)
-u string
  	URL to make requests to. If not passed, a local built-in server is created and requests are sent to that.
```

## Example

```shell
benchpress -c 2 -n 10 -r 5
```
```
Number of clients:  2
Number of requests per client:  10
Total requests:  20
Requests per second per client:  5
URL flag not passed, using built-in server.
Making requests...

Starting server on port 8080
Response counts by status code:
	200: 20/20 (100.00%)

Success: 20/20
Time taken: 2.000799553s
Successful requests per second: 9.996004
Total requests per second: 9.996004
```

## To do

- [ ] Add support for defining the HTTP method.
- [ ] Add support for defining the request body.
- [ ] Add support for defining the request headers.
- [ ] Add support for defining the request timeout.
- [ ] Add support for machine readable output.
- [ ] Add support for co-ordinating across multiple instances.