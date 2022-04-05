import { ApolloClient, InMemoryCache } from "@apollo/client";

export const NewBlendbaseClient = (jwtToken = "") => {
  return new ApolloClient({
    uri: "http://localhost:8080/omni/query",
    cache: new InMemoryCache(),
    fetchOptions: {
      mode: "no-cors"
    },
    headers: {
      Origin: "http://localhost:3000",
      Authorization: ["Bearer", jwtToken].join(" ")
    }
  });
};
