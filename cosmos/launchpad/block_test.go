package launchpad

import (
	"context"
	"testing"
	"time"

	"github.com/antihax/optional"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	cosmosclient "github.com/tendermint/cosmos-rosetta-gateway/cosmos/launchpad/client/sdk/generated"
	cosmosmocks "github.com/tendermint/cosmos-rosetta-gateway/cosmos/launchpad/client/sdk/mocks"
	tendermintclient "github.com/tendermint/cosmos-rosetta-gateway/cosmos/launchpad/client/tendermint/generated"
	tendermintmocks "github.com/tendermint/cosmos-rosetta-gateway/cosmos/launchpad/client/tendermint/mocks"
	"github.com/tendermint/cosmos-rosetta-gateway/rosetta"
)

func TestLaunchpad_Block(t *testing.T) {
	var (
		mt = &tendermintmocks.TendermintInfoAPI{}
		mc = &cosmosmocks.CosmosTransactionsAPI{}
	)
	defer mt.AssertExpectations(t)
	defer mc.AssertExpectations(t)

	ti, err := time.Parse(time.RFC3339, "2019-04-22T17:01:51Z")
	require.NoError(t, err)

	mt.
		On("Block", mock.Anything, &tendermintclient.BlockOpts{
			Height: optional.NewFloat32(float32(1)),
		}).
		Return(tendermintclient.BlockResponse{
			Result: tendermintclient.BlockComplete{
				BlockId: tendermintclient.BlockId{
					Hash: "11",
				},
				Block: tendermintclient.Block{
					Header: tendermintclient.BlockHeader{
						Height: "2",
						Time:   ti.Format(time.RFC3339),
						LastBlockId: tendermintclient.BlockId{
							Hash: "12",
						},
					},
				},
			},
		}, nil, nil).
		Once()

	var opts *tendermintclient.TxSearchOpts
	mt.
		On("TxSearch", mock.Anything, `"tx.height=2"`, opts).
		Return(tendermintclient.TxSearchResponse{
			Result: tendermintclient.TxSearchResponseResult{
				Txs: []tendermintclient.TxSearchResponseResultTxs{
					{
						Hash: "3",
					},
					{
						Hash: "4",
					},
				},
			},
		}, nil, nil).
		Once()

	mc.
		On("TxsHashGet", mock.Anything, "3").
		Return(cosmosclient.TxQuery{
			Txhash: "3",
			Tx: cosmosclient.StdTx{
				Value: cosmosclient.StdTxValue{
					Msg: []cosmosclient.Msg{
						{
							Type: "5",
							Value: cosmosclient.MsgValue{
								Creator: "6",
							},
						},
					},
				},
			},
		}, nil, nil).
		Once()

	mc.
		On("TxsHashGet", mock.Anything, "4").
		Return(cosmosclient.TxQuery{
			Txhash: "4",
			Tx: cosmosclient.StdTx{
				Value: cosmosclient.StdTxValue{
					Msg: []cosmosclient.Msg{
						{
							Type: "7",
							Value: cosmosclient.MsgValue{
								FromAddress: "8",
								Amount: []cosmosclient.Coin{
									{
										Amount: "9",
										Denom:  "10",
									},
								},
							},
						},
					},
				},
			},
		}, nil, nil).
		Once()

	properties := rosetta.NetworkProperties{
		Blockchain: "TheBlockchain",
		Network:    "TheNetwork",
		SupportedOperations: []string{
			"Transfer",
			"Reward",
		},
	}

	adapter := NewLaunchpad(TendermintAPI{Info: mt}, CosmosAPI{Transactions: mc}, properties)

	var h int64 = 1
	block, blockErr := adapter.Block(context.Background(), &types.BlockRequest{
		BlockIdentifier: &types.PartialBlockIdentifier{
			Index: &h,
		},
	})
	require.Nil(t, blockErr)
	require.NotNil(t, block)

	require.Equal(t, &types.BlockResponse{
		Block: &types.Block{
			BlockIdentifier: &types.BlockIdentifier{
				Index: int64(2),
				Hash:  "11",
			},
			Transactions: []*types.Transaction{
				&types.Transaction{
					TransactionIdentifier: &types.TransactionIdentifier{
						Hash: "4",
					},
					Operations: []*types.Operation{
						&types.Operation{
							OperationIdentifier: &types.OperationIdentifier{},
							Type:                "7",
							Status:              "TODO",
							Account: &types.AccountIdentifier{
								Address: "8",
							},
							Amount: &types.Amount{
								Value: "9",
								Currency: &types.Currency{
									Symbol: "10",
								},
							},
						},
					},
				},
				&types.Transaction{
					TransactionIdentifier: &types.TransactionIdentifier{
						Hash: "3",
					},
					Operations: []*types.Operation{
						&types.Operation{
							OperationIdentifier: &types.OperationIdentifier{},
							Type:                "5",
							Status:              "TODO",
							Account: &types.AccountIdentifier{
								Address: "6",
							},
						},
					},
				},
			},
			Timestamp: ti.UnixNano() / 1000000,
			ParentBlockIdentifier: &types.BlockIdentifier{
				Index: int64(1),
				Hash:  "12",
			},
		},
	}, block)

}