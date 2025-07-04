#!/bin/bash
GETH="./go-ethereum-1.11.6/build/bin/geth"
ETH_DATA=./v5data

rm -rf "$ETH_DATA"/geth/chaindata "$ETH_DATA"/geth/lightchaindata "$ETH_DATA"/geth/transactions.rlp

$GETH --datadir=$ETH_DATA init genesis.json

$GETH --datadir=$ETH_DATA --port 11004 --authrpc.port '8554' --v5disc --http --http.addr '0.0.0.0' --http.port '18544' --http.api="db,eth,net,web3,personal,miner" --maxpeers 32 --networkid 20191003 --syncmode "full" --verbosity 5 --txpool.nolocals --allow-insecure-unlock 2>&1 |tee -a v5_output.log
