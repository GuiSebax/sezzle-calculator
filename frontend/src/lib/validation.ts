// Plain decimal or scientific notation: "5", "-3.2", ".5", "1e3".
// Stricter than Number(), which also accepts "0x10", "Infinity" and "".
const NUMBER_PATTERN = /^[+-]?(\d+\.?\d*|\.\d+)([eE][+-]?\d+)?$/;

export type ParsedOperand = { ok: true; value: number } | { ok: false; error: string };

export function parseOperand(raw: string): ParsedOperand {
  const text = raw.trim();
  if (text === '') {
    return { ok: false, error: 'This field is required.' };
  }
  const value = Number(text);
  if (!NUMBER_PATTERN.test(text) || !Number.isFinite(value)) {
    return { ok: false, error: 'Enter a valid number.' };
  }
  return { ok: true, value };
}
