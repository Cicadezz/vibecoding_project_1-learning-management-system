import { type CSSProperties, useMemo, useState } from 'react';

import { CheckinPage } from './pages/CheckinPage';
import { DashboardPage } from './pages/DashboardPage';
import { SettingsPage } from './pages/SettingsPage';
import { StudyPage } from './pages/StudyPage';
import { TasksPage } from './pages/TasksPage';

type View = 'dashboard' | 'tasks' | 'study' | 'checkin' | 'settings';

type NavItem = {
  key: View;
  label: string;
};

const navItems: NavItem[] = [
  { key: 'dashboard', label: 'Dashboard' },
  { key: 'tasks', label: 'Tasks' },
  { key: 'study', label: 'Study' },
  { key: 'checkin', label: 'Checkin' },
  { key: 'settings', label: 'Settings' },
];

export default function App() {
  const [view, setView] = useState<View>('dashboard');

  const currentPage = useMemo(() => {
    switch (view) {
      case 'dashboard':
        return <DashboardPage />;
      case 'tasks':
        return <TasksPage />;
      case 'study':
        return <StudyPage />;
      case 'checkin':
        return <CheckinPage />;
      case 'settings':
        return <SettingsPage />;
      default:
        return null;
    }
  }, [view]);

  return (
    <main style={styles.shell}>
      <section style={styles.frame}>
        <header style={styles.header}>
          <div>
            <p style={styles.kicker}>Learning Growth MVP</p>
            <h1 style={styles.brand}>成长看板</h1>
          </div>
          <nav aria-label="Main navigation" style={styles.nav}>
            {navItems.map((item) => (
              <button
                key={item.key}
                type="button"
                onClick={() => setView(item.key)}
                style={view === item.key ? styles.activeNavButton : styles.navButton}
              >
                {item.label}
              </button>
            ))}
          </nav>
        </header>

        <section style={styles.content}>{currentPage}</section>
      </section>
    </main>
  );
}

const styles: Record<string, CSSProperties> = {
  shell: {
    minHeight: '100vh',
    padding: '24px',
    background:
      'radial-gradient(circle at top left, rgba(56, 189, 248, 0.22), transparent 32%), radial-gradient(circle at top right, rgba(124, 58, 237, 0.18), transparent 28%), linear-gradient(180deg, #06111f 0%, #0a1527 100%)',
    color: '#e5eefb',
    fontFamily: 'system-ui, sans-serif',
  },
  frame: {
    maxWidth: '1200px',
    margin: '0 auto',
    display: 'grid',
    gap: '24px',
  },
  header: {
    display: 'flex',
    flexWrap: 'wrap',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: '16px',
  },
  kicker: {
    margin: 0,
    color: '#7dd3fc',
    fontSize: '12px',
    letterSpacing: '0.16em',
    textTransform: 'uppercase',
  },
  brand: {
    margin: '6px 0 0',
    fontSize: '30px',
    lineHeight: 1.1,
  },
  nav: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: '10px',
  },
  navButton: {
    border: '1px solid rgba(148, 163, 184, 0.22)',
    borderRadius: '999px',
    padding: '10px 14px',
    background: 'rgba(15, 23, 42, 0.72)',
    color: '#d7e3f4',
    fontSize: '14px',
    cursor: 'pointer',
  },
  activeNavButton: {
    border: '1px solid rgba(56, 189, 248, 0.6)',
    borderRadius: '999px',
    padding: '10px 14px',
    background: 'linear-gradient(135deg, #38bdf8 0%, #6366f1 100%)',
    color: '#eff6ff',
    fontSize: '14px',
    cursor: 'pointer',
    boxShadow: '0 10px 24px rgba(56, 189, 248, 0.2)',
  },
  content: {
    minHeight: 'calc(100vh - 180px)',
  },
};
