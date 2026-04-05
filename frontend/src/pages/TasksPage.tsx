import { type CSSProperties } from 'react';

export function TasksPage() {
  const tasks = ['完成英语单词复习', '提交数学练习', '整理今日错题'];

  return (
    <main style={styles.page}>
      <h1 style={styles.title}>Tasks</h1>
      <p style={styles.description}>最小任务页，用于展示待办列表。</p>
      <ul style={styles.list}>
        {tasks.map((task) => (
          <li key={task} style={styles.item}>
            {task}
          </li>
        ))}
      </ul>
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
  list: {
    margin: 0,
    paddingLeft: '20px',
    display: 'grid',
    gap: '10px',
  },
  item: {
    color: '#f8fbff',
  },
};
