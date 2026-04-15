import { type CSSProperties } from 'react';

type CheckinPageProps = {
  streak: number;
  onCheckin: () => Promise<void>;
};

export function CheckinPage({ streak, onCheckin }: CheckinPageProps) {
  return (
    <main style={styles.page}>
      <h1 style={styles.title}>每日打卡</h1>
      <p style={styles.description}>当前连续打卡：{streak} 天</p>
      <button type="button" style={styles.button} onClick={() => void onCheckin()}>
        今日打卡
      </button>
    </main>
  );
}

const styles: Record<string, CSSProperties> = {
  page: { display: 'grid', gap: '12px', color: '#e5eefb' },
  title: { margin: 0, fontSize: '28px' },
  description: { margin: 0, color: '#b7c6de' },
  button: { width: 'fit-content', border: 'none', borderRadius: '12px', padding: '10px 14px', background: '#22c55e', color: '#052e16', fontWeight: 700, cursor: 'pointer' },
};
