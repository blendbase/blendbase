import jwt from "jsonwebtoken";

export const CreateToken = (consumerId = null) => {
  const currentTime = Math.floor(Date.now() / 1000);
  var claim = {
    iat: currentTime,
    exp: currentTime + 60 * 60
  };

  if (consumerId !== null) {
    claim.consumer_id = consumerId;
  }

  const token = jwt.sign(claim, process.env.BLENDBASE_AUTH_SECRET, {
    algorithm: "HS256"
  });

  return token;
};
