import { type CSSProperties } from 'react';

export type WeeklyTrendPoint = {
  label: string;
  minutes: number;
};

type WeeklyTrendChartProps = {
  data?: WeeklyTrendPoint[];
};

const defaultData: WeeklyTrendPoint[] = [
  { label: '周一', minutes: 40 },
  { label: '周二', minutes: 72 },
  { label: '周三', minutes: 68 },
  { label: '周四', minutes: 90 },
  { label: '周五', minutes: 120 },
  { label: '周六', minutes: 84 },
  { label: '周日', minutes: 100 },
];

export function WeeklyTrendChart({ data = defaultData }: WeeklyTrendChartProps) {
  const maxMinutes = Math.max(...data.map((point) => point.minutes), 1);

  return (
    <section style={styles.card} aria-label="本周学习趋势">
      <header style={styles.header}>
        <h2 style={styles.title}>本周趋势</h2>
        <p style={styles.subtitle}>按天展示学习时长</p>
      </header>

      <div style={styles.chart} role="img" aria-label="每周学习趋势图">
        {data.map((point) => {
          const height = Math.max((point.minutes / maxMinutes) * 100, 12);

          return (
            <div key={point.label} style={styles.barColumn}>
              <div style={styles.barTrack}>
                <div style={{ ...styles.barFill, height: `${height}%` }} />
              </div>
              <span style={styles.barLabel}>{point.label}</span>
              <span style={styles.barValue}>{point.minutes} 分</span>
            </div>
          );
        })}
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
  chart: {
    display: 'grid',
    gridTemplateColumns: 'repeat(7, minmax(0, 1fr))',
    alignItems: 'end',
    gap: '10px',
    minHeight: '220px',
  },
  barColumn: {
    display: 'grid',
    justifyItems: 'center',
    gap: '8px',
    height: '100%',
  },
  barTrack: {
    position: 'relative',
    width: '100%',
    minHeight: '150px',
    borderRadius: '999px',
    background: 'rgba(148, 163, 184, 0.12)',
    display: 'flex',
    alignItems: 'end',
    overflow: 'hidden',
  },
  barFill: {
    width: '100%',
    borderRadius: '999px 999px 0 0',
    background: 'linear-gradient(180deg, #38bdf8 0%, #6366f1 100%)',
  },
  barLabel: {
    fontSize: '12px',
    color: '#d7e3f4',
  },
  barValue: {
    fontSize: '12px',
    color: '#8fbaf8',
  },
};
