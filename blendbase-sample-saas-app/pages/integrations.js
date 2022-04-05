import { useMemo, useEffect } from "react";
import { useQuery } from "@apollo/client";
import { useRouter } from "next/router";

import Layout from "../components/Layout";
import Integration from "../components/Integration";
import { LIST_INTEGRATIONS } from "../components/queries";
import { NewBlendbaseClient } from "../utils/blendbaseClient";
import { CreateToken } from "../utils/jwt";

import { SuccessCallout, ErrorCallout } from "../components/callouts";

export async function getServerSideProps() {
  return {
    props: {
      blendbaseJwtToken: CreateToken(process.env.CONSUMER_ID)
    }
  };
}

function Integrations({ blendbaseJwtToken }) {
  const { query } = useRouter();
  const blendbaseErrorMessage = query.blendbaseErrorMessage;
  const blendbaseSuccessMessage = query.blendbaseSuccessMessage;

  const blendbaseClient = useMemo(() => NewBlendbaseClient(blendbaseJwtToken), [blendbaseJwtToken]);

  const { loading, error, data } = useQuery(LIST_INTEGRATIONS, {
    client: blendbaseClient
  });

  if (loading) return "Loading...";
  if (error) return `Error fetching data from Blendbase: ${error.message}`;

  const integrations = data.connect.integrations;

  return (
    <Layout>
      <div className="App">
        <div className="mx-auto w-full p-4 sm:px-6">
          <div className="sm:flex sm:items-center">
            <div className="sm:flex-auto">
              <h1 className="text-3xl font-semibold text-slate-800">Settings: integrations</h1>
              <p className="mt-2 text-sm text-gray-700">A list of all available integrations in your account.</p>
            </div>
            <div className="mt-4 sm:mt-0 sm:ml-16 sm:flex-none">
              {/* <button
                type="button"
                className="inline-flex items-center justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:w-auto"
              >
                Add contact
              </button> */}
            </div>
          </div>
          <div className="mt-8 space-y-8">
            {!!blendbaseSuccessMessage && <SuccessCallout message={blendbaseSuccessMessage} />}
            {!!blendbaseErrorMessage && <ErrorCallout message={blendbaseErrorMessage} />}
            {integrations.map((integration) => (
              <Integration integration={integration} key={integration.serviceCode} blendbaseClient={blendbaseClient} />
            ))}
          </div>
        </div>
      </div>
    </Layout>
  );
}

export default Integrations;
