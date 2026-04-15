import { request } from './api';

export type AuthPayload = {
  username: string;
  password: string;
};

export type AuthUser = {
  id: number;
  username: string;
};

export type AuthResponse = {
  token: string;
  user: AuthUser;
};

export function login(payload: AuthPayload) {
  return request<AuthResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function register(payload: AuthPayload) {
  return request<AuthResponse>('/auth/register', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}
