package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

const (
	BService = "BlockchainRPCAPI"
)

func (s *Service) CallCreateTransaction(t *TransactionRequest) error {
	peers := s.host.Peerstore().Peers()

	replies := make([]*Response, len(peers))

	errs := s.rpcClient.MultiCall(
		Ctxts(len(peers)),
		peers,
		BService,
		"CreateTransaction",
		*t,
		ToInterfaces(replies),
	)

	for _, err := range errs {
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return err
		}
	}
	return nil
}

func (s *Service) CallConsensus() bool {
	peers := s.host.Peerstore().Peers()

	replies := make([]*Response, len(peers))

	req := Chain{Block: s.blockchain.Chain()}
	reqData, _ := json.Marshal(req)

	errs := s.rpcClient.MultiCall(
		Ctxts(len(peers)),
		peers,
		BService,
		"Consensus",
		Request{reqData},
		ToInterfaces(replies),
	)

	var errCount int
	for _, err := range errs {
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			errCount++
		}
	}

	if len(peers) == errCount {
		return false
	}

	s.CallResolveConflicts(Request{reqData})

	return true
}

func (s *Service) CallResolveConflicts(req Request) bool {
	log.Println("ResolveConflicts!!!")

	peers := s.host.Peerstore().Peers()
	chains := make([]*Response, len(peers))

	errs := s.rpcClient.MultiCall(
		Ctxts(len(peers)),
		peers,
		BService,
		"SetChain",
		req,
		ToInterfaces(chains),
	)

	for _, err := range errs {
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}

	return true
}

func Ctxts(n int) []context.Context {
	ctxs := make([]context.Context, n)
	for i := 0; i < n; i++ {
		ctxs[i] = context.Background()
	}
	return ctxs
}

func ToInterfaces[T any](in []*T) []interface{} {
	ifaces := make([]interface{}, len(in))
	for i := range in {
		var t T
		in[i] = &t
		ifaces[i] = in[i]
	}
	return ifaces
}
