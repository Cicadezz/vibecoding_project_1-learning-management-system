import { useState, type CSSProperties } from 'react';

import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';

type View = 'login' | 'register';

export default function App() {
  const [view, setView] = useState<View>('login');

  return (
    <main style={styles.shell}>
      <section style={styles.card}>
        <div style={styles.switcher}>
          <button
            type="button"
            onClick={() => setView('login')}
            style={view === 'login' ? styles.activeTab : styles.tab}
          >
            登录
          </button>
          <button
            type="button"
            onClick={() => setView('register')}
            style={view === 'register' ? styles.activeTab : styles.tab}
          >
            注册
          </button>
        </div>
        {view === 'login' ? (
          <LoginPage onSwitchToRegister={() => setView('register')} />
        ) : (
          <RegisterPage onSwitchToLogin={() => setView('login')} />
        )}
      </section>
    </main>
  );
}

const styles: Record<string, CSSProperties> = {
  shell: {
    minHeight: '100vh',
    display: 'grid',
    placeItems: 'center',
    padding: '24px',
    background: '#f5f7fb',
    color: '#102033',
    fontFamily: 'system-ui, sans-serif',
  },
  card: {
    width: '100%',
    maxWidth: '420px',
    background: '#ffffff',
    borderRadius: '16px',
    padding: '24px',
    boxShadow: '0 12px 40px rgba(16, 32, 51, 0.12)',
  },
  switcher: {
    display: 'grid',
    gridTemplateColumns: '1fr 1fr',
    gap: '8px',
    marginBottom: '24px',
  },
  tab: {
    border: '1px solid #c8d2e0',
    background: '#f8fafc',
    color: '#102033',
    borderRadius: '10px',
    padding: '10px 12px',
    fontSize: '14px',
    cursor: 'pointer',
  },
  activeTab: {
    border: '1px solid #1d4ed8',
    background: '#1d4ed8',
    color: '#ffffff',
    borderRadius: '10px',
    padding: '10px 12px',
    fontSize: '14px',
    cursor: 'pointer',
  },
};
