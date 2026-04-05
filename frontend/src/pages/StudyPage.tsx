import { type CSSProperties } from 'react';

export function StudyPage() {
  return (
    <main style={styles.page}>
      <h1 style={styles.title}>Study</h1>
      <p style={styles.description}>专注学习页，展示当前学习内容和时段记录。</p>
      <section style={styles.card}>
        <p style={styles.cardLabel}>当前课程</p>
        <p style={styles.cardValue}>高等数学 - 极限与连续</p>
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
  cardLabel: {
    margin: 0,
    color: '#93a4bf',
    fontSize: '13px',
  },
  cardValue: {
    margin: '10px 0 0',
    fontSize: '20px',
    fontWeight: 600,
  },
};
