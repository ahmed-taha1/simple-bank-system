package db

import (
	"context"
	"database/sql"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) db.Account {
	arg := db.CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoeny(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	/*account := */ CreateRandomAccount(t)

	// clean up
	// err := testQueries.DeleteAccount(context.Background(), account.ID)
	// require.NoError(t, err)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	arg := db.UpdateAccountBalanceParams{
		ID:      account1.ID,
		Balance: util.RandomMoeny(),
	}

	account2, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	// Create 10 random accounts
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	arg := db.ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
