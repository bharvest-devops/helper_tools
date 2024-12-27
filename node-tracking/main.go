package main

import (
	"context"
	"fmt"
	"time"

	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/types"
)

func main() {
	client, err := rpchttp.New("${RPC_ADDRESS}", "/websocket")
	if err != nil {
		panic(err)
	}

	err = client.Start()
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Check Block Producing
	go func() {
		query := "tm.event = 'CompleteProposal'"
		txs, err := client.Subscribe(ctx, "client", query)
		if err != nil {
			panic(err)
		}
		for e := range txs {
			// {Height, Round, Step, BlockID:{Total(Block내의 tx개수로 추정됨), Hash(tx hash)}}
			fmt.Println("Check Block Producing:", e.Data.(types.EventDataCompleteProposal))
		}
	}()

	// Check Voting
	go func() {
		query := "tm.event = 'Vote'"
		txs, err := client.WSEvents.Subscribe(ctx, "client", query)
		if err != nil {
			panic(err)
		}
		for e := range txs {
			d := e.Data.(types.EventDataVote)
			// fmt.Println(d.Vote.Height, " / ", d.Vote.Round, " / ", d.Vote.Type, "\t", d.Vote.ValidatorAddress)
			if d.Vote.ValidatorAddress.String() == "${VALIDATOR_ADDRESS}" {
				fmt.Printf("%+v\n", d.Vote)
			}
		}
	}()

	// Check Block Proposer
	query := "tm.event = 'NewRound'"
	txs, err := client.Subscribe(ctx, "client", query)
	if err != nil {
		panic(err)
	}
	for e := range txs {
		// {Height, Round, Step, Proposer:{Address(Proposer validator address), Index(Proposer validator순위)}}
		fmt.Println("Check Block Proposer:", e.Data.(types.EventDataNewRound))
	}

	// Gov Submit
	go func() {
		query := "tm.event = 'Tx' AND message.action = '/cosmos.gov.v1beta1.MsgSubmitProposal'"
		txs, err := client.Subscribe(ctx, "client", query)
		if err != nil {
			panic(err)
		}
		for e := range txs {
			if e.Events["submit_proposal.proposal_type"][0] == "ConsumerAddition" {
				fmt.Println("🛑 Proposal Notice 🛑\n\n\"Consumer Addition Proposal\" is submitted!\nCheck GoC Provider chain.")
			} else if e.Events["submit_proposal.proposal_type"][0] == "ConsumerRemoval" {
				fmt.Println("🛑 Proposal Notice 🛑\n\n\"Consumer Removal Proposal\" is submitted!\nCheck GoC Provider chain.")
			}
		}
	}()
}
