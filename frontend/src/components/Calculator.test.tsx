import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { jsonResponse } from '../test/http';
import { Calculator } from './Calculator';

function setup() {
  const user = userEvent.setup();
  render(<Calculator />);
  return {
    user,
    first: () => screen.getByLabelText('First number'),
    second: () => screen.queryByLabelText('Second number'),
    operation: () => screen.getByLabelText('Operation'),
    button: () => screen.getByRole('button', { name: /calculat/i }),
  };
}

describe('Calculator', () => {
  it('calculates and shows the result', async () => {
    const fetchMock = vi.fn(async () => jsonResponse(200, { operation: 'multiply', result: 15 }));
    vi.stubGlobal('fetch', fetchMock);
    const { user, first, second, operation, button } = setup();

    await user.type(first(), '5');
    await user.selectOptions(operation(), 'multiply');
    await user.type(second()!, '3');
    await user.click(button());

    expect(await screen.findByText('15')).toBeInTheDocument();
    expect(screen.getByText('5 × 3 =')).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalledWith(
      'http://localhost:8080/api/v1/calculate/multiply',
      expect.objectContaining({ body: JSON.stringify({ a: 5, b: 3 }) }),
    );
  });

  it('submits with the Enter key', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => jsonResponse(200, { operation: 'add', result: 3 })));
    const { user, first, second } = setup();

    await user.type(first(), '1');
    await user.type(second()!, '2{Enter}');

    expect(await screen.findByText('3')).toBeInTheDocument();
  });

  it('shows inline validation messages and disables the button while invalid', async () => {
    const fetchMock = vi.fn();
    vi.stubGlobal('fetch', fetchMock);
    const { user, first, second, button } = setup();

    expect(button()).toBeDisabled();

    await user.type(first(), 'abc');
    expect(screen.getByText('Enter a valid number.')).toBeInTheDocument();
    expect(first()).toHaveAttribute('aria-invalid', 'true');

    await user.clear(first());
    expect(screen.getByText('This field is required.')).toBeInTheDocument();

    await user.type(first(), '4');
    await user.type(second()!, '2');
    expect(screen.queryByText(/required|valid number/)).not.toBeInTheDocument();
    expect(button()).toBeEnabled();
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it('shows the backend error message', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => jsonResponse(422, { error: 'division by zero' })));
    const { user, first, second, operation, button } = setup();

    await user.type(first(), '5');
    await user.selectOptions(operation(), 'divide');
    await user.type(second()!, '0');
    await user.click(button());

    const alert = await screen.findByRole('alert');
    expect(alert).toHaveTextContent('Math error');
    expect(alert).toHaveTextContent('Division by zero');
  });

  it('shows a friendly message when the backend is down', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => Promise.reject(new TypeError('Failed to fetch'))));
    const { user, first, second, button } = setup();

    await user.type(first(), '1');
    await user.type(second()!, '2');
    await user.click(button());

    expect(await screen.findByRole('alert')).toHaveTextContent('Could not reach the calculator service');
  });

  it('shows a loading state while the request is in flight', async () => {
    let respond!: (r: Response) => void;
    vi.stubGlobal('fetch', vi.fn(() => new Promise<Response>((resolve) => (respond = resolve))));
    const { user, first, second, button } = setup();

    await user.type(first(), '1');
    await user.type(second()!, '2');
    await user.click(button());

    expect(button()).toHaveTextContent('Calculating…');
    expect(button()).toBeDisabled();

    respond(jsonResponse(200, { operation: 'add', result: 3 }));
    expect(await screen.findByText('3')).toBeInTheDocument();
    expect(button()).toHaveTextContent('Calculate');
  });

  it('hides the second input for sqrt and sends only "a"', async () => {
    const fetchMock = vi.fn(async () => jsonResponse(200, { operation: 'sqrt', result: 4 }));
    vi.stubGlobal('fetch', fetchMock);
    const { user, second, operation, button } = setup();

    await user.selectOptions(operation(), 'sqrt');

    expect(second()).not.toBeInTheDocument();
    await user.type(screen.getByLabelText('Number'), '16');
    expect(button()).toBeEnabled();
    await user.click(button());

    expect(await screen.findByText('√16 =')).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/sqrt'),
      expect.objectContaining({ body: JSON.stringify({ a: 16 }) }),
    );
  });
});
