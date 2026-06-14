import { useEffect } from 'react';

import { wsClient } from '../ws/WebSocketClient';
import { useAppStore } from '../store/appStore';

export default function WSStatusPanel() {
  const wsStatus = useAppStore((state) => state.wsStatus);
  const wsMessages = useAppStore((state) => state.wsMessages);
  const setWSStatus = useAppStore((state) => state.setWSStatus);
  const pushWSMessage = useAppStore((state) => state.pushWSMessage);
  const clearWSMessages = useAppStore((state) => state.clearWSMessages);

  useEffect(() => {
    const offStatus = wsClient.onStatus(setWSStatus);
    const offMessage = wsClient.onMessage(pushWSMessage);

    return () => {
      offStatus();
      offMessage();
    };
  }, [setWSStatus, pushWSMessage]);

  function handleConnect() {
    wsClient.connect();
  }

  function handleClose() {
    wsClient.close();
  }

  async function handlePing() {
    try {
      await wsClient.request({
        type: 'ping',
        payload: {
          from: 'frontend',
          now: Date.now()
        }
      });
    } catch (err) {
      pushWSMessage({
        type: 'client_error',
        payload: {
          message: err instanceof Error ? err.message : 'ping failed'
        }
      });
    }
  }

  return (
    <div className="ws-panel">
      <div className="ws-panel-header">
        <strong>WebSocket</strong>
        <span className={`ws-status ws-status-${wsStatus}`}>{wsStatus}</span>
      </div>

      <div className="ws-actions">
        <button type="button" onClick={handleConnect}>
          连接
        </button>
        <button type="button" onClick={handlePing}>
          Ping
        </button>
        <button type="button" onClick={handleClose}>
          关闭
        </button>
        <button type="button" onClick={clearWSMessages}>
          清空消息
        </button>
      </div>

      <div className="ws-messages">
        {wsMessages.length === 0 ? (
          <p className="muted">暂无消息</p>
        ) : (
          wsMessages.map((msg, index) => (
            <pre key={index}>{JSON.stringify(msg, null, 2)}</pre>
          ))
        )}
      </div>
    </div>
  );
}
