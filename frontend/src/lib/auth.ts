import { request } from './api';

export type AuthPayload = {
  username: string;
  password: string;
};

export function login(payload: AuthPayload) {
  return request<void>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function register(payload: AuthPayload) {
  return request<void>('/auth/register', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}
