import { type CSSProperties, useCallback, useEffect, useMemo, useState } from 'react';

import { type SubjectBreakdown } from './components/charts/SubjectPieChart';
import { type WeeklyTrendPoint } from './components/charts/WeeklyTrendChart';
import { request } from './lib/api';
import { login, register, type AuthPayload, type AuthUser } from './lib/auth';
import { CheckinPage } from './pages/CheckinPage';
import { type DashboardOverview, DashboardPage } from './pages/DashboardPage';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import { SettingsPage } from './pages/SettingsPage';
import { type SubjectOption, StudyPage } from './pages/StudyPage';
import { type TaskItem, TasksPage } from './pages/TasksPage';

type View = 'dashboard' | 'tasks' | 'study' | 'checkin' | 'settings';

type AuthMode = 'login' | 'register';

type NavItem = {
  key: View;
  label: string;
};

const TOKEN_KEY = 'learning_growth_token';

const navItems: NavItem[] = [
  { key: 'dashboard', label: '总览' },
  { key: 'tasks', label: '任务' },
  { key: 'study', label: '学习记录' },
  { key: 'checkin', label: '打卡' },
  { key: 'settings', label: '设置' },
];

const subjectColors = ['#7c3aed', '#2563eb', '#0f766e', '#f59e0b', '#db2777', '#0891b2'];

function getTodayDateString() {
  const now = new Date();
  const year = now.getFullYear();
  const month = `${now.getMonth() + 1}`.padStart(2, '0');
  const day = `${now.getDate()}`.padStart(2, '0');
  return `${year}-${month}-${day}`;
}

function getAuthHeader(token: string) {
  return {
    Authorization: `Bearer ${token}`,
  };
}

function toLocalISOString(value: Date) {
  const pad = (num: number) => `${Math.trunc(Math.abs(num))}`.padStart(2, '0');
  const year = value.getFullYear();
  const month = pad(value.getMonth() + 1);
  const day = pad(value.getDate());
  const hour = pad(value.getHours());
  const minute = pad(value.getMinutes());
  const second = pad(value.getSeconds());

  const offsetMinutes = -value.getTimezoneOffset();
  const sign = offsetMinutes >= 0 ? '+' : '-';
  const offsetHour = pad(Math.floor(Math.abs(offsetMinutes) / 60));
  const offsetMinute = pad(Math.abs(offsetMinutes) % 60);

  return `${year}-${month}-${day}T${hour}:${minute}:${second}${sign}${offsetHour}:${offsetMinute}`;
}

function pickString(record: Record<string, unknown>, keys: string[], fallback = '') {
  for (const key of keys) {
    const value = record[key];
    if (typeof value === 'string') {
      return value;
    }
  }
  return fallback;
}

function pickNumber(record: Record<string, unknown>, keys: string[], fallback = 0) {
  for (const key of keys) {
    const value = record[key];
    if (typeof value === 'number') {
      return value;
    }
    if (typeof value === 'string') {
      const parsed = Number(value);
      if (!Number.isNaN(parsed)) {
        return parsed;
      }
    }
  }
  return fallback;
}

