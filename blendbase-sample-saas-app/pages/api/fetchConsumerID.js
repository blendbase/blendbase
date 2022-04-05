import { CREATE_CONSUMER_QUERY } from "../../components/queries";
import { NewBlendbaseClient } from "../../utils/blendbaseClient";
import { CreateToken } from "../../utils/jwt";

export default async function handler(req, res) {
  const token = CreateToken();
  const client = NewBlendbaseClient(token);

  const result = await client.mutate({
    mutation: CREATE_CONSUMER_QUERY
  });

  const consumerID = result.data.createConsumer;

  res
    .status(200)
    .send(
      `Created new consumer for you.\nExecute the next command in console to update consumer ID:\n\necho \"CONSUMER_ID=${consumerID}" >> .env\n\n`
    );
}
