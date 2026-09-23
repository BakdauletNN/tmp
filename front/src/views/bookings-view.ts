import { ApiError, deleteBooking, getBookings, getRoom } from '../api';
import type { Booking, Room } from '../types';
import { emptyStateHtml } from '../components/empty-state';
import { showToast } from '../components/toast';

export function renderBookingsView(container: HTMLElement): void {
  container.innerHTML = `
    <div class="page">
      <header class="page__header">
        <div>
          <h1>Мои брони</h1>
          <p class="page__subtitle">Все ваши бронирования в одном месте.</p>
        </div>
      </header>
      <div id="bookings-list" class="stack stack--bookings">
        <p class="loading-note">Загружаем брони…</p>
      </div>
    </div>
  `;

  const listEl = container.querySelector<HTMLElement>('#bookings-list')!;
  void load(listEl);
}

async function load(listEl: HTMLElement): Promise<void> {
  listEl.innerHTML = '<p class="loading-note">Загружаем брони…</p>';
  try {
    const bookings = await getBookings();
    await renderBookings(listEl, bookings);
  } catch (err) {
    const message = err instanceof ApiError ? err.message : 'Не удалось загрузить брони.';
    showToast(message);
    listEl.innerHTML = emptyStateHtml({ title: 'Не получилось загрузить брони', message });
  }
}

async function renderBookings(listEl: HTMLElement, bookings: Booking[]): Promise<void> {
  if (bookings.length === 0) {
    listEl.innerHTML = emptyStateHtml({
      title: 'Броней пока нет',
      message: 'Выберите коворкинг и комнату, чтобы сделать первую бронь.',
    });
    return;
  }

  // Брони отдают только room_id, поэтому подтягиваем комнаты параллельно,
  // а если какая-то не найдётся — просто показываем её номер.
  const rooms = await Promise.all(
    bookings.map((booking) => getRoom(booking.room_id).catch(() => null)),
  );

  listEl.innerHTML = bookings.map((booking, i) => bookingRowHtml(booking, rooms[i])).join('');

  listEl.querySelectorAll<HTMLButtonElement>('[data-action="cancel"]').forEach((btn) => {
    btn.addEventListener('click', () => void cancelBooking(listEl, Number(btn.dataset.bookingId)));
  });
}

async function cancelBooking(listEl: HTMLElement, bookingId: number): Promise<void> {
  const btn = listEl.querySelector<HTMLButtonElement>(`[data-booking-id="${bookingId}"]`);
  if (btn) {
    btn.disabled = true;
    btn.textContent = 'Отменяем…';
  }
  try {
    const res = await deleteBooking(bookingId);
    showToast(res.message || 'Бронь отменена.', 'success');
    await load(listEl);
  } catch (err) {
    const message = err instanceof ApiError ? err.message : 'Не удалось отменить бронь.';
    showToast(message);
    if (btn) {
      btn.disabled = false;
      btn.textContent = 'Отменить';
    }
  }
}

function bookingRowHtml(booking: Booking, room: Room | null): string {
  const typeLabel = room
    ? room.type === 'meeting_room'
      ? 'Переговорная'
      : room.type === 'open_space'
        ? 'Open space'
        : room.type
    : `Комната #${booking.room_id}`;

  return `
    <div class="booking-row">
      <div>
        <p class="booking-row__title">${escapeHtml(typeLabel)}</p>
        <p class="booking-row__time">${formatRange(booking.start_time, booking.end_time)}</p>
      </div>
      <button type="button" class="btn btn--ghost" data-action="cancel" data-booking-id="${booking.id}">
        Отменить
      </button>
    </div>
  `;
}

function formatRange(start: string, end: string): string {
  const format = (iso: string) =>
    new Date(iso).toLocaleString('ru-RU', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' });
  return `${format(start)} → ${format(end)}`;
}

function escapeHtml(value: string): string {
  const div = document.createElement('div');
  div.textContent = value;
  return div.innerHTML;
}
