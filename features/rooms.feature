Feature: Manage rooms on the API
  Background: Create accounts and locations we can create rooms in
    Given data loaded from the seed file "seed-data/rooms-background.json"

  Scenario: Able to create, get, update, and delete a room in a location on an account
    Given a JSON "chacerapp.v1.CreateRoomRequest"
      """
        {
          "parent": "accounts/default/locations/default",
          "room": {
            "displayName": "Front Desk",
            "description": "This is my description"
          },
          "room_id": "front-desk"
        }
      """
     When calling the "chacerapp.v1.Rooms/CreateRoom" RPC
     Then I will receive a successful response
     And the response value "name" will be "accounts/default/locations/default/rooms/front-desk"
     And the response value "displayName" will be "Front Desk"
     And the response value "description" will be "This is my description"
     When calling the "chacerapp.v1.Rooms/CreateRoom" RPC
     Then I will receive an error with code "ALREADY_EXISTS"
    Given a JSON "chacerapp.v1.GetRoomRequest"
      """
        { "name": "accounts/default/locations/default/rooms/front-desk" }
      """
     When calling the "chacerapp.v1.Rooms/GetRoom" RPC
     Then I will receive a successful response
     And the response value "displayName" will be "Front Desk"
     And the response value "description" will be "This is my description"
    Given a JSON "chacerapp.v1.UpdateRoomRequest"
      """
        {
          "room": {
            "name": "accounts/default/locations/default/rooms/front-desk",
            "displayName": "Default #2",
            "description": "This is my updated description"
          }
        }
      """
     When calling the "chacerapp.v1.Rooms/UpdateRoom" RPC
     Then I will receive a successful response
     And the response value "displayName" will be "Default #2"
     And the response value "description" will be "This is my updated description"
    Given a JSON "chacerapp.v1.DeleteRoomRequest"
      """
        { "name": "accounts/default/locations/default/rooms/front-desk" }
      """
     When calling the "chacerapp.v1.Rooms/DeleteRoom" RPC
     Then I will receive a successful response
    Given a JSON "chacerapp.v1.GetRoomRequest"
      """
        { "name": "accounts/default/locations/default/rooms/front-desk" }
      """
     When calling the "chacerapp.v1.Rooms/GetRoom" RPC
     Then I will receive an error with code "NOT_FOUND"

  Scenario: Verify that a partial update with a field mask will work correctly
    Given a JSON "chacerapp.v1.CreateRoomRequest"
      """
        {
          "parent": "accounts/default/locations/default",
          "room": {
            "displayName": "Default",
            "description": "This is my default room"
          },
          "room_id": "default"
        }
      """
     And calling the "chacerapp.v1.Rooms/CreateRoom" RPC
     And I will receive a successful response
     And a JSON "chacerapp.v1.UpdateRoomRequest"
      """
        {
          "room": {
            "name": "accounts/default/locations/default/rooms/default",
            "displayName": "Default #2"
          },
          "updateMask": {
            "paths": [ "displayName" ]
          }
        }
      """
     When calling the "chacerapp.v1.Rooms/UpdateRoom" RPC
     Then I will receive a successful response
     And the response value "displayName" will be "Default #2"
      # Verify the description remains unchanged
     And the response value "description" will be "This is my default room"

  Scenario: Able to correctly list out the rooms for a specific account
    # Create several rooms across the two accounts and two locations
    Given a JSON "chacerapp.v1.ListRoomsRequest"
      """
        { "parent": "accounts/default/locations/default", "pageSize": 3 }
      """
     And data loaded from the seed file "seed-data/rooms-list.json"
     When calling the "chacerapp.v1.Rooms/ListRooms" RPC
     Then I will receive a successful response
     And the response value "rooms" will have a length of 2
     And the response value "rooms[0].name" will be "accounts/default/locations/default/rooms/default"
     And the response value "nextPageToken" will be ""

  Scenario: Able to list out the Rooms for a location without specifying the account
    Given a JSON "chacerapp.v1.ListRoomsRequest"
      """
        { "parent": "accounts/-/locations/default", "pageSize": 5 }
      """
     And data loaded from the seed file "seed-data/rooms-list.json"
     When calling the "chacerapp.v1.Rooms/ListRooms" RPC
     Then I will receive a successful response
     And the response value "rooms" will have a length of 4
     And the response value "rooms[0].name" will be "accounts/default/locations/default/rooms/default"
     And the response value "nextPageToken" will be ""

  Scenario: Able to paginate the list of rooms across accounts
    Given a JSON "chacerapp.v1.ListRoomsRequest"
      """
        { "parent": "accounts/-/locations/default", "pageSize": 3 }
      """
     And data loaded from the seed file "seed-data/rooms-list.json"
     When calling the "chacerapp.v1.Rooms/ListRooms" RPC
     Then I will receive a successful response
     And the response value "rooms" will have a length of 3
     And stashing the next page token from the response
    Given a JSON "chacerapp.v1.ListRoomsRequest"
      """
        { "parent": "accounts/-/locations/default", "pageSize": 3 }
      """
     And using the stashed next page token
     When calling the "chacerapp.v1.Rooms/ListRooms" RPC
     Then I will receive a successful response
     And the response value "nextPageToken" will be ""
     And the response value "rooms" will have a length of 1

  Scenario: Able to paginate the list of rooms across all accounts
    Given a JSON "chacerapp.v1.ListRoomsRequest"
      """
        { "parent": "accounts/-/locations/-", "pageSize": 3 }
      """
     And data loaded from the seed file "seed-data/rooms-list.json"
     When calling the "chacerapp.v1.Rooms/ListRooms" RPC
     Then I will receive a successful response
     And the response value "rooms" will have a length of 3
     And stashing the next page token from the response
    Given a JSON "chacerapp.v1.ListRoomsRequest"
      """
        { "parent": "accounts/-/locations/-", "pageSize": 3 }
      """
     And using the stashed next page token
     When calling the "chacerapp.v1.Rooms/ListRooms" RPC
     Then I will receive a successful response
     And the response value "rooms" will have a length of 3
     And stashing the next page token from the response
    Given a JSON "chacerapp.v1.ListRoomsRequest"
      """
        { "parent": "accounts/-/locations/-", "pageSize": 3 }
      """
     And using the stashed next page token
     When calling the "chacerapp.v1.Rooms/ListRooms" RPC
     Then I will receive a successful response
     And the response value "rooms" will have a length of 2
     And the response value "nextPageToken" will be ""

  Scenario: Verify that changing the URL and using a previous page token will fail
    Given a JSON "chacerapp.v1.ListRoomsRequest"
      """
        { "parent": "accounts/-/locations/-", "pageSize": 3 }
      """
     And data loaded from the seed file "seed-data/rooms-list.json"
     When calling the "chacerapp.v1.Rooms/ListRooms" RPC
     Then I will receive a successful response
     And the response value "rooms" will have a length of 3
     And stashing the next page token from the response
    Given a JSON "chacerapp.v1.ListRoomsRequest"
      """
        { "parent": "accounts/-/locations/default", "pageSize": 3 }
      """
     And using the stashed next page token
     When calling the "chacerapp.v1.Rooms/ListRooms" RPC
     Then I will receive an error with code "INVALID_ARGUMENT"

  Scenario: Verify that all endpoints return a not found error when the account, location, or room do not exist
    Given a JSON "chacerapp.v1.GetRoomRequest"
      """
        { "name": "accounts/non-existant/locations/some-location/rooms/non-existant" }
      """
     When calling the "chacerapp.v1.Rooms/GetRoom" RPC
     Then I will receive an error with code "NOT_FOUND"
    Given a JSON "chacerapp.v1.CreateRoomRequest"
      """
        {
          "parent": "accounts/non-existant/locations/does-not-exist",
          "room": {
            "displayName": "Default",
            "description": "This is my default account"
          }
        }
      """
     When calling the "chacerapp.v1.Rooms/CreateRoom" RPC
     Then I will receive an error with code "NOT_FOUND"
    Given a JSON "chacerapp.v1.UpdateRoomRequest"
      """
        {
          "room": {
            "name": "accounts/non-existant/locations/does-not-exist/rooms/does-not-exist",
            "displayName": "Default"
          }
        }
      """
     When calling the "chacerapp.v1.Rooms/UpdateRoom" RPC
     Then I will receive an error with code "NOT_FOUND"
    Given a JSON "chacerapp.v1.DeleteRoomRequest"
      """
        { "name": "accounts/non-existant/locations/default/rooms/does-not-exist" }
      """
     When calling the "chacerapp.v1.Rooms/DeleteRoom" RPC
     Then I will receive an error with code "NOT_FOUND"
