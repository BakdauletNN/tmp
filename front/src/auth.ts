import * as api from './api';

// Токен намеренно хранится только в памяти модуля, а не в localStorage —
// это учебное демо, обновление страницы разлогинивает пользователя.
let token: string | null = null;
let currentEmail: string | null = null;

type Listener = () => void;
const listeners = new Set<Listener>();

function notify(): void {
  listeners.forEach((fn) => fn());
}

export function onAuthChange(fn: Listener): () => void {
  listeners.add(fn);
  return () => listeners.delete(fn);
}

export function getToken(): string | null {
  return token;
}

export function isAuthenticated(): boolean {
  return token !== null;
}

/** /login не возвращает данные пользователя — показываем email, под которым вошли. */
export function getCurrentEmail(): string | null {
  return currentEmail;
}

export async function login(email: string, password: string): Promise<void> {
  const res = await api.login(email, password);
  token = res.token;
  currentEmail = email;
  notify();
}

export async function register(name: string, email: string, password: string): Promise<void> {
  await api.register(name, email, password);
  // Бэкенд не логинит автоматически после регистрации — токена тут нет и не будет.
}

export function logout(): void {
  token = null;
  currentEmail = null;
  notify();
}
