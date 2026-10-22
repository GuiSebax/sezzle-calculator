import { describe, expect, it, vi } from 'vitest';
import { jsonResponse } from '../test/http';
import { ApiError, createCalculatorClient, NetworkError } from './calculator';

function stubFetch(impl: () => Promise<Response>) {
  const fetchMock = vi.fn<typeof fetch>(impl);
  vi.stubGlobal('fetch', fetchMock);
  return fetchMock;
}

const client = createCalculatorClient('http://api.test');

describe('calculatorClient.calculate', () => {
  it('posts both operands and returns the result', async () => {
    const fetchMock = stubFetch(async () => jsonResponse(200, { operation: 'add', result: 8 }));

    await expect(client.calculate('add', 5, 3)).resolves.toEqual({ operation: 'add', result: 8 });

    expect(fetchMock).toHaveBeenCalledWith('http://api.test/api/v1/calculate/add', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ a: 5, b: 3 }),
    });
  });

  it('sends only "a" for sqrt', async () => {
    const fetchMock = stubFetch(async () => jsonResponse(200, { operation: 'sqrt', result: 4 }));

    await client.calculate('sqrt', 16, 99);

    expect(fetchMock.mock.calls[0]?.[1]?.body).toBe(JSON.stringify({ a: 16 }));
  });

  it('ignores trailing slashes in the base URL', async () => {
    const fetchMock = stubFetch(async () => jsonResponse(200, { operation: 'add', result: 1 }));

    await createCalculatorClient('http://api.test//').calculate('add', 0, 1);

    expect(fetchMock.mock.calls[0]?.[0]).toBe('http://api.test/api/v1/calculate/add');
  });

  it.each([
    { status: 422, error: 'division by zero', kind: 'math' },
    { status: 400, error: 'field "b" is required', kind: 'invalid-input' },
    { status: 413, error: 'request body must not exceed 1024 bytes', kind: 'invalid-input' },
    { status: 500, error: 'internal error', kind: 'server' },
  ])('maps $status to an ApiError of kind $kind', async ({ status, error, kind }) => {
    stubFetch(async () => jsonResponse(status, { error }));

    const err = await client.calculate('divide', 5, 0).catch((e: unknown) => e);

    expect(err).toBeInstanceOf(ApiError);
    expect(err).toMatchObject({ status, kind, message: error });
  });

  it('falls back to a generic message when the error body is not JSON', async () => {
    stubFetch(async () => new Response('Method Not Allowed', { status: 405 }));

    await expect(client.calculate('add', 1, 2)).rejects.toMatchObject({
      status: 405,
      message: 'Request failed with status 405',
    });
  });

  it('rejects a success response without a numeric result', async () => {
    stubFetch(async () => jsonResponse(200, { operation: 'add' }));

    await expect(client.calculate('add', 1, 2)).rejects.toBeInstanceOf(ApiError);
  });

  it('throws NetworkError when the backend is unreachable', async () => {
    const cause = new TypeError('Failed to fetch');
    stubFetch(async () => {
      throw cause;
    });

    const err = await client.calculate('add', 1, 2).catch((e: unknown) => e);

    expect(err).toBeInstanceOf(NetworkError);
    expect(err).toMatchObject({ cause });
  });
});
