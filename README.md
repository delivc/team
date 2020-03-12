# Team - Account management for APIs

Used mainly to organize "Company" structures within delivc services.

## Endpoints

Team exposes the following endpoints:

**All endpoints except Health requires Authentication with the [@delivc/identity](https://github.com/delivc/identity) service**

* **GET /health**

  Returns the publicly available healthcheck for this service.



  ```json
    {
        "alloc": "1 MiB",
        "cached_items": "1",
        "description": "Team is a management Service for Delivc Teams",
        "garbage_collector_runs": "3",
        "name": "Team",
        "start_time": "2020-03-12 04:50:07.786883 +0100 CET m=+0.005006555",
        "sys": "68 MiB",
        "total_alloc": "5 MiB",
        "version": "dev"
    }
  ```

* **GET /permissions**

  Returns a list of all available Permissions

  ```json
    {
        "permissions": [
            {
                "id": "85bbbd1b-2a68-4241-9245-26ac7ab3f594",
                "name": "account-destroy"
            },
            [...],
        ]
    }
  ```

* **GET /accounts**
  
  Returns a list of available Accounts, if User is SuperAdmin it returns all accounts

  ```json
    {
        "accounts": [
            {
                "id": "b11516e8-1d1d-4c05-82de-6f15c9e4a0cd",
                "aud": "app.delivc.com",
                "name": "Your Test Company",
                "owner_ids": {
                    "0": "1dffa867-718b-4488-b07e-f838ef7b01e4"
                },
                "createdAt": "2020-03-11T07:57:12Z",
                "updatedAt": "2020-03-11T07:57:12Z",
                "roles": [
                    {
                        "account_id": "b11516e8-1d1d-4c05-82de-6f15c9e4a0cd",
                        "id": "54c2fccf-b8be-4a83-98c7-2440de0239bb",
                        "name": "Admin",
                        "createdAt": "2020-03-11T07:57:12Z",
                        "updatedAt": "2020-03-11T07:57:12Z"
                    }
                ]
            }
        ],
        "aud": "app.delivc.com"
    }
  ```

* **GET /accounts/{id}**

    Returns:
    ```json
    {
        "id": "263aa240-8bb1-4f27-8926-a14b16e69936",
        "aud": "app.delivc.com",
        "name": "Awesome Team",
        "owner_ids": {
            "0": "1dffa867-718b-4488-b07e-f838ef7b01e4"
        },
        "createdAt": "2020-03-11T08:57:33Z",
        "updatedAt": "2020-03-11T08:57:33Z",
        "roles": [
            {
                "id": "9e5ba411-364b-4757-ad9f-890b87eeb157",
                "name": "Admin",
                "createdAt": "2020-03-11T08:57:33Z",
                "updatedAt": "2020-03-11T08:57:33Z"
            }
        ],
        "users": [
            {
                "user_id": "1dffa867-718b-4488-b07e-f838ef7b01e4",
                "role_id": "9e5ba411-364b-4757-ad9f-890b87eeb157",
                "invited_by": "00000000-0000-0000-0000-000000000000"
            }
        ]
    }
    ```

* **PUT /accounts/{id}**

  Updates the given Account. Only the Fields are getting updated from request

  This Requests needs one of the following Permissions
  * `account-edit`
  * `isOwner`
  * `isSuperAdmin`

  ```json
    {
        "name": "Awesome Team Updated",
        "billing_name": "Delivc GmbH",
        "billing_email": "me@julian.pro"
    }
  ```

  Returns:
  ```json
    {
        "id": "263aa240-8bb1-4f27-8926-a14b16e69936",
        "aud": "app.delivc.com",
        "name": "Awesome Team Updated",
        "billing_name": "Delivc GmbH",
        "billing_email": "me@julian.pro",
        "owner_ids": {
            "0": "1dffa867-718b-4488-b07e-f838ef7b01e4"
        },
        "createdAt": "2020-03-11T08:57:33Z",
        "updatedAt": "2020-03-12T06:17:47.072033+01:00",
        "roles": [
            {
                "id": "9e5ba411-364b-4757-ad9f-890b87eeb157",
                "name": "Admin",
                "createdAt": "2020-03-11T08:57:33Z",
                "updatedAt": "2020-03-11T08:57:33Z"
            }
        ],
        "users": [
            {
                "user_id": "1dffa867-718b-4488-b07e-f838ef7b01e4",
                "role_id": "9e5ba411-364b-4757-ad9f-890b87eeb157",
                "invited_by": "00000000-0000-0000-0000-000000000000"
            }
        ]
    }
  ```

* **POST /accounts**
  
  Create a new account by given Name
  ```json
  {
    "name": "Awesome Team"
  }
  ```

  Returns:

    ```json
    {
        "id": "5989f172-c967-4969-84bd-25c00932daa3",
        "aud": "app.delivc.com",
        "name": "Awesome Team",
        "owner_ids": {
            "0": "1dffa867-718b-4488-b07e-f838ef7b01e4"
        },
        "createdAt": "2020-03-11T09:00:31.799286+01:00",
        "updatedAt": "2020-03-11T09:00:31.79929+01:00",
        "roles": [
            {
                "account_id": "5989f172-c967-4969-84bd-25c00932daa3",
                "id": "68bc6f1f-0083-4aed-bf77-1076fd0155b6",
                "name": "Admin",
                "createdAt": "2020-03-11T09:00:31.802765+01:00",
                "updatedAt": "2020-03-11T09:00:31.802768+01:00",
                "permissions": [
                    {
                        "id": "00a15670-7d4f-4a50-ba44-0578707fe123",
                        "name": "spaces-create"
                    },
                    {
                        "id": "088a16c0-2e8f-435c-b33d-5bde70775130",
                        "name": "spaces-read-apikeys"
                    },
                    [...]
                ]
            }
        ]
    }
    ```

* **DELETE /accounts/{id}**
  
  Delete given Account. User MUST be SuperAdmin or Owner of given Account

  Returns
  ```json
  {}
  ```

* **GET /accounts/{id}/role**
  
  Returns a list of related Roles.

  ```json
    {
        "roles": [
            {
                "id": "9e5ba411-364b-4757-ad9f-890b87eeb157",
                "name": "Admin",
                "createdAt": "2020-03-11T08:57:33Z",
                "updatedAt": "2020-03-11T08:57:33Z",
                "permissions": [
                    {
                        "id": "00a15670-7d4f-4a50-ba44-0578707fe123",
                        "name": "spaces-create"
                    },
                    {
                        "id": "088a16c0-2e8f-435c-b33d-5bde70775130",
                        "name": "spaces-read-apikeys"
                    },
                    [...]
                ]
            }
        ]
    }
  ```

* **GET /accounts/{id}/role/{roleID}**

  Returns given Role by ID if exists

  ```json
    {
        "id": "9e5ba411-364b-4757-ad9f-890b87eeb157",
        "name": "Admin",
        "createdAt": "2020-03-11T08:57:33Z",
        "updatedAt": "2020-03-11T08:57:33Z",
        "permissions": [
            {
                "id": "00a15670-7d4f-4a50-ba44-0578707fe123",
                "name": "spaces-create"
            },
            [...]
        ]
    }
  ```

* **POST /accounts/{id}/role**

  Creates a new Role, with Permissions if defined
  Accepts:
  ```json
    {
	    "name": "MyNewRole",
	    "permissions": ["account-edit", "does-not-exists", "account-destroy"]
    }
  ```

  Returns:
  ```json
    {
        "id": "97a2588b-69ea-4006-b7f4-d0d6d84870e8",
        "name": "MyNewRole",
        "createdAt": "2020-03-12T08:59:36.141506+01:00",
        "updatedAt": "2020-03-12T08:59:36.141509+01:00",
        "permissions": [
            {
                "id": "85bbbd1b-2a68-4241-9245-26ac7ab3f594",
                "name": "account-destroy"
            },
            {
                "id": "904841b0-802f-4dff-ba24-416b4363f7b9",
                "name": "account-edit"
            }
        ]
    }
    ```

* **PUT /accounts/{id}/role/{roleId}**
  
  Updates Role with Permissions
  User MUST be SuperAdmin or Owner or have `account-role-create` permission of given Account

  Accepts:
  ```json
    {
	    "name": "newRoleName123",
	    "permissions": ["account-role-update", "account-role-destroy"]
    }
  ```

  Returns
  ```json
    {
        "id": "f76912cf-a5e2-4faa-a80e-763250194620",
        "name": "newRoleName123",
        "createdAt": "2020-03-12T07:57:52Z",
        "updatedAt": "2020-03-12T10:55:42.427786+01:00",
        "permissions": [
            {
                "id": "ddd4bb5e-0b82-45e3-ba00-a8f8de459e3b",
                "name": "account-role-destroy"
            },
            {
                "id": "f61528a1-7396-4a6d-ae03-5d0d7153cb63",
                "name": "account-role-update"
            }
        ]
    }
  ```

* **DELETE /accounts/{id}/role/{roleId}**
  
  Delete given Account.
  User MUST be SuperAdmin or Owner or have account-role-destroy permission of given Account

  Returns
  ```json
  {}
  ```