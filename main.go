package main

import (
	"flag"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ethereum/go-ethereum/log"

	"github.com/qiaopengjun5162/crypto-utxo-gateway/chaindispatcher"
	"github.com/qiaopengjun5162/crypto-utxo-gateway/config"
	"github.com/qiaopengjun5162/crypto-utxo-gateway/rpc/utxo"
)

// main starts a grpc server which provides wallet utxo service.
//
// The server will listen on port specified in config file.
// The default config file path is "config.yml", but you can specify
// another path by using "-c" flag. For example:
//
//	./crypto-utxo-gateway -c path/to/your/config.yml
//
// The server will panic if any error occurs during setup process.
func main() {
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	dispatcher, err := chaindispatcher.New(conf)
	if err != nil {
		log.Error("Setup dispatcher failed", "err", err)
		panic(err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(dispatcher.Interceptor))
	defer grpcServer.GracefulStop()

	utxo.RegisterWalletUtxoServiceServer(grpcServer, dispatcher)

	listen, err := net.Listen("tcp", ":"+conf.Server.Port)
	if err != nil {
		log.Error("net listen failed", "err", err)
		panic(err)
	}
	reflection.Register(grpcServer)

	log.Info("crypto-utxo-gateway wallet rpc services start success", "port", conf.Server.Port)

	if err := grpcServer.Serve(listen); err != nil {
		log.Error("grpc server serve failed", "err", err)
		panic(err)
	}
}
