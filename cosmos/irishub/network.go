package irishub

import (
	"context"
	"strconv"
	"time"

	"github.com/tendermint/cosmos-rosetta-gateway/cosmos/irishub/client/tendermint"

	"github.com/coinbase/rosetta-sdk-go/types"
	"golang.org/x/sync/errgroup"
)

func (l launchpad) NetworkList(context.Context, *types.MetadataRequest) (*types.NetworkListResponse, *types.Error) {
	if l.properties.OfflineMode {
		return nil, ErrEndpointDisabledOfflineMode
	}

	return &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{
			{
				Blockchain: l.properties.Blockchain,
				Network:    l.properties.Network,
			},
		},
	}, nil
}

func (l launchpad) NetworkOptions(ctx context.Context, _ *types.NetworkRequest) (*types.NetworkOptionsResponse, *types.Error) {
	if l.properties.OfflineMode {
		return nil, ErrEndpointDisabledOfflineMode
	}

	resp, err := l.tendermint.Status()
	if err != nil {
		return nil, ErrNodeConnection
	}

	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: "1.2.5",
			NodeVersion:    resp.NodeInfo.Version,
		},
		Allow: &types.Allow{
			OperationStatuses: []*types.OperationStatus{
				{
					Status:     StatusSuccess,
					Successful: true,
				},
				{
					Status:     StatusReverted,
					Successful: false,
				},
			},
			OperationTypes: []string{OperationTransfer},
		},
	}, nil
}

func (l launchpad) NetworkStatus(ctx context.Context, _ *types.NetworkRequest) (*types.NetworkStatusResponse, *types.Error) {
	if l.properties.OfflineMode {
		return nil, ErrEndpointDisabledOfflineMode
	}

	var (
		latestBlock  tendermint.BlockResponse
		genesisBlock tendermint.BlockResponse
		netInfo      tendermint.NetInfoResponse
	)

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		latestBlock, err = l.tendermint.Block(0)
		return
	})
	g.Go(func() (err error) {
		genesisBlock, err = l.tendermint.Block(1)
		return
	})
	g.Go(func() (err error) {
		netInfo, err = l.tendermint.NetInfo()
		return
	})
	if err := g.Wait(); err != nil {
		return nil, ErrNodeConnection
	}

	var peers []*types.Peer
	for _, p := range netInfo.Peers {
		peers = append(peers, &types.Peer{
			PeerID: p.NodeInfo.Id,
		})
	}

	height, err := strconv.ParseUint(latestBlock.Block.Header.Height, 10, 64)
	if err != nil {
		return nil, ErrInterpreting
	}

	t, err := time.Parse(time.RFC3339Nano, latestBlock.Block.Header.Time)
	if err != nil {
		return nil, ErrInterpreting
	}

	return &types.NetworkStatusResponse{
		CurrentBlockIdentifier: &types.BlockIdentifier{
			Index: int64(height),
			Hash:  latestBlock.BlockMeta.BlockId.Hash,
		},
		CurrentBlockTimestamp: t.UnixNano() / 1000000,
		GenesisBlockIdentifier: &types.BlockIdentifier{
			Index: 1,
			Hash:  genesisBlock.BlockMeta.BlockId.Hash,
		},
		Peers: peers,
	}, nil
}