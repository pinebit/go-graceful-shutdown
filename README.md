# go-graceful-shutdown

This project demonstrates a pattern for graceful server shutdown in Go. It was inspired by this
article: https://rudderstack.com/blog/implementing-graceful-shutdown-in-go

The sample server maintains:

- a connection to postgres database
- an HTTP webserver
- a few long-running services using PG connection

Upon startup:

- hooks on signals: SIGINT and SIGTERM to initiate a shutdown flow
- creates a global context that will be cancelled in case of shutdown
- the server connects to postgres database "test"
- starts HTTP webserver
- spawns a few goroutines that represent long-running tasks
- main thread remains blocked until one of the signals caught

Once one of the signals caught:

- it cancels the global context
- all child components will finish executing because of ctx.Done()
- the main thread unblocks receiving the error (if any)
- database is closing as the last step of shutdown

Edge cases:

- if a termination signal received upon startup sequence, it cancels the global context, and then newly creating
  components should start and stop pretty much immediately.
- if a termination signal received during shutdown sequence, it will be ignored.
- because SIGKILL cannot be handled, killing the server process will result in inconsistent state.

## How to build and run

- have a postgres instance up and running on your machine
- create an empty DB "test"
- `go run .`
- try terminating the app by Ctrl+C (or by sending SIGINT or SIGTERM from terminal)
