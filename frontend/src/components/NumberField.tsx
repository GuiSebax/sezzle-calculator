import styles from './NumberField.module.css';

interface NumberFieldProps {
  id: string;
  label: string;
  value: string;
  error?: string | undefined;
  onChange: (value: string) => void;
  onBlur: () => void;
}

export function NumberField({ id, label, value, error, onChange, onBlur }: NumberFieldProps) {
  const errorId = `${id}-error`;

  return (
    <div className={styles.field}>
      <label htmlFor={id}>{label}</label>
      {/* type="text" + inputMode keeps the raw text so we can validate it;
          type="number" reports "" for invalid input, hiding the reason. */}
      <input
        id={id}
        type="text"
        inputMode="decimal"
        autoComplete="off"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onBlur={onBlur}
        aria-invalid={error ? true : undefined}
        aria-describedby={error ? errorId : undefined}
      />
      {error && (
        <p id={errorId} className={styles.error}>
          {error}
        </p>
      )}
    </div>
  );
}
