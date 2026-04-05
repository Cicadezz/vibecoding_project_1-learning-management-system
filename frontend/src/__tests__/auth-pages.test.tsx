import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { LoginPage } from '../pages/LoginPage';
import { RegisterPage } from '../pages/RegisterPage';
import { login, register } from '../lib/auth';

vi.mock('../lib/auth', () => ({
  login: vi.fn(),
  register: vi.fn(),
}));

describe('auth pages', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('submits username and password from the login page', async () => {
    render(<LoginPage />);

    fireEvent.change(screen.getByLabelText('用户名'), {
      target: { value: 'alice' },
    });
    fireEvent.change(screen.getByLabelText('密码'), {
      target: { value: 'secret123' },
    });
    fireEvent.click(screen.getByRole('button', { name: '登录' }));

    await waitFor(() => {
      expect(login).toHaveBeenCalledWith({
        username: 'alice',
        password: 'secret123',
      });
    });
  });

  it('submits username and password from the register page', async () => {
    render(<RegisterPage />);

    fireEvent.change(screen.getByLabelText('用户名'), {
      target: { value: 'bob' },
    });
    fireEvent.change(screen.getByLabelText('密码'), {
      target: { value: 'pass456' },
    });
    fireEvent.click(screen.getByRole('button', { name: '注册' }));

    await waitFor(() => {
      expect(register).toHaveBeenCalledWith({
        username: 'bob',
        password: 'pass456',
      });
    });
  });
});
