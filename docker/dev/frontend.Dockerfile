FROM node:22-bookworm

WORKDIR /workspace

COPY package.json package-lock.json ./
COPY apps/web/package.json ./apps/web/package.json
COPY packages/ui/package.json ./packages/ui/package.json
COPY packages/sdk-runtime/package.json ./packages/sdk-runtime/package.json
COPY packages/feature-classifications/package.json ./packages/feature-classifications/package.json
COPY packages/feature-integrations/package.json ./packages/feature-integrations/package.json
COPY packages/feature-marketplaces/package.json ./packages/feature-marketplaces/package.json
COPY packages/feature-products/package.json ./packages/feature-products/package.json
COPY packages/feature-simulator/package.json ./packages/feature-simulator/package.json

RUN npm ci

EXPOSE 5174