export default function App() {
  const [view, setView] = useState<View>('dashboard');
  const [authMode, setAuthMode] = useState<AuthMode>('login');
  const [token, setToken] = useState<string>(() => localStorage.getItem(TOKEN_KEY) ?? '');
  const [user, setUser] = useState<AuthUser | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const [overview, setOverview] = useState<DashboardOverview>({
    todayMinutes: 0,
    weekMinutes: 0,
    doneTasks: 0,
    streak: 0,
  });
  const [trend, setTrend] = useState<WeeklyTrendPoint[]>([]);
  const [subjectDistribution, setSubjectDistribution] = useState<SubjectBreakdown[]>([]);
  const [tasks, setTasks] = useState<TaskItem[]>([]);
  const [subjects, setSubjects] = useState<SubjectOption[]>([]);
  const [streak, setStreak] = useState(0);

  const fetchDashboardData = useCallback(async (currentToken: string) => {
    const headers = getAuthHeader(currentToken);

    const [overviewResp, trendResp, subjectResp, tasksResp, subjectsResp, streakResp] = await Promise.all([
      request<{ overview: { today_minutes: number; week_minutes: number; done_tasks: number; streak: number } }>('/stats/overview', {
        headers,
      }),
      request<{ weekly_trend: Array<{ date: string; minutes: number }> }>('/stats/weekly-trend', {
        headers,
      }),
      request<{ subject_distribution: Array<{ subject_id: number; subject_name: string; minutes: number }> }>('/stats/subject-distribution', {
        headers,
      }),
      request<{ tasks: Array<Record<string, unknown>> }>(`/tasks?date=${getTodayDateString()}`, {
        headers,
      }),
      request<{ subjects: Array<Record<string, unknown>> }>('/subjects', {
        headers,
      }),
      request<{ streak: number }>('/checkin/streak', {
        headers,
      }),
    ]);

    setOverview({
      todayMinutes: overviewResp.overview.today_minutes,
      weekMinutes: overviewResp.overview.week_minutes,
      doneTasks: overviewResp.overview.done_tasks,
      streak: overviewResp.overview.streak,
    });

    setTrend(
      trendResp.weekly_trend.map((item) => ({
        label: item.date.slice(5),
        minutes: item.minutes,
      })),
    );

    setSubjectDistribution(
      subjectResp.subject_distribution.map((item, index) => ({
        label: item.subject_name || `科目${item.subject_id}`,
        minutes: item.minutes,
        color: subjectColors[index % subjectColors.length],
      })),
    );

    setTasks(
      tasksResp.tasks
        .map((task) => ({
          id: pickNumber(task, ['id', 'ID']),
          title: pickString(task, ['title', 'Title']),
          status: pickString(task, ['status', 'Status'], 'PENDING'),
          priority: pickString(task, ['priority', 'Priority'], 'MEDIUM'),
          planDate: pickString(task, ['plan_date', 'PlanDate']),
        }))
        .filter((task) => task.id > 0),
    );

    setSubjects(
      subjectsResp.subjects
        .map((subject) => ({
          id: pickNumber(subject, ['id', 'ID']),
          name: pickString(subject, ['name', 'Name']),
        }))
        .filter((subject) => subject.id > 0 && subject.name.trim().length > 0),
    );
    setStreak(streakResp.streak);
  }, []);

  const loadSession = useCallback(async () => {
    if (!token) {
      return;
    }

    setLoading(true);
    setError('');
    try {
      const me = await request<{ user: AuthUser }>('/auth/me', {
        headers: getAuthHeader(token),
      });
      setUser(me.user);
      await fetchDashboardData(token);
    } catch (e) {
      setToken('');
      setUser(null);
      localStorage.removeItem(TOKEN_KEY);
      setError(e instanceof Error ? e.message : '会话已失效，请重新登录');
    } finally {
      setLoading(false);
    }
  }, [fetchDashboardData, token]);

  useEffect(() => {
    void loadSession();
  }, [loadSession]);

  async function handleAuthSubmit(mode: AuthMode, payload: AuthPayload) {
    setLoading(true);
    setError('');
    try {
      const resp = mode === 'login' ? await login(payload) : await register(payload);
      localStorage.setItem(TOKEN_KEY, resp.token);
      setToken(resp.token);
      setUser(resp.user);
      await fetchDashboardData(resp.token);
      setView('dashboard');
    } catch (e) {
      setError(e instanceof Error ? e.message : '认证失败');
      throw e;
    } finally {
      setLoading(false);
    }
  }

  async function createTask(input: { title: string; priority: string }) {
    if (!token) return;
    const planDate = new Date();
    planDate.setHours(12, 0, 0, 0);

    await request('/tasks', {
      method: 'POST',
      headers: getAuthHeader(token),
      body: JSON.stringify({
        title: input.title,
        priority: input.priority,
        plan_date: toLocalISOString(planDate),
        status: 'PENDING',
      }),
    });

    await fetchDashboardData(token);
  }

  async function toggleTaskDone(task: TaskItem) {
    if (!token) return;

    const nextStatus = task.status === 'DONE' ? 'PENDING' : 'DONE';
    const payload: Record<string, unknown> = { status: nextStatus };

    if (nextStatus === 'DONE') {
      payload.completed_at = toLocalISOString(new Date());
    }

    await request(`/tasks/${task.id}`, {
      method: 'PUT',
      headers: getAuthHeader(token),
      body: JSON.stringify(payload),
    });

    await fetchDashboardData(token);
  }

  async function createStudy(input: { subjectId: number; minutes: number; note: string }) {
    if (!token) return;
    const endAt = new Date();
    const startAt = new Date(endAt.getTime() - input.minutes * 60 * 1000);

    await request('/study/sessions', {
      method: 'POST',
      headers: getAuthHeader(token),
      body: JSON.stringify({
        subject_id: input.subjectId,
        start_at: toLocalISOString(startAt),
        end_at: toLocalISOString(endAt),
        note: input.note || undefined,
      }),
    });

    await fetchDashboardData(token);
  }

  async function createSubject(input: { name: string }) {
    if (!token) return;

    await request('/subjects', {
      method: 'POST',
      headers: getAuthHeader(token),
      body: JSON.stringify({
        name: input.name,
      }),
    });

    await fetchDashboardData(token);
  }

  async function checkinToday() {
    if (!token) return;

    await request('/checkin/today', {
      method: 'POST',
      headers: getAuthHeader(token),
    });

    await fetchDashboardData(token);
  }

  function handleLogout() {
    localStorage.removeItem(TOKEN_KEY);
    setToken('');
    setUser(null);
    setTasks([]);
    setSubjects([]);
    setTrend([]);
    setSubjectDistribution([]);
    setOverview({ todayMinutes: 0, weekMinutes: 0, doneTasks: 0, streak: 0 });
    setStreak(0);
    setView('dashboard');
  }

  const currentPage = useMemo(() => {
    switch (view) {
      case 'dashboard':
        return <DashboardPage overview={overview} trend={trend} subjects={subjectDistribution} />;
      case 'tasks':
        return <TasksPage tasks={tasks} onCreateTask={createTask} onToggleDone={toggleTaskDone} loading={loading} />;
      case 'study':
        return <StudyPage subjects={subjects} onCreateStudy={createStudy} onCreateSubject={createSubject} />;
      case 'checkin':
        return <CheckinPage streak={streak} onCheckin={checkinToday} />;
      case 'settings':
        return <SettingsPage username={user?.username ?? '未登录'} onLogout={handleLogout} />;
      default:
        return null;
    }
  }, [checkinToday, createStudy, createSubject, createTask, loading, overview, streak, subjectDistribution, subjects, tasks, toggleTaskDone, trend, user?.username, view]);

  if (!token || !user) {
    return (
      <main style={styles.authShell}>
        <section style={styles.authCard}>
          <h1 style={styles.authTitle}>学习成长管理平台</h1>
          <p style={styles.authDescription}>请先登录后再进入业务页面。</p>
          {error ? <p style={styles.error}>{error}</p> : null}
          {authMode === 'login' ? (
            <LoginPage onSubmit={(payload) => handleAuthSubmit('login', payload)} onSwitchToRegister={() => setAuthMode('register')} />
          ) : (
            <RegisterPage onSubmit={(payload) => handleAuthSubmit('register', payload)} onSwitchToLogin={() => setAuthMode('login')} />
          )}
        </section>
      </main>
    );
  }

  return (
    <main style={styles.shell}>
      <section style={styles.frame}>
        <header style={styles.header}>
          <div>
            <p style={styles.kicker}>学习成长平台</p>
            <h1 style={styles.brand}>成长看板</h1>
            <p style={styles.user}>当前用户：{user.username}</p>
          </div>
          <nav aria-label="主导航" style={styles.nav}>
            {navItems.map((item) => (
              <button
                key={item.key}
                type="button"
                onClick={() => setView(item.key)}
                style={view === item.key ? styles.activeNavButton : styles.navButton}
              >
                {item.label}
              </button>
            ))}
          </nav>
        </header>

        {error ? <p style={styles.error}>{error}</p> : null}
        <section style={styles.content}>{currentPage}</section>
      </section>
    </main>
  );
}

