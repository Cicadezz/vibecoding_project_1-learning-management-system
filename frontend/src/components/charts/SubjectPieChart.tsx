import { type CSSProperties } from 'react';

export type SubjectBreakdown = {
  label: string;
  minutes: number;
  color: string;
};

type SubjectPieChartProps = {
  data?: SubjectBreakdown[];
};

const defaultData: SubjectBreakdown[] = [
  { label: '数学', minutes: 240, color: '#7c3aed' },
  { label: '英语', minutes: 180, color: '#2563eb' },
  { label: '编程', minutes: 140, color: '#0f766e' },
  { label: '其他', minutes: 80, color: '#f59e0b' },
];

export function SubjectPieChart({ data = defaultData }: SubjectPieChartProps) {
  const totalMinutes = data.reduce((sum, item) => sum + item.minutes, 0) || 1;
  let accumulated = 0;

  const segments = data.map((item) => {
    const start = accumulated;
    accumulated += (item.minutes / totalMinutes) * 360;
    return {
      ...item,
      start,
      end: accumulated,
    };
  });

  const background = segments
    .map((segment) => `${segment.color} ${segment.start}deg ${segment.end}deg`)
    .join(', ');

  return (
    <section style={styles.card} aria-label="学科分布">
      <header style={styles.header}>
        <h2 style={styles.title}>学科占比</h2>
        <p style={styles.subtitle}>按学科拆分学习时长</p>
      </header>

      <div style={styles.layout}>
        <div style={{ ...styles.ring, background: `conic-gradient(${background})` }} role="img" aria-label="学科时长分布图">
          <div style={styles.innerRing}>
            <span style={styles.totalLabel}>{totalMinutes} 分</span>
          </div>
        </div>

        <ul style={styles.legend}>
          {data.map((item) => (
            <li key={item.label} style={styles.legendItem}>
              <span style={{ ...styles.swatch, background: item.color }} />
              <span style={styles.legendText}>{item.label}</span>
              <span style={styles.legendValue}>{item.minutes} 分</span>
            </li>
          ))}
        </ul>
      </div>
    </section>
  );
}

const styles: Record<string, CSSProperties> = {
  card: {
    background: 'rgba(12, 20, 38, 0.9)',
    border: '1px solid rgba(148, 163, 184, 0.16)',
    borderRadius: '18px',
    padding: '18px',
  },
  header: {
    display: 'grid',
    gap: '4px',
    marginBottom: '16px',
  },
  title: {
    margin: 0,
    fontSize: '18px',
  },
  subtitle: {
    margin: 0,
    color: '#b7c6de',
    fontSize: '13px',
  },
  layout: {
    display: 'grid',
    gridTemplateColumns: 'minmax(160px, 220px) 1fr',
    gap: '18px',
    alignItems: 'center',
  },
  ring: {
    position: 'relative',
    width: '100%',
    aspectRatio: '1',
    borderRadius: '50%',
    padding: '16px',
  },
  innerRing: {
    position: 'absolute',
    inset: '16%',
    borderRadius: '50%',
    display: 'grid',
    placeItems: 'center',
    background: '#0c1426',
    border: '1px solid rgba(148, 163, 184, 0.16)',
  },
  totalLabel: {
    fontSize: '28px',
    fontWeight: 700,
    color: '#f8fbff',
  },
  legend: {
    listStyle: 'none',
    margin: 0,
    padding: 0,
    display: 'grid',
    gap: '10px',
  },
  legendItem: {
    display: 'grid',
    gridTemplateColumns: '12px 1fr auto',
    alignItems: 'center',
    gap: '10px',
  },
  swatch: {
    width: '12px',
    height: '12px',
    borderRadius: '999px',
  },
  legendText: {
    color: '#d7e3f4',
    fontSize: '14px',
  },
  legendValue: {
    color: '#8fbaf8',
    fontSize: '14px',
  },
};
