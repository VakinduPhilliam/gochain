package cross

import (
	"context"
	"fmt"
	"math/big"

	"github.com/gochain/gochain/v3/accounts/abi/bind"
	"github.com/gochain/gochain/v3/common"
)

type ConfirmationRequest struct {
	BlockNum  *big.Int
	LogIndex  *big.Int
	EventHash [32]byte
}

// ConfirmationsVoters returns the set of voters from confs at the given block.
func ConfirmationsVoters(ctx context.Context, block *big.Int, confs *Confirmations) (map[common.Address]struct{}, error) {
	opts := &bind.CallOpts{Context: ctx, BlockNumber: block}
	l, err := confs.VotersLength(opts)
	if err != nil {
		return nil, err
	}
	voters := make(map[common.Address]struct{})
	var i int64
	for i = 0; i < l.Int64(); i++ {
		v, err := confs.GetVoter(opts, big.NewInt(i))
		if err != nil {
			return nil, err
		}
		voters[v] = struct{}{}
	}
	return voters, nil
}

// ConfirmationsSigners returns the set of signers from confs at the given block.
func ConfirmationsSigners(ctx context.Context, block *big.Int, confs *Confirmations) (map[common.Address]struct{}, error) {
	opts := &bind.CallOpts{Context: ctx, BlockNumber: block}
	l, err := confs.SignersLength(opts)
	if err != nil {
		return nil, err
	}
	signers := make(map[common.Address]struct{})
	var i int64
	for i = 0; i < l.Int64(); i++ {
		s, err := confs.GetSigner(opts, big.NewInt(i))
		if err != nil {
			return nil, err
		}
		signers[s] = struct{}{}
	}
	return signers, nil
}

// difference returns a slice of addresses which are in a and not in b.
func difference(a, b map[common.Address]struct{}) []common.Address {
	var d []common.Address
	for x := range a {
		if _, ok := b[x]; !ok {
			d = append(d, x)
		}
	}
	return d
}

// pendingRequests returns the pending list of confirmation requests.
func pendingRequests(confs *Confirmations, num *big.Int) ([]ConfirmationRequest, error) {
	confirmedOpts := &bind.CallOpts{BlockNumber: num}
	pll, err := confs.PendingListLength(confirmedOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending list length: %v", err)
	}
	max := pll.Uint64()
	reqs := make([]ConfirmationRequest, max)
	var i uint64
	var arg big.Int
	for ; i < max; i++ {
		_ = arg.SetUint64(i)
		raw, err := confs.PendingList(confirmedOpts, &arg)
		if err != nil {
			return nil, fmt.Errorf("failed to get pending request %d: %v", i, err)
		}
		reqs[i] = raw
	}
	return reqs, nil
}
