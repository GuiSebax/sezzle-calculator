import type { CalculationState } from '../hooks/useCalculation';
import { formatResult } from '../lib/operations';
import styles from './ResultPanel.module.css';

export function ResultPanel({ state }: { state: CalculationState }) {
  return (
    // The live region is always mounted so screen readers announce updates.
    <section className={styles.panel} aria-live="polite" aria-busy={state.status === 'loading'} aria-label="Result">
      {state.status === 'idle' && <p className={styles.muted}>Enter values and press Calculate.</p>}
      {state.status === 'loading' && <p className={styles.muted}>Calculating…</p>}
      {state.status === 'success' && (
        <p className={styles.success}>
          <span className={styles.expression}>{state.expression} =</span>{' '}
          <output className={styles.result}>{formatResult(state.result)}</output>
        </p>
      )}
      {state.status === 'error' && (
        <div role="alert" className={styles.error}>
          <strong>{state.title}</strong>
          <p>{state.message}</p>
        </div>
      )}
    </section>
  );
}
