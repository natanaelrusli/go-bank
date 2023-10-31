package db

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/natanaelrusli/go-bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	arg := CreateTransferParams{
		FromAccountID: util.RandomInt(0, 100),
		ToAccountID:   util.RandomInt(0, 100),
		Amount:        util.RandomInt(0, 10000),
	}

	rows := sqlmock.NewRows([]string{
		"id",
		"from_account_id",
		"to_account_id",
		"amount",
		"created_at",
	}).AddRow(util.RandomInt(0, 100), arg.FromAccountID, arg.ToAccountID, arg.Amount, time.Now())

	mocksql.ExpectQuery("INSERT INTO transfers").
		WithArgs(arg.FromAccountID, arg.ToAccountID, arg.Amount).
		WillReturnRows(rows)

	transfer, err := testMockQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NoError(t, mocksql.ExpectationsWereMet())

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestGetTransfer(t *testing.T) {
	arg := CreateTransferParams{
		FromAccountID: util.RandomInt(0, 100),
		ToAccountID:   util.RandomInt(0, 100),
		Amount:        util.RandomInt(0, 10000),
	}

	mockId := util.RandomInt(0, 100)

	rows := sqlmock.NewRows([]string{
		"id",
		"from_account_id",
		"to_account_id",
		"amount",
		"created_at",
	}).AddRow(mockId, arg.FromAccountID, arg.ToAccountID, arg.Amount, time.Now())

	mocksql.ExpectQuery("SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers").
		WithArgs(mockId).
		WillReturnRows(rows)

	transfer, err := testMockQueries.GetTransfer(context.Background(), mockId)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, mockId, transfer.ID)
}

func TestListTransfers(t *testing.T) {
	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	numOfRows := 5

	rows := sqlmock.NewRows([]string{
		"id",
		"from_account_id",
		"to_account_id",
		"amount",
		"created_at",
	})

	for i := 0; i < numOfRows; i++ {
		row := []driver.Value{
			util.RandomInt(0, 100),
			util.RandomInt(0, 100),
			util.RandomInt(0, 100),
			util.RandomMoney(),
			time.Now(),
		}
		rows = rows.AddRow(row...)
	}

	mocksql.ExpectQuery("SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers ORDER BY id").
		WithArgs(arg.Limit, arg.Offset).
		WillReturnRows(rows)

	transfers, err := testMockQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, int(arg.Limit))
}
