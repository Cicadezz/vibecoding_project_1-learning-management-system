import { type CSSProperties } from 'react';

import { SubjectPieChart, type SubjectBreakdown } from '../components/charts/SubjectPieChart';
import { WeeklyTrendChart, type WeeklyTrendPoint } from '../components/charts/WeeklyTrendChart';

type Metric = {
  label: string;
  value: string;
  hint: string;
};

const metrics: Metric[] = [
  {
    label: '今日学习总时长',
    value: '120 分钟',
    hint: '比昨日多 15 分钟',
  },
  {
    label: '本周总时长',
    value: '640 分钟',
    hint: '完成率 80%',
  },
  {
    label: '完成任务数',
    value: '8 个',
    hint: '还有 2 个待完成',
  },
  {
    label: '连续打卡',
    value: '5 天',
    hint: '保持当前节奏',
  },
];

const weeklyTrend: WeeklyTrendPoint[] = [
  { label: '周一', minutes: 55 },
  { label: '周二', minutes: 80 },
  { label: '周三', minutes: 70 },
  { label: '周四', minutes: 95 },
  { label: '周五', minutes: 120 },
  { label: '周六', minutes: 110 },
  { label: '周日', minutes: 110 },
];

const subjectBreakdown: SubjectBreakdown[] = [
  { label: '数学', minutes: 240, color: '#7c3aed' },
  { label: '英语', minutes: 180, color: '#2563eb' },
  { label: '编程', minutes: 140, color: '#0f766e' },
  { label: '其他', minutes: 80, color: '#f59e0b' },
];

export function DashboardPage() {
  return (
    <main style={styles.page}>
      <header style={styles.header}>
        <p style={styles.kicker}>学习成长看板</p>
        <h1 style={styles.title}>Dashboard</h1>
        <p style={styles.description}>查看今天的学习进度、任务状态和打卡趋势。</p>
      </header>

      <section aria-label="核心指标" style={styles.metricGrid}>
        {metrics.map((metric) => (
          <article key={metric.label} style={styles.metricCard}>
            <p style={styles.metricLabel}>{metric.label}</p>
            <p style={styles.metricValue}>{metric.value}</p>
            <p style={styles.metricHint}>{metric.hint}</p>
            {metric.label === '连续打卡' ? (
              <p style={styles.metricSummary}>{metric.label} {metric.value}</p>
            ) : null}
          </article>
        ))}
      </section>

      <section style={styles.grid}>
        <WeeklyTrendChart data={weeklyTrend} />
        <SubjectPieChart data={subjectBreakdown} />
      </section>
    </main>
  );
}

const styles: Record<string, CSSProperties> = {
  page: {
    display: 'grid',
    gap: '20px',
    color: '#e5eefb',
  },
  header: {
    display: 'grid',
    gap: '8px',
  },
  kicker: {
    margin: 0,
    textTransform: 'uppercase',
    letterSpacing: '0.16em',
    fontSize: '12px',
    color: '#7dd3fc',
  },
  title: {
    margin: 0,
    fontSize: '32px',
    lineHeight: 1.1,
  },
  description: {
    margin: 0,
    color: '#b7c6de',
  },
  metricGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))',
    gap: '14px',
  },
  metricCard: {
    background: 'rgba(12, 20, 38, 0.9)',
    border: '1px solid rgba(148, 163, 184, 0.16)',
    borderRadius: '18px',
    padding: '18px',
    boxShadow: '0 18px 42px rgba(3, 7, 18, 0.22)',
  },
  metricLabel: {
    margin: 0,
    color: '#93a4bf',
    fontSize: '13px',
  },
  metricValue: {
    margin: '10px 0 6px',
    fontSize: '28px',
    fontWeight: 700,
    color: '#f8fbff',
  },
  metricHint: {
    margin: 0,
    fontSize: '13px',
    color: '#c9d5e7',
  },
  metricSummary: {
    margin: '8px 0 0',
    fontSize: '13px',
    color: '#7dd3fc',
    fontWeight: 600,
  },
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
    gap: '14px',
  },
};
