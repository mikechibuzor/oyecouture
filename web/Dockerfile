# Nuxt storefront. The build stage produces .output and the runtime stage serves it.
FROM node:22-alpine AS build

WORKDIR /app

# Dependencies first so this layer caches between source-only changes.
COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# `npm run build` passes --dotenv ../.env, which does not exist inside the image.
# Call nuxt directly and let compose supply the environment instead.
RUN npx nuxt build

# Runtime carries only the built output, with no source or dev dependencies.
FROM node:22-alpine

WORKDIR /app

COPY --from=build /app/.output ./.output

ENV NODE_ENV=production \
    HOST=0.0.0.0 \
    PORT=3000

USER node
EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]
