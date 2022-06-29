# erigon-bug
Showcase for a bug with erigon log subscriptions

```
make build
bin/main [-erc20] [-usdt]
```

This program runs an ethereum log subscription with options to filter by logs from the usdt address (using the `-usdt` flag) and/or to filter by logs with the erc20 transfer topic 0 hash (using the `-erc20` flag). When running this program multiple times concurrently with different flags (resulting in different filter queries), certain combinations result in not all logs being received as expected:
- Running one with no flags and one with the `-erc20` flag, neither receive any logs
- Running one with no flags and one with the `-usdt` flag, neither receive any logs
- Running one with no flags, one with the `-erc20` flag, and one with both flags, the former two only receive logs from the usdt address
- Running one with no flags, one with the `-usdt` flag, and one with both flags, the former two only receive logs with the erc20 transfer topic 0 hash
- Running one with no flags, one with the `-usdt` flag, one with the `-erc20` flag, and one with both flags, all four only receive logs with the erc20 transfer topic 0 hash from the usdt address
- Running one with the `-usdt` flag, one with the `-erc20` flag, and one with both flags, all behave as expected
- Running one with no flags, and one with both flags, both behave as expected

This is running erigon on the devel branch (commit `aa7985341e194cf33ad51afaf09252b1558f2599`) using
```
~/erigon/build/bin/erigon --datadir=/erigon --chain=mainnet --port=30303 --http.port=8545 --private.api.addr=127.0.0.1:9091 --torrent.port=42069 --http --metrics --ws --http.addr=0.0.0.0 --http.api=eth,debug,net,trace,web3,erigon
```

Note that these issues do not occur when running with a fiews node