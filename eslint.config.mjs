// ESLint, flat config. The rule set is GATE-TOPOLOGY.md §2.3a and nothing else.
//
// Two properties are deliberate and both are easy to lose by accident:
//
//   1. Correctness only. Not one stylistic rule. Formatting belongs to Prettier,
//      and the moment a linter has an opinion about quotes the config becomes a
//      preference argument instead of a defect finder.
//
//   2. No `recommended` preset. The presets are a moving target across minor
//      versions -- an upgrade silently adds rules, the ratchet baseline jumps,
//      and the number stops meaning "defects we have not paid down yet". Every
//      rule below was chosen against a failure mode named in the audit, so every
//      rule is written out.
//
// The type-aware rules resolve against `tsconfig.eslint.json`, which covers all
// 220 tracked .ts/.tsx files. See that file for why it exists.

import tseslint from "typescript-eslint";
import reactHooks from "eslint-plugin-react-hooks";

export default [
  {
    // ESLint's flat config does not read `.gitignore`, so everything git already
    // ignores has to be named here or it gets linted. Measured 2026-08-08 on the
    // first run: 399 of 734 files walked were parse failures, every one of them
    // from a tree git ignores -- `scripts/.runs/` snapshots (396) and a vendored
    // VS Code extension inside `.gomodcache/` (3). They fail rather than pass
    // because they belong to no tsconfig, which is the same signal the type-aware
    // rules depend on, so leaving them in would drown it.
    //
    // Also generated output: the authority for the SDK and the API client is the
    // generator, and a finding there is a finding about a schema that is not
    // actionable in the file it is reported in.
    ignores: [
      "**/node_modules/**",
      "**/dist/**",
      "**/build/**",
      "**/.next/**",
      "**/coverage/**",
      "**/.vite/**",
      ".gate/**",
      "**/.gomodcache/**",
      "**/.gocache/**",
      "scripts/.runs/**",
      ".worktrees/**",
      ".mnfs/**",
      // tsc emit alongside source, gitignored for the reason recorded in
      // `.gitignore`: Vite resolves `.js`/`.d.ts` before `.tsx`.
      "apps/web/src/**/*.d.ts",
      "packages/*/src/**/*.d.ts",
      "apps/web/src/api/generated/**",
      "packages/sdk-runtime/src/generated/**",
    ],
  },
  {
    files: ["**/*.ts", "**/*.tsx"],
    languageOptions: {
      parser: tseslint.parser,
      parserOptions: {
        project: ["./tsconfig.eslint.json"],
        tsconfigRootDir: import.meta.dirname,
      },
    },
    plugins: {
      "@typescript-eslint": tseslint.plugin,
      "react-hooks": reactHooks,
    },
    rules: {
      // An unawaited promise is a silent failure with no stack. Needs type
      // information, which is why tsc alone does not catch it.
      "@typescript-eslint/no-floating-promises": "error",

      // An async function passed where a void callback is expected -- including
      // React event handlers, where the rejection has nowhere to go.
      "@typescript-eslint/no-misused-promises": "error",

      // Awaiting a non-promise is usually a refactor that lost a call.
      "@typescript-eslint/await-thenable": "error",

      // Same reason as Go's `unused`: dead code keeps a retired concept alive.
      // `argsIgnorePattern` allows the leading-underscore convention already used
      // in this tree for deliberately unused parameters.
      "@typescript-eslint/no-unused-vars": [
        "error",
        {
          argsIgnorePattern: "^_",
          varsIgnorePattern: "^_",
          caughtErrors: "all",
          caughtErrorsIgnorePattern: "^_",
        },
      ],

      // Stale-closure bugs. Not findable by review at reading speed.
      "react-hooks/rules-of-hooks": "error",
      "react-hooks/exhaustive-deps": "error",
    },
  },
];
