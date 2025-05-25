package db

import (
	"context"
	"database/sql"
	"fmt"
)

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Store provides all functions to execute SQL queries and transactions.
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store instance with the given database connection and queries.
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rollBackError := tx.Rollback()
		if rollBackError != nil {
			return fmt.Errorf("transaction error: %v, roll back error: %v", err, rollBackError)
		}
		return err
	}

	tx.Commit()
	return nil
}

// TransferTx performs a money transfer from one account to another.
// It creates a transfer record, add account entries and update accounts' balance within a single database transaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// ? condition to avoid deadlock
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = transferMoney(
				ctx,
				q,
				arg.FromAccountID, -arg.Amount,
				arg.ToAccountID, arg.Amount,
			)
		} else {
			result.ToAccount, result.FromAccount, err = transferMoney(
				ctx,
				q,
				arg.ToAccountID, arg.Amount,
				arg.FromAccountID, -arg.Amount,
			)
		}
		if err != nil {
			return fmt.Errorf("failed to transfer money: %w", err)
		}
		return nil
	})

	return result, err
}

func transferMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddBalanceToAccountById(ctx, AddBalanceToAccountByIdParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	account2, err = q.AddBalanceToAccountById(ctx, AddBalanceToAccountByIdParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return account1, account2, err
}
