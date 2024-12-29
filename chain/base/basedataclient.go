package base

import (
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/chain-explorer-api/common/chain"
	"github.com/dapplink-labs/chain-explorer-api/common/gas_fee"
	"github.com/dapplink-labs/chain-explorer-api/common/transaction"
	"github.com/dapplink-labs/chain-explorer-api/explorer/oklink"
)

type DataClient struct {
	ChainShortName string
	ExplorerName   string
	BaseDataCli    *oklink.ChainExplorerAdaptor
}

func NewBaseDataClient(baseUrl, apiKey, chainShortName, explorerName string) (*DataClient, error) {
	ltcCli, err := oklink.NewChainExplorerAdaptor(apiKey, baseUrl, false, time.Second*15)
	if err != nil {
		log.Error("New base client fail", "err", err)
		return nil, err
	}
	return &DataClient{
		ChainShortName: chainShortName,
		ExplorerName:   explorerName,
		BaseDataCli:    ltcCli,
	}, nil
}

func (bdc *DataClient) GetFee() (*gas_fee.GasEstimateFeeResponse, error) {
	gasFeeRequest := &gas_fee.GasEstimateFeeRequest{
		ChainShortName: bdc.ChainShortName,
		ExplorerName:   oklink.ChainExplorerName,
	}
	okResp, err := bdc.BaseDataCli.GetEstimateGasFee(gasFeeRequest)
	if err != nil {
		log.Error("get estimate gas fee fail", "err", err)
		return nil, err
	}
	return okResp, nil
}

func (bdc *DataClient) GetAccountBalance(address string) (*account.AccountBalanceResponse, error) {
	accountItem := []string{address}
	contractAddress := []string{"0x00"}
	balanceRequest := &account.AccountBalanceRequest{
		ChainShortName:  bdc.ChainShortName,
		ExplorerName:    oklink.ChainExplorerName,
		Account:         accountItem,
		ContractAddress: contractAddress,
	}
	balanceResponse, err := bdc.BaseDataCli.GetAccountBalance(balanceRequest)
	if err != nil {
		log.Error("get balance response fail", "err", err)
		return nil, err
	}
	return balanceResponse, nil
}

func (bdc *DataClient) GetAccountUtxoList(address string) ([]account.AccountUtxoResponse, error) {
	utxoRequest := &account.AccountUtxoRequest{
		ChainShortName: bdc.ChainShortName,
		ExplorerName:   oklink.ChainExplorerName,
		Address:        address,
		Page:           "",
		Limit:          "",
	}
	utxoResponse, err := bdc.BaseDataCli.GetAccountUtxo(utxoRequest)
	if err != nil {
		log.Error("get account utxo fail", "err", err)
		return nil, err
	}
	return utxoResponse, nil
}

func (bdc *DataClient) GetTxListByAddress(address string, page, pageSize uint64) (*account.TransactionResponse[account.AccountTxResponse], error) {
	txRequest := &account.AccountTxRequest{
		ChainShortName: bdc.ChainShortName,
		ExplorerName:   oklink.ChainExplorerName,
		Action:         account.OkLinkActionUtxo,
		Address:        address,
		PageRequest: chain.PageRequest{
			Page:  page,
			Limit: pageSize,
		},
	}
	txListResponse, err := bdc.BaseDataCli.GetTxByAddress(txRequest)
	if err != nil {
		log.Error("get tx by address fail", "err", err)
		return nil, err
	}
	log.Info("tx list response success", "transactionList Length", len(txListResponse.TransactionList))
	return txListResponse, nil
}

func (bdc *DataClient) GetTxByHash(txId string) (*transaction.TxResponse, error) {
	txRequest := &transaction.TxRequest{
		ChainShortName: bdc.ChainShortName,
		ExplorerName:   oklink.ChainExplorerName,
		Txid:           txId,
	}
	txResponse, err := bdc.BaseDataCli.GetTxByHash(txRequest)
	if err != nil {
		log.Error("get tx by address fail", "err", err)
		return nil, err
	}
	return txResponse, nil
}
