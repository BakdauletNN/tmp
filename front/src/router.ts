type RouteHandler = (params: Record<string, string>) => void;

interface Route {
  pattern: RegExp;
  keys: string[];
  handler: RouteHandler;
}

const routes: Route[] = [];
let fallbackPath = '/';

function compile(path: string): { pattern: RegExp; keys: string[] } {
  const keys: string[] = [];
  const source = path
    .split('/')
    .map((segment) => {
      if (segment.startsWith(':')) {
        keys.push(segment.slice(1));
        return '([^/]+)';
      }
      return segment.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    })
    .join('/');
  return { pattern: new RegExp(`^${source}$`), keys };
}

export function addRoute(path: string, handler: RouteHandler): void {
  const { pattern, keys } = compile(path);
  routes.push({ pattern, keys, handler });
}

export function setFallback(path: string): void {
  fallbackPath = path;
}

export function resolve(): void {
  const path = window.location.pathname || '/';
  for (const route of routes) {
    const match = path.match(route.pattern);
    if (match) {
      const params: Record<string, string> = {};
      route.keys.forEach((key, i) => {
        params[key] = decodeURIComponent(match[i + 1]);
      });
      route.handler(params);
      return;
    }
  }
  if (path !== fallbackPath) {
    navigate(fallbackPath);
  }
}

export function navigate(path: string): void {
  if (window.location.pathname !== path) {
    window.history.pushState({}, '', path);
  }
  resolve();
}

window.addEventListener('popstate', resolve);

// Перехватываем клики по ссылкам с data-link, чтобы навигация оставалась
// клиентской (через History API), а не перезагружала страницу.
document.addEventListener('click', (event) => {
  const target = (event.target as HTMLElement).closest('a[data-link]');
  if (!(target instanceof HTMLAnchorElement)) return;
  if (event.defaultPrevented || event.button !== 0) return;
  if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;

  const href = target.getAttribute('href');
  if (!href || !href.startsWith('/')) return;

  event.preventDefault();
  navigate(href);
});
