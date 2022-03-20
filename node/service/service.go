package service

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"blockchain"
	"library/response"
	"library/sign"
	"wallet_server/wallet"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	rpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/rs/cors"
)

const (
	BService = "BlockchainRPCAPI"
)

type Service struct {
	rpcServer *rpc.Server
	rpcClient *rpc.Client
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
	s.rpcServer = rpc.NewServer(s.host, s.protocol)

	api := BlockchainRPCAPI{service: s}
	err := s.rpcServer.Register(&api)
	if err != nil {
		return err
	}

	s.rpcClient = rpc.NewClientWithServer(s.host, s.protocol, s.rpcServer)
	return nil
}

func (s *Service) SetBlockchain() error {
	if s.blockchain == nil {
		log.Println("NEW WALLET")
		minersWallet := wallet.NewWallet()
		s.blockchain = blockchain.NewBlockchain(minersWallet.BlockchainAddress(), uint16(s.port))

		log.Printf("%s", strings.Repeat("-", 40))
		log.Printf("private_key          %v", minersWallet.PrivateKeyStr())
		log.Printf("public_key           %v", minersWallet.PublicKeyStr())
		log.Printf("blockchain_address   %v", minersWallet.BlockchainAddress())
		log.Printf("%s", strings.Repeat("-", 40))
	}

	return nil
}

func (s *Service) Run(port int) {
	port = port + 10

	mux := http.NewServeMux()
	mux.HandleFunc("/transactions", s.Transaction)
	mux.HandleFunc("/chain", s.GetChains)
	mux.HandleFunc("/mine", s.Mine)
	mux.HandleFunc("/amount", s.Amount)
	// mux.HandleFunc("/start/amount", s.StartMine)

	handler := cors.Default().Handler(mux)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func (s *Service) CreateTransaction(sender, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, sig *sign.Signature) bool {
	isTransacted := s.blockchain.AddTransaction(sender, recipient, value, senderPublicKey, sig)
	return isTransacted
}

func (s *Service) GetChains(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:

		m, _ := s.blockchain.MarshalJSON()
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (s *Service) ValidChain(chain []*blockchain.Block) bool {
	if len(chain) == 0 {
		return false
	}
	return s.blockchain.ValidChain(chain)
}

func (s *Service) GetBlockchain() *blockchain.Blockchain {
	if s.blockchain == nil {
		s.SetBlockchain()
	}
	return s.blockchain
}

func (s *Service) Transaction(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:
		log.Println("Transactions Get now")
		w.Header().Add("Content-Type", "application/json")
		tr := s.blockchain.TransactionPool()

		m, _ := json.Marshal(struct {
			Transactions []*blockchain.Transaction `json:"transactions"`
			Length       int                       `json:"length"`
		}{
			Transactions: tr,
			Length:       len(tr),
		})

		io.WriteString(w, string(m[:]))

	case http.MethodPost:
		log.Println("Transactions Post now")

		decoder := json.NewDecoder(req.Body)
		var t TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(response.JsonMessage("fail")))
		}
		if !t.Validate() {
			log.Println("ERROR: Missing field(s)")
			io.WriteString(w, string(response.JsonMessage("fail")))
			return
		}
		w.Header().Add("Content-Type", "application/json")

		var m []byte
		if err := s.CallCreateTransaction(&t); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			m = response.JsonMessage("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			m = response.JsonMessage("success")
		}
		io.WriteString(w, string(m))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (s *Service) Mine(w http.ResponseWriter, req *http.Request) {
	log.Println("mining now")

	switch req.Method {
	case http.MethodGet:
		isMinded := s.Mining()

		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isMinded {
			w.WriteHeader(http.StatusBadRequest)
			m = response.JsonMessage("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			m = response.JsonMessage("success")
		}
		io.WriteString(w, string(m))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (s *Service) Mining() bool {
	isMined := s.blockchain.Mining()

	if isMined {
		return s.CallConsensus()
	}

	return false
}

func (s *Service) Amount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		amount := s.blockchain.CalculateTotal(blockchainAddress)

		ar := &blockchain.AmountResponse{amount}
		m, _ := ar.MarshalJSON()

		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m[:]))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}
