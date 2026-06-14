import { create } from 'zustand';

import type { User } from '../api/auth';
import type { WSMessage, WSStatus } from '../ws/types';

interface AppState {
  appName: string;
  user: User | null;

  wsStatus: WSStatus;
  wsMessages: WSMessage[];

  setUser: (user: User | null) => void;
  setWSStatus: (status: WSStatus) => void;
  pushWSMessage: (msg: WSMessage) => void;
  clearWSMessages: () => void;
}

export const useAppStore = create<AppState>((set) => ({
  appName: 'Whiteboard Dashboard',
  user: null,

  wsStatus: 'idle',
  wsMessages: [],

  setUser: (user) => set({ user }),

  setWSStatus: (status) => set({ wsStatus: status }),

  pushWSMessage: (msg) =>
    set((state) => ({
      wsMessages: [msg, ...state.wsMessages].slice(0, 20)
    })),

  clearWSMessages: () => set({ wsMessages: [] })
}));

