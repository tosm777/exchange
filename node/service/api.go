package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"blockchain"
	"library/sign"
)

type BlockchainRPCAPI struct {
	service *Service
}

type TransactionRequest struct {
	SenderAddress    *string  `json:"sender_address"`
	RecipientAddress *string  `json:"recipient_address"`
	SenderPublicKey  *string  `json:"sender_public_key"`
	Value            *float32 `json:"value"`
	Signature        *string  `json:"signature"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderAddress == nil ||
		tr.RecipientAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil ||
		tr.Signature == nil {
		return false
	}
	return true
}

func (bcs *BlockchainRPCAPI) CreateTransaction(ctx context.Context, t TransactionRequest, out *TransactionRequest) error {
	publicKey := sign.String2PublicKey(*t.SenderPublicKey)
	signature := sign.String2Signature(*t.Signature)

	isCreated := bcs.service.CreateTransaction(
		*t.SenderAddress,
		*t.RecipientAddress,
		*t.Value,
		publicKey,
		signature,
	)

	if !isCreated {
		return errors.New("Failed to create Transaction")
	}

	return nil
}

func (bcs *BlockchainRPCAPI) Mine(ctx context.Context, t TransactionRequest, out *TransactionRequest) error {

	bcs.service.Mining()

	return nil
}

func (bcs *BlockchainRPCAPI) Consensus(ctx context.Context, t BlockResponse, out *TransactionRequest) error {

	var resp ResponseData
	_ = json.Unmarshal(t.Data, &resp)

	maxLength := len(bcs.service.blockchain.Chain())
	if len(resp.Block) < maxLength {
		return errors.New("Less Than This Node")
	}

	if !bcs.service.ValidChain(resp.Block) {
		return errors.New("No Valid Block")
	}

	return nil
}

type BlockResponse struct {
	Data []byte
}

type ResponseData struct {
	Block []*blockchain.Block `json:"block"`
}

func (b *ResponseData) Print() {
	for _, v := range b.Block {
		v.Print()
	}
}

func (bcs *BlockchainRPCAPI) GetChain(ctx context.Context, t TransactionRequest, out *BlockResponse) error {

	bk := ResponseData{Block: bcs.service.blockchain.Chain()}
	m, _ := json.Marshal(bk)

	out.Data = m

	return nil
}

func (bcs *BlockchainRPCAPI) SetChain(ctx context.Context, t BlockResponse, out *BlockResponse) error {

	var resp ResponseData
	_ = json.Unmarshal(t.Data, &resp)

	if len(resp.Block) > 0 {
		bcs.service.blockchain.SetChain(resp.Block)
		bcs.service.blockchain.ClearTransactionPool()
		log.Printf("Resovle confilicts replaced")
	}

	return nil
}
