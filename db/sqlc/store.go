package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
	Querier
}

type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParam struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *SQLStore) TransferTx(ctx context.Context, param TransferTxParam) (TransferTxResult, error) {
	var result = TransferTxResult{}
	var err error
	error := s.execTx(ctx, func(q *Queries) error {
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: param.FromAccountId,
			ToAccountID:   param.ToAccountId,
			Amount:        param.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: param.FromAccountId,
			Amount:    -param.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: param.ToAccountId,
			Amount:    param.Amount,
		})
		if err != nil {
			return err
		}

		//update account balance

		//result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		//	ID:     param.FromAccountId,
		//	Amount: -param.Amount,
		//})
		//
		//result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		//	ID:     param.ToAccountId,
		//	Amount: param.Amount,
		//})
		if param.FromAccountId < param.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, param.FromAccountId, -param.Amount, param.ToAccountId, param.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, param.ToAccountId, param.Amount, param.FromAccountId, -param.Amount)
		}

		return nil
	})

	return result, error
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
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
