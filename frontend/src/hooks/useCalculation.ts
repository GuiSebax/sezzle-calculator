import { useCallback, useState } from 'react';
import { ApiError, calculatorClient, NetworkError, type Operation } from '../api';
import { formatExpression } from '../lib/operations';

export type CalculationState =
  | { status: 'idle' }
  | { status: 'loading' }
  | { status: 'success'; expression: string; result: number }
  | { status: 'error'; title: string; message: string };

/** Runs a calculation against the backend and tracks its lifecycle. */
export function useCalculation() {
  const [state, setState] = useState<CalculationState>({ status: 'idle' });

  const calculate = useCallback(async (operation: Operation, a: number, b: number) => {
    setState({ status: 'loading' });
    try {
      const { result } = await calculatorClient.calculate(operation, a, b);
      setState({ status: 'success', expression: formatExpression(operation, a, b), result });
    } catch (err) {
      setState({ status: 'error', ...describeError(err) });
    }
  }, []);

  return { state, calculate };
}

function describeError(err: unknown): { title: string; message: string } {
  if (err instanceof NetworkError) {
    return { title: 'Service unavailable', message: err.message };
  }
  if (err instanceof ApiError) {
    const titles = { 'invalid-input': 'Invalid input', math: 'Math error', server: 'Server error' } as const;
    return { title: titles[err.kind], message: capitalize(err.message) };
  }
  return { title: 'Unexpected error', message: 'Something went wrong. Please try again.' };
}

function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}
