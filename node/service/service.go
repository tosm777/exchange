package service

import (
	"log"
	"strings"

	"blockchain"
	"wallet"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
)

type Service struct {
	rpcServer *gorpc.Server
	rpcClient *gorpc.Client
	host      host.Host
	protocol  protocol.ID
	counter   int

	blockchain *blockchain.Blockchain
	port       int
}

func NewService(host host.Host, protocol protocol.ID, port int) *Service {
	return &Service{
		host:       host,
		protocol:   protocol,
		blockchain: nil,
		port:       port,
	}
}

func (s *Service) SetupRPC() error {
	s.rpcServer = gorpc.NewServer(s.host, s.protocol)

	api := BlockchainRPCAPI{service: s}
	err := s.rpcServer.Register(&api)
	if err != nil {
		return err
	}

	s.rpcClient = gorpc.NewClientWithServer(s.host, s.protocol, s.rpcServer)
	return nil
}

func (s *Service) SetBlockchain() error {
	if s.blockchain == nil {
		minersWallet := wallet.NewWallet()
		s.blockchain = blockchain.NewBlockchain(minersWallet.BlockchainAddress(), uint16(s.port))

		log.Printf("%s", strings.Repeat("-", 40))
		log.Println("NEW WALLET")
		log.Printf("private_key          %v", minersWallet.PrivateKeyStr())
		log.Printf("public_key           %v", minersWallet.PublicKeyStr())
		log.Printf("blockchain_address   %v", minersWallet.BlockchainAddress())
		log.Printf("%s", strings.Repeat("-", 40))
	}

	return nil
}
