import { type CSSProperties } from 'react';

export function CheckinPage() {
  return (
    <main style={styles.page}>
      <h1 style={styles.title}>Checkin</h1>
      <p style={styles.description}>打卡页，确认今日学习是否完成。</p>
      <button type="button" style={styles.button}>
        今日已打卡
      </button>
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
  button: {
    width: 'fit-content',
    border: 'none',
    borderRadius: '12px',
    padding: '10px 14px',
    background: '#22c55e',
    color: '#052e16',
    fontWeight: 700,
    cursor: 'pointer',
  },
};
