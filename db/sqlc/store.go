package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// store providers all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLstore providers all functions to execute SQL DB queries and transactions
type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

// Newstore creates a new store instance\
func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// start the transaction with default options
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("falied to begin transactions: %v ", err)
	}

	q := New(tx) //Initialize with the transaction not the pool
	err = fn(q)  //Execute the transactional function
	if err != nil {
		// attempt to rollback transactions on error
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("transaction err: %v, rollback err: %v", rbErr, err)
		}
		return err
	}

	return tx.Commit(ctx)
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the trasnfer transactions
type TransferTxResult struct {
	Transfer      Transfer `json:"transfer"`
	FromAccountID Account  `json:"from_account"`
	ToAccountID   Account  `json:"to_account"`
	FromEntry     Entry    `json:"from_entry"`
	ToEntry       Entry    `json:"to_entry"`
}

// type contextKey string
// const txKey contextKey = "txName"

// TransferTx performs a money transfer from one account to the other
// it create a transfer record, add account entries and update accounts balance within a single database trnx
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		//this is used in testing the db deadlock
		// txName := ctx.Value(txKey)

		// Create transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))

		if err != nil {
			return err
		}

		// create entry from account sender
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// create entry from account reciever
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// getting the account and updating it
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccountID, result.ToAccountID, err = addMoney(ctx, q,  arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil{
				return err
			}
		} else {
			result.ToAccountID, result.FromAccountID, err = addMoney(ctx, q,  arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil{
				return err
			}
		}

		// getting the account and updating it ends here

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID2,
		Amount: amount2,
	})
	return
}
