package store

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/chacerapp/apiserver/name"
	"github.com/chacerapp/apiserver/server/serverpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/lib/pq"
)

// Account provides an interface that can be used for managing accounts within storage
type Account interface {
	// GetAccount will retrieve an account by name from storage
	//
	// This function will return a nil account when an account does not
	// exist with the given name. An error will only be returned when
	// the account failed to be retrieved.
	GetAccount(ctx context.Context, name string) (*serverpb.Account, error)
	ListAccounts(ctx context.Context, opts ...ListOption) ([]*serverpb.Account, error)
	CreateAccount(ctx context.Context, account *serverpb.Account) (*serverpb.Account, error)
	UpdateAccount(ctx context.Context, account *serverpb.Account, opts ...UpdateOption) (*serverpb.Account, error)
	UpdateAccountStatus(ctx context.Context, accountName string, status *serverpb.AccountStatus) (*serverpb.AccountStatus, error)
	UpdateAccountQuotas(ctx context.Context, accountName string, quotas *serverpb.AccountQuotas) (*serverpb.AccountQuotas, error)
	DeleteAccount(ctx context.Context, name string) (*serverpb.Account, error)
}

func (s *store) GetAccount(ctx context.Context, name string) (*serverpb.Account, error) {
	return doGetAccount(ctx, s.db, name)
}

// ListAccounts will list all of the accounts in storage
func (s *store) ListAccounts(ctx context.Context, opts ...ListOption) ([]*serverpb.Account, error) {
	options := getListOptions(opts...)
	rows, err := s.db.Query(paginateQuery(selectAccountBaseQuery+" ORDER BY name", options.pageInfo, options.pageSize))
	if err != nil {
		return nil, err
	}

	// Close the rows once we are done retrieving results
	defer rows.Close()

	var accounts []*serverpb.Account
	for rows.Next() {
		account, err := scanAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// CreateAccount will create a new account in storage
//
// Only settable fields are respected when creating an account. All other fields
// will be discarded or overwritten. The Quotas and Status of the returned account
// will be guaranteed to be set. If an account with the provided name already exists
// a nil account will be returned.
func (s *store) CreateAccount(ctx context.Context, account *serverpb.Account) (*serverpb.Account, error) {
	var newAccount *serverpb.Account

	accountName, err := name.ParseAccount(account.Name)
	if err != nil {
		return nil, err
	}

	// Run in a transaction so we can atomically check if the account already exists
	err = doTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		// Check that the account doesn't already exists, when it does
		// then we should return without returning an account.
		if existing, err := doGetAccount(ctx, tx, account.Name); err != nil || existing != nil {
			return err
		}

		// Create the new account with all the defaults that should be set
		newAccount = &serverpb.Account{
			Name:        account.Name,
			SelfLink:    serviceName + account.Name,
			CreateTime:  ptypes.TimestampNow(),
			DisplayName: account.DisplayName,
			Quotas: &serverpb.AccountQuotas{
				Name: account.Name + "/quotas",
			},
			Status: &serverpb.AccountStatus{
				Phase: serverpb.AccountPhase_ACCOUNT_PHASE_ACTIVE,
			},
		}

		status, err := protoMarshaller.MarshalToString(newAccount.Status)
		if err != nil {
			return err
		}
		quotas, err := protoMarshaller.MarshalToString(newAccount.Quotas)
		if err != nil {
			return err
		}
		created, err := ptypes.Timestamp(newAccount.CreateTime)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			accountInsertQuery,
			accountName,
			newAccount.DisplayName,
			status,
			quotas,
			created,
		)
		return err
	})

	if err != nil {
		return nil, err
	}

	return newAccount, nil
}

func (s *store) UpdateAccount(ctx context.Context, account *serverpb.Account, opts ...UpdateOption) (*serverpb.Account, error) {
	return s.doUpdateAccount(ctx, account.Name, func(existing *serverpb.Account) error {
		existing.DisplayName = account.DisplayName
		return nil
	})
}

func (s *store) UpdateAccountStatus(ctx context.Context, name string, status *serverpb.AccountStatus) (*serverpb.AccountStatus, error) {
	updated, err := s.doUpdateAccount(ctx, name, func(existing *serverpb.Account) error {
		existing.Status.Phase = status.Phase
		existing.Status.Reason = status.Reason
		existing.Status.Message = status.Message
		return nil
	})

	if err != nil || updated == nil {
		return nil, err
	}

	return updated.Status, nil
}

