// Singleton WASM bridge. Spawns one Worker for the page and exposes a
// typed runCommand(...) that returns a Promise<string>.

import type { Resume } from './schema';

type Pending = {
  resolve: (output: string) => void;
  reject: (err: Error) => void;
};

let workerPromise: Promise<Worker> | null = null;
const pending = new Map<string, Pending>();

function spawnWorker(): Promise<Worker> {
  // Classic (non-module) worker — `importScripts` (used inside the
  // worker to load wasm_exec.js) is only available in classic workers.
  const worker = new Worker(new URL('./wasm.worker.ts', import.meta.url));

  return new Promise<Worker>((resolve, reject) => {
    const ready = (event: MessageEvent) => {
      if (event.data?.type === 'ready') {
        worker.removeEventListener('message', ready);
        resolve(worker);
      }
    };
    worker.addEventListener('message', ready);
    worker.addEventListener('error', (e) => reject(new Error(e.message)));

    worker.addEventListener('message', (event: MessageEvent) => {
      const { type, requestId, output, error } = event.data ?? {};
      if (type !== 'response' || !requestId) return;
      const cb = pending.get(requestId);
      if (!cb) return;
      pending.delete(requestId);
      if (error) cb.reject(new Error(error));
      else cb.resolve(output ?? '');
    });
  });
}

export async function runCommand(command: string): Promise<string> {
  workerPromise ??= spawnWorker();
  const worker = await workerPromise;

  return new Promise<string>((resolve, reject) => {
    const requestId = crypto.randomUUID();
    pending.set(requestId, { resolve, reject });
    worker.postMessage({ type: 'command', requestId, command });
  });
}

/**
 * Fetches the resume.json by running the same kubectl-style command the
 * user would type in the terminal, then parses the wrapping array.
 */
export async function fetchResume(): Promise<Resume> {
  const raw = await runCommand('kubectl get resume main-resume -o json');
  const parsed = JSON.parse(raw) as Resume[] | Resume;
  return Array.isArray(parsed) ? parsed[0] : parsed;
}
