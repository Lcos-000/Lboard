export type WSStatus =
  | 'idle'
  | 'connecting'
  | 'connected'
  | 'reconnecting'
  | 'closed'
  | 'error';

export interface WSMessage<TPayload = unknown> {
  type: string;
  requestId?: string;
  roomId?: string;
  payload?: TPayload;
}

export interface PendingRequest {
  resolve: (msg: WSMessage) => void;
  reject: (err: Error) => void;
  timeoutId: number;
}
