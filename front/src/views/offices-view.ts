import { ApiError, getOffices } from '../api';
import type { Office } from '../types';
import { emptyStateHtml } from '../components/empty-state';
import { showToast } from '../components/toast';

export function renderOfficesView(container: HTMLElement): void {
  let addressFilter = '';

  container.innerHTML = `
    <div class="page">
      <header class="page__header">
        <div>
          <h1>Коворкинги</h1>
          <p class="page__subtitle">Выберите пространство и переходите к комнатам.</p>
        </div>
        <form class="filter-bar" id="office-filter">
          <input type="text" name="address" placeholder="Фильтр по адресу…" aria-label="Фильтр по адресу" />
          <button type="submit" class="btn btn--secondary">Найти</button>
        </form>
      </header>
      <div id="offices-list" class="grid grid--offices">
        <p class="loading-note">Загружаем коворкинги…</p>
      </div>
    </div>
  `;

  const listEl = container.querySelector<HTMLElement>('#offices-list')!;
  const filterForm = container.querySelector<HTMLFormElement>('#office-filter')!;

  filterForm.addEventListener('submit', (event) => {
    event.preventDefault();
    addressFilter = new FormData(filterForm).get('address')?.toString().trim() ?? '';
    void load();
  });

  async function load(): Promise<void> {
    listEl.innerHTML = '<p class="loading-note">Загружаем коворкинги…</p>';
    try {
      const offices = await getOffices(addressFilter ? { address: addressFilter } : undefined);
      renderList(listEl, offices, addressFilter);
    } catch (err) {
      const message = err instanceof ApiError ? err.message : 'Не удалось загрузить список коворкингов.';
      showToast(message);
      listEl.innerHTML = emptyStateHtml({
        title: 'Не получилось загрузить список',
        message,
      });
    }
  }

  void load();
}

function renderList(listEl: HTMLElement, offices: Office[], addressFilter: string): void {
  if (offices.length === 0) {
    listEl.innerHTML = emptyStateHtml(
      addressFilter
        ? { title: 'Ничего не нашлось', message: `По адресу «${addressFilter}» коворкингов пока нет.` }
        : { title: 'Коворкингов пока нет', message: 'Как только появятся новые пространства, они отобразятся здесь.' },
    );
    return;
  }

  listEl.innerHTML = offices.map(officeCardHtml).join('');
}

function officeCardHtml(office: Office): string {
  return `
    <a href="/offices/${office.id}" data-link class="card card--office">
      <div class="card__top">
        <h2 class="card__title">${escapeHtml(office.name)}</h2>
        <span class="tag tag--star">★ ${office.star}</span>
      </div>
      <p class="card__address">${escapeHtml(office.address)}</p>
      <dl class="card__meta">
        <div><dt>Часы работы</dt><dd>${escapeHtml(office.time_range_work)}</dd></div>
        <div><dt>Метро</dt><dd>${escapeHtml(office.metro_near)}</dd></div>
      </dl>
      ${office.has_kitchen ? '<span class="tag tag--amenity">Кухня</span>' : ''}
    </a>
  `;
}

function escapeHtml(value: string): string {
  const div = document.createElement('div');
  div.textContent = value;
  return div.innerHTML;
}
