import { ApiError, createBooking } from '../api';
import { showToast } from './toast';

interface BookingFormOptions {
  roomId: number;
  onBooked: () => void;
}

/**
 * Рисует форму бронирования комнаты целиком (desk_id всегда null — карта
 * отдельных мест появится позже, когда на бэкенде будет эндпоинт desks).
 */
export function renderBookingForm(container: HTMLElement, { roomId, onBooked }: BookingFormOptions): void {
  container.innerHTML = `
    <form class="booking-form" data-room-id="${roomId}">
      <div class="field">
        <label for="start-${roomId}">Начало</label>
        <input id="start-${roomId}" name="start" type="datetime-local" required />
      </div>
      <div class="field">
        <label for="end-${roomId}">Конец</label>
        <input id="end-${roomId}" name="end" type="datetime-local" required />
      </div>
      <button type="submit" class="btn btn--primary">Забронировать</button>
    </form>
  `;

  const form = container.querySelector('form')!;
  form.addEventListener('submit', async (event) => {
    event.preventDefault();

    const startInput = form.querySelector<HTMLInputElement>('input[name="start"]')!;
    const endInput = form.querySelector<HTMLInputElement>('input[name="end"]')!;
    const start = toRfc3339(startInput.value);
    const end = toRfc3339(endInput.value);

    if (!start || !end) {
      showToast('Укажите дату и время начала и конца брони.');
      return;
    }
    if (new Date(end) <= new Date(start)) {
      showToast('Время окончания должно быть позже времени начала.');
      return;
    }

    const submitBtn = form.querySelector<HTMLButtonElement>('button[type="submit"]')!;
    submitBtn.disabled = true;
    submitBtn.textContent = 'Бронируем…';

    try {
      const res = await createBooking({ room_id: roomId, desk_id: null, start_time: start, end_time: end });
      showToast(res.message || 'Бронь создана.', 'success');
      form.reset();
      onBooked();
    } catch (err) {
      const message = err instanceof ApiError ? err.message : 'Не удалось создать бронь.';
      showToast(message);
    } finally {
      submitBtn.disabled = false;
      submitBtn.textContent = 'Забронировать';
    }
  });
}

/** datetime-local ("YYYY-MM-DDTHH:mm", локальное время) -> RFC3339 в UTC. */
function toRfc3339(localValue: string): string | null {
  if (!localValue) return null;
  const date = new Date(localValue);
  if (Number.isNaN(date.getTime())) return null;
  return date.toISOString();
}
