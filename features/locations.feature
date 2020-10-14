Feature: Manage devices on the API
  Background: Create an account we can create locations in
    Given these resources are created:
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

  Scenario: Able to create, get, update, and delete a location on an account
    Given a JSON "chacerapp.v1.CreateLocationRequest"
      """
        {
          "parent": "accounts/default-account",
          "location": {
            "displayName": "Default",
            "description": "This is my default location"
          }
        }
      """
     When calling the "chacerapp.v1.Locations/CreateLocation" RPC
     Then I will receive a successful response
      And the response value "displayName" will be "Default"
      And the response value "description" will be "This is my default location"
     When calling the "chacerapp.v1.Locations/CreateLocation" RPC
     Then I will receive an error with code "ALREADY_EXISTS"
    Given a JSON "chacerapp.v1.GetLocationRequest"
      """
        { "name": "accounts/default-account/locations/default" }
      """
     When calling the "chacerapp.v1.Locations/GetLocation" RPC
     Then I will receive a successful response
      And the response value "displayName" will be "Default"
      And the response value "description" will be "This is my default location"
    Given a JSON "chacerapp.v1.UpdateLocationRequest"
      """
        {
          "location": {
            "name": "accounts/default-account/locations/default",
            "displayName": "Default #2",
            "description": "This is my default location #2"
          }
        }
      """
     When calling the "chacerapp.v1.Locations/UpdateLocation" RPC
     Then I will receive a successful response
      And the response value "displayName" will be "Default #2"
      And the response value "description" will be "This is my default location #2"
    Given a JSON "chacerapp.v1.DeleteLocationRequest"
      """
        { "name": "accounts/default-account/locations/default" }
      """
     When calling the "chacerapp.v1.Locations/DeleteLocation" RPC
     Then I will receive a successful response
    Given a JSON "chacerapp.v1.GetLocationRequest"
      """
        { "name": "accounts/default-account/locations/default" }
      """
     When calling the "chacerapp.v1.Locations/GetLocation" RPC
     Then I will receive an error with code "NOT_FOUND"

  Scenario: Verify that a partial update with a field mask will work correctly
    Given a JSON "chacerapp.v1.CreateLocationRequest"
      """
        {
          "parent": "accounts/default-account",
          "location": {
            "displayName": "Default",
            "description": "This is my default location"
          }
        }
      """
     When calling the "chacerapp.v1.Locations/CreateLocation" RPC
     Then I will receive a successful response
    Given a JSON "chacerapp.v1.UpdateLocationRequest"
      """
        {
          "location": {
            "name": "accounts/default-account/locations/default",
            "displayName": "Default #2"
          },
          "updateMask": {
            "paths": [
              "displayName"
            ]
          }
        }
      """
     When calling the "chacerapp.v1.Locations/UpdateLocation" RPC
     Then I will receive a successful response
      And the response value "displayName" will be "Default #2"
      And the response value "description" will be "This is my default location"

  Scenario: Able to correctly list out the locations
    # Create several devices across the two accounts
    Given a JSON "chacerapp.v1.CreateLocationRequest"
      """
        {
          "parent": "accounts/default-account",
          "location": {
            "displayName": "Default",
            "description": "This is my default location"
          }
        }
      """
     When calling the "chacerapp.v1.Locations/CreateLocation" RPC
     Then I will receive a successful response
    Given a JSON "chacerapp.v1.CreateLocationRequest"
      """
        {
          "parent": "accounts/default-account",
          "location": {
            "displayName": "Secondary",
            "description": "This is my secondary location"
          }
        }
      """
     When calling the "chacerapp.v1.Locations/CreateLocation" RPC
     Then I will receive a successful response
    Given a JSON "chacerapp.v1.CreateLocationRequest"
      """
        {
          "parent": "accounts/secondary-account",
          "location": {
            "displayName": "Default",
            "description": "This is my default location"
          }
        }
      """
     When calling the "chacerapp.v1.Locations/CreateLocation" RPC
     Then I will receive a successful response
    Given a JSON "chacerapp.v1.CreateLocationRequest"
      """
        {
          "parent": "accounts/secondary-account",
          "location": {
            "displayName": "Secondary",
            "description": "This is my secondary location"
          }
        }
      """
     When calling the "chacerapp.v1.Locations/CreateLocation" RPC
     Then I will receive a successful response
    Given a JSON "chacerapp.v1.ListLocationsRequest"
      """
        { "parent": "accounts/default-account", "pageSize": 3 }
      """
     When calling the "chacerapp.v1.Locations/ListLocations" RPC
     Then I will receive a successful response
      And the response value "locations" will have a length of 2
      And the response value "locations[0].name" will be "accounts/default-account/locations/default"
      And the response value "nextPageToken" will be ""
    Given a JSON "chacerapp.v1.ListLocationsRequest"
      """
        { "parent": "accounts/-", "pageSize": 5 }
      """
     When calling the "chacerapp.v1.Locations/ListLocations" RPC
     Then I will receive a successful response
      And the response value "locations" will have a length of 4
      And the response value "locations[0].name" will be "accounts/default-account/locations/default"
      And the response value "nextPageToken" will be ""
    Given a JSON "chacerapp.v1.ListLocationsRequest"
      """
        { "parent": "accounts/-", "pageSize": 3 }
      """
     When calling the "chacerapp.v1.Locations/ListLocations" RPC
     Then I will receive a successful response
      And the response value "locations" will have a length of 3
      And stashing the next page token from the response
    Given a JSON "chacerapp.v1.ListLocationsRequest"
      """
        { "parent": "accounts/-", "pageSize": 3 }
      """
      And using the stashed next page token
     When calling the "chacerapp.v1.Locations/ListLocations" RPC
     Then I will receive a successful response
      And the response value "nextPageToken" will be ""
      And the response value "locations" will have a length of 1
    # Verify that when the URL is changed, the previous next page token cannot be used
    Given a JSON "chacerapp.v1.ListLocationsRequest"
      """
        { "parent": "accounts/-", "pageSize": 3 }
      """
     When calling the "chacerapp.v1.Locations/ListLocations" RPC
     Then I will receive a successful response
      And stashing the next page token from the response
    Given a JSON "chacerapp.v1.ListLocationsRequest"
      """
        { "parent": "accounts/default", "pageSize": 3 }
      """
      And using the stashed next page token
     When calling the "chacerapp.v1.Locations/ListLocations" RPC
     Then I will receive an error with code "INVALID_ARGUMENT"

  Scenario: Verify that all endpoints return a not found error when the account or location do not exist
    Given a JSON "chacerapp.v1.GetLocationRequest"
      """
        { "name": "accounts/non-existant-account/locations/some-location" }
      """
     When calling the "chacerapp.v1.Locations/GetLocation" RPC
     Then I will receive an error with code "NOT_FOUND"
    Given a JSON "chacerapp.v1.CreateLocationRequest"
      """
        {
          "parent": "accounts/non-existant-account",
          "location": {
            "displayName": "Default",
            "description": "This is my default account"
          }
        }
      """
     When calling the "chacerapp.v1.Locations/CreateLocation" RPC
     Then I will receive an error with code "NOT_FOUND"
    Given a JSON "chacerapp.v1.UpdateLocationRequest"
      """
        {
          "location": {
            "name": "accounts/non-existant-account/locations/does-not-exist",
            "displayName": "Default"
          }
        }
      """
     When calling the "chacerapp.v1.Locations/UpdateLocation" RPC
     Then I will receive an error with code "NOT_FOUND"
    Given a JSON "chacerapp.v1.DeleteLocationRequest"
      """
        { "name": "accounts/non-existant-account/locations/does-not-exist" }
      """
     When calling the "chacerapp.v1.Locations/DeleteLocation" RPC
     Then I will receive an error with code "NOT_FOUND"
