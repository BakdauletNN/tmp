# CoworkGo — web client

React + TypeScript + Vite. A responsive catalog of coworking spaces, rooms, sign-in, and booking management.

```bash
npm install
npm run dev        # http://localhost:5173
npm run typecheck  # strict type checking
npm run build      # typecheck + production bundle → dist/
```

API address: the `VITE_API_URL` environment variable (see `.env.example`), defaulting to `http://localhost:8080`.

The session token is stored in `localStorage`. The API address is set via `VITE_API_URL`, defaulting to `http://localhost:8080`.
