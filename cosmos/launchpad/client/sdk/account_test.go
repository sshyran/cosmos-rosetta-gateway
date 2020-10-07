package sdk

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/cosmos-rosetta-gateway/cosmos/launchpad/clienttest"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestAuthAccountClient(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx, cancel := clienttest.Ctx()
	t.Cleanup(cancel)
	e, err := clienttest.NewLaunchpad(ctx, "crgapp")
	require.NoError(t, err)
	t.Cleanup(e.Cleanup)

	client := NewClient(e.SDKAddr)

	accountsb := &bytes.Buffer{}
	require.NoError(t, cmdrunner.
		New().
		Run(ctx, step.New(
			step.Exec(
				e.Appcli(), "keys", "list",
			),
			step.Stdout(accountsb),
		)))
	var accounts []struct {
		Address string `json:"address"`
	}
	require.NoError(t, json.NewDecoder(accountsb).Decode(&accounts))
	require.True(t, len(accounts) != 0)

	addr := accounts[0].Address
	res, err := client.GetAuthAccount(ctx, addr, 0)
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Greater(t, res.Height, int64(0))
	require.Equal(t, addr, res.Result.Value.Address)
	require.Equal(t, "2", res.Result.Value.AccountNumber)
	require.Equal(t, "1", res.Result.Value.Sequence)
	require.Equal(t, int64(1000), res.Result.Value.Coins[0].Amount.Int64())
}
