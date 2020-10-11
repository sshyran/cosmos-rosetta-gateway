package launchpad

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"

	"github.com/coinbase/rosetta-sdk-go/types"
)

const (
	// Metadata Keys
	ChainIdKey       = "chain_id"
	SequenceKey      = "sequence"
	AccountNumberKey = "account_number"
	GasKey           = "gas"
)

type PayloadReqMetadata struct {
	ChainId       string
	Sequence      uint64
	AccountNumber uint64
	Gas           uint64
	Memo          string
	Fee           sdk.Coin
}

// GetMetadataFromPayloadReq obtains the metadata from the request to /construction/payloads endpoint.
func GetMetadataFromPayloadReq(req *types.ConstructionPayloadsRequest) (*PayloadReqMetadata, error) {
	chainId, ok := req.Metadata[ChainIdKey].(string)
	if !ok {
		return nil, fmt.Errorf("chain_id metadata was not provided")
	}

	sequence, ok := req.Metadata[SequenceKey]
	if !ok {
		return nil, fmt.Errorf("sequence metadata was not provided")
	}
	seqStr, ok := sequence.(string)
	if !ok {
		return nil, fmt.Errorf("invalid sequence value")
	}
	seqNum, err := strconv.Atoi(seqStr)
	if err != nil {
		return nil, fmt.Errorf("error converting sequence num to int")
	}

	accountNum, ok := req.Metadata[AccountNumberKey]
	if !ok {
		return nil, fmt.Errorf("account_number metadata was not provided")
	}
	accStr, ok := accountNum.(string)
	if !ok {
		return nil, fmt.Errorf("invalid account_number value")
	}
	accNum, err := strconv.Atoi(accStr)
	if err != nil {
		return nil, fmt.Errorf("error converting account num to int")
	}

	gasNum, ok := req.Metadata[GasKey]
	if !ok {
		return nil, fmt.Errorf("gas metadata was not provided")
	}
	gasF64, ok := gasNum.(float64)
	if !ok {
		return nil, fmt.Errorf("invalid gas value")
	}

	memo, ok := req.Metadata[OptionMemo]
	if !ok {
		return nil, fmt.Errorf("memo metadata was not provided")
	}
	memoStr, ok := memo.(string)
	if !ok {
		return nil, fmt.Errorf("invalid memo")
	}

	_, ok = req.Metadata[OptionMemo]
	if !ok {
		return nil, fmt.Errorf("fee metadata was not provided")
	}
	return &PayloadReqMetadata{
		ChainId:       chainId,
		Sequence:      uint64(seqNum),
		AccountNumber: uint64(accNum),
		Gas:           uint64(gasF64),
		Memo:          memoStr,
		Fee:           sdk.NewCoin("tokens", sdk.NewInt(0.001)),
	}, nil
}
