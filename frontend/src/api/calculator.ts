export const OPERATIONS = [
  'add',
  'subtract',
  'multiply',
  'divide',
  'power',
  'percentage',
  'sqrt',
] as const;

export type Operation = (typeof OPERATIONS)[number];

/** Operations that only take the first operand. */
export function isUnary(operation: Operation): boolean {
  return operation === 'sqrt';
}

export interface CalculationResult {
  operation: Operation;
  result: number;
}

/** What went wrong, derived from the HTTP status. */
export type ApiErrorKind = 'invalid-input' | 'math' | 'server';

/** The backend answered with a non-2xx status (or an unusable body). */
export class ApiError extends Error {
  readonly kind: ApiErrorKind;

  constructor(
    message: string,
    readonly status: number,
  ) {
    super(message);
    this.name = 'ApiError';
    this.kind = kindForStatus(status);
  }
}

/** The backend could not be reached at all. */
export class NetworkError extends Error {
  constructor(options?: ErrorOptions) {
    super('Could not reach the calculator service. Please check that it is running and try again.', options);
    this.name = 'NetworkError';
  }
}

function kindForStatus(status: number): ApiErrorKind {
  if (status === 422) return 'math';
  if (status >= 400 && status < 500) return 'invalid-input';
  return 'server';
}

export interface CalculatorClient {
  calculate(operation: Operation, a: number, b?: number): Promise<CalculationResult>;
}

export function createCalculatorClient(baseUrl: string): CalculatorClient {
  const root = baseUrl.replace(/\/+$/, '');

  return {
    async calculate(operation, a, b) {
      const body = isUnary(operation) ? { a } : { a, b };

      let response: Response;
      try {
        response = await fetch(`${root}/api/v1/calculate/${operation}`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
      } catch (err) {
        throw new NetworkError({ cause: err });
      }

      const payload: unknown = await response.json().catch(() => undefined);

      if (!response.ok) {
        const error = field(payload, 'error');
        const message = typeof error === 'string' ? error : `Request failed with status ${response.status}`;
        throw new ApiError(message, response.status);
      }

      const result = field(payload, 'result');
      if (typeof result !== 'number') {
        throw new ApiError('Unexpected response from the calculator service', response.status);
      }
      return { operation, result };
    },
  };
}

/** Reads a property from an untrusted JSON value. */
function field(value: unknown, key: string): unknown {
  return typeof value === 'object' && value !== null ? (value as Record<string, unknown>)[key] : undefined;
}
