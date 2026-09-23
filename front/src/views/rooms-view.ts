import { ApiError, getOffice, getRooms } from '../api';
import { isAuthenticated } from '../auth';
import type { Office, Room, RoomFilters } from '../types';
import { emptyStateHtml } from '../components/empty-state';
import { renderBookingForm } from '../components/booking-form';
import { showToast } from '../components/toast';

export function renderRoomsView(container: HTMLElement, officeId: number): void {
  container.innerHTML = `
    <div class="page">
      <header class="page__header page__header--office" id="office-header">
        <p class="loading-note">Загружаем коворкинг…</p>
      </header>
      <form class="filter-bar filter-bar--rooms" id="room-filter">
        <select name="type" aria-label="Тип комнаты">
          <option value="">Любой тип</option>
          <option value="meeting_room">Переговорная</option>
          <option value="open_space">Open space</option>
        </select>
        <input type="number" name="max_price_hour" min="0" placeholder="Цена/час, до" aria-label="Максимальная цена в час" />
        <input type="number" name="min_qty_person" min="1" placeholder="Мест, от" aria-label="Минимальная вместимость" />
        <label class="checkbox">
          <input type="checkbox" name="has_air_conditioner" />
          С кондиционером
        </label>
        <button type="submit" class="btn btn--secondary">Применить</button>
      </form>
      <div id="rooms-list" class="grid grid--rooms">
        <p class="loading-note">Загружаем комнаты…</p>
      </div>
    </div>
  `;

  const headerEl = container.querySelector<HTMLElement>('#office-header')!;
  const listEl = container.querySelector<HTMLElement>('#rooms-list')!;
  const filterForm = container.querySelector<HTMLFormElement>('#room-filter')!;

  loadOffice(headerEl, officeId);

  filterForm.addEventListener('submit', (event) => {
    event.preventDefault();
    void loadRooms(listEl, officeId, readFilters(filterForm));
  });

  void loadRooms(listEl, officeId, {});
}

async function loadOffice(headerEl: HTMLElement, officeId: number): Promise<void> {
  try {
    const office = await getOffice(officeId);
    headerEl.innerHTML = officeHeaderHtml(office);
  } catch (err) {
    const message = err instanceof ApiError ? err.message : 'Не удалось загрузить коворкинг.';
    headerEl.innerHTML = `<h1>Комнаты</h1><p class="page__subtitle">${escapeHtml(message)}</p>`;
  }
}

function readFilters(form: HTMLFormElement): RoomFilters {
  const data = new FormData(form);
  const type = data.get('type')?.toString() ?? '';
  const maxPrice = data.get('max_price_hour')?.toString() ?? '';
  const minPeople = data.get('min_qty_person')?.toString() ?? '';
  return {
    type: type || undefined,
    max_price_hour: maxPrice ? Number(maxPrice) : undefined,
    min_qty_person: minPeople ? Number(minPeople) : undefined,
    has_air_conditioner: form.querySelector<HTMLInputElement>('[name="has_air_conditioner"]')!.checked || undefined,
  };
}

async function loadRooms(listEl: HTMLElement, officeId: number, filters: RoomFilters): Promise<void> {
  listEl.innerHTML = '<p class="loading-note">Загружаем комнаты…</p>';
  try {
    const rooms = await getRooms({ ...filters, office_id: officeId });
    renderRooms(listEl, rooms, officeId);
  } catch (err) {
    const message = err instanceof ApiError ? err.message : 'Не удалось загрузить список комнат.';
    showToast(message);
    listEl.innerHTML = emptyStateHtml({ title: 'Не получилось загрузить комнаты', message });
  }
}

function renderRooms(listEl: HTMLElement, rooms: Room[], officeId: number): void {
  if (rooms.length === 0) {
    listEl.innerHTML = emptyStateHtml({
      title: 'Комнат пока нет',
      message: 'В этом коворкинге ещё не добавили комнаты под фильтр, попробуйте изменить условия.',
    });
    return;
  }

  listEl.innerHTML = rooms.map(roomCardHtml).join('');

  rooms.forEach((room) => {
    const bookingSlot = listEl.querySelector<HTMLElement>(`[data-booking-slot="${room.id}"]`);
    if (!bookingSlot) return;

    if (!isAuthenticated()) {
      bookingSlot.innerHTML = `<a href="/login" data-link class="hint-link">Войдите, чтобы забронировать</a>`;
      return;
    }

    renderBookingForm(bookingSlot, {
      roomId: room.id,
      onBooked: () => void loadRooms(listEl, officeId, {}),
    });
  });
}

function roomCardHtml(room: Room): string {
  const typeLabel = room.type === 'meeting_room' ? 'Переговорная' : room.type === 'open_space' ? 'Open space' : room.type;
  return `
    <div class="card card--room">
      <div class="card__top">
        <h2 class="card__title">${escapeHtml(typeLabel)}</h2>
        <span class="tag tag--price">${room.price_hour} / час</span>
      </div>
      <dl class="card__meta">
        <div><dt>Мест</dt><dd>${room.qty_person}</dd></div>
        <div><dt>Столов</dt><dd>${room.qty_desks}</dd></div>
        <div><dt>Код доступа</dt><dd class="mono">${escapeHtml(room.access_code)}</dd></div>
      </dl>
      <div class="card__tags">
        ${room.has_air_conditioner ? '<span class="tag tag--amenity">Кондиционер</span>' : ''}
        ${room.has_prayer_room ? '<span class="tag tag--amenity">Молитвенная комната</span>' : ''}
      </div>
      <div class="card__booking" data-booking-slot="${room.id}"></div>
    </div>
  `;
}

function officeHeaderHtml(office: Office): string {
  return `
    <div>
      <a href="/offices" data-link class="back-link">← Все коворкинги</a>
      <h1>${escapeHtml(office.name)}</h1>
      <p class="page__subtitle">${escapeHtml(office.address)} · ${escapeHtml(office.time_range_work)}</p>
    </div>
  `;
}

function escapeHtml(value: string): string {
  const div = document.createElement('div');
  div.textContent = value;
  return div.innerHTML;
}
