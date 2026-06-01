import { Link } from 'react-router-dom';

import { useAppStore } from '../store/appStore';

export default function DashboardPage() {
  const appName = useAppStore((state) => state.appName);

  return (
    <section className="page-card">
      <h1>{appName}</h1>
      <p>Phase 0：工作台页面占位。</p>

      <div className="room-list">
        <div className="room-card">
          <h2>Demo Room</h2>
          <p>用于验证前端路由和白板页面是否正常。</p>
          <Link to="/board/demo-room">进入白板</Link>
        </div>
      </div>
    </section>
  );
}
