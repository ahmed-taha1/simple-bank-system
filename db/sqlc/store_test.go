package db

import (
	"context"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	senderAccount := createRandomAccount(t)
	reciverAccount := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: senderAccount.ID,
				ToAccountID:   reciverAccount.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check results
	isExisted := make(map[int]bool)

	for range n {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, senderAccount.ID, transfer.FromAccountID)
		require.Equal(t, reciverAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check accounts
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, senderAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, reciverAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, senderAccount.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, reciverAccount.ID, toAccount.ID)

		// check accounts' balance
		outMoneyFromAccount1 := senderAccount.Balance - fromAccount.Balance
		inMoneyToAccount2 := toAccount.Balance - reciverAccount.Balance
		require.Equal(t, outMoneyFromAccount1, inMoneyToAccount2)
		require.True(t, outMoneyFromAccount1 > 0)
		require.True(t, outMoneyFromAccount1%amount == 0)

		k := int(outMoneyFromAccount1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, isExisted, k)
		isExisted[k] = true
	}

	// check the final updated balances
	senderUpdatedAccount, err := testQueries.GetAccountById(context.Background(), senderAccount.ID)
	require.NoError(t, err)

	receiverUpdatedAccount, err := testQueries.GetAccountById(context.Background(), reciverAccount.ID)
	require.NoError(t, err)

	require.Equal(t, senderAccount.Balance-int64(n)*amount, senderUpdatedAccount.Balance)
	require.Equal(t, reciverAccount.Balance+int64(n)*amount, receiverUpdatedAccount.Balance)
}
