package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Printf(">> Before >> account1 Balance: %v, account2 balance: %v\n", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		//this is used for testing the deadlock in our db
		// txName := fmt.Sprintf("tx : %d", i+1)
		go func() {
			// ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)
	// check results from the channel
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check the transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// checking if the transfer record is created in the DB
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check the account entries of the result for account1 (sender)
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// checking if the entries record is created in the DB
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check the account entries of the result for account2 (reciever)
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// checking if the entries record is created in the DB
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts (sender)
		fromAccount := result.FromAccountID
		require.NotEmpty(t, fromAccount)
		require.NotZero(t, fromAccount.ID)
		require.Equal(t, account1.ID, fromAccount.ID)

		// check accounts (receiver)
		toAccount := result.ToAccountID
		require.NotEmpty(t, toAccount)
		require.NotZero(t, toAccount.ID)
		require.Equal(t, account2.ID, toAccount.ID)

		// check the account's balance
		fmt.Printf(">> Tx >> account1 Balance: %v, account2 balance: %v\n", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // amount, 2 * amount, 3 * amount,....n * amount
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check the final update balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">>Updated>> account1 Balance: %v, account2 balance: %v", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Printf(">> Before >> account1 Balance: %v, account2 balance: %v\n", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i % 2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		//this is used for testing the deadlock in our db
		// txName := fmt.Sprintf("tx : %d", i+1)
		go func() {
			// ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	// check results from the channel
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	//check the final update balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">>Updated>> account1 Balance: %v, account2 balance: %v", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
} 
