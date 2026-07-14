#!/usr/bin/env node

import { spawn } from "node:child_process";
import { mkdtemp, readFile, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { createInterface } from "node:readline";

const SERVER_NAME = "marketplace-central-codex-bridge";
const SERVER_VERSION = "0.1.0";
const DEFAULT_TIMEOUT_MS = 10 * 60 * 1000;

function parseArgs(argv) {
  const options = {
    taskId: process.env.CODEX_PORTFOLIO_TASK_ID ?? "",
    timeoutMs: Number(process.env.CODEX_BRIDGE_TIMEOUT_MS ?? DEFAULT_TIMEOUT_MS),
    check: false,
  };

  for (let index = 0; index < argv.length; index += 1) {
    const argument = argv[index];
    if (argument === "--task-id") options.taskId = argv[++index] ?? "";
    else if (argument === "--timeout-ms") options.timeoutMs = Number(argv[++index]);
    else if (argument === "--check") options.check = true;
    else throw new Error(`Unknown argument: ${argument}`);
  }

  if (!/^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$/i.test(options.taskId)) {
    throw new Error("A valid Codex task ID is required via --task-id or CODEX_PORTFOLIO_TASK_ID.");
  }
  if (!Number.isFinite(options.timeoutMs) || options.timeoutMs < 1_000) {
    throw new Error("--timeout-ms must be at least 1000.");
  }
  return options;
}

function runProcess(command, args, { input, timeoutMs }) {
  return new Promise((resolve, reject) => {
    const child = spawn(command, args, {
      cwd: process.cwd(),
      env: process.env,
      shell: false,
      stdio: ["pipe", "pipe", "pipe"],
      windowsHide: true,
    });
    let stdout = "";
    let stderr = "";
    let settled = false;

    const timer = setTimeout(async () => {
      if (settled) return;
      settled = true;
      if (process.platform === "win32" && child.pid) {
        await new Promise((done) => {
          const killer = spawn("taskkill", ["/PID", String(child.pid), "/T", "/F"], {
            windowsHide: true,
            stdio: "ignore",
          });
          killer.on("error", done);
          killer.on("close", done);
        });
      } else {
        child.kill("SIGTERM");
      }
      reject(new Error(`Codex timed out after ${timeoutMs} ms.`));
    }, timeoutMs);

    child.stdout.setEncoding("utf8");
    child.stderr.setEncoding("utf8");
    child.stdout.on("data", (chunk) => { stdout += chunk; });
    child.stderr.on("data", (chunk) => { stderr += chunk; });
    child.on("error", (error) => {
      if (settled) return;
      settled = true;
      clearTimeout(timer);
      reject(error);
    });
    child.on("close", (code) => {
      if (settled) return;
      settled = true;
      clearTimeout(timer);
      if (code === 0) resolve({ stdout, stderr });
      else reject(new Error(`Codex exited with code ${code}: ${stderr.trim() || stdout.trim()}`));
    });
    child.stdin.end(input ?? "");
  });
}

async function askCodex(taskId, message, timeoutMs) {
  const directory = await mkdtemp(join(tmpdir(), "codex-bridge-"));
  const outputPath = join(directory, "last-message.txt");
  const codexCommand = process.env.CODEX_COMMAND || "codex";
  const prompt = [
    "External advisory request from Claude Code.",
    "Respond with analysis or guidance only unless the message explicitly authorizes repository changes.",
    "Do not create, fork, or hand off another task.",
    "",
    message,
  ].join("\n");

  try {
    await runProcess(
      codexCommand,
      [
        "-s", "read-only",
        "-a", "never",
        "exec", "resume",
        "--all",
        "--output-last-message", outputPath,
        taskId,
        "-",
      ],
      { input: prompt, timeoutMs },
    );
    const response = (await readFile(outputPath, "utf8")).trim();
    if (!response) throw new Error("Codex completed without a final message.");
    return response;
  } finally {
    await rm(directory, { recursive: true, force: true });
  }
}

function writeResponse(payload) {
  process.stdout.write(`${JSON.stringify(payload)}\n`);
}

function toolDefinition(taskId) {
  return {
    name: "ask_codex_portfolio",
    description: `Send one advisory message to the existing Codex Portfolio task ${taskId} and return its final response.`,
    inputSchema: {
      type: "object",
      properties: {
        message: { type: "string", description: "The complete question or handoff for the Portfolio task." },
      },
      required: ["message"],
      additionalProperties: false,
    },
  };
}

async function handleRequest(request, options) {
  if (request.method === "initialize") {
    return {
      protocolVersion: request.params?.protocolVersion ?? "2024-11-05",
      capabilities: { tools: {} },
      serverInfo: { name: SERVER_NAME, version: SERVER_VERSION },
    };
  }
  if (request.method === "ping") return {};
  if (request.method === "tools/list") return { tools: [toolDefinition(options.taskId)] };
  if (request.method === "tools/call") {
    if (request.params?.name !== "ask_codex_portfolio") {
      throw new Error(`Unknown tool: ${request.params?.name ?? "<missing>"}`);
    }
    const message = request.params?.arguments?.message;
    if (typeof message !== "string" || message.trim() === "") {
      throw new Error("message must be a non-empty string.");
    }
    const response = await askCodex(options.taskId, message.trim(), options.timeoutMs);
    return { content: [{ type: "text", text: response }] };
  }
  return null;
}

async function check(options) {
  const result = await runProcess(process.env.CODEX_COMMAND || "codex", ["--version"], {
    timeoutMs: 10_000,
  });
  process.stdout.write(`bridge: ok\ntask: ${options.taskId}\n${result.stdout.trim()}\n`);
}

async function main() {
  const options = parseArgs(process.argv.slice(2));
  if (options.check) return check(options);

  const lines = createInterface({ input: process.stdin, crlfDelay: Infinity });
  for await (const line of lines) {
    if (!line.trim()) continue;
    let request;
    try {
      request = JSON.parse(line);
      const result = await handleRequest(request, options);
      if (request.id !== undefined && result !== null) {
        writeResponse({ jsonrpc: "2.0", id: request.id, result });
      }
    } catch (error) {
      if (request?.id !== undefined) {
        writeResponse({
          jsonrpc: "2.0",
          id: request.id,
          error: { code: -32000, message: error instanceof Error ? error.message : String(error) },
        });
      } else {
        process.stderr.write(`${error instanceof Error ? error.message : String(error)}\n`);
      }
    }
  }
}

main().catch((error) => {
  process.stderr.write(`${error instanceof Error ? error.message : String(error)}\n`);
  process.exitCode = 1;
});