const styles: Record<string, CSSProperties> = {
  shell: {
    minHeight: '100vh',
    padding: '24px',
    background:
      'radial-gradient(circle at top left, rgba(56, 189, 248, 0.22), transparent 32%), radial-gradient(circle at top right, rgba(124, 58, 237, 0.18), transparent 28%), linear-gradient(180deg, #06111f 0%, #0a1527 100%)',
    color: '#e5eefb',
    fontFamily: 'system-ui, sans-serif',
  },
  frame: {
    maxWidth: '1200px',
    margin: '0 auto',
    display: 'grid',
    gap: '24px',
  },
  header: {
    display: 'flex',
    flexWrap: 'wrap',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: '16px',
  },
  kicker: {
    margin: 0,
    color: '#7dd3fc',
    fontSize: '12px',
    letterSpacing: '0.16em',
    textTransform: 'uppercase',
  },
  brand: {
    margin: '6px 0 0',
    fontSize: '30px',
    lineHeight: 1.1,
  },
  user: {
    margin: '8px 0 0',
    fontSize: '14px',
    color: '#c7d2fe',
  },
  nav: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: '10px',
  },
  navButton: {
    border: '1px solid rgba(148, 163, 184, 0.22)',
    borderRadius: '999px',
    padding: '10px 14px',
    background: 'rgba(15, 23, 42, 0.72)',
    color: '#d7e3f4',
    fontSize: '14px',
    cursor: 'pointer',
  },
  activeNavButton: {
    border: '1px solid rgba(56, 189, 248, 0.6)',
    borderRadius: '999px',
    padding: '10px 14px',
    background: 'linear-gradient(135deg, #38bdf8 0%, #6366f1 100%)',
    color: '#eff6ff',
    fontSize: '14px',
    cursor: 'pointer',
    boxShadow: '0 10px 24px rgba(56, 189, 248, 0.2)',
  },
  content: {
    minHeight: 'calc(100vh - 180px)',
  },
  authShell: {
    minHeight: '100vh',
    display: 'grid',
    placeItems: 'center',
    background: 'linear-gradient(180deg, #06111f 0%, #0a1527 100%)',
    fontFamily: 'system-ui, sans-serif',
    padding: '24px',
  },
  authCard: {
    width: 'min(480px, 100%)',
    background: '#ffffff',
    color: '#0f172a',
    borderRadius: '16px',
    padding: '24px',
    boxShadow: '0 20px 50px rgba(2, 6, 23, 0.35)',
  },
  authTitle: { margin: 0, fontSize: '26px' },
  authDescription: { margin: '8px 0 16px', color: '#475569' },
  error: { margin: 0, color: '#ef4444' },
};
