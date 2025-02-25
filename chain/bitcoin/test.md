# Bitcoin GRPC Test

## getFee
- gRPCurl
```bash
grpcurl -plaintext -d '{
  "chain": "Bitcoin",
  "coin": "Btc",
  "network": "mainnet"
}' 127.0.0.1:8289 dapplink.utxo.WalletUtxoService.getFee
```

- request:
```bash
{
  "chain": "Bitcoin",
  "coin": "Btc",
  "network": "mainnet"
}
```

- response:
```bash
{
  "code": "SUCCESS",
  "msg": "get fee success",
  "best_fee": "0.00003101",
  "best_fee_sat": "3101",
  "slow_fee": "",
  "normal_fee": "",
  "fast_fee": ""
}
```