func (s *store) UpdateAccountQuotas(ctx context.Context, name string, Quotas *serverpb.AccountQuotas) (*serverpb.AccountQuotas, error) {
	updated, err := s.doUpdateAccount(ctx, name, func(existing *serverpb.Account) error {
		existing.Quotas.Devices = Quotas.Devices
		existing.Quotas.Locations = Quotas.Locations
		return nil
	})

	if err != nil || updated == nil {
		return nil, err
	}

	return updated.Quotas, nil
}

// DeleteAccount will delete an account from storage.
//
// If the requested account does not exist a nil account will be returned. Otherwise,
// the returned account will be the account at the time of deletion. This operation
// can not be undone.
func (s *store) DeleteAccount(ctx context.Context, name string) (*serverpb.Account, error) {
	var account *serverpb.Account

	// Run in a transaction so we can atomically check if the account already exists
	err := doTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		// Check if the account exists
		if account, err = doGetAccount(ctx, tx, name); err != nil {
			return err
		} else if account == nil {
			// return nil here so we can indicate the account does not exist in the system
			return nil
		}

		_, err = tx.Exec(accountDeleteQuery, name)
		return err
	})

	if err != nil {
		return nil, err
	}

	return account, nil
}

func doGetAccount(ctx context.Context, query retriever, accountName string) (*serverpb.Account, error) {
	accountName, err := name.ParseAccount(accountName)
	if err != nil {
		return nil, err
	}

	rows := query.QueryRowContext(ctx, selectAccountBaseQuery+` WHERE name = $1`, accountName)
	account, err := scanAccount(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *store) doUpdateAccount(ctx context.Context, account string, updater func(existing *serverpb.Account) error) (*serverpb.Account, error) {
	var existing *serverpb.Account

	accountName, err := name.ParseAccount(account)
	if err != nil {
		return nil, err
	}

	err = doTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		if existing, err = doGetAccount(ctx, tx, account); err != nil || existing == nil {
			return err
		}

		// Override the values in the existing account
		existing.UpdateTime = ptypes.TimestampNow()
		if err := updater(existing); err != nil {
			return err
		}

		status, err := protoMarshaller.MarshalToString(existing.Status)
		if err != nil {
			return err
		}
		quotas, err := protoMarshaller.MarshalToString(existing.Quotas)
		if err != nil {
			return err
		}
		updated, err := ptypes.Timestamp(existing.UpdateTime)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, updateAccountQuery, existing.DisplayName, status, quotas, updated, accountName)
		return err
	})

	return existing, nil
}

func scanAccount(scan scanner) (*serverpb.Account, error) {
	// Allocate all the variables we will need to scan
	var uid, accountName, displayName, status, quotas string
	var createdTime time.Time
	var updateTime pq.NullTime
	// Scan the row from the database
	if err := scan.Scan(&uid, &accountName, &displayName, &status, &quotas, &createdTime, &updateTime); err != nil {
		return nil, err
	}

	created, err := ptypes.TimestampProto(createdTime)
	if err != nil {
		return nil, err
	}

	var updated *timestamp.Timestamp
	if updateTime.Valid {
		updated, err = ptypes.TimestampProto(updateTime.Time)
		if err != nil {
			return nil, err
		}
	}

	// Unmarshal the quotas and status columns
	statusProtobuf := &serverpb.AccountStatus{}
	if err := protoUnmarshaller.Unmarshal(strings.NewReader(status), statusProtobuf); err != nil {
		return nil, err
	}
	quotasProtobuf := &serverpb.AccountQuotas{}
	if err := protoUnmarshaller.Unmarshal(strings.NewReader(quotas), quotasProtobuf); err != nil {
		return nil, err
	}

	fqName := name.BuildAccount(accountName)

	return &serverpb.Account{
		Name:        fqName,
		Uid:         uid,
		SelfLink:    serviceName + fqName,
		CreateTime:  created,
		UpdateTime:  updated,
		DisplayName: displayName,
		Status:      statusProtobuf,
		Quotas:      quotasProtobuf,
	}, nil
}

const selectAccountBaseQuery = `
SELECT id, name, display_name, status, quotas, created_time, updated_time FROM account`

const accountInsertQuery = `
INSERT INTO account (name, display_name, status, quotas, created_time, updated_time)
VALUES ($1, $2, $3, $4, $5, NULL)`

const accountDeleteQuery = `
DELETE FROM account WHERE name = $1`

const updateAccountQuery = `
UPDATE account SET display_name = $1, status = $2, quotas = $3, updated_time = $4 WHERE name = $5`
