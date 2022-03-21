package service

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"blockchain"
	"library/response"
	"library/sign"

	"github.com/rs/cors"
)

const (
	MINING_TIMER_SEC = 10
)

func (s *Service) Run(port int) {
	port = port + 10

	s.StartMining()

	mux := http.NewServeMux()
	mux.HandleFunc("/transactions", s.Transaction)
	mux.HandleFunc("/chain", s.GetChains)
	mux.HandleFunc("/amount", s.Amount)

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

func (s *Service) StartMining() {
	s.Mining()
	_ = time.AfterFunc(time.Second*MINING_TIMER_SEC, s.StartMining)
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
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}
