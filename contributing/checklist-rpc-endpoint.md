# New RPC Endpoint Checklist

Prefer adding a new message to changing any existing RPC messages.

## Code

* [ ] `Request` struct and `RequestType` constant in `nomad/structs/structs.go`
* [ ] In `nomad/fsm.go`, add a dispatch case to the switch statement in `Apply`
  * `*nomadFSM` method to decode the request and call the state method
* [ ] State method for modifying objects in a `Txn` in `nomad/state/state_store.go`
  * `nomad/state/state_store_test.go`
* [ ] Handler for the request in `nomad/foo_endpoint.go`
  * RPCs are resolved by matching the method name for bound structs
	[net/rpc](https://golang.org/pkg/net/rpc/)
  * Wrapper for the HTTP request in `command/agent/foo_endpoint.go`
* [ ] `nomad/core_sched.go` sends many RPCs
  * `ServersMeetMinimumVersion` asserts that the server cluster is
    upgraded, so use this to gaurd sending the new RPC, else send the old RPC
  * Version must match the actual release version!

## Docs

* [ ] Changelog
