import { type CSSProperties, useEffect, useState } from 'react';

export type SubjectOption = {
  id: number;
  name: string;
};

type StudyPageProps = {
  subjects: SubjectOption[];
  onCreateStudy: (input: { subjectId: number; minutes: number; note: string }) => Promise<void>;
  onCreateSubject: (input: { name: string }) => Promise<void>;
};

export function StudyPage({ subjects, onCreateStudy, onCreateSubject }: StudyPageProps) {
  const [subjectId, setSubjectId] = useState<number>(subjects[0]?.id ?? 0);
  const [minutes, setMinutes] = useState(45);
  const [note, setNote] = useState('');
  const [newSubjectName, setNewSubjectName] = useState('');
  const [status, setStatus] = useState('');

  useEffect(() => {
    if (subjects.length === 0) {
      setSubjectId(0);
      return;
    }

    const exists = subjects.some((subject) => subject.id === subjectId);
    if (!exists) {
      setSubjectId(subjects[0].id);
    }
  }, [subjectId, subjects]);

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setStatus('');
    if (!subjectId) {
      setStatus('请先创建科目');
      return;
    }
    try {
      await onCreateStudy({ subjectId, minutes, note });
      setStatus('学习记录已保存');
      setNote('');
    } catch (e) {
      setStatus(e instanceof Error ? e.message : '保存失败');
    }
  }

  async function handleCreateSubject(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setStatus('');

    const name = newSubjectName.trim();
    if (!name) {
      setStatus('请输入科目名称');
      return;
    }

    try {
      await onCreateSubject({ name });
      setNewSubjectName('');
      setStatus('科目已创建');
    } catch (e) {
      setStatus(e instanceof Error ? e.message : '科目创建失败');
    }
  }

  return (
    <main style={styles.page}>
      <h1 style={styles.title}>学习记录</h1>
      <p style={styles.description}>手动记录本次学习时长。</p>

      <form onSubmit={handleCreateSubject} style={styles.subjectForm}>
        <input
          aria-label="科目名称"
          value={newSubjectName}
          onChange={(e) => setNewSubjectName(e.target.value)}
          placeholder="新增科目，例如：英语"
          style={styles.input}
        />
        <button type="submit" style={styles.secondaryButton}>新增科目</button>
      </form>

      <form onSubmit={handleSubmit} style={styles.form}>
        <select aria-label="科目" value={subjectId} onChange={(e) => setSubjectId(Number(e.target.value))} style={styles.input}>
          {subjects.length === 0 ? <option value={0}>暂无科目</option> : null}
          {subjects.map((s) => (
            <option key={s.id} value={s.id}>{s.name}</option>
          ))}
        </select>
        <input aria-label="学习分钟" type="number" min={1} value={minutes} onChange={(e) => setMinutes(Number(e.target.value))} style={styles.input} />
        <input aria-label="备注" value={note} onChange={(e) => setNote(e.target.value)} placeholder="可选备注" style={styles.input} />
        <button type="submit" style={styles.button}>保存记录</button>
      </form>

      {status ? <p style={styles.description}>{status}</p> : null}
    </main>
  );
}

const styles: Record<string, CSSProperties> = {
  page: { display: 'grid', gap: '12px', color: '#e5eefb' },
  title: { margin: 0, fontSize: '28px' },
  description: { margin: 0, color: '#b7c6de' },
  subjectForm: { display: 'grid', gridTemplateColumns: '1fr 120px', gap: '10px' },
  form: { display: 'grid', gridTemplateColumns: '1fr 120px 1fr 120px', gap: '10px' },
  input: { border: '1px solid rgba(148, 163, 184, 0.35)', borderRadius: '10px', background: '#0f172a', color: '#e2e8f0', padding: '10px 12px' },
  button: { border: 'none', borderRadius: '10px', background: '#0ea5e9', color: '#fff', padding: '10px 12px', cursor: 'pointer' },
  secondaryButton: { border: '1px solid rgba(125, 211, 252, 0.3)', borderRadius: '10px', background: '#082f49', color: '#e0f2fe', padding: '10px 12px', cursor: 'pointer' },
};
