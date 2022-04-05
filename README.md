# Overview

Blendbase provides single unified GraphQL API to access CRMs for SaaS solutions builders.
Instead of building and maintaining various integrations with third-party apps (e.g. Salesforce, HubSpot CRM, PipeDrive, etc.) you can use Blendbase CRM API to access all of them through a single unified interface regardless which CRM your users use. Blendbase also manages the complexity of authentication and authorization with CRMs.

## Supported CRMs

- Salesforce
- HubSpot

## Configuring Blendbase

1. Copy .env.sample `cp .env.sample .env`
2. Copy .env.sample `connect-fullstack-webapp-sample/.env.sample connect-fullstack-webapp-sample/.env`
3. Generate a secret used for encrypting sensitive information and store it in `SECRET_ENCRYPTION_KEY` in `.env`: `go run main.go gen-enc-key`
4. Blendbase is using JWT tokens to authenticate both Omni and Connect APIs (both client side and server-side). Generate a secret that is going to be used for encrypting JWT tokens by running `openssl rand -hex 32` and store it in

- `BLENDBASE_AUTH_SECRET` in `~/.env.local`
- `BLENDBASE_AUTH_SECRET` in `~/connect-fullstack-webapp-sample/.env`

5. `docker-compose build`
6. `docker-compose up`
7. Go to http://localhost:3000/ to see the sample Connect app, follow the instructions below on configuring CRM integrations
8. Run `go run main.go gen-auth-token --consumer-id c6a82fd9-7e22-40c2-8bf2-db58a40839a9` to obtain an authentication token (the used consumer ID is preconfigured for test purposes)
9. Go to http://localhost:8080/ and configure HTTP headers (replace $token with the value from the previous step)

```
{
  "Authorization": "Bearer $token"
}
```

10. Test Omni API with the query

```
query {
  crm {
    contacts {
      edges {
        node {
          id
          name
        }
      }
    }
  }
}
```

### Connecting to Salesforce

1. Login to your instances of Salesforce or create a new one at https://login.salesforce.com/ as an **admin**
2. Go to Settings > Setup > Apps > App Manager
3. Click on "New Connected App"
   1. Fill out the name and contact email
   2. Enable "Enable OAuth Settings"
   3. Fill out "Callback URL" with `https://example.com` - we will change it later
   4. In the "Selected OAuth Scopes" select:
      1. Manage user data via APIs (api)
      2. Perform Requests at any time (…)
   5. Click "Save"
4. Set "Consumer Key" and "Consumer Secret" to h at http://localhost:3000/ (`SALESFORCE_CLIENT_ID` and `SALESFORCE_CLIENT_SECRET` in .env for development and testing).

### Connecting to HubSpot

1. Log in to or create a new instance of HubSpot at https://app.hubspot.com/login as **admin**
2. Go to Settings (cog icon in the top-left)
3. Sidebar. Integrations > Private Apps.
4. Click "Create a private app"
5. Give it a name, e.g. "Blendbase app"
6. Switch to "Scopes" in the header and select:
   1. `crm.objects.companies` - Read & Write
   2. `crm.objects.contacts` - Read & Write
   3. `crm.objects.deals` - Read & Write
7. Click "Create app"
8. In app view page navigate to "Access token" section and copy the `Access token`, store it in "Integration Secret" at http://localhost:3000/ (`HUBSPOT_ACCESS_TOKEN` in the .env file for development and tests).

## API

APIs:

- Connect API for managing your consumers and integrations with CRMs
- Omni API for interacting with CRM objects like contacts, notes, deals, etc.

API authentication is done via the Authorization header which should have a JWT token encoded with the value of the `BLENDBASE_AUTH_SECRET` environment variable. The JWT token should have the `cunsomer_id` claim that represents the current `Consumer` on behalf of whom CRM is being called. See [jwt.js](connect-fullstack-webapp-sample/utils/jwt.js) for an example.

# Development

### DB Setup

1. Set up Postgres database

```
createuser -d blendbase
createdb -O blendbase blendbase
```

2. Run `go run main.go db:migrate` to migrate the database
3. Run `go run main.go db:seed` to init DB with test data

### Running the server

1. `go run main.go server`
2. Go to http://localhost:8080/ to use GraphQL Playground
3. Execute the query

### GraphQL Generation

Blendbase is using `gqlgen` to generate the code based on the schema located at `/graph/schema.graphqls`.
After updating the schema make sure to run `go run github.com/99designs/gqlgen generate`

## Sample connect React application

Sample React connect app is used to demonstrate how blendbase can be integrated into a React app.

### Setup

1. `cd connect-fullstack-webapp-sample`
2. `cp .env.sample .env`
3. Make sure to assigned a value for `BLENDBASE_AUTH_SECRET` (see API Authentication section above)
4. `yarn install` - to install the dependencies
5. Run `curl http://localhost:3000/api/fetchConsumerID` in directory of the React app. That will call Blendbase API and will create and new consumer. Follow the instruction in the output of the API call.
6. `yarn run dev` - to run the app

## Testing

> make sure you have `.env.test` file in your root directory. You can copy .env.sample if you are doing this for the first time:

```shell
cp .env.sample .env.test
```

To run all the tests:

```shell
go test -v ./...
```

Run a specific file:

```shell
go test -v blendbase/connectors/salesforce
```

Disable cache:

```shell
go test -count 1 -v ...
```

## License

Blendbase monorepo uses multiple licenses.

The license for a particular work is defined with following prioritized rules:

1. License directly present in the file
2. LICENSE file in the same directory as the work
3. First LICENSE found when exploring parent directories up to the project top level directory
4. Defaults to Elastic License 2.0
