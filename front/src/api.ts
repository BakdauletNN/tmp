import type { AdminUser, Booking, CreateBookingBody, Desk, Office, Room } from './types';

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
  }
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  headers.set('Content-Type', 'application/json');
  const token = localStorage.getItem('coworkgo-token');
  if (token) headers.set('Authorization', `Bearer ${token}`);
  let response: Response;
  try {
    response = await fetch(`${API_BASE}${path}`, { ...init, headers });
  } catch {
    throw new ApiError(0, `Could not connect to the server ${API_BASE}`);
  }
  const data: unknown = await response.json().catch(() => null);
  if (!response.ok) {
    const raw = (data as { error?: string } | null)?.error;
    const messages: Record<string, string> = {
      'place already booked': 'This time slot is already taken. Choose a different interval.',
      unauthorized: 'Please sign in to continue.',
      'invalid token': 'Your session has ended. Please sign in again.',
    };
    throw new ApiError(response.status, messages[raw ?? ''] ?? raw ?? `Error ${response.status}`);
  }
  return data as T;
}

function list<T>(data: T[] | { offices?: T[]; rooms?: T[] } | null, key: 'offices' | 'rooms'): T[] {
  if (Array.isArray(data)) return data;
  return data?.[key] ?? [];
}

function query(values: Record<string, string | number | boolean | undefined>): string {
  const params = new URLSearchParams();
  Object.entries(values).forEach(([key, value]) => {
    if (value !== undefined && value !== '' && value !== false) params.set(key, String(value));
  });
  const encoded = params.toString();
  return encoded ? `?${encoded}` : '';
}

export const api = {
  offices: async (filters: Record<string, string | number | boolean | undefined> = {}) =>
    list(await request<Office[] | { offices?: Office[] }>('/offices' + query(filters)), 'offices'),
  office: (id: number) => request<Office>(`/offices/${id}`),
  rooms: async (filters: Record<string, string | number | boolean | undefined> = {}) =>
    list(await request<Room[] | { rooms?: Room[] }>('/rooms' + query(filters)), 'rooms'),
  room: (id: number) => request<Room>(`/rooms/${id}`),
  desks: (roomId: number, start: string, end: string) => request<Desk[]>(
    `/rooms/${roomId}/desks` + query({ start_time: start, end_time: end }),
  ),
  bookings: async () => (await request<Booking[] | null>('/bookings')) ?? [],
  createBooking: (body: CreateBookingBody) =>
    request<{ booking: Booking }>('/create_booking', { method: 'POST', body: JSON.stringify(body) }),
  cancelBooking: (id: number) => request<{ message: string }>(`/bookings/${id}`, { method: 'DELETE' }),
  login: (email: string, password: string) =>
    request<{ token: string }>('/login', { method: 'POST', body: JSON.stringify({ email, password }) }),
  register: (name: string, email: string, password: string) =>
    request<{ id: number; name: string; email: string }>('/register', { method: 'POST', body: JSON.stringify({ name, email, password }) }),
  forgotPassword: (email: string) =>
    request<{ message: string }>('/forgot-password', { method: 'POST', body: JSON.stringify({ email }) }),
  resetPassword: (token: string, password: string) =>
    request<{ message: string }>('/reset-password', { method: 'POST', body: JSON.stringify({ token, password }) }),
  adminOffices: async () => (await request<Office[] | null>('/admin/offices')) ?? [],
  adminBookings: async () => (await request<Booking[] | null>('/admin/bookings')) ?? [],
  adminUsers: async () => (await request<AdminUser[] | null>('/admin/users')) ?? [],
  adminCancelBooking: (id: number) => request<{ message: string }>(`/admin/bookings/${id}`, { method: 'DELETE' }),
};
