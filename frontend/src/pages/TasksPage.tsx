import { type CSSProperties, useState } from 'react';

export type TaskItem = {
  id: number;
  title: string;
  status: string;
  priority: string;
  planDate: string;
};

type TasksPageProps = {
  tasks: TaskItem[];
  onCreateTask: (input: { title: string; priority: string }) => Promise<void>;
  onToggleDone: (task: TaskItem) => Promise<void>;
  loading?: boolean;
};

export function TasksPage({ tasks, onCreateTask, onToggleDone, loading = false }: TasksPageProps) {
  const [title, setTitle] = useState('');
  const [priority, setPriority] = useState('MEDIUM');
  const [error, setError] = useState('');

  async function handleCreate(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError('');
    if (!title.trim()) {
      setError('请输入任务标题');
      return;
    }
    try {
      await onCreateTask({ title: title.trim(), priority });
      setTitle('');
    } catch (e) {
      setError(e instanceof Error ? e.message : '创建任务失败');
    }
  }

  return (
    <main style={styles.page}>
      <h1 style={styles.title}>今日任务</h1>
      <form onSubmit={handleCreate} style={styles.form}>
        <input aria-label="任务标题" value={title} onChange={(e) => setTitle(e.target.value)} placeholder="新增任务，例如：刷数学题" style={styles.input} />
        <select aria-label="优先级" value={priority} onChange={(e) => setPriority(e.target.value)} style={styles.input}>
          <option value="HIGH">高</option>
          <option value="MEDIUM">中</option>
          <option value="LOW">低</option>
        </select>
        <button type="submit" style={styles.button}>添加任务</button>
      </form>
      {error ? <p style={styles.error}>{error}</p> : null}
      {loading ? <p style={styles.description}>加载中...</p> : null}
      <ul style={styles.list}>
        {tasks.map((task) => (
          <li key={task.id} style={styles.item}>
            <span>{task.title} ({task.priority})</span>
            <button type="button" onClick={() => void onToggleDone(task)} style={styles.smallButton}>
              {task.status === 'DONE' ? '标记待办' : '标记完成'}
            </button>
          </li>
        ))}
      </ul>
    </main>
  );
}

const styles: Record<string, CSSProperties> = {
  page: { display: 'grid', gap: '12px', color: '#e5eefb' },
  title: { margin: 0, fontSize: '28px' },
  description: { margin: 0, color: '#b7c6de' },
  form: { display: 'grid', gridTemplateColumns: '1fr 120px 120px', gap: '10px' },
  input: { border: '1px solid rgba(148, 163, 184, 0.35)', borderRadius: '10px', background: '#0f172a', color: '#e2e8f0', padding: '10px 12px' },
  button: { border: 'none', borderRadius: '10px', background: '#2563eb', color: '#fff', padding: '10px 12px', cursor: 'pointer' },
  smallButton: { border: 'none', borderRadius: '8px', background: '#334155', color: '#fff', padding: '8px 10px', cursor: 'pointer' },
  list: { margin: 0, paddingLeft: '20px', display: 'grid', gap: '10px' },
  item: { color: '#f8fbff', display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '12px' },
  error: { margin: 0, color: '#fca5a5' },
};
