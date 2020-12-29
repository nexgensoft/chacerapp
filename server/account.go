package server

import (
	"context"

	"github.com/chacerapp/apiserver/name"
	"github.com/chacerapp/apiserver/server/serverpb"
	"github.com/chacerapp/apiserver/store"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var validAccountSuspendedReasons = []string{
	"Billing", "Fraud",
}

func (s *server) ListAccounts(ctx context.Context, req *serverpb.ListAccountsRequest) (*serverpb.ListAccountsResponse, error) {
	// Validate the pagination request
	pageInfo, err := s.validatePageableRequest(req)
	if err != nil {
		return nil, err
	}

	accounts, err := s.store.ListAccounts(ctx, store.WithPageSize(req.PageSize), store.WithPageInfo(pageInfo))
	if err != nil {
		return nil, err
	}

	var nextPageToken string
	// The next page token should only be generated when the number
	// of results being returned is equal to the page size. The lack
	// of a next page token is used to determine if a next page exists.
	if len(accounts) == int(req.PageSize) {
		nextPageToken, err = s.store.GenerateNextPageToken(pageInfo, req.PageSize)
		if err != nil {
			return nil, err
		}
	}

	return &serverpb.ListAccountsResponse{
		Accounts:      accounts,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *server) CreateAccount(ctx context.Context, req *serverpb.CreateAccountRequest) (*serverpb.Account, error) {
	if err := validateCreateAccount(req); err != nil {
		return nil, err
	}

	account := proto.Clone(req.Account).(*serverpb.Account)
	account.Name = name.BuildAccount(req.AccountId)

	if account, err := s.store.CreateAccount(ctx, account); err != nil {
		return nil, err
	} else if account == nil {
		return nil, errAlreadyExists
	} else {
		return account, nil
	}
}

func (s *server) GetAccount(ctx context.Context, req *serverpb.GetAccountRequest) (*serverpb.Account, error) {
	if _, err := name.ParseAccount(req.Name); err != nil {
		return nil, err
	}

	if account, err := s.store.GetAccount(ctx, req.Name); err != nil {
		return nil, err
	} else if account == nil {
		return nil, errNotFound
	} else {
		return account, nil
	}
}

func (s *server) UpdateAccount(ctx context.Context, req *serverpb.UpdateAccountRequest) (*serverpb.Account, error) {
	if err := validateUpdateAccount(req); err != nil {
		return nil, err
	}

	if account, err := s.store.UpdateAccount(ctx, req.Account, store.WithUpdateMask(req.UpdateMask)); err != nil {
		return nil, err
	} else if account == nil {
		return nil, errNotFound
	} else {
		return account, nil
	}
}

func (s *server) ActivateAccount(ctx context.Context, req *serverpb.ActivateAccountRequest) (*serverpb.ActivateAccountResponse, error) {
	if _, err := name.ParseAccount(req.Name); err != nil {
		return nil, err
	}

	if account, err := s.store.GetAccount(ctx, req.Name); err != nil {
		return nil, err
	} else if account == nil {
		return nil, errNotFound
	} else if account.Status.Phase == serverpb.AccountPhase_ACCOUNT_PHASE_ACTIVE {
		return nil, errFailedPrecondition("account is already active")
	}

	status, err := s.store.UpdateAccountStatus(ctx, req.Name, &serverpb.AccountStatus{
		Phase: serverpb.AccountPhase_ACCOUNT_PHASE_ACTIVE,
	})
	if err != nil {
		return nil, err
	} else if status == nil {
		return nil, errNotFound
	}
	return &serverpb.ActivateAccountResponse{}, nil
}

func (s *server) SuspendAccount(ctx context.Context, req *serverpb.SuspendAccountRequest) (*serverpb.SuspendAccountResponse, error) {
	if err := validateSuspendAccount(req); err != nil {
		return nil, err
	}

	if account, err := s.store.GetAccount(ctx, req.Name); err != nil {
		return nil, err
	} else if account == nil {
		return nil, errNotFound
	} else if account.Status.Phase == serverpb.AccountPhase_ACCOUNT_PHASE_SUSPENDED {
		return nil, errFailedPrecondition("account is already suspended")
	}

	status, err := s.store.UpdateAccountStatus(ctx, req.Name, &serverpb.AccountStatus{
		Phase:   serverpb.AccountPhase_ACCOUNT_PHASE_SUSPENDED,
		Reason:  req.Reason,
		Message: req.Message,
	})
	if err != nil {
		return nil, err
	} else if status == nil {
		return nil, errNotFound
	}
	return &serverpb.SuspendAccountResponse{}, nil
}

func (s *server) DeleteAccount(ctx context.Context, req *serverpb.DeleteAccountRequest) (*serverpb.Account, error) {
	if _, err := name.ParseAccount(req.Name); err != nil {
		return nil, err
	}

	existing, err := s.store.DeleteAccount(ctx, req.Name)
	if err != nil {
		return nil, err
	} else if existing == nil {
		return nil, errNotFound
	}
	return existing, nil
}

func (s *server) GetAccountStatus(ctx context.Context, req *serverpb.GetAccountStatusRequest) (*serverpb.AccountStatus, error) {
	if _, err := name.ParseAccount(req.Name); err != nil {
		return nil, err
	}

	if account, err := s.store.GetAccount(ctx, req.Name); err != nil {
		return nil, err
	} else if account == nil {
		return nil, errNotFound
	} else {
		return account.GetStatus(), nil
	}
}

func (s *server) UpdateAccountStatus(ctx context.Context, req *serverpb.UpdateAccountStatusRequest) (*serverpb.AccountStatus, error) {
	if _, err := name.ParseAccount(req.AccountStatus.Name); err != nil {
		return nil, err
	}

	if updatedStatus, err := s.store.UpdateAccountStatus(ctx, req.AccountStatus.Name, req.AccountStatus); err != nil {
		return nil, err
	} else if updatedStatus == nil {
		return nil, errNotFound
	} else {
		return updatedStatus, nil
	}
}

func (s *server) GetAccountQuotas(ctx context.Context, req *serverpb.GetAccountQuotasRequest) (*serverpb.AccountQuotas, error) {
	if _, err := name.ParseAccount(req.Name); err != nil {
		return nil, err
	}

	if account, err := s.store.GetAccount(ctx, req.Name); err != nil {
		return nil, err
	} else if account == nil {
		return nil, errNotFound
	} else {
		return account.GetQuotas(), nil
	}
}

func (s *server) UpdateAccountQuotas(ctx context.Context, req *serverpb.UpdateAccountQuotasRequest) (*serverpb.AccountQuotas, error) {
	if _, err := name.ParseAccount(req.AccountQuotas.Name); err != nil {
		return nil, err
	}

	if updatedQuotas, err := s.store.UpdateAccountQuotas(ctx, req.AccountQuotas.Name, req.AccountQuotas); err != nil {
		return nil, err
	} else if updatedQuotas == nil {
		return nil, errNotFound
	} else {
		return updatedQuotas, nil
	}
}

func validateSuspendAccount(req *serverpb.SuspendAccountRequest) error {
	var errs field.ErrorList
	if _, err := name.ParseAccount(req.Name); err != nil {
		if s, ok := status.FromError(err); ok {
			errs = append(errs, field.Invalid(field.NewPath("name"), req.Name, s.Message()))
		} else {
			return err
		}
	}

	validReason := false
	for i := range validAccountSuspendedReasons {
		if req.Reason == validAccountSuspendedReasons[i] {
			validReason = true
			break
		}
	}
	if !validReason {
		errs = append(errs, field.NotSupported(field.NewPath("reason"), req.Reason, validAccountSuspendedReasons))
	}
	if len(req.Message) > 1024 {
		errs = append(errs, field.Invalid(field.NewPath("message"), req.Message, "message must be between 0 and 1024 characters"))
	}

	return convertErrorList(errs)
}

func validateCreateAccount(req *serverpb.CreateAccountRequest) error {
	var errs field.ErrorList

	if req.AccountId == "" {
		errs = append(errs, field.Required(field.NewPath("account_id"), "account_id is required"))
	} else if !name.ValidResourceID(req.AccountId) {
		errs = append(errs, field.Invalid(field.NewPath("account_id"), req.AccountId, "invalid account ID"))
	}

	errs = append(errs, validateAccount(req.Account, true)...)

	return convertErrorList(errs)
}

func validateUpdateAccount(req *serverpb.UpdateAccountRequest) error {
	return convertErrorList(validateAccount(req.Account, false))
}

// Validates the account has all of its values set correctly. This will return
// a list of errors for each field that is missing or invalid.
func validateAccount(account *serverpb.Account, creating bool) field.ErrorList {
	accountPath := field.NewPath("account")
	if account == nil {
		return field.ErrorList{
			field.Required(accountPath, "account is required"),
		}
	}

	var errs field.ErrorList
	errs = append(errs, validateLabels(accountPath, account)...)
	errs = append(errs, validateAnnotations(accountPath, account)...)
	errs = append(errs, validateDisplayName(accountPath, account)...)

	// Specific checks to perform only when updating a resource
	if !creating {
		// Verify the account name is valid
		if account.Name == "" {
			errs = append(errs, field.Required(accountPath.Child("name"), "name is required"))
		} else if _, err := name.ParseAccount(account.Name); err != nil {
			errs = append(errs, field.Invalid(accountPath.Child("name"), account.Name, err.Error()))
		}
	}

	return errs
}
