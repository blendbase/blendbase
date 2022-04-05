import { useMemo } from "react";
import { useQuery } from "@apollo/client";

import Layout from "../components/Layout";
import { LIST_CONTACTS } from "../components/queries";
import { NewBlendbaseClient } from "../utils/blendbaseClient";
import { CreateToken } from "../utils/jwt";

export async function getServerSideProps() {
  const token = CreateToken(process.env.CONSUMER_ID);

  return {
    props: {
      blendbaseJwtToken: token
    }
  };
}

function Contacts({ blendbaseJwtToken }) {
  const blendbaseClient = useMemo(() => NewBlendbaseClient(blendbaseJwtToken), [blendbaseJwtToken]);

  const { loading, error, data } = useQuery(LIST_CONTACTS, {
    client: blendbaseClient
  });

  if (loading) return null;
  if (error) return `Error fetching data from Blendbase: ${error.message}`;

  const edges = data.crm.contacts.edges;

  return (
    <Layout>
      <div className="App">
        <div className="mx-auto w-full p-4">
          <div className="sm:flex sm:items-center">
            <div className="sm:flex-auto">
              <h1 className="text-3xl font-semibold text-slate-800">CRM: Contacts</h1>
              <p className="mt-2 text-sm text-gray-700">
                A list of all the contacts in your account including their name, title and email.
              </p>
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
          <div className="-mx-4 mt-8 overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:-mx-6 md:mx-0 md:rounded-lg">
            <table className="min-w-full divide-y divide-gray-300">
              <thead className="bg-gray-50">
                <tr>
                  <th scope="col" className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6">
                    ID
                  </th>
                  <th scope="col" className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6">
                    Company
                  </th>
                  <th
                    scope="col"
                    className="hidden px-3 py-3.5 text-left text-sm font-semibold text-gray-900 sm:table-cell"
                  >
                    Name
                  </th>
                  <th
                    scope="col"
                    className="hidden px-3 py-3.5 text-left text-sm font-semibold text-gray-900 lg:table-cell"
                  >
                    Email
                  </th>
                  <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                    Phone
                  </th>
                  <th scope="col" className="relative py-3.5 pl-3 pr-4 sm:pr-6">
                    <span className="sr-only">Edit</span>
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 bg-white">
                {edges.map((edge) => (
                  <tr key={edge.node.email}>
                    <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6">
                      {edge.node.id}
                    </td>
                    <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6">
                      {edge.node.companyName}
                    </td>
                    <td className="hidden whitespace-nowrap px-3 py-4 text-sm text-gray-500 sm:table-cell">
                      {edge.node.name}
                    </td>
                    <td className="hidden whitespace-nowrap px-3 py-4 text-sm text-gray-500 lg:table-cell">
                      {edge.node.email}
                    </td>
                    <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{edge.node.phone}</td>
                    <td className="whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                      {/* <a href="#" className="text-indigo-600 hover:text-indigo-900">
                        Edit
                      </a> */}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </Layout>
  );
}

export default Contacts;
