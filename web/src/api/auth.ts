import { request, setToken } from './http';

export interface User {
  id: string;
  username: string;
  email: string;
  createdAt: string;
  updatedAt: string;
}

export interface AuthResult {
  user: User;
  token: string;
}

export async function register(input: {
  username: string;
  email: string;
  password: string;
}) {
  const result = await request<AuthResult>('/auth/register', {
    method: 'POST',
    body: JSON.stringify(input)
  });

  setToken(result.token);
  return result;
}

export async function login(input: {
  email: string;
  password: string;
}) {
  const result = await request<AuthResult>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(input)
  });

  setToken(result.token);
  return result;
}

export async function me() {
  return request<{ user: User }>('/auth/me');
}
