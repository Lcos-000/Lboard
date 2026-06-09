import { request } from './http';

export interface Room {
  id: string;
  name: string;
  description: string;
  ownerId: string;
  createdAt: string;
  updatedAt: string;
}

export async function createRoom(input: {
  name: string;
  description: string;
}) {
  return request<{ room: Room }>('/rooms', {
    method: 'POST',
    body: JSON.stringify(input)
  });
}

export async function listRooms() {
  return request<{ rooms: Room[] }>('/rooms');
}

export async function getRoom(roomId: string) {
  return request<{ room: Room }>(`/rooms/${roomId}`);
}

export async function inviteMember(
  roomId: string,
  input: {
    email: string;
    role: string;
  }
) {
  return request<{ status: string }>(`/rooms/${roomId}/members`, {
    method: 'POST',
    body: JSON.stringify(input)
  });
}
