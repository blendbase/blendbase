import { gql } from "@apollo/client";

export const LIST_INTEGRATIONS = gql`
  query {
    connect {
      integrations {
        id
        type
        serviceCode
        serviceName
        enabled
        description
        authType
        loginURL
        callbackURL

        oauth2Metadata {
          clientCredentialsSet
          tokensSet
        }
      }
    }
  }
`;

export const LIST_CONTACTS = gql`
  query {
    crm {
      contacts {
        edges {
          node {
            id
            name
            email
            phone
            companyName
          }
        }
      }
    }
  }
`;

export const ENABLE_INTEGRATION_QUERY = gql`
  mutation enableIntegration($serviceCode: String!, $enabled: Boolean!) {
    enableConsumerIntegration(serviceCode: $serviceCode, enabled: $enabled)
  }
`;

export const CONFIGURE_CONSUMER_INTEGRATION_OAUTH_QUERY = gql`
  mutation addIntegration($consumerIntegrationID: String!, $input: OAuth2ConfigurationInput!) {
    configureConsumerIntegrationOAuth(consumerIntegrationID: $consumerIntegrationID, input: $input)
  }
`;

export const SET_CONSUMER_INTEGRATION_SECRET_QUERY = gql`
  mutation setIntegrationSecret($consumerIntegrationID: String!, $secret: String!) {
    setConsumerIntegrationSecret(consumerIntegrationID: $consumerIntegrationID, secret: $secret)
  }
`;

export const CREATE_CONSUMER_QUERY = gql`
  mutation CreateConsumer {
    createConsumer
  }
`;
