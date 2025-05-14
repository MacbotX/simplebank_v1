package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/MacbotX/simplebank_v1/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntries(t *testing.T) Entry {
	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	// making sure its not zero
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntries(t *testing.T) {
	createRandomEntries(t)
}

func TestGetEntires(t *testing.T) {
	entry1 := createRandomEntries(t)
	entry, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry1.ID, entry.ID)
	require.Equal(t, entry1.AccountID, entry.AccountID)
	require.Equal(t, entry1.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

}

func TestGetUpdateEntries(t *testing.T) {
	entry1 := createRandomEntries(t)
	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomMoney(),
	}

	entry, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.ID, entry.ID)
	require.Equal(t, entry1.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntries(t)
	}

	arg := ListEntryParams{
		Limit:  5,
		Offset: 5,
	}

	entry, err := testQueries.ListEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entry, 5)

	for _, entries := range entry {
		require.NotEmpty(t, entries)
		require.NotEmpty(t, entries.AccountID)
		require.NotEmpty(t, entries.CreatedAt)
	}

}

func TestDeleteEntries(t *testing.T) {
	entry1 := createRandomEntries(t)
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	// checking if the account entries is been deleted successfully
	entry, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.True(t, errors.Is(err, sql.ErrNoRows))
	require.Empty(t, entry)

}
