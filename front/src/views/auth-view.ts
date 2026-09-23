import { ApiError } from '../api';
import { login, register } from '../auth';
import { navigate } from '../router';
import { showToast } from '../components/toast';

type Mode = 'login' | 'register';

export function renderAuthView(container: HTMLElement): void {
  let mode: Mode = 'login';

  const draw = () => {
    container.innerHTML = `
      <div class="auth-screen">
        <aside class="auth-screen__panel" aria-hidden="true">
          <svg viewBox="0 0 200 200" class="auth-screen__blueprint">
            <rect x="10" y="10" width="180" height="180" fill="none" stroke="currentColor" stroke-width="1"/>
            <rect x="30" y="30" width="60" height="45" fill="none" stroke="currentColor" stroke-width="1"/>
            <rect x="110" y="30" width="60" height="45" fill="none" stroke="currentColor" stroke-width="1"/>
            <rect x="30" y="95" width="140" height="30" fill="none" stroke="currentColor" stroke-width="1"/>
            <rect x="30" y="140" width="35" height="35" fill="none" stroke="currentColor" stroke-width="1"/>
            <rect x="80" y="140" width="35" height="35" fill="none" stroke="currentColor" stroke-width="1"/>
            <rect x="130" y="140" width="40" height="35" fill="none" stroke="currentColor" stroke-width="1"/>
          </svg>
          <p class="auth-screen__caption">Планировка рабочих пространств</p>
        </aside>
        <div class="auth-screen__form-wrap">
          <div class="auth-tabs" role="tablist">
            <button type="button" class="auth-tabs__tab ${mode === 'login' ? 'is-active' : ''}" data-mode="login" role="tab">Вход</button>
            <button type="button" class="auth-tabs__tab ${mode === 'register' ? 'is-active' : ''}" data-mode="register" role="tab">Регистрация</button>
          </div>
          <form class="auth-form" id="auth-form">
            ${
              mode === 'register'
                ? `<div class="field">
                     <label for="name">Имя</label>
                     <input id="name" name="name" type="text" autocomplete="name" required />
                   </div>`
                : ''
            }
            <div class="field">
              <label for="email">Email</label>
              <input id="email" name="email" type="email" autocomplete="email" required />
            </div>
            <div class="field">
              <label for="password">Пароль</label>
              <input id="password" name="password" type="password" autocomplete="${mode === 'login' ? 'current-password' : 'new-password'}" required minlength="6" />
            </div>
            <button type="submit" class="btn btn--primary btn--block">
              ${mode === 'login' ? 'Войти' : 'Создать аккаунт'}
            </button>
          </form>
        </div>
      </div>
    `;

    container.querySelectorAll<HTMLButtonElement>('.auth-tabs__tab').forEach((tab) => {
      tab.addEventListener('click', () => {
        mode = tab.dataset.mode as Mode;
        draw();
      });
    });

    const form = container.querySelector<HTMLFormElement>('#auth-form')!;
    form.addEventListener('submit', (event) => handleSubmit(event, mode, form, () => {
      mode = 'login';
      draw();
    }));
  };

  draw();
}

async function handleSubmit(event: SubmitEvent, mode: Mode, form: HTMLFormElement, switchToLogin: () => void): Promise<void> {
  event.preventDefault();
  const data = new FormData(form);
  const email = String(data.get('email') ?? '').trim();
  const password = String(data.get('password') ?? '');

  const submitBtn = form.querySelector<HTMLButtonElement>('button[type="submit"]')!;
  submitBtn.disabled = true;

  try {
    if (mode === 'login') {
      await login(email, password);
      navigate('/offices');
    } else {
      const name = String(data.get('name') ?? '').trim();
      await register(name, email, password);
      showToast('Регистрация прошла успешно — теперь войдите в аккаунт.', 'success');
      switchToLogin();
    }
  } catch (err) {
    const message = err instanceof ApiError ? err.message : 'Что-то пошло не так, попробуйте ещё раз.';
    showToast(message);
  } finally {
    submitBtn.disabled = false;
  }
}
