package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/MacbotX/simplebank_v1/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfers(t *testing.T) Transfer {
	from_account1 := createRandomAccount(t)
	to_account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: from_account1.ID,
		ToAccountID:   to_account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfers(t)
}

func TestGetTransfers(t *testing.T) {
	from_account1 := createRandomTransfers(t)
	transfer, err := testQueries.GetTransfer(context.Background(), from_account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, from_account1.FromAccountID, transfer.FromAccountID)
	require.Equal(t, from_account1.ToAccountID, transfer.ToAccountID)
	require.Equal(t, from_account1.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestGetListTrasfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfers(t)

	}
	arg := ListTransferParams{
		Limit:  5,
		Offset: 5,
	}

	transfer, err := testQueries.ListTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfer, 5)

	for _, ts := range transfer {
		require.NotEmpty(t, ts.FromAccountID)
		require.NotEmpty(t, ts.ToAccountID)
		require.NotEmpty(t, ts.Amount)
	}

}

func TestDeleteTransfer(t *testing.T)  {
	from_account1 := createRandomTransfers(t)
	err := testQueries.DeleteTransfer(context.Background(), from_account1.ID)
	require.NoError(t, err)

	transfer, err := testQueries.GetTransfer(context.Background(), from_account1.ID)
	require.Error(t, err)
	require.Empty(t, transfer)
	require.True(t, errors.Is(err , sql.ErrNoRows))
}
