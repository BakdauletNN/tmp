import './styles/main.css';
import { addRoute, navigate, resolve, setFallback } from './router';
import { isAuthenticated } from './auth';
import { renderNavbar } from './components/navbar';
import { renderAuthView } from './views/auth-view';
import { renderOfficesView } from './views/offices-view';
import { renderRoomsView } from './views/rooms-view';
import { renderBookingsView } from './views/bookings-view';

const app = document.getElementById('app');
if (!app) {
  throw new Error('#app не найден в index.html');
}

app.innerHTML = `
  <header id="navbar"></header>
  <main id="view" class="view"></main>
`;

const navEl = document.getElementById('navbar')!;
const viewEl = document.getElementById('view')!;

renderNavbar(navEl);

setFallback('/offices');

addRoute('/', () => navigate('/offices'));
addRoute('/login', () => renderAuthView(viewEl));
addRoute('/offices', () => renderOfficesView(viewEl));
addRoute('/offices/:id', (params) => renderRoomsView(viewEl, Number(params.id)));
addRoute('/bookings', () => {
  if (!isAuthenticated()) {
    navigate('/login');
    return;
  }
  renderBookingsView(viewEl);
});

resolve();
