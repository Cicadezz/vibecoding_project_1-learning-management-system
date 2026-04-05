import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';

import { DashboardPage } from '../pages/DashboardPage';

describe('DashboardPage', () => {
  it('renders the key dashboard metrics', () => {
    render(<DashboardPage />);

    expect(screen.getByText('今日学习总时长')).toBeTruthy();
    expect(screen.getByText('120 分钟')).toBeTruthy();
    expect(screen.getByText('连续打卡 5 天')).toBeTruthy();
  });
});
