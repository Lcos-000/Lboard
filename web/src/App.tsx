import { Link, Navigate, Route, Routes } from 'react-router-dom';

import LoginPage from './routes/LoginPage';
import RegisterPage from './routes/RegisterPage';
import DashboardPage from './routes/DashboardPage';
import BoardPage from './routes/BoardPage';

export default function App() {
  return (
    <div className="app">
      <header className="app-header">
        <Link to="/" className="brand">
          Whiteboard
        </Link>

        <nav className="nav">
          <Link to="/login">登录</Link>
          <Link to="/register">注册</Link>
          <Link to="/dashboard">工作台</Link>
          <Link to="/board/demo-room">白板 Demo</Link>
        </nav>
      </header>

      <main className="app-main">
        <Routes>
          {/* 路由配置，这里重定向到工作台 */}
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/board/:roomId" element={<BoardPage />} />
        </Routes>
      </main>
    </div>
  );
}
