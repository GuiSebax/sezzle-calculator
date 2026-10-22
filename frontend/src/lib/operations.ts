import type { Operation } from '../api';

export const OPERATION_LABELS: Record<Operation, string> = {
  add: 'Add (+)',
  subtract: 'Subtract (−)',
  multiply: 'Multiply (×)',
  divide: 'Divide (÷)',
  power: 'Power (^)',
  percentage: 'Percentage (a% of b)',
  sqrt: 'Square root (√)',
};

/** Human-readable expression, e.g. "5 + 3" or "√16". */
export function formatExpression(operation: Operation, a: number, b: number): string {
  switch (operation) {
    case 'add':
      return `${a} + ${b}`;
    case 'subtract':
      return `${a} − ${b}`;
    case 'multiply':
      return `${a} × ${b}`;
    case 'divide':
      return `${a} ÷ ${b}`;
    case 'power':
      return `${a} ^ ${b}`;
    case 'percentage':
      return `${a}% of ${b}`;
    case 'sqrt':
      return `√${a}`;
  }
}

/** Rounds away binary floating-point noise for display (0.1 + 0.2 → 0.3). */
export function formatResult(value: number): string {
  return String(Number(value.toPrecision(12)));
}
