package dash

import "github.com/qiaopengjun5162/crypto-utxo-gateway/chain/base"

const ChainName = "Dash"

type ChainAdaptor struct {
	dashClient     *base.Client
	dashDataClient *base.DataClient
}
