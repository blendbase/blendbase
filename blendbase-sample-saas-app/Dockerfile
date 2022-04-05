FROM node:lts as dependencies
WORKDIR /connect-app
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

FROM node:lts as builder
WORKDIR /connect-app
COPY . .
COPY --from=dependencies /connect-app/node_modules ./node_modules
RUN yarn build

FROM node:lts as runner
WORKDIR /connect-app
ENV NODE_ENV production
# If you are using a custom next.config.js file, uncomment this line.
# COPY --from=builder /my-project/next.config.js ./
COPY --from=builder /connect-app/public ./public
COPY --from=builder /connect-app/.next ./.next
COPY --from=builder /connect-app/node_modules ./node_modules
COPY --from=builder /connect-app/package.json ./package.json

EXPOSE 3000
CMD ["yarn", "start"]
