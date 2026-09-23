interface EmptyStateOptions {
  title: string;
  message: string;
}

/** Рисует аккуратную заглушку вместо пустого списка (без белого экрана). */
export function emptyStateHtml({ title, message }: EmptyStateOptions): string {
  return `
    <div class="empty-state">
      <svg class="empty-state__glyph" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
        <rect x="8" y="14" width="48" height="36" rx="2" stroke="currentColor" stroke-width="2"/>
        <path d="M8 24H56" stroke="currentColor" stroke-width="2"/>
        <path d="M20 32H44M20 40H36" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
      <p class="empty-state__title">${title}</p>
      <p class="empty-state__message">${message}</p>
    </div>
  `;
}
