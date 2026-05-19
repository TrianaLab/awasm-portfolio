/// <reference lib="webworker" />
// Web Worker that owns the Go WASM runtime. It runs commands sent from
// the main thread and replies with the stdout string, correlated via
// the requestId field so multiple callers can share one worker.

interface GoRuntime {
  importObject: WebAssembly.Imports;
  run(instance: WebAssembly.Instance): Promise<void>;
}

// The wasm_exec.js shim is loaded via importScripts; it attaches Go to
// the worker's global scope and main.go publishes executeCommand later.
declare const Go: new () => GoRuntime;
const workerSelf = self as DedicatedWorkerGlobalScope & {
  executeCommand?: (cmd: string) => string;
};

interface CommandRequest {
  type: 'command';
  requestId: string;
  command: string;
}

interface CommandResponse {
  type: 'response';
  requestId: string;
  output?: string;
  error?: string;
}

interface ReadyMessage {
  type: 'ready';
}

let ready = false;
const queue: CommandRequest[] = [];

async function bootstrap() {
  // wasm_exec.js is shipped at the same origin under /scripts/wasm_exec.js
  importScripts('/scripts/wasm_exec.js');

  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(fetch('/assets/app.wasm'), go.importObject);
  void go.run(result.instance);

  // Wait one tick for Go's main to register js.Global().Set("executeCommand", ...)
  await new Promise<void>((resolve) => setTimeout(resolve, 0));

  ready = true;
  postMessage({ type: 'ready' } satisfies ReadyMessage);
  for (const req of queue.splice(0)) {
    handle(req);
  }
}

function handle(req: CommandRequest) {
  try {
    const output = workerSelf.executeCommand?.(req.command) ?? '';
    postMessage({ type: 'response', requestId: req.requestId, output } satisfies CommandResponse);
  } catch (err) {
    postMessage({
      type: 'response',
      requestId: req.requestId,
      error: err instanceof Error ? err.message : String(err),
    } satisfies CommandResponse);
  }
}

self.addEventListener('message', (event: MessageEvent<CommandRequest>) => {
  if (event.data.type !== 'command') return;
  if (!ready) {
    queue.push(event.data);
    return;
  }
  handle(event.data);
});

void bootstrap();
