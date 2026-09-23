import { getCurrentEmail, isAuthenticated, logout, onAuthChange } from '../auth';
import { navigate } from '../router';

export function renderNavbar(container: HTMLElement): void {
  const draw = () => {
    const authed = isAuthenticated();
    container.innerHTML = `
      <nav class="navbar">
        <a href="/offices" data-link class="navbar__logo">
          <span class="navbar__logo-mark" aria-hidden="true"></span>
          CoworkGo
        </a>
        <div class="navbar__links">
          <a href="/offices" data-link class="navbar__link">Коворкинги</a>
          ${authed ? '<a href="/bookings" data-link class="navbar__link">Мои брони</a>' : ''}
        </div>
        <div class="navbar__session">
          ${
            authed
              ? `<span class="navbar__email">${escapeHtml(getCurrentEmail() ?? '')}</span>
                 <button type="button" class="btn btn--ghost" data-action="logout">Выйти</button>`
              : `<a href="/login" data-link class="btn btn--primary">Войти</a>`
          }
        </div>
      </nav>
    `;

    const logoutBtn = container.querySelector<HTMLButtonElement>('[data-action="logout"]');
    logoutBtn?.addEventListener('click', () => {
      logout();
      navigate('/offices');
    });
  };

  draw();
  onAuthChange(draw);
}

function escapeHtml(value: string): string {
  const div = document.createElement('div');
  div.textContent = value;
  return div.innerHTML;
}
