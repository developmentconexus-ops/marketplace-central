import { appendFileSync, existsSync, readFileSync, readdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { spawnSync } from "node:child_process";

const [kind, ...args] = process.argv.slice(2);
const logPath = process.env.HARNESS_POSTGRES_PROBE_LOG;

if (kind === "real-go") {
  const goPath = process.env.HARNESS_REAL_GO;
  if (!goPath) process.exit(64);
  if (args.includes("test")) process.exit(17);
  const result = spawnSync(goPath, args, {
    cwd: process.cwd(),
    env: process.env,
    stdio: "inherit",
  });
  process.exit(result.status ?? 1);
}

if (!logPath || !["docker", "go"].includes(kind)) {
  process.stderr.write("probe configuration invalid");
  process.exit(64);
}

function operationForDocker(argv) {
  if (argv[0] === "version") return "daemon";
  if (argv[0] === "image" && argv[1] === "inspect") return "image";
  if (argv[0] === "run") return "start";
  if (argv[0] === "inspect") return "ownership";
  if (argv[0] === "port") return "port";
  if (argv[0] === "rm") return "remove";
  if (argv[0] === "ps") {
    const filterIndex = argv.indexOf("--filter");
    const filter = filterIndex >= 0 ? argv[filterIndex + 1] : "";
    if (filter.startsWith("name=")) return state.startedEver ? "inventory-name" : "preflight-name";
    if (filter.startsWith("label="))
      return state.startedEver ? "inventory-label" : "preflight-label";
  }
  const text = argv.join(" ");
  if (text.includes("pg_isready")) return "ready";
  if (text.includes("CREATE DATABASE")) return "create";
  if (text.includes("DROP DATABASE")) return "drop";
  if (text.includes("SELECT 1 FROM pg_database")) return "create-verify";
  if (text.includes("FROM pg_database")) return "drop-verify";
  if (text.includes("FROM pg_stat_activity")) return "held-check";
  return "docker-unknown";
}

function operationForGo(argv) {
  if (argv.includes("test")) return "tests";
  if (argv.includes("run")) return "migrate";
  return "go-unknown";
}

function canonicalMigrationCount() {
  try {
    const count = readdirSync(join(process.cwd(), "migrations"), { withFileTypes: true }).filter(
      (entry) => entry.isFile() && entry.name.endsWith(".sql"),
    ).length;
    if (count > 0) return count;
  } catch {
    // The harness must fail rather than manufacture a migration inventory.
  }
  process.stderr.write("HPG_MIGRATION_INVENTORY_INVALID");
  process.exit(64);
}

const statePath = `${logPath}.state.json`;
let state = existsSync(statePath)
  ? JSON.parse(readFileSync(statePath, "utf8"))
  : {
      running: false,
      startedEver: false,
      container: "",
      runId: "",
      migrateCalls: 0,
      readyCalls: 0,
      databaseExists: false,
    };

const operation = kind === "docker" ? operationForDocker(args) : operationForGo(args);
const secret = process.env.POSTGRES_PASSWORD ?? "";
appendFileSync(
  logPath,
  `${JSON.stringify({
    kind,
    operation,
    args,
    cwd: process.cwd(),
    hasPostgresPassword: secret.length > 0,
    secretInArgs: secret.length > 0 && args.some((arg) => arg.includes(secret)),
  })}\n`,
);

if (operation === "start") {
  const nameIndex = args.indexOf("--name");
  state = {
    running: true,
    startedEver: true,
    container: nameIndex >= 0 ? args[nameIndex + 1] : "probe-container",
    runId: (args[args.indexOf("--label") + 1] ?? "").split("=").at(-1),
    migrateCalls: 0,
    readyCalls: 0,
    databaseExists: false,
  };
  writeFileSync(statePath, JSON.stringify(state));
}

const failures = (process.env.HARNESS_POSTGRES_PROBE_FAIL_OPERATIONS ?? "")
  .split(",")
  .map((value) => value.trim())
  .filter(Boolean);
const readyFailures = Number.parseInt(process.env.HARNESS_POSTGRES_PROBE_READY_FAILURES ?? "0", 10);
if (operation === "ready") {
  state.readyCalls += 1;
  writeFileSync(statePath, JSON.stringify(state));
  const hangMilliseconds = Number.parseInt(
    process.env.HARNESS_POSTGRES_PROBE_READY_HANG_MILLISECONDS ?? "0",
    10,
  );
  if (hangMilliseconds > 0)
    Atomics.wait(new Int32Array(new SharedArrayBuffer(4)), 0, 0, hangMilliseconds);
  if (state.readyCalls <= readyFailures) {
    process.stderr.write(`probe readiness pending attempt=${state.readyCalls}`);
    process.exit(29);
  }
}
if (failures.includes(operation)) {
  const arbitraryOnly = process.env.HARNESS_POSTGRES_PROBE_ARBITRARY_TEST_FAILURE === "1";
  // The shape `go test -v` actually prints when a subtest fails on a constraint,
  // with the same untrusted noise around it the other fixtures carry. It exists so
  // the subtest/at/constraint tokens are proved against real output, not against a
  // string written to match the regex.
  const richFailure = [
    "=== RUN   TestProbeMustFail",
    '    probe_mustfail_test.go:41: rejected as designed: violates check constraint "probe_only_one_open_row"',
    "--- PASS: TestProbeMustFail (0.00s)",
    "=== RUN   TestProbeContract",
    "=== RUN   TestProbeContract/seeded_rows_survive_the_walk",
    "    probe_contract_test.go:290: seed listing: ERROR: duplicate key value violates unique constraint " +
      '"listings_pkey" (SQLSTATE 23505)',
    "    --- FAIL: TestProbeContract/seeded_rows_survive_the_walk (0.01s)",
    "--- FAIL: TestProbeContract (0.02s)",
    "FAIL\tmarketplace-central/apps/server_core/tests/integration\t0.30s",
    "C:\\private\\customer.txt person@example.test",
  ].join("\n");
  process.stderr.write(
    operation === "tests"
      ? arbitraryOnly
        ? "C:\\private\\customer.txt person@example.test"
        : process.env.HARNESS_POSTGRES_PROBE_RICH_TEST_FAILURE === "1"
          ? richFailure
          : "probe failure reason=HPG_TEST_FAILED_SENTINEL C:\\private\\customer.txt person@example.test"
      : `probe failure operation=${operation}`,
  );
  process.exit(operation === "tests" ? 17 : 29);
}

if (operation === "tests") {
  // `go test -v` prints one `=== RUN` per test and one `--- PASS` per result, and
  // the lifecycle counts those lines. A fixture that printed nothing modelled a
  // run that executed no tests -- which is what this probe did before, and what
  // the lifecycle now refuses to call a pass. Zero is reachable on purpose via
  // HARNESS_POSTGRES_PROBE_TEST_COUNT=0, which is the vacuity fixture.
  const count = Number.parseInt(process.env.HARNESS_POSTGRES_PROBE_TEST_COUNT ?? "2", 10);
  // Skips are their own shape: `=== RUN` is printed, then `--- SKIP` instead of
  // `--- PASS`. A suite where every test skipped still exits 0, which is the
  // vacuous green HARNESS-PROFILE.md records as RUN 27 / PASS 1 / SKIP 26.
  const skipped = Number.parseInt(process.env.HARNESS_POSTGRES_PROBE_TEST_SKIPPED ?? "0", 10);
  for (let index = 1; index <= count; index += 1) {
    const verdict = index <= skipped ? "SKIP" : "PASS";
    process.stdout.write(
      `=== RUN   TestProbe${index}\n--- ${verdict}: TestProbe${index} (0.00s)\n`,
    );
  }
  process.stdout.write(
    count > 0
      ? "PASS\nok  \tprobe\t0.01s\n"
      : "testing: warning: no tests to run\nPASS\nok  \tprobe\t0.01s [no tests to run]\n",
  );
}

if (operation === "preflight-name" && process.env.HARNESS_POSTGRES_PROBE_NAME_CONFLICT === "1") {
  process.stdout.write("conflicting-name\n");
}
if (operation === "preflight-label" && process.env.HARNESS_POSTGRES_PROBE_LABEL_CONFLICT === "1") {
  process.stdout.write("conflicting-label\n");
}
if (operation === "ownership") {
  process.stdout.write(`/${state.container}|${state.runId}\n`);
}
if (operation === "create") {
  state.databaseExists = true;
  writeFileSync(statePath, JSON.stringify(state));
}
if (
  operation === "drop" &&
  !failures.includes(operation) &&
  process.env.HARNESS_POSTGRES_PROBE_DROP_REMAINS !== "1"
) {
  state.databaseExists = false;
  writeFileSync(statePath, JSON.stringify(state));
}
if (operation === "create-verify" && state.databaseExists) process.stdout.write("1\n");
if (operation === "drop-verify" && state.databaseExists)
  process.stdout.write("database-still-exists\n");

if (operation === "remove") {
  state.running = false;
  writeFileSync(statePath, JSON.stringify(state));
}
if (operation === "migrate") {
  const applied = state.migrateCalls === 0 ? canonicalMigrationCount() : 0;
  state.migrateCalls += 1;
  writeFileSync(statePath, JSON.stringify(state));
  process.stdout.write(`applied ${applied} migration(s)\n`);
}
if (operation === "port") process.stdout.write("127.0.0.1:49152\n");
if (operation === "held-check") process.stdout.write("1\n");
if (["inventory-name", "inventory-label"].includes(operation) && state.running)
  process.stdout.write(`${state.container}\n`);
