import { useState, type CSSProperties, type FormEvent } from 'react';

import { login, type AuthPayload } from '../lib/auth';

type LoginPageProps = {
  onSubmit?: (payload: AuthPayload) => Promise<void>;
  onSwitchToRegister?: () => void;
};

type LoginFormState = {
  username: string;
  password: string;
};

const initialState: LoginFormState = {
  username: '',
  password: '',
};

export function LoginPage({ onSubmit, onSwitchToRegister }: LoginPageProps = {}) {
  const [formState, setFormState] = useState<LoginFormState>(initialState);
  const [status, setStatus] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsSubmitting(true);
    setStatus('');

    try {
      if (onSubmit) {
        await onSubmit(formState);
      } else {
        await login(formState);
      }
      setStatus('登录成功');
    } catch (error) {
      setStatus(error instanceof Error ? error.message : '登录失败');
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <section>
      <h1 style={styles.title}>登录</h1>
      <form onSubmit={handleSubmit} style={styles.form}>
        <label style={styles.field}>
          <span>用户名</span>
          <input
            aria-label="用户名"
            value={formState.username}
            onChange={(event) =>
              setFormState((current) => ({ ...current, username: event.target.value }))
            }
            autoComplete="username"
            style={styles.input}
          />
        </label>
        <label style={styles.field}>
          <span>密码</span>
          <input
            aria-label="密码"
            type="password"
            value={formState.password}
            onChange={(event) =>
              setFormState((current) => ({ ...current, password: event.target.value }))
            }
            autoComplete="current-password"
            style={styles.input}
          />
        </label>
        <button type="submit" disabled={isSubmitting} style={styles.button}>
          {isSubmitting ? '提交中' : '登录'}
        </button>
      </form>
      {onSwitchToRegister ? (
        <button type="button" onClick={onSwitchToRegister} style={styles.linkButton}>
          去注册
        </button>
      ) : null}
      {status ? (
        <p aria-live="polite" style={styles.status}>
          {status}
        </p>
      ) : null}
    </section>
  );
}

const styles: Record<string, CSSProperties> = {
  title: {
    margin: '0 0 16px',
    fontSize: '24px',
  },
  form: {
    display: 'grid',
    gap: '14px',
  },
  field: {
    display: 'grid',
    gap: '8px',
    fontSize: '14px',
  },
  input: {
    border: '1px solid #c8d2e0',
    borderRadius: '10px',
    padding: '10px 12px',
    fontSize: '14px',
  },
  button: {
    border: 'none',
    borderRadius: '10px',
    padding: '10px 12px',
    background: '#1d4ed8',
    color: '#ffffff',
    fontSize: '14px',
    cursor: 'pointer',
  },
  linkButton: {
    marginTop: '12px',
    border: 'none',
    padding: 0,
    background: 'transparent',
    color: '#1d4ed8',
    fontSize: '14px',
    cursor: 'pointer',
  },
  status: {
    marginTop: '12px',
    fontSize: '14px',
    color: '#334155',
  },
};
