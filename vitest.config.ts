import { defineConfig } from "vitest/config";
import path from "path";

// The single vitest entry point for the repository (GATE-TOPOLOGY §2.3):
// discovery is glob-based over every workspace, with no filename pins, so a
// new test file is reachable the moment it exists. Two projects because the
// environments genuinely differ -- sdk-runtime is a node library, everything
// else renders DOM. The gate's test-web lane runs this config and its census
// guard compares the executed file count against the tree.
export default defineConfig({
  test: {
    projects: [
      {
        // root pinned to the repository so the include globs mean the same
        // thing from any cwd -- workspace `npm run test` scripts included.
        root: __dirname,
        esbuild: {
          jsx: "automatic",
        },
        test: {
          name: "web",
          environment: "jsdom",
          globals: true,
          setupFiles: [
            path.resolve(__dirname, "node_modules/@testing-library/jest-dom/dist/vitest.mjs"),
          ],
          include: ["apps/web/src/**/*.test.{ts,tsx}", "packages/*/src/**/*.test.{ts,tsx}"],
          exclude: ["**/node_modules/**", "**/.worktrees/**", "packages/sdk-runtime/**"],
        },
      },
      {
        root: __dirname,
        test: {
          name: "sdk-runtime",
          environment: "node",
          globals: true,
          include: ["packages/sdk-runtime/src/**/*.test.ts"],
          exclude: ["**/node_modules/**", "**/.worktrees/**"],
        },
      },
    ],
  },
});
