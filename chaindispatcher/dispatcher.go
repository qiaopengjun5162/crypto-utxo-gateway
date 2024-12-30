package chaindispatcher

import (
	"context"

	"github.com/qiaopengjun5162/crypto-utxo-gateway/chain/bitcoincash"

	"runtime/debug"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ethereum/go-ethereum/log"

	"github.com/qiaopengjun5162/crypto-utxo-gateway/chain"
	"github.com/qiaopengjun5162/crypto-utxo-gateway/chain/bitcoin"
	"github.com/qiaopengjun5162/crypto-utxo-gateway/config"
	"github.com/qiaopengjun5162/crypto-utxo-gateway/rpc/common"
	"github.com/qiaopengjun5162/crypto-utxo-gateway/rpc/utxo"
)

type CommonRequest interface {
	GetChain() string
}

type CommonReply = utxo.SupportChainsResponse

type ChainType = string

type ChainDispatcher struct {
	registry map[ChainType]chain.IChainAdaptor
}

func New(conf *config.Config) (*ChainDispatcher, error) {
	dispatcher := ChainDispatcher{
		registry: make(map[ChainType]chain.IChainAdaptor),
	}
	chainAdaptorFactoryMap := map[string]func(conf *config.Config) (chain.IChainAdaptor, error){
		bitcoin.ChainName:     bitcoin.NewChainAdaptor,
		bitcoincash.ChainName: bitcoincash.NewChainAdaptor,
		//dash.ChainName:        dash.NewChainAdaptor,
		//litecoin.ChainName:    litecoin.NewChainAdaptor,
	}
	supportedChains := []string{
		bitcoin.ChainName,
		bitcoincash.ChainName,
		//dash.ChainName,
		//litecoin.ChainName,
	}
	for _, c := range conf.Chains {
		if factory, ok := chainAdaptorFactoryMap[c]; ok {
			adaptor, err := factory(conf)
			if err != nil {
				log.Crit("failed to setup chain", "chain", c, "error", err)
			}
			dispatcher.registry[c] = adaptor
		} else {
			log.Error("unsupported chain", "chain", c, "supportedChains", supportedChains)
		}
	}
	return &dispatcher, nil
}

func (d *ChainDispatcher) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic error", "msg", e)
			log.Debug(string(debug.Stack()))
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	pos := strings.LastIndex(info.FullMethod, "/")
	method := info.FullMethod[pos+1:]

	chainName := req.(CommonRequest).GetChain()
	log.Info(method, "chain", chainName, "req", req)

	resp, err = handler(ctx, req)
	log.Debug("Finish handling", "resp", resp, "err", err)
	return
}

func (d *ChainDispatcher) preHandler(req interface{}) (resp *CommonReply) {
	chainName := req.(CommonRequest).GetChain()
	if _, ok := d.registry[chainName]; !ok {
		return &CommonReply{
			Code:    common.ReturnCode_ERROR,
			Msg:     config.UnsupportedOperation,
			Support: false,
		}
	}
	return nil
}

// GetSupportChains retrieves the list of supported chains by the gateway.
// It first calls the preHandler to validate the request and then calls the
// corresponding method on the registered chain adaptor.
//
// ctx: The context for the RPC call.
// request: The SupportChainsRequest containing the chain name.
//
// Returns:
// - A SupportChainsResponse containing the list of supported chains and any error encountered.
// - An error if the preHandler or the chain adaptor method returns an error.
func (d *ChainDispatcher) GetSupportChains(ctx context.Context, request *utxo.SupportChainsRequest) (*utxo.SupportChainsResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.SupportChainsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
		}, nil
	}
	return d.registry[request.Chain].GetSupportChains(request)
}

// ConvertAddress implements the IChainAdaptor interface.
// It converts a public key to a wallet address according to the specified format.
// It first calls the preHandler to validate the request and then calls the
// corresponding method on the registered chain adaptor.
//
// ctx: The context for the RPC call.
// request: The ConvertAddressRequest containing the public key, format and chain.
//
// Returns:
// - A ConvertAddressResponse containing the wallet address and any error encountered.
// - An error if the preHandler or the chain adaptor method returns an error.
func (d *ChainDispatcher) ConvertAddress(ctx context.Context, request *utxo.ConvertAddressRequest) (*utxo.ConvertAddressResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.ConvertAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "covert address fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].ConvertAddress(request)
}

func (d *ChainDispatcher) ValidAddress(ctx context.Context, request *utxo.ValidAddressRequest) (*utxo.ValidAddressResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.ValidAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "valid address error at pre handle",
		}, nil
	}
	return d.registry[request.Chain].ValidAddress(request)
}

func (d *ChainDispatcher) GetFee(ctx context.Context, request *utxo.FeeRequest) (*utxo.FeeResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.FeeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get fee fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetFee(request)
}

