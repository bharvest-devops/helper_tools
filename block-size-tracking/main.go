package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	types "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"berachain-testnet-grpc.polkachu.com:25490",
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cmtClient := types.NewServiceClient(conn)

	height := 2450333
	for i := 0; i < 100; i++ {
		resp, err := cmtClient.GetBlockByHeight(
			ctx,
			&types.GetBlockByHeightRequest{
				Height: int64(height - (1*i)),
			},
		)
		if err != nil {
			panic(err)
		}

		fmt.Println(resp.GetBlock().Size() / 1024)
	}
}
