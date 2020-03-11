# Team - Account management for APIs

Used mainly to organize "Company" structures within delivc services.

## Endpoints

Team exposes the following endpoints:

* **GET /health**

  Returns the publicly available healthcheck for this service.

  All endpoints except Health requires Authentication with the [@delivc/identity](https://github.com/delivc/identity) service

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