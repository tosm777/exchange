package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

func (s *Service) CallCreateTransaction(t *TransactionRequest) error {
	peers := s.host.Peerstore().Peers()

	replies := make([]*TransactionRequest, len(peers))

	log.Printf("peers length: %d\n", len(peers))

	errs := s.rpcClient.MultiCall(
		Ctxts(len(peers)),
		peers,
		BService,
		"CreateTransaction",
		*t,
		ToInterfaces(replies),
	)

	for i, err := range errs {
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return err
		} else {
			fmt.Printf("SUCCESS: %v\n", replies[i])
		}
	}
	return nil
}

func (s *Service) CallConsensus() bool {
	peers := s.host.Peerstore().Peers()

	replies := make([]*TransactionRequest, len(peers))

	log.Printf("peers length: %d\n", len(peers))

	req := ResponseData{Block: s.blockchain.Chain()}
	reqData, _ := json.Marshal(req)

	errs := s.rpcClient.MultiCall(
		Ctxts(len(peers)),
		peers,
		BService,
		"Consensus",
		BlockResponse{reqData},
		ToInterfaces(replies),
	)

	var errCount int
	for _, err := range errs {
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			errCount++
		}
	}

	if len(peers)-2 > errCount {
		return false
	}

	s.CallResolveConflicts(BlockResponse{reqData})

	return true
}

func (s *Service) CallResolveConflicts(req BlockResponse) bool {
	log.Println("ResolveConflicts!!!")

	peers := s.host.Peerstore().Peers()
	chains := make([]*BlockResponse, len(peers))

	log.Printf("peers length: %d\n", len(peers))

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
