import { createCalculatorClient } from './calculator';

export * from './calculator';

const DEFAULT_API_URL = 'http://localhost:8080';

/** Client configured from VITE_API_URL; components should use this. */
export const calculatorClient = createCalculatorClient(import.meta.env.VITE_API_URL || DEFAULT_API_URL);
