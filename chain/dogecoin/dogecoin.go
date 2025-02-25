// doge coin
package dogecoin

import "github.com/qiaopengjun5162/crypto-utxo-gateway/chain/base"

const ChainName = "Doge"

type ChainAdaptor struct {
	dogeClient           *base.Client
	dogeDataClientClient *base.DataClient
}
