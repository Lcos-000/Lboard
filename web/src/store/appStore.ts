import { create } from 'zustand';

import type { User } from '../api/auth';

interface AppState {
  appName: string;
  user: User | null;
  setUser: (user: User | null) => void;
}

export const useAppStore = create<AppState>((set) => ({
  appName: 'Whiteboard Dashboard',
  user: null,
  setUser: (user) => set({ user })
}));
