import { getToken } from '../api/http';
import type { PendingRequest, WSMessage, WSStatus } from './types';

type MessageListener = (msg: WSMessage) => void;
type StatusListener = (status: WSStatus) => void;

export class WebSocketClient {
  private ws: WebSocket | null = null;
  private status: WSStatus = 'idle';

  private reconnectTimer: number | null = null;
  private reconnectAttempt = 0;
  private manuallyClosed = false;

  private pending = new Map<string, PendingRequest>();
  private messageListeners = new Set<MessageListener>();
  private statusListeners = new Set<StatusListener>();

  connect() {
    // 如果已经连接或正在连接中，不重复创建
    if (this.ws) {
        if (this.ws.readyState === WebSocket.OPEN) {
        console.log('WebSocket already connected, skip');
        return;
        }
        if (this.ws.readyState === WebSocket.CONNECTING) {
        console.log('WebSocket already connecting, skip');
        return;
        }
        // readyState 是 CLOSING 或 CLOSED 时，继续往下走，重新连接
    }
    
    // 清理 pending 请求
    this.rejectAllPending(new Error('重新连接'));
    
    // 清除重连定时器
    if (this.reconnectTimer !== null) {
        window.clearTimeout(this.reconnectTimer);
        this.reconnectTimer = null;
    }

    const token = getToken();

    if (!token) {
      this.setStatus('error');
      return;
    }

    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return;
    }

    this.manuallyClosed = false;

    if (this.status === 'connected') {
      return;
    }

    this.setStatus(this.reconnectAttempt > 0 ? 'reconnecting' : 'connecting');

    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const host = window.location.host;

    this.ws = new WebSocket(`${protocol}://${host}/ws?token=${encodeURIComponent(token)}`);

    this.ws.onopen = () => {
      this.reconnectAttempt = 0;
      this.setStatus('connected');
    };

    this.ws.onmessage = (event) => {
      this.handleMessage(event.data);
    };

    this.ws.onerror = () => {
      this.setStatus('error');
    };

    this.ws.onclose = () => {
      this.ws = null;

      if (this.manuallyClosed) {
        this.setStatus('closed');
        return;
      }

      this.scheduleReconnect();
    };
  }

  close() {
    this.manuallyClosed = true;

    if (this.reconnectTimer !== null) {
      window.clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    this.rejectAllPending(new Error('websocket closed'));

    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }

    this.setStatus('closed');
  }

  send(msg: WSMessage) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error('websocket is not connected');
    }

    this.ws.send(JSON.stringify(msg));
  }

  request(msg: WSMessage, timeoutMs = 5000): Promise<WSMessage> {
    if (!msg.requestId) {
      msg.requestId = this.newRequestId();
    }

    this.send(msg);

    return new Promise((resolve, reject) => {
      const timeoutId = window.setTimeout(() => {
        this.pending.delete(msg.requestId!);
        reject(new Error('websocket request timeout'));
      }, timeoutMs);

      this.pending.set(msg.requestId!, {
        resolve,
        reject,
        timeoutId
      });
    });
  }

  onMessage(listener: MessageListener) {
    this.messageListeners.add(listener);

    return () => {
      this.messageListeners.delete(listener);
    };
  }

  onStatus(listener: StatusListener) {
    this.statusListeners.add(listener);

    listener(this.status);

    return () => {
      this.statusListeners.delete(listener);
    };
  }

  getStatus() {
    return this.status;
  }

  private handleMessage(raw: string) {
    let msg: WSMessage;

    try {
      msg = JSON.parse(raw) as WSMessage;
    } catch {
      return;
    }

    if (msg.requestId) {
      const pending = this.pending.get(msg.requestId);
      if (pending) {
        window.clearTimeout(pending.timeoutId);
        this.pending.delete(msg.requestId);

        if (msg.type === 'error') {
          const payload = msg.payload as { message?: string } | undefined;
          pending.reject(new Error(payload?.message ?? 'websocket request failed'));
        } else {
          pending.resolve(msg);
        }
      }
    }

    for (const listener of this.messageListeners) {
      listener(msg);
    }
  }

  private scheduleReconnect() {
    this.reconnectAttempt += 1;
    this.setStatus('reconnecting');

    const delay = Math.min(1000 * this.reconnectAttempt, 5000);

    this.reconnectTimer = window.setTimeout(() => {
      this.connect();
    }, delay);
  }

  private setStatus(status: WSStatus) {
    this.status = status;

    for (const listener of this.statusListeners) {
      listener(status);
    }
  }

  private rejectAllPending(err: Error) {
    for (const pending of this.pending.values()) {
      window.clearTimeout(pending.timeoutId);
      pending.reject(err);
    }

    this.pending.clear();
  }

  private newRequestId() {
    return `req-${Date.now()}-${Math.random().toString(16).slice(2)}`;
  }
}

export const wsClient = new WebSocketClient();
