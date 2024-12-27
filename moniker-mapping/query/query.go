package query

import (
	"context"
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	//"github.com/initia-labs/initia/x/mstaking/types"
)

func Validators(ctx context.Context, target string, secureConnection bool) ([]types.Validator, error) {
	options := []grpc.DialOption{grpc.WithBlock()}
	if secureConnection {
		options = append(
			options,
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		)
	} else {
		options = append(
			options,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	}

	conn, err := grpc.DialContext(
		ctx,
		target,
		options...,
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := types.NewQueryClient(conn)

	page := &query.PageRequest{}
	validators := []types.Validator{}
	for {
		resp, err := client.Validators(
			ctx,
			&types.QueryValidatorsRequest{
				Pagination: page,
			},
		)
		if err != nil {
			return nil, err
		}
		validators = append(validators, resp.Validators...)

		if len(resp.Pagination.NextKey) == 0 {
			break
		}
		page.Key = resp.Pagination.NextKey
	}

	return validators, nil
}
