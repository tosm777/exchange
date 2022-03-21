package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"blockchain"
	"library/response"
	"library/sign"
	"wallet"

	"github.com/rs/cors"
)

const tempDir = "templates/view"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port: port, gateway: gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		m, _ := myWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		var t wallet.TransactionRequest

		err := json.NewDecoder(req.Body).Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(response.JsonMessage("fail")))
			return
		}
		if !t.Validate() {
			log.Println("ERROR: Missing field(s)")
			io.WriteString(w, string(response.JsonMessage("fail")))
			return
		}

		publicKey := sign.String2PublicKey(*t.SenderPublicKey)
		privateKey := sign.String2PrivateKey(*t.SenderPrivateKey, publicKey)
		value64, err := t.Value.Float64()
		if err != nil {
			log.Println("ERROR: parse error")
			io.WriteString(w, string(response.JsonMessage("fail")))
			return
		}
		value32 := float32(value64)

		w.Header().Add("Content-Type", "application/json")

		transaction := wallet.NewTransaction(privateKey, publicKey, *t.SenderAddress, *t.RecipientAddress, value32)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		bt := &blockchain.TransactionRequest{
			SenderAddress:    t.SenderAddress,
			RecipientAddress: t.RecipientAddress,
			SenderPublicKey:  t.SenderPublicKey,
			Value:            &value32,
			Signature:        &signatureStr,
		}

		log.Println("posting....")
		m, _ := json.Marshal(bt)
		resp, err := http.Post("http://node:7000/transactions", "application/json", bytes.NewBuffer(m))

		if resp.StatusCode == 201 {
			io.WriteString(w, string(response.JsonMessage("success")))
			return
		}
		io.WriteString(w, string(response.JsonMessage("fail")))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) WalletAmount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		endpoint := "http://node:7000/amount"

		client := http.DefaultClient
		bcsReq, _ := http.NewRequest("GET", endpoint, nil)

		q := bcsReq.URL.Query()
		q.Add("blockchain_address", blockchainAddress)
		bcsReq.URL.RawQuery = q.Encode()

		bcsResp, err := client.Do(bcsReq)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(response.JsonMessage("fail")))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if bcsResp.StatusCode == 200 {
			decoder := json.NewDecoder(bcsResp.Body)
			var bar blockchain.AmountResponse

			err := decoder.Decode(&bar)
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, string(response.JsonMessage("fail")))
				return
			}

			m, _ := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "success",
				Amount:  bar.Amount,
			})
			io.WriteString(w, string(m[:]))
		} else {
			io.WriteString(w, string(response.JsonMessage("fail")))
		}
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (ws *WalletServer) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/wallet", ws.Wallet)
	mux.HandleFunc("/wallet/amount", ws.WalletAmount)
	mux.HandleFunc("/transaction", ws.CreateTransaction)

	handler := cors.Default().Handler(mux)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), handler))
}
