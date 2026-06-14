import { FormEvent, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

import WSStatusPanel from '../WSStatusPanel';
import { getRoom, inviteMember, type Room } from '../../api/rooms';
import { wsClient } from '../../ws/WebSocketClient';
import { useAppStore } from '../../store/appStore';

export default function BoardPage() {
  const { roomId } = useParams();

  const pushWSMessage = useAppStore((state) => state.pushWSMessage);

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

  useEffect(() => {
    if (!roomId) return;

    wsClient.connect();

    const offStatus = wsClient.onStatus(async (status) => {
      if (status !== 'connected') return;

      try {
        const ack = await wsClient.request({
          type: 'join_room',
          roomId
        });

        pushWSMessage(ack);
      } catch (err) {
        pushWSMessage({
          type: 'client_error',
          roomId,
          payload: {
            message: err instanceof Error ? err.message : 'join room failed'
          }
        });
      }
    });

    return () => {
      offStatus();
    };
  }, [roomId, pushWSMessage]);

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

        <hr />

        <WSStatusPanel />
      </aside>

      <div className="board-placeholder">
        <h1>白板页面占位</h1>
        <p>Phase 2：WebSocket Gateway 已接入。</p>
        <p>当前阶段已支持连接、鉴权、心跳、自动重连、join room。</p>
        <p>Phase 3 才会接入 Room Actor 和房间内广播。</p>
      </div>
    </section>
  );
}
