# syntax=docker/dockerfile:1
# Production image for the web SPA: vite build -> static files behind nginx.

FROM node:22-bookworm AS build

WORKDIR /src

COPY package.json package-lock.json ./
COPY apps/web/package.json ./apps/web/package.json
COPY packages/ui/package.json ./packages/ui/package.json
COPY packages/sdk-runtime/package.json ./packages/sdk-runtime/package.json
COPY packages/web-query/package.json ./packages/web-query/package.json
COPY packages/feature-classifications/package.json ./packages/feature-classifications/package.json
COPY packages/feature-inventory/package.json ./packages/feature-inventory/package.json
COPY packages/feature-products/package.json ./packages/feature-products/package.json
COPY packages/feature-simulator/package.json ./packages/feature-simulator/package.json

RUN npm ci

COPY tsconfig.base.json ./
COPY apps/web ./apps/web
COPY packages ./packages

RUN npm run build --workspace @marketplace-central/web


FROM nginx:1.27-alpine

COPY docker/prod/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /src/apps/web/dist /usr/share/nginx/html

EXPOSE 80
