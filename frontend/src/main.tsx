import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { Calculator } from './components/Calculator';
import './index.css';

const root = document.getElementById('root');
if (!root) throw new Error('Root element #root not found');

createRoot(root).render(
  <StrictMode>
    <main>
      <Calculator />
    </main>
  </StrictMode>,
);
