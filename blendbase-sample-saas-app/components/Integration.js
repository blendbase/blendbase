import { useState } from "react";
import { useMutation } from "@apollo/client";
import Toggle from "./Toggle";
import {
  LIST_INTEGRATIONS,
  ENABLE_INTEGRATION_QUERY,
  CONFIGURE_CONSUMER_INTEGRATION_OAUTH_QUERY,
  SET_CONSUMER_INTEGRATION_SECRET_QUERY
} from "./queries";

const SERVICE_CODE_CRM_SALESFORCE = "crm_salesforce";

function CheckIcon() {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
      <path
        fillRule="evenodd"
        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
        clipRule="evenodd"
      />
    </svg>
  );
}

function Integration({ integration, blendbaseClient }) {
  const [clientId, setClientId] = useState();
  const [clientSecret, setClientSecret] = useState();
  const [salesforceInstanceSubdomain, setSalesforceInstanceSubdomain] = useState();
  const [errorMessage, setErrorMessage] = useState();
  const [successMessage, setSuccessMessage] = useState();
  const [clientCredentialsSet, setClientCredentialsSet] = useState(integration.oauth2Metadata?.clientCredentialsSet);
  const [integrationSecret, setIntegrationSecret] = useState();
  const [tokensSet, setTokensSet] = useState(integration.oauth2Metadata?.tokensSet);

  const [enableMutation] = useMutation(ENABLE_INTEGRATION_QUERY, {
    client: blendbaseClient,
    refetchQueries: [
      LIST_INTEGRATIONS, // DocumentNode object parsed with gql
      "ListIntegrations" // Query name
    ]
  });
  const [configureOAuthMutation] = useMutation(CONFIGURE_CONSUMER_INTEGRATION_OAUTH_QUERY, { client: blendbaseClient });
  const [setIntegrationSecretMutation] = useMutation(SET_CONSUMER_INTEGRATION_SECRET_QUERY, {
    client: blendbaseClient
  });

  async function enableIntegration(newValue) {
    const response = await enableMutation({ variables: { serviceCode: integration.serviceCode, enabled: newValue } });
    if (!!response.data && response.data.enableConsumerIntegration) {
      setSuccessMessage("Integration successfully toggled");
    } else {
      setErrorMessage("Error enabling integration");
    }
  }

  async function updateCredentials() {
    cleanMessages();
    if (!clientId || !clientSecret) {
      setErrorMessage("Client ID and Client Secret are required");
      return;
    }

    const response = await configureOAuthMutation({
      variables: {
        consumerIntegrationID: integration.id,
        input: {
          clientID: clientId,
          clientSecret: clientSecret,
          salesforceInstanceSubdomain: salesforceInstanceSubdomain
        }
      }
    });

    if (!!response && !!response.data && response.data.configureConsumerIntegrationOAuth) {
      setSuccessMessage("Credentials updated successfully");
      setClientId("");
      setClientSecret("");
      setSalesforceInstanceSubdomain("");
      setClientCredentialsSet(true);
    } else {
      setErrorMessage("Error updating credentials");
    }
  }

  async function updateIntegrationSecret() {
    cleanMessages();
    if (!integrationSecret) {
      setErrorMessage("Secret cannot be empty");
      return;
    }

    const response = await setIntegrationSecretMutation({
      variables: {
        consumerIntegrationID: integration.id,
        secret: integrationSecret
      }
    });

    if (!!response && !!response.data && response.data.setConsumerIntegrationSecret) {
      setSuccessMessage("Secret was updated successfully");
      setIntegrationSecret("");
    } else {
      setErrorMessage("Error updating secret");
    }
  }

  function onClientIdChanged(event) {
    setClientId(event.target.value);
  }

  function onClientSecretChanged(event) {
    setClientSecret(event.target.value);
  }

  function onSalesforceInstanceSubdomainChanged(event) {
    setSalesforceInstanceSubdomain(event.target.value);
  }

  function cleanMessages() {
    setErrorMessage(null);
    setSuccessMessage(null);
  }

  return (
    <div className="overflow-hidden border-gray-200 bg-white shadow sm:rounded-lg">
      <div className="px-4 py-5 sm:px-6">
        <h3 className="text-lg font-medium leading-6 text-gray-900">{integration.serviceName}</h3>
        <p className="mt-1 max-w-2xl text-sm text-gray-500">{integration.description}</p>
      </div>
      {!!errorMessage && (
        <div className="relative bg-red-500  px-6 py-2 text-sm text-white" role="alert">
          {errorMessage}
        </div>
      )}

      {!!successMessage && (
        <div className="relative bg-green-100  px-6 py-2 text-sm text-green-800" role="alert">
          {successMessage}
        </div>
      )}
      <div className="border-t border-gray-200">
        <dl>
          <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500">Enabled</dt>
            <dd className="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
              <Toggle enabled={integration.enabled} setEnabled={enableIntegration} />
            </dd>
          </div>
          <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500">Integration ID</dt>
            <dd className="mt-1 text-sm text-gray-400 sm:col-span-2 sm:mt-0">{integration.id}</dd>
          </div>
          <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500">Integration Type</dt>
            <dd className="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">{integration.type}</dd>
          </div>
          <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500">Authentication Type</dt>
            <dd className="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">{integration.authType}</dd>
          </div>

          {integration.enabled && integration.authType == "secret" && (
            <div className="bg-gray-50 px-4 py-5 text-sm sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
              <dt className="font-medium text-gray-500">Integration Secret</dt>
              <dd className="mt-1 sm:col-span-2 sm:mt-0">
                <div className="space-y-3">
                  <input
                    value={clientId}
                    onChange={(event) => setIntegrationSecret(event.target.value)}
                    type="text"
                    placeholder="Secret"
                    name="integration-secret"
                    id="integration-secret"
                    className="block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                  />

                  <button
                    className="focus:shadow-outline-indigo-500 inline-flex items-center rounded-md border border-transparent bg-indigo-100 px-3 py-1.5 text-sm font-medium leading-5 text-indigo-600 transition duration-150 ease-in-out hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 active:bg-indigo-200"
                    onClick={updateIntegrationSecret}
                  >
                    Update secret
                  </button>
                </div>
              </dd>
            </div>
          )}

          {integration.enabled && integration.authType == "oauth2" && (
            <div className="bg-gray-50 px-4 py-5 text-sm sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
              <dt className="font-medium text-gray-500">Client Credentials</dt>
              <dd className="mt-1 sm:col-span-2 sm:mt-0">
                {clientCredentialsSet && (
                  <div className="flex space-x-2">
                    <div className="flex-0 text-green-600">
                      <CheckIcon />
                    </div>
                    <div className="flex-auto">
                      <button
                        className="text-sm text-indigo-500 underline underline-offset-1 hover:text-indigo-900"
                        onClick={() => setClientCredentialsSet(false)}
                      >
                        Update credentials
                      </button>
                    </div>
                  </div>
                )}
                {!clientCredentialsSet && (
                  <div className="space-y-3">
                    <div>
                      <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                        Callback URL{" "}
                        <span className="text-xs text-gray-400">
                          (copy and paste into the callback URL field in the CRM)
                        </span>
                      </label>
                      <input
                        value={integration.callbackURL}
                        type="text"
                        disabled
                        placeholder="Callback URL"
                        name="callback-url"
                        id="callback-url"
                        className="block w-full rounded-md border-gray-500 px-3 py-2 text-gray-400 shadow-sm"
                      />
                    </div>
                    <input
                      value={clientId}
                      onChange={onClientIdChanged}
                      type="text"
                      placeholder="Client ID"
                      name="client-id"
                      id="client-id"
                      className="block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                    />
                    <input
                      value={clientSecret}
                      onChange={onClientSecretChanged}
                      type="password"
                      placeholder="Client Secret"
                      name="client-secret"
                      id="client-secret"
                      className="block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                    />
                    {integration.serviceCode == SERVICE_CODE_CRM_SALESFORCE && (
                      <input
                        value={salesforceInstanceSubdomain}
                        onChange={onSalesforceInstanceSubdomainChanged}
                        type="text"
                        placeholder="Instance Subdomain"
                        name="salesforce-instance-subdomain"
                        id="salesforce-instance-subdomain"
                        className="block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                      />
                    )}

                    <button
                      className="focus:shadow-outline-indigo-500 inline-flex items-center rounded-md border border-transparent bg-indigo-100 px-3 py-1.5 text-sm font-medium leading-5 text-indigo-600 transition duration-150 ease-in-out hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 active:bg-indigo-200"
                      onClick={updateCredentials}
                    >
                      Update credentials
                    </button>
                  </div>
                )}
              </dd>
            </div>
          )}

          {integration.enabled && integration.authType == "oauth2" && (
            <div className="bg-white px-4 py-5 text-sm sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
              <dt className="font-medium text-gray-500">Access Token</dt>
              <dd className="mt-1 sm:col-span-2 sm:mt-0">
                <div className="flex space-x-2">
                  {tokensSet && (
                    <div className="text-green-600">
                      <CheckIcon />
                    </div>
                  )}
                  <div className="flex-auto">
                    <a
                      className="text-sm text-indigo-500 underline underline-offset-1 hover:text-indigo-900"
                      href={integration.loginURL}
                    >
                      Refresh access token
                    </a>
                  </div>
                </div>
              </dd>
            </div>
          )}
        </dl>
      </div>
    </div>
  );
}

export default Integration;
