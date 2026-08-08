import { spawn } from "node:child_process";
import { writeFileSync } from "node:fs";

const [mode, ...args] = process.argv.slice(2);

switch (mode) {
  case "inspect":
    process.stdout.write(
      JSON.stringify({
        args,
        cwd: process.cwd(),
        env: Object.fromEntries(Object.entries(process.env).sort(([a], [b]) => a.localeCompare(b))),
      }),
    );
    break;
  case "streams": {
    const size = Number(args[0] ?? 0);
    process.stdout.write("o".repeat(size));
    process.stderr.write("e".repeat(size));
    break;
  }
  case "exit":
    process.exit(Number(args[0] ?? 0));
    break;
  case "redact": {
    const candidate = process.env.HARNESS_TEST_CANDIDATE ?? "";
    const encoded = encodeURIComponent(candidate);
    process.stdout.write(
      `${candidate}\nKEY=${candidate}\nprefix-${candidate}-suffix\npostgresql://user:${encoded}@host/db\n`,
    );
    process.stderr.write(`SECRET=${candidate}\nhttps://user:${encoded}@host/path\n`);
    break;
  }
  case "tree": {
    const marker = args[0];
    const childScript = `setTimeout(() => require("node:fs").writeFileSync(${JSON.stringify(marker)}, "late"), 2000); setTimeout(() => {}, 20000);`;
    spawn(process.execPath, ["-e", childScript], { stdio: "ignore" });
    setTimeout(() => {}, 20000);
    break;
  }
  case "marker":
    writeFileSync(args[0], "started");
    break;
  default:
    process.stderr.write("unsupported probe mode");
    process.exit(64);
}
