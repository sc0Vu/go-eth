# subgraph

The golang client implementation to query the popular subgraph on Ethereum, eg Uniswap.

# what is subgraph?

[The Graph](https://thegraph.com/) is a cool indexing protocol for querying blockchain networks like Ethereum.

Subgraph is config that define data to collect and save in the graph protoocl.

# test local
If you already install make tools, simply execute this command to test:
```BASH
$ make
```

If you don't install make tools, type this command in your terminal:
```BASH
$ go test ./...
```

# examples

* tokenlon
Fetch and print transaction history from Tokenlon.
```BASH
$ go run examples/tokenlon/*.go
```

* uniswap v2
Fetch and print ethereum price and pair data from UniswapV2.
```BASH
$ go run examples/uniswapv2/*.go
```

# license
MIT
