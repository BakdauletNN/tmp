let hideTimer: number | undefined;

function getToastEl(): HTMLElement {
  let el = document.getElementById('toast');
  if (!el) {
    el = document.createElement('div');
    el.id = 'toast';
    el.setAttribute('role', 'status');
    document.body.appendChild(el);
  }
  return el;
}

export function showToast(message: string, kind: 'error' | 'success' = 'error'): void {
  const el = getToastEl();
  el.textContent = message;
  el.className = `toast toast--${kind} toast--visible`;

  window.clearTimeout(hideTimer);
  hideTimer = window.setTimeout(() => {
    el.classList.remove('toast--visible');
  }, 4500);
}
