import { describe, expect, it } from 'vitest';
import type { Operation } from '../api';
import { formatExpression, formatResult } from './operations';

describe('formatExpression', () => {
  it.each<[Operation, string]>([
    ['add', '5 + 3'],
    ['subtract', '5 − 3'],
    ['multiply', '5 × 3'],
    ['divide', '5 ÷ 3'],
    ['power', '5 ^ 3'],
    ['percentage', '5% of 3'],
    ['sqrt', '√5'],
  ])('formats %s', (operation, expected) => {
    expect(formatExpression(operation, 5, 3)).toBe(expected);
  });
});

describe('formatResult', () => {
  it.each([
    [0.1 + 0.2, '0.3'],
    [8, '8'],
    [-2.5, '-2.5'],
    [1 / 3, '0.333333333333'],
  ])('formats %d as %s', (value, expected) => {
    expect(formatResult(value)).toBe(expected);
  });
});
