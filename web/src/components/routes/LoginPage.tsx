import { FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { login } from '../../api/auth';
import { useAppStore } from '../../store/appStore';

export default function LoginPage() {
  const navigate = useNavigate();
  const setUser = useAppStore((state) => state.setUser);

  const [email, setEmail] = useState('alice@example.com');
  const [password, setPassword] = useState('123456');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();

    setLoading(true);
    setError('');

    try {
      const result = await login({ email, password });
      setUser(result.user);
      navigate('/dashboard');
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败');
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="page-card">
      <h1>登录</h1>
      <p>使用已注册账号登录。</p>

      <form className="form" onSubmit={handleSubmit}>
        <input
          placeholder="邮箱"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />

        <input
          placeholder="密码"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />

        {error && <p className="error-text">{error}</p>}

        <button type="submit" disabled={loading}>
          {loading ? '登录中...' : '登录'}
        </button>
      </form>
    </section>
  );
}

