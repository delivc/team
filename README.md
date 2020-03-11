# Team - Account management for APIs

Used mainly to organize "Company" structures within delivc services.

## Endpoints

Team exposes the following endpoints:

**All endpoints except Health requires Authentication with the [@delivc/identity](https://github.com/delivc/identity) service**

* **GET /health**

  Returns the publicly available healthcheck for this service.



  ```json
  {
    "description": "Team is a management Service for Delivc Teams",
    "name": "Team",
    "version": "dev"
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
                        "name": "spaces-create",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "088a16c0-2e8f-435c-b33d-5bde70775130",
                        "name": "spaces-read-apikeys",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "2ae680e6-f421-45ab-8ee9-426afafb96f2",
                        "name": "account-users-remove",
                        "createdAt": "2020-03-10T10:50:53Z",
                        "updatedAt": "2020-03-10T10:50:53Z"
                    },
                    {
                        "id": "2b4cc127-4e17-480c-acf3-2ca4f6099002",
                        "name": "spaces-create-assets",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "2e44dcf4-adbf-4c6c-a7ea-ca1dbf993435",
                        "name": "account-users-invite",
                        "createdAt": "2020-03-10T10:50:53Z",
                        "updatedAt": "2020-03-10T10:50:53Z"
                    },
                    {
                        "id": "3f87f515-45c5-4451-9a1a-afed6f8ccdbc",
                        "name": "spaces-destroy-assets",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "42a1c9ba-4d61-4dd3-bc0c-f8404c94181b",
                        "name": "spaces-create-content",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "47a257f2-6a40-497c-9e4c-089ed6ac3b7e",
                        "name": "spaces-create-apikeys",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "57e09a63-62d2-4708-b3c1-b145eece0f89",
                        "name": "spaces-destroy-models",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "59f7637a-9b78-4571-8e4c-77ba619c64c4",
                        "name": "spaces-destroy-content",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "722fae9f-32fe-4755-852d-df1d6f15a0b0",
                        "name": "spaces-edit-models",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "75d50643-50cd-4a3d-97eb-53f6c2411c36",
                        "name": "spaces-delete",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "85bbbd1b-2a68-4241-9245-26ac7ab3f594",
                        "name": "account-destroy",
                        "createdAt": "2020-03-10T10:51:24Z",
                        "updatedAt": "2020-03-10T10:51:24Z"
                    },
                    {
                        "id": "8a8d26a5-f733-480f-9129-e165b676f98c",
                        "name": "spaces-edit",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "904841b0-802f-4dff-ba24-416b4363f7b9",
                        "name": "account-edit",
                        "createdAt": "2020-03-10T10:50:53Z",
                        "updatedAt": "2020-03-10T10:50:53Z"
                    },
                    {
                        "id": "c1ce5f7e-f541-41ef-b54c-117049fc4f21",
                        "name": "spaces-destroy-apikeys",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "f9d3d8f8-69d0-41f4-ae31-cbdcf6a4ca87",
                        "name": "spaces-edit-content",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    },
                    {
                        "id": "fe280ed7-e5c0-49d3-8cd4-c47281cbd7d4",
                        "name": "spaces-create-models",
                        "createdAt": "2020-03-10T10:18:25Z",
                        "updatedAt": "2020-03-10T10:18:25Z"
                    }
                ]
            }
        ]
    }
    ```
