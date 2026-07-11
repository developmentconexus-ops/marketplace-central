import { appendFileSync, existsSync, readFileSync, writeFileSync } from "node:fs";

const [kind, ...args] = process.argv.slice(2);
const logPath = process.env.HARNESS_POSTGRES_PROBE_LOG;

if (!logPath || !["docker", "go"].includes(kind)) {
  process.stderr.write("probe configuration invalid");
  process.exit(64);
}

function operationForDocker(argv) {
  if (argv[0] === "run") return "start";
  if (argv[0] === "port") return "port";
  if (argv[0] === "rm") return "remove";
  if (argv[0] === "ps") return "inventory";
  const text = argv.join(" ");
  if (text.includes("pg_isready")) return "ready";
  if (text.includes("CREATE DATABASE")) return "create";
  if (text.includes("DROP DATABASE")) return "drop";
  return "docker-unknown";
}

function operationForGo(argv) {
  if (argv.includes("test")) return "tests";
  if (argv.includes("run")) return "migrate";
  return "go-unknown";
}

const operation = kind === "docker" ? operationForDocker(args) : operationForGo(args);
const secret = process.env.POSTGRES_PASSWORD ?? "";
appendFileSync(logPath, `${JSON.stringify({
  kind,
  operation,
  args,
  cwd: process.cwd(),
  hasPostgresPassword: secret.length > 0,
  secretInArgs: secret.length > 0 && args.some((arg) => arg.includes(secret)),
})}\n`);

const statePath = `${logPath}.state.json`;
let state = existsSync(statePath)
  ? JSON.parse(readFileSync(statePath, "utf8"))
  : { running: false, container: "", migrateCalls: 0 };

if (operation === "start") {
  const nameIndex = args.indexOf("--name");
  state = {
    running: true,
    container: nameIndex >= 0 ? args[nameIndex + 1] : "probe-container",
    migrateCalls: 0,
  };
  writeFileSync(statePath, JSON.stringify(state));
}

const failures = (process.env.HARNESS_POSTGRES_PROBE_FAIL_OPERATIONS ?? "")
  .split(",")
  .map((value) => value.trim())
  .filter(Boolean);
if (failures.includes(operation)) {
  process.stderr.write(`probe failure operation=${operation}`);
  process.exit(operation === "tests" ? 17 : 29);
}

if (operation === "remove") {
  state.running = false;
  writeFileSync(statePath, JSON.stringify(state));
}
if (operation === "migrate") {
  const applied = state.migrateCalls === 0 ? 32 : 0;
  state.migrateCalls += 1;
  writeFileSync(statePath, JSON.stringify(state));
  process.stdout.write(`applied ${applied} migration(s)\n`);
}
if (operation === "port") process.stdout.write("127.0.0.1:49152\n");
if (operation === "inventory" && state.running) process.stdout.write(`${state.container}\n`);
