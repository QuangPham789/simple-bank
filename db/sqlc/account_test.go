package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "Tom",
		Balance:  1000,
		Currency: "USD",
	}

	acc, err := testQueries.CreateAccount(t.Context(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.Equal(t, arg.Owner, acc.Owner)
	require.Equal(t, arg.Balance, acc.Balance)
	require.Equal(t, arg.Currency, acc.Currency)

	require.NotZero(t, acc.ID)
}
