Feature: Account Management
  In order to manage accounts
  As a user of the system
  I need to be able to manage any accounts

  Scenario: List accounts will paginate correctly
    Given a JSON "chacerapp.v1.ListAccountsRequest"
      """
        { "pageSize": 1 }
      """
    And these resources are created:
      """
        {
          "resources": [
            {
              "@type": "chacerapp.v1.CreateAccountRequest",
              "account": { "displayName": "My Testing Account" },
              "account_id": "my-testing-account"
            },
            {
              "@type": "chacerapp.v1.CreateAccountRequest",
              "account": { "displayName": "My Second Testing Account" },
              "account_id": "my-second-testing-account"
            }
          ]
        }
      """
     When calling the "chacerapp.v1.Accounts/ListAccounts" RPC
     Then I will receive a successful response
      And the response value "accounts" will have a length of 1
      And the response value "accounts[0].name" will be "accounts/my-second-testing-account"
      And stashing the next page token from the response
    Given a JSON "chacerapp.v1.ListAccountsRequest"
      """
        { "pageSize": 1 }
      """
      And using the stashed next page token
     When calling the "chacerapp.v1.Accounts/ListAccounts" RPC
      And the response value "accounts" will have a length of 1
      And the response value "accounts[0].name" will be "accounts/my-testing-account"
      And stashing the next page token from the response
    Given a JSON "chacerapp.v1.ListAccountsRequest"
      """
        { "pageSize": 1 }
      """
      And using the stashed next page token
     When calling the "chacerapp.v1.Accounts/ListAccounts" RPC
      And the response value "accounts" will have a length of 0
      And the response value "nextPageToken" will be ""

  Scenario: Create an account fails with an invalid account parameters
    Given a JSON "chacerapp.v1.CreateAccountRequest"
      """
        {
          "account": {
            "displayName": "This is an extremely long description for an account. This is an extremely long description for an account. This is an extremely long description for an account. This is an extremely long description for an account. This is an extremely long description for an account."
          },
          "account_id": "This is an invalid name, with characters"
        }
      """
    When calling the "chacerapp.v1.Accounts/CreateAccount" RPC
    Then I will receive an error with code "INVALID_ARGUMENT"
    And the BadRequest error details will be for the following fields
      | account_id          | invalid account ID                                  |
      | account.displayName | display name must not be longer than 255 characters |

  Scenario: Create an account will succeed and the new account can be retrieved
    Given a JSON "chacerapp.v1.CreateAccountRequest"
      """
        {
          "account": { "displayName": "My Testing Account" },
          "account_id": "my-testing-account"
        }
      """
    When calling the "chacerapp.v1.Accounts/CreateAccount" RPC
    Then I will receive a successful response
     And the response value "name" will be "accounts/my-testing-account"
     And the response value "displayName" will be "My Testing Account"
     And the response value "status.phase" will be "ACCOUNT_PHASE_ACTIVE"
     And the response value "selfLink" will be "//chacerappapis.com/accounts/my-testing-account"
    # Test that calling the error will result in an error message
    When calling the "chacerapp.v1.Accounts/CreateAccount" RPC
    Then I will receive an error with code "ALREADY_EXISTS"
    # Verify we can retrieve the newly created account
    Given a JSON "chacerapp.v1.GetAccountRequest"
      """
        { "name": "accounts/my-testing-account" }
      """
    When calling the "chacerapp.v1.Accounts/GetAccount" RPC
    Then I will receive a successful response
     And the response value "name" will be "accounts/my-testing-account"
     And the response value "displayName" will be "My Testing Account"
     And the response value "status.phase" will be "ACCOUNT_PHASE_ACTIVE"
     And the response value "selfLink" will be "//chacerappapis.com/accounts/my-testing-account"

  Scenario: Verify an account can be activated and suspended
    Given a JSON "chacerapp.v1.SuspendAccountRequest"
      """
        {
          "name": "accounts/my-testing-account",
          "reason": "Billing"
        }
      """
     And these resources are created:
      """
        {
          "resources": [{
            "@type": "chacerapp.v1.CreateAccountRequest",
            "account": { "displayName": "My Testing Account" },
            "account_id": "my-testing-account"
          }]
        }
      """
    When calling the "chacerapp.v1.Accounts/SuspendAccount" RPC
    Then I will receive a successful response
    # Call the endpoint again and it should fail because the account
    # is no longer active.
    When calling the "chacerapp.v1.Accounts/SuspendAccount" RPC
    Then I will receive an error with code "FAILED_PRECONDITION"
    # Reactivate the account
    Given a JSON "chacerapp.v1.ActivateAccountRequest"
      """
        { "name": "accounts/my-testing-account" }
      """
    When calling the "chacerapp.v1.Accounts/ActivateAccount" RPC
    Then I will receive a successful response
    # Call it again to ensure that an activate account can not be activated
    When calling the "chacerapp.v1.Accounts/ActivateAccount" RPC
    Then I will receive an error with code "FAILED_PRECONDITION"

  Scenario: Account endpoints return a NotFound error when the account does not exist
    Given a JSON "chacerapp.v1.GetAccountRequest"
      """
        { "name": "accounts/this-account-does-not-exist" }
      """
     When calling the "chacerapp.v1.Accounts/GetAccount" RPC
     Then I will receive an error with code "NOT_FOUND"
    Given a JSON "chacerapp.v1.UpdateAccountRequest"
      """
        {
          "account": { "name": "accounts/this-account-does-not-exist" }
        }
      """
     When calling the "chacerapp.v1.Accounts/UpdateAccount" RPC
     Then I will receive an error with code "NOT_FOUND"
    Given a JSON "chacerapp.v1.DeleteAccountRequest"
      """
        { "name": "accounts/this-account-does-not-exist" }
      """
     When calling the "chacerapp.v1.Accounts/DeleteAccount" RPC
     Then I will receive an error with code "NOT_FOUND"
    Given a JSON "chacerapp.v1.SuspendAccountRequest"
      """
        { "name": "accounts/this-account-does-not-exist", "reason": "Billing" }
      """
     When calling the "chacerapp.v1.Accounts/SuspendAccount" RPC
     Then I will receive an error with code "NOT_FOUND"
    Given a JSON "chacerapp.v1.ActivateAccountRequest"
      """
        { "name": "accounts/this-account-does-not-exist" }
      """
     When calling the "chacerapp.v1.Accounts/ActivateAccount" RPC
     Then I will receive an error with code "NOT_FOUND"
    Given a JSON "chacerapp.v1.UpdateAccountQuotasRequest"
      """
        {
          "accountQuotas": { "name": "accounts/this-account-does-not-exist" }
        }
      """
     When calling the "chacerapp.v1.Accounts/UpdateAccountQuotas" RPC
     Then I will receive an error with code "NOT_FOUND"
