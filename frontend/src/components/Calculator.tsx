import { useState, type FormEvent } from 'react';
import { isUnary, OPERATIONS, type Operation } from '../api';
import { useCalculation } from '../hooks/useCalculation';
import { OPERATION_LABELS } from '../lib/operations';
import { parseOperand, type ParsedOperand } from '../lib/validation';
import styles from './Calculator.module.css';
import { NumberField } from './NumberField';
import { ResultPanel } from './ResultPanel';

export function Calculator() {
  const [operation, setOperation] = useState<Operation>('add');
  const [a, setA] = useState('');
  const [b, setB] = useState('');
  // Errors are shown only once the user has interacted with a field.
  const [touched, setTouched] = useState({ a: false, b: false });
  const { state, calculate } = useCalculation();

  const unary = isUnary(operation);
  const parsedA = parseOperand(a);
  const parsedB: ParsedOperand = unary ? { ok: true, value: 0 } : parseOperand(b);
  const loading = state.status === 'loading';
  const canSubmit = parsedA.ok && parsedB.ok && !loading;

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!parsedA.ok || !parsedB.ok || loading) return;
    void calculate(operation, parsedA.value, parsedB.value);
  }

  const touch = (field: 'a' | 'b') => setTouched((t) => ({ ...t, [field]: true }));

  return (
    <form className={styles.card} onSubmit={handleSubmit} noValidate>
      <h1 className={styles.title}>Calculator</h1>

      <div className={unary ? `${styles.inputs} ${styles.unary}` : styles.inputs}>
        <NumberField
          id="operand-a"
          label={unary ? 'Number' : 'First number'}
          value={a}
          error={touched.a && !parsedA.ok ? parsedA.error : undefined}
          onChange={(value) => {
            setA(value);
            touch('a');
          }}
          onBlur={() => touch('a')}
        />

        <div className={styles.operation}>
          <label htmlFor="operation">Operation</label>
          <select id="operation" value={operation} onChange={(e) => setOperation(e.target.value as Operation)}>
            {OPERATIONS.map((op) => (
              <option key={op} value={op}>
                {OPERATION_LABELS[op]}
              </option>
            ))}
          </select>
        </div>

        {!unary && (
          <NumberField
            id="operand-b"
            label="Second number"
            value={b}
            error={touched.b && !parsedB.ok ? parsedB.error : undefined}
            onChange={(value) => {
              setB(value);
              touch('b');
            }}
            onBlur={() => touch('b')}
          />
        )}
      </div>

      <button type="submit" className={styles.submit} disabled={!canSubmit}>
        {loading ? 'Calculating…' : 'Calculate'}
      </button>

      <ResultPanel state={state} />
    </form>
  );
}
