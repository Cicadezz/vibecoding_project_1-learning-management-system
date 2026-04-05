type JsonRecord = Record<string, unknown>;

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? '/api';

export async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers = new Headers(options.headers);

  if (options.body != null && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const response = await fetch(joinUrl(apiBaseUrl, path), {
    ...options,
    headers,
  });

  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }

  if (response.status === 204) {
    return undefined as T;
  }

  const text = await response.text();

  if (text.length === 0) {
    return undefined as T;
  }

  const contentType = response.headers.get('content-type') ?? '';

  if (contentType.includes('application/json')) {
    return JSON.parse(text) as T;
  }

  return text as T;
}

function joinUrl(baseUrl: string, path: string) {
  const normalizedBaseUrl = baseUrl.replace(/\/$/, '');
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;

  return `${normalizedBaseUrl}${normalizedPath}`;
}

async function parseErrorMessage(response: Response) {
  const body = await response.text();
  const contentType = response.headers.get('content-type') ?? '';

  if (contentType.includes('application/json') && body) {
    try {
      const payload = JSON.parse(body) as JsonRecord;
      const error = payload.error;

      if (typeof error === 'string' && error.trim()) {
        return error;
      }

      const message = payload.message;

      if (typeof message === 'string' && message.trim()) {
        return message;
      }
    } catch {
      // Fall through to the plain-text handling below.
    }
  }

  return body.trim() || `Request failed with status ${response.status}`;
}

