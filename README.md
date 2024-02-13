# txstreet-bulk-eth-api

> A rewrite of
> [txstreet/tmp-bulk-geth-api](https://github.com/txstreet/tmp-bulk-geth-api)

An HTTP API for making RPC calls in bulk to an Ethereum execution client.
Limited to those needed by
[txstreet/processor](https://github.com/txstreet/processor).

Compared to the original version:

* Configurable inbound/outbound concurrency and timeout
  * Via options or by adjustments in the source code
* Ingestable logs ([logfmt][logfmt] format)

[logfmt]: https://brandur.org/logfmt

Usage:

```
Usage of bulk-eth-api:
  -addr string
    	Listen interface and port (default ":8081")
  -concurrency int
    	Max concurrency of RPC calls (default 8)
  -eth string
    	RPC URL of Ethereum execution client (default "http://localhost:8545")
  -log value
    	Log level (DEBUG, INFO, WARN, ERROR) (default DEBUG)
```

Example logs:

```
time=2024-02-13T11:42:34.806Z level=INFO msg=Started requestId=2cJL5tiZjoT1fTrvg7wN0OrzIYF method=POST url=/nonces
time=2024-02-13T11:42:34.807Z level=INFO msg=Started requestId=2cJL5slXQVwwaNhCmQCfhTuh3vW method=POST url=/nonces
time=2024-02-13T11:42:34.807Z level=INFO msg=Started requestId=2cJL5t9ua4SRUzbnIkj6Edrxtvz method=POST url=/nonces
time=2024-02-13T11:42:34.808Z level=INFO msg=Finished requestId=2cJL5tiZjoT1fTrvg7wN0OrzIYF method=POST url=/nonces duration=1.239781ms
time=2024-02-13T11:42:34.808Z level=INFO msg=Finished requestId=2cJL5slXQVwwaNhCmQCfhTuh3vW method=POST url=/nonces duration=751.517µs
time=2024-02-13T11:42:34.808Z level=INFO msg=Finished requestId=2cJL5t9ua4SRUzbnIkj6Edrxtvz method=POST url=/nonces duration=901.097µs
time=2024-02-13T11:42:34.808Z level=INFO msg=Started requestId=2cJL5rljlr7DyIItDtsrQwZdFzL method=POST url=/nonces
time=2024-02-13T11:42:34.809Z level=INFO msg=Finished requestId=2cJL5rljlr7DyIItDtsrQwZdFzL method=POST url=/nonces duration=595.585µs
time=2024-02-13T11:42:34.915Z level=INFO msg=Started requestId=2cJL5tiZ5aoiynrz1ZxLu4svthv method=POST url=/contract-codes
time=2024-02-13T11:42:34.965Z level=INFO msg=Finished requestId=2cJL5tiZ5aoiynrz1ZxLu4svthv method=POST url=/contract-codes duration=50.321516ms
time=2024-02-13T11:42:34.970Z level=INFO msg=Started requestId=2cJL5vquzfFSrWPQHAIVaTDMtRz method=POST url=/transaction-receipts
time=2024-02-13T11:42:34.994Z level=INFO msg=Finished requestId=2cJL5vquzfFSrWPQHAIVaTDMtRz method=POST url=/transaction-receipts duration=24.007902ms
```
