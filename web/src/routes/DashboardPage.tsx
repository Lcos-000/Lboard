import { FormEvent, useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

import { createRoom, listRooms, type Room } from '../api/rooms';
import { me } from '../api/auth';
import { useAppStore } from '../store/appStore';

export default function DashboardPage() {
  const appName = useAppStore((state) => state.appName);
  const user = useAppStore((state) => state.user);
  const setUser = useAppStore((state) => state.setUser);

  const [rooms, setRooms] = useState<Room[]>([]);
  const [name, setName] = useState('我的白板');
  const [description, setDescription] = useState('Phase 1 创建的房间');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function loadRooms() {
    try {
      const meResult = await me();
      setUser(meResult.user);

      const result = await listRooms();
      setRooms(result.rooms ?? []);
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载失败，请先登录');
    }
  }

  useEffect(() => {
    void loadRooms();
  }, []);

  async function handleCreateRoom(e: FormEvent) {
    e.preventDefault();

    setLoading(true);
    setError('');

    try {
      await createRoom({ name, description });
      await loadRooms();
    } catch (err) {
      setError(err instanceof Error ? err.message : '创建房间失败');
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="page-card">
      <h1>{appName}</h1>

      {user ? (
        <p>
          当前用户：{user.username} / {user.email}
        </p>
      ) : (
        <p>
          未登录，请先 <Link to="/login">登录</Link> 或{' '}
          <Link to="/register">注册</Link>。
        </p>
      )}

      <form className="create-room-form" onSubmit={handleCreateRoom}>
        <h2>创建房间</h2>

        <input
          placeholder="房间名称"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />

        <input
          placeholder="房间描述"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />

        <button type="submit" disabled={loading}>
          {loading ? '创建中...' : '创建房间'}
        </button>
      </form>

      {error && <p className="error-text">{error}</p>}

      <div className="room-list">
        {rooms.length === 0 ? (
          <p>暂无房间。</p>
        ) : (
          rooms.map((room) => (
            <div className="room-card" key={room.id}>
              <h2>{room.name}</h2>
              <p>{room.description || '暂无描述'}</p>
              <p className="muted">Room ID: {room.id}</p>
              <Link to={`/board/${room.id}`}>进入白板</Link>
            </div>
          ))
        )}
      </div>
    </section>
  );
}
