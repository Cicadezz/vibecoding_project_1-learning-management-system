import { type CSSProperties } from 'react';

export function SettingsPage() {
  return (
    <main style={styles.page}>
      <h1 style={styles.title}>Settings</h1>
      <p style={styles.description}>设置页，用于管理基础偏好。</p>
      <section style={styles.card}>
        <p style={styles.row}>学习提醒：已开启</p>
        <p style={styles.row}>深色主题：启用</p>
      </section>
    </main>
  );
}

const styles: Record<string, CSSProperties> = {
  page: {
    display: 'grid',
    gap: '12px',
    color: '#e5eefb',
  },
  title: {
    margin: 0,
    fontSize: '28px',
  },
  description: {
    margin: 0,
    color: '#b7c6de',
  },
  card: {
    background: 'rgba(12, 20, 38, 0.9)',
    border: '1px solid rgba(148, 163, 184, 0.16)',
    borderRadius: '18px',
    padding: '18px',
  },
  row: {
    margin: '0 0 8px',
  },
};
