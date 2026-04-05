import { afterEach, describe, expect, it, vi } from 'vitest';

import { request } from '../lib/api';

describe('request error parsing', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  it('prefers error over message when parsing json failures', async () => {
    const response = new Response(JSON.stringify({ error: 'boom', message: 'fallback' }), {
      status: 400,
      headers: {
        'content-type': 'application/json',
      },
    });

    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(response));

    await expect(request('/test')).rejects.toThrow('boom');
  });
});
