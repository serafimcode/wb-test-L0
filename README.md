# wb-test-L0

All commands are expected to run from workspace root
- Build subscriber:
`go build -o sub subscriber/main.go`
- Run subscriber:
`./sub`
- Build publisher:
`go build -o pub publisher/main.go`
- Run publisher:
`./pub`
- To start nats-streaming:
`docker run -p 4223:4223 -p 8223:8223 nats-streaming -p 4223 -m 8223 -cid wb-cluster`
