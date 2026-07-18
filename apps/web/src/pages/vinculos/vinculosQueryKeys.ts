export const vinculosQueryKeys = {
  queue: (installationId: string, filters: Record<string, unknown> = {}) =>
    ["product-links", "queue", { installation_id: installationId, filters }] as const,
  resolved: (installationId: string) =>
    ["product-links", "resolved", { installation_id: installationId }] as const,
};
