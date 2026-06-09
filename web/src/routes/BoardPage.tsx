import { FormEvent, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

import { getRoom, inviteMember, type Room } from '../api/rooms';

export default function BoardPage() {
  const { roomId } = useParams();

  const [room, setRoom] = useState<Room | null>(null);
  const [inviteEmail, setInviteEmail] = useState('bob@example.com');
  const [inviteRole, setInviteRole] = useState('editor');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  async function loadRoom() {
    if (!roomId) return;

    try {
      const result = await getRoom(roomId);
      setRoom(result.room);
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载房间失败');
    }
  }

  useEffect(() => {
    void loadRoom();
  }, [roomId]);

  async function handleInvite(e: FormEvent) {
    e.preventDefault();

    if (!roomId) return;

    setError('');
    setMessage('');

    try {
      await inviteMember(roomId, {
        email: inviteEmail,
        role: inviteRole
      });

      setMessage('邀请成功');
    } catch (err) {
      setError(err instanceof Error ? err.message : '邀请失败');
    }
  }

  return (
    <section className="board-page">
      <aside className="toolbar">
        <h2>房间信息</h2>

        {room ? (
          <>
            <p>名称：{room.name}</p>
            <p>描述：{room.description || '暂无描述'}</p>
            <p className="muted">Room ID: {room.id}</p>
          </>
        ) : (
          <p>加载中...</p>
        )}

        <hr />

        <h2>邀请成员</h2>

        <form className="invite-form" onSubmit={handleInvite}>
          <input
            placeholder="成员邮箱"
            value={inviteEmail}
            onChange={(e) => setInviteEmail(e.target.value)}
          />

          <select
            value={inviteRole}
            onChange={(e) => setInviteRole(e.target.value)}
          >
            <option value="viewer">viewer</option>
            <option value="editor">editor</option>
            <option value="admin">admin</option>
          </select>

          <button type="submit">邀请</button>
        </form>

        {message && <p className="success-text">{message}</p>}
        {error && <p className="error-text">{error}</p>}
      </aside>

      <div className="board-placeholder">
        <h1>白板页面占位</h1>
        <p>Phase 1：已完成房间权限校验。</p>
        <p>Phase 2 才会接入 WebSocket。</p>
      </div>
    </section>
  );
}
