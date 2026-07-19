// Loads @testing-library/jest-dom's matcher augmentations for vitest's
// `expect` under tsc. Runtime registration lives in vitest.config.ts
// setupFiles; this ambient reference only supplies the types so
// `toBeInTheDocument`, `toHaveTextContent`, etc. resolve during typecheck.
/// <reference types="@testing-library/jest-dom/vitest" />