func (d *ChainDispatcher) GetAccount(ctx context.Context, request *utxo.AccountRequest) (*utxo.AccountResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.AccountResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get account information fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetAccount(request)
}

func (d *ChainDispatcher) GetUnspentOutputs(ctx context.Context, request *utxo.UnspentOutputsRequest) (*utxo.UnspentOutputsResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.UnspentOutputsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get un spend out fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetUnspentOutputs(request)
}

// GetBlockByNumber retrieves a block from the blockchain by block number.
// It first calls the preHandler to validate the request and then calls the
// corresponding method on the registered chain adaptor.
//
// ctx: The context for the RPC call.
// request: The BlockNumberRequest containing the block number.
//
// Returns:
// - A BlockResponse containing the block and any error encountered.
// - An error if the preHandler or the chain adaptor method returns an error.
func (d *ChainDispatcher) GetBlockByNumber(ctx context.Context, request *utxo.BlockNumberRequest) (*utxo.BlockResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block by number fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetBlockByNumber(request)
}

func (d *ChainDispatcher) GetBlockByHash(ctx context.Context, request *utxo.BlockHashRequest) (*utxo.BlockResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block by hash fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetBlockByHash(request)
}

func (d *ChainDispatcher) GetBlockHeaderByHash(ctx context.Context, request *utxo.BlockHeaderHashRequest) (*utxo.BlockHeaderResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block header by hash fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetBlockHeaderByHash(request)
}

func (d *ChainDispatcher) GetBlockHeaderByNumber(ctx context.Context, request *utxo.BlockHeaderNumberRequest) (*utxo.BlockHeaderResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block header by number fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetBlockHeaderByNumber(request)
}

// SendTx sends a transaction request to the appropriate chain adaptor.
// It first calls the preHandler to validate the request and then calls the
// corresponding SendTx method on the registered chain adaptor.
//
// ctx: The context for the RPC call.
// request: The SendTxRequest containing the transaction details.
//
// Returns:
// - A SendTxResponse containing the result of the transaction and any error encountered.
// - An error if the preHandler or the chain adaptor method returns an error.
func (d *ChainDispatcher) SendTx(ctx context.Context, request *utxo.SendTxRequest) (*utxo.SendTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "send tx fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].SendTx(request)
}

func (d *ChainDispatcher) GetTxByAddress(ctx context.Context, request *utxo.TxAddressRequest) (*utxo.TxAddressResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx by address fail pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetTxByAddress(request)
}

func (d *ChainDispatcher) GetTxByHash(ctx context.Context, request *utxo.TxHashRequest) (*utxo.TxHashResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx by hash fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetTxByHash(request)
}

func (d *ChainDispatcher) CreateUnSignTransaction(ctx context.Context, request *utxo.UnSignTransactionRequest) (*utxo.UnSignTransactionResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.UnSignTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get un sign tx fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].CreateUnSignTransaction(request)
}

// BuildSignedTransaction builds a signed transaction.
//
// It first calls the preHandler to validate the request and then calls the
// corresponding BuildSignedTransaction method on the registered chain adaptor.
//
// ctx: The context for the RPC call.
// request: The SignedTransactionRequest containing the transaction details to be signed.
//
// Returns:
// - A SignedTransactionResponse containing the signed transaction and any error encountered.
// - An error if the preHandler or the chain adaptor method returns an error.
func (d *ChainDispatcher) BuildSignedTransaction(ctx context.Context, request *utxo.SignedTransactionRequest) (*utxo.SignedTransactionResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "signed tx fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].BuildSignedTransaction(request)
}

// DecodeTransaction decodes a transaction.
//
// It first calls the preHandler to validate the request and then calls the
// corresponding DecodeTransaction method on the registered chain adaptor.
//
// ctx: The context for the RPC call.
// request: The DecodeTransactionRequest containing the transaction details.
//
// Returns:
// - A DecodeTransactionResponse containing the decoded transaction and any error encountered.
// - An error if the preHandler or the chain adaptor method returns an error.
func (d *ChainDispatcher) DecodeTransaction(ctx context.Context, request *utxo.DecodeTransactionRequest) (*utxo.DecodeTransactionResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.DecodeTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "decode tx fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].DecodeTransaction(request)
}

// VerifySignedTransaction verifies a signed transaction.
//
// It first calls the preHandler to validate the request and then calls the
// corresponding VerifySignedTransaction method on the registered chain adaptor.
//
// ctx: The context for the RPC call.
// request: The VerifyTransactionRequest containing the signed transaction details.
//
// Returns:
// - A VerifyTransactionResponse containing the result of the verification and any error encountered.
// - An error if the preHandler or the chain adaptor method returns an error.
func (d *ChainDispatcher) VerifySignedTransaction(ctx context.Context, request *utxo.VerifyTransactionRequest) (*utxo.VerifyTransactionResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &utxo.VerifyTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "verify tx fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].VerifySignedTransaction(request)
}
