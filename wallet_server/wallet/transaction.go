package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"

	"library/sign"
)

type Transaction struct {
	senderPrivateKey *ecdsa.PrivateKey
	senderPublicKey  *ecdsa.PublicKey
	senderAddress    string
	recipentAddress  string
	value            float32
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender, recipient string, value float32) *Transaction {
	return &Transaction{
		senderPrivateKey: privateKey,
		senderPublicKey:  publicKey,
		senderAddress:    sender,
		recipentAddress:  recipient,
		value:            value,
	}
}

func (t *Transaction) GenerateSignature() *sign.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &sign.Signature{R: r, S: s}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_address"`
		Recipient string  `json:"recipient_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderAddress,
		Recipient: t.recipentAddress,
		Value:     t.value,
	})
}

type TransactionRequest struct {
	SenderPrivateKey *string      `json:"sender_private_key"`
	SenderPublicKey  *string      `json:"sender_public_key"`
	SenderAddress    *string      `json:"sender_address"`
	RecipientAddress *string      `json:"recipient_address"`
	Value            *json.Number `json:"value"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderAddress == nil ||
		tr.SenderPrivateKey == nil ||
		tr.RecipientAddress == nil ||
		tr.SenderPrivateKey == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil {
		return false
	}

	if *tr.SenderAddress == "" ||
		*tr.SenderPrivateKey == "" ||
		*tr.RecipientAddress == "" ||
		*tr.SenderPrivateKey == "" ||
		*tr.SenderPublicKey == "" {
		return false
	}

	return true
}
