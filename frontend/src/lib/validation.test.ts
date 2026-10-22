import { describe, expect, it } from 'vitest';
import { parseOperand } from './validation';

describe('parseOperand', () => {
  it.each([
    ['5', 5],
    ['-3.25', -3.25],
    ['+7', 7],
    ['.5', 0.5],
    ['5.', 5],
    ['1e3', 1000],
    ['  42  ', 42],
  ])('accepts %j', (raw, value) => {
    expect(parseOperand(raw)).toEqual({ ok: true, value });
  });

  it.each(['', '   '])('requires a value for %j', (raw) => {
    expect(parseOperand(raw)).toEqual({ ok: false, error: 'This field is required.' });
  });

  it.each(['abc', '1,5', '0x10', 'Infinity', '1e999', '--1', '.', '1e'])('rejects %j', (raw) => {
    expect(parseOperand(raw)).toEqual({ ok: false, error: 'Enter a valid number.' });
  });
});
