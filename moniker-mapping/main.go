package main

import (
	"context"
	"encoding/json"
	"os"

	"moniker-mapping/query"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
)

func main() {
	ctx := context.Background()

	validators, err := query.Validators(
		ctx,
		"injective-grpc.polkachu.com:14390",
		false,
	)
	if err != nil {
		panic(err)
	}

	output := make(map[string]string)

	for _, validator := range validators {
		pk := ed25519.PubKey{}
		err := pk.Unmarshal(validator.ConsensusPubkey.Value)
		if err != nil {
			panic(err)
		}

		output[validator.Description.Moniker] = pk.Address().String()
	}

	outputMarshaled, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("validators.json", outputMarshaled, 0644)
	if err != nil {
		panic(err)
	}
}
