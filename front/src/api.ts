import { getToken } from './auth';
import type {
  AuthResponse,
  Booking,
  CreateBookingPayload,
  CreateBookingResponse,
  Office,
  OfficeFilters,
  Room,
  RoomFilters,
  User,
} from './types';

const API_BASE = 'http://localhost:8080';

/** Ошибка запроса к API. Сообщение — ровно то, что вернул бэкенд в поле "error". */
export class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

function buildQuery(params: Record<string, unknown> | undefined): string {
  if (!params) return '';
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null || value === '') continue;
    search.set(key, String(value));
  }
  const query = search.toString();
  return query ? `?${query}` : '';
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getToken();
  const headers = new Headers(options.headers);
  headers.set('Content-Type', 'application/json');
  if (token) headers.set('Authorization', `Bearer ${token}`);

  let response: Response;
  try {
    response = await fetch(`${API_BASE}${path}`, { ...options, headers });
  } catch {
    throw new ApiError('Не удалось соединиться с сервером. Проверьте, что бэкенд запущен.', 0);
  }

  const raw = await response.text();
  const data = raw ? safeJsonParse(raw) : null;

  if (!response.ok) {
    const message =
      data && typeof data === 'object' && typeof (data as { error?: unknown }).error === 'string'
        ? (data as { error: string }).error
        : `Ошибка запроса (${response.status})`;
    throw new ApiError(message, response.status);
  }

  return data as T;
}

function safeJsonParse(raw: string): unknown {
  try {
    return JSON.parse(raw);
  } catch {
    return null;
  }
}

// ---- Авторизация ----

export function register(name: string, email: string, password: string): Promise<User> {
  return request<User>('/register', {
    method: 'POST',
    body: JSON.stringify({ name, email, password }),
  });
}

export function login(email: string, password: string): Promise<AuthResponse> {
  return request<AuthResponse>('/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
}

// ---- Коворкинги ----

export function getOffices(filters?: OfficeFilters): Promise<Office[]> {
  return request<Office[]>(`/offices${buildQuery(filters as Record<string, unknown>)}`);
}

export function getOffice(id: number): Promise<Office> {
  return request<Office>(`/offices/${id}`);
}

// ---- Комнаты ----

export function getRooms(filters?: RoomFilters): Promise<Room[]> {
  return request<Room[]>(`/rooms${buildQuery(filters as Record<string, unknown>)}`);
}

export function getRoom(id: number): Promise<Room> {
  return request<Room>(`/rooms/${id}`);
}

// ---- Брони (требуют JWT — токен подставляется в request() автоматически) ----

export function getBookings(): Promise<Booking[]> {
  return request<Booking[]>('/bookings');
}

export function createBooking(payload: CreateBookingPayload): Promise<CreateBookingResponse> {
  return request<CreateBookingResponse>('/create_booking', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function deleteBooking(id: number): Promise<{ message: string }> {
  return request<{ message: string }>(`/bookings/${id}`, { method: 'DELETE' });
}
