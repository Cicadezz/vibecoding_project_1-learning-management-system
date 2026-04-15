import { type CSSProperties } from 'react';

type SettingsPageProps = {
  username: string;
  onLogout: () => void;
};

export function SettingsPage({ username, onLogout }: SettingsPageProps) {
  return (
    <main style={styles.page}>
      <h1 style={styles.title}>设置</h1>
      <p style={styles.description}>当前登录用户：{username}</p>
      <section style={styles.card}>
        <p style={styles.row}>学习提醒：已开启</p>
        <p style={styles.row}>统计周起点：周一</p>
      </section>
      <button type="button" onClick={onLogout} style={styles.button}>退出登录</button>
    </main>
  );
}

const styles: Record<string, CSSProperties> = {
  page: { display: 'grid', gap: '12px', color: '#e5eefb' },
  title: { margin: 0, fontSize: '28px' },
  description: { margin: 0, color: '#b7c6de' },
  card: { background: 'rgba(12, 20, 38, 0.9)', border: '1px solid rgba(148, 163, 184, 0.16)', borderRadius: '18px', padding: '18px' },
  row: { margin: '0 0 8px' },
  button: { width: 'fit-content', border: 'none', borderRadius: '10px', padding: '10px 14px', background: '#ef4444', color: '#fff', cursor: 'pointer' },
};
