# CoworkGo — фронтенд

HTML + CSS + TypeScript на Vite (vanilla-ts, без фреймворков). SPA с
клиентским роутингом через History API.

## Запуск

```bash
npm install
npm run dev
```

Откроется на `http://localhost:5173`. Бэкенд должен быть поднят отдельно на
`http://localhost:8080` (см. `API_BASE` в `src/api.ts`, если у вас другой
адрес — поменяйте там).

## Структура

- `src/api.ts` — обёртки над `fetch` под каждый эндпоинт бэкенда.
- `src/auth.ts` — токен в памяти модуля (не localStorage) + простая
  подписка (`onAuthChange`), чтобы навбар и вьюхи реагировали на вход/выход.
- `src/router.ts` — самописный роутер на History API: `addRoute`,
  `navigate`, перехват кликов по `<a data-link>`.
- `src/views/*` — по одному файлу на экран (auth, offices, rooms,
  bookings), каждый рендерится в `#view`.
- `src/components/*` — переиспользуемые кусочки: навбар, тост с ошибками,
  пустые состояния, форма бронирования.
- `src/styles/main.css` — вся дизайн-система на CSS-переменных.

## На будущее

Код специально разложен так, чтобы легко дописать:
- **Карту мест (desks)** — когда на бэкенде появится
  `GET /rooms/:id/desks`, добавьте `getDesks()` в `api.ts` и новый view
  `desks-view.ts`; сейчас бронирование всегда идёт на уровне комнаты
  (`desk_id: null`).
- **Админку** — новые вьюхи (`admin-offices.ts`, `admin-rooms.ts`) и
  роуты `addRoute('/admin/...', ...)`, без изменений в существующих
  файлах.
