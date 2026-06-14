import { FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { register } from '../../api/auth';
import { useAppStore } from '../../store/appStore';

export default function RegisterPage() {
  const navigate = useNavigate();
  const setUser = useAppStore((state) => state.setUser);

  const [username, setUsername] = useState('alice');
  const [email, setEmail] = useState('alice@example.com');
  const [password, setPassword] = useState('123456');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();

    setLoading(true);
    setError('');

    try {
      const result = await register({
        username,
        email,
        password
      });

      setUser(result.user);
      navigate('/dashboard');
    } catch (err) {
      setError(err instanceof Error ? err.message : '注册失败');
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="page-card">
      <h1>注册</h1>
      <p>创建一个新用户。</p>

      <form className="form" onSubmit={handleSubmit}>
        <input
          placeholder="用户名"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />

        <input
          placeholder="邮箱"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />

        <input
          placeholder="密码，至少 6 位"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />

        {error && <p className="error-text">{error}</p>}

        <button type="submit" disabled={loading}>
          {loading ? '注册中...' : '注册'}
        </button>
      </form>
    </section>
  );
}
