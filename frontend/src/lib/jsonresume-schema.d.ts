// Minimal typings for the official JSON Resume validator, which ships as
// untyped CommonJS. We use the package itself as the source of truth — this
// only describes its surface, it does not restate the schema.
declare module '@jsonresume/schema' {
  export function validate(
    resume: unknown,
    callback: (errors: unknown[] | null, valid: boolean) => void,
  ): void;
  export const schema: {
    properties: Record<string, unknown>;
    [key: string]: unknown;
  };
  export const jobSchema: unknown;
}
