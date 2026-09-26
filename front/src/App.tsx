import { useEffect, useState, type FormEvent } from 'react';
import {
  ArrowLeft, ArrowRight, Armchair, Building2, CalendarDays, Check, ChevronRight,
  CircleHelp, Clock3, Coffee, Headset, MapPin, Menu, Monitor, Search, Send,
  ShieldCheck, SlidersHorizontal, Star, TrainFront, Users, Wind, X,
} from 'lucide-react';
import { FaTelegramPlane } from 'react-icons/fa';
import { ApiError, api } from './api';
import type { AdminUser, Booking, Desk, Office, Room } from './types';

type BookingSchedule = { date: string; time: string; duration: number };
type Route = { page: 'catalog' | 'office' | 'bookings' | 'login' | 'forgot-password' | 'reset-password' | 'admin'; officeId?: number; next?: string; resetToken?: string; schedule?: BookingSchedule };
const officePhotos = ['photo-1497366754035-f200968a6e72', 'photo-1497366216548-37526070297c', 'photo-1497366811353-6870744d04b2', 'photo-1524758631624-e2822e304c36'];
const TELEGRAM_BOT_USERNAME = (import.meta.env.VITE_TELEGRAM_BOT_USERNAME || 'CoworkGoBot').replace(/^@/, '');

function getRoute(): Route {
  const [path, query = ''] = (window.location.hash.slice(1) || '/').split('?');
  const params = new URLSearchParams(query);
  if (path === '/bookings') return { page: 'bookings' };
  if (path === '/admin') return { page: 'admin' };
  if (path === '/login') return { page: 'login', next: params.get('next') ?? undefined };
  if (path === '/forgot-password') return { page: 'forgot-password' };
  if (path === '/reset-password') return { page: 'reset-password', resetToken: params.get('token') ?? '' };
  const match = path.match(/^\/offices\/(\d+)$/);
  if (!match) return { page: 'catalog' };
  const date = params.get('date');
  const time = params.get('time');
  const duration = Number(params.get('duration'));
  return {
    page: 'office',
    officeId: Number(match[1]),
    schedule: date && time && duration > 0 ? { date, time, duration } : undefined,
  };
}
function navigate(path: string) { window.location.hash = path; }
function money(value: number) { return new Intl.NumberFormat('en-US', { maximumFractionDigits: 0 }).format(value) + ' ₸'; }
function roomName(type: string) { return ({ private: 'Private room', open: 'Open workspace', meeting: 'Meeting room' } as Record<string, string>)[type] ?? type; }
function formatRange(start: string, end: string) {
  const date = new Intl.DateTimeFormat('en-GB', { day: 'numeric', month: 'long' });
  const time = new Intl.DateTimeFormat('en-GB', { hour: '2-digit', minute: '2-digit' });
  const from = new Date(start);
  return `${date.format(from)}, ${time.format(from)}–${time.format(new Date(end))}`;
}
function localDateTime(value: Date) { return new Date(value.getTime() - value.getTimezoneOffset() * 60_000).toISOString().slice(0, 16); }
function localDate(value: Date) { return localDateTime(value).slice(0, 10); }
function defaultStartTime() {
  const value = new Date();
  value.setHours(value.getHours() + 1, 0, 0, 0);
  return localDateTime(value).slice(11, 16);
}
function scheduleInterval(schedule?: BookingSchedule) {
  const fallback = new Date();
  fallback.setHours(fallback.getHours() + 1, 0, 0, 0);
  const start = schedule ? new Date(`${schedule.date}T${schedule.time}:00`) : fallback;
  const safeStart = Number.isNaN(start.getTime()) ? fallback : start;
  const end = new Date(safeStart.getTime() + (schedule?.duration || 2) * 3_600_000);
  return { start: safeStart, end };
}
function tokenRole(token: string): string {
  try {
    const payload = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
    return JSON.parse(atob(payload)).role ?? '';
  } catch {
    return '';
  }
}

function App() {
  const [route, setRoute] = useState<Route>(getRoute);
  const [token, setToken] = useState(() => localStorage.getItem('coworkgo-token') ?? '');
  const [toast, setToast] = useState('');
  const isAdmin = tokenRole(token) === 'admin';
  useEffect(() => {
    const update = () => setRoute(getRoute());
    window.addEventListener('hashchange', update);
    return () => window.removeEventListener('hashchange', update);
  }, []);
  useEffect(() => {
    if (!toast) return;
    const timer = window.setTimeout(() => setToast(''), 3200);
    return () => window.clearTimeout(timer);
  }, [toast]);
  function logout() {
    localStorage.removeItem('coworkgo-token');
    setToken('');
    setToast('You have signed out');
    navigate('/');
  }
  const needsLogin = route.page === 'bookings' && !token;
  return <div className="app-shell">
    <header className="topbar">
      <a className="brand" href="#/" aria-label="CoworkGo, go to homepage"><span className="brand-mark"><Armchair size={19} strokeWidth={2.2} /></span><span>Cowork<span className="brand-accent">Go</span></span></a>
      <nav className="desktop-nav" aria-label="Main navigation"><a className={route.page === 'catalog' || route.page === 'office' ? 'active' : ''} href="#/">Spaces</a>{token && <a className={route.page === 'bookings' ? 'active' : ''} href="#/bookings">My bookings</a>}{isAdmin && <a className={route.page === 'admin' ? 'active' : ''} href="#/admin">Admin</a>}</nav>
      <div className="top-actions">{token ? <button className="avatar-button" onClick={logout} title="Sign out" aria-label="Sign out">A</button> : <a className="login-link" href="#/login">Sign in <ArrowRight size={15} /></a>}
        <button className="mobile-menu" onClick={() => navigate(isAdmin ? (route.page === 'admin' ? '/' : '/admin') : token ? '/bookings' : '/login')} aria-label={isAdmin ? 'Admin' : token ? 'My bookings' : 'Sign in'}><Menu size={20} /></button></div>
    </header>
    <main className="main-content">
      {route.page === 'catalog' && <Catalog onOpen={(id, schedule) => navigate(`/offices/${id}?date=${schedule.date}&time=${schedule.time}&duration=${schedule.duration}`)} />}
      {route.page === 'office' && route.officeId !== undefined && <OfficePage key={route.officeId} officeId={route.officeId} token={token} onBack={() => navigate('/')}
        schedule={route.schedule} onLogin={() => navigate(`/login?next=${encodeURIComponent(`/offices/${route.officeId}`)}`)} onToast={setToast} />}
      {route.page === 'bookings' && token && <BookingsPage onBack={() => navigate('/')} onToast={setToast} />}
      {route.page === 'admin' && token && isAdmin && <AdminPage />}
      {route.page === 'admin' && token && !isAdmin && <div className="empty-state"><strong>No access</strong><p>This panel is available to admins only.</p><button className="text-button" onClick={() => navigate('/')}>Back to spaces <ArrowRight size={15} /></button></div>}
      {route.page === 'forgot-password' && <ForgotPasswordPage />}
      {route.page === 'reset-password' && <ResetPasswordPage token={route.resetToken ?? ''} />}
      {(route.page === 'login' || needsLogin) && <AuthPage onSuccess={(newToken) => {
        localStorage.setItem('coworkgo-token', newToken); setToken(newToken); navigate(route.next || '/bookings');
      }} />}
    </main>
    <footer className="footer"><a className="brand footer-brand" href="#/">Cowork<span className="brand-accent">Go</span></a><span>A workspace that fits your day.</span><span>Almaty · 2026</span></footer>
    <div className={`toast ${toast ? 'is-visible' : ''}`} role="status" aria-live="polite">{toast}</div>
    <SupportWidget />
  </div>;
}

function Catalog({ onOpen }: { onOpen: (id: number, schedule: BookingSchedule) => void }) {
  const [offices, setOffices] = useState<Office[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [address, setAddress] = useState('');
  const [metro, setMetro] = useState(false);
  const [kitchen, setKitchen] = useState(false);
  const [minStar, setMinStar] = useState('');
  const [retry, setRetry] = useState(0);
  const [bookingDate, setBookingDate] = useState(() => localDate(new Date()));
  const [bookingTime, setBookingTime] = useState(defaultStartTime);
  const [duration, setDuration] = useState(2);
  useEffect(() => {
    let active = true;
    const timer = window.setTimeout(() => {
      setLoading(true); setError('');
      api.offices({ address: address.trim(), metro_near: metro, has_kitchen: kitchen, min_star: minStar || undefined })
        .then((items) => { if (active) setOffices(items); })
        .catch((reason: unknown) => { if (active) setError(reason instanceof Error ? reason.message : 'Could not load spaces'); })
        .finally(() => { if (active) setLoading(false); });
    }, 180);
    return () => { active = false; window.clearTimeout(timer); };
  }, [address, metro, kitchen, minStar, retry]);
  return <>
    <section className="catalog-intro"><div><p className="eyebrow"><span className="live-dot" /> ALMATY COWORKING SPACES</p><h1>Find your<br /><em>place to work.</em></h1><p className="intro-copy">A quiet corner to focus, a spacious studio for your team, or a meeting room for a couple of hours.</p></div>
      <div className="intro-photo" role="img" aria-label="Bright workspace with a large table"><div className="photo-caption"><span>01 — 08</span><span>Spaces around the city</span></div></div></section>
    <form className="booking-search" onSubmit={(event) => { event.preventDefault(); document.getElementById('spaces-results')?.scrollIntoView({ behavior: 'smooth' }); }}>
      <label><span>Date</span><input type="date" min={localDate(new Date())} required value={bookingDate} onChange={(event) => setBookingDate(event.target.value)} /></label>
      <label><span>Start</span><input type="time" required value={bookingTime} onChange={(event) => setBookingTime(event.target.value)} /></label>
      <label><span>Duration</span><select value={duration} onChange={(event) => setDuration(Number(event.target.value))}><option value={2}>2 hours</option><option value={4}>4 hours</option><option value={8}>Full day · 8 hours</option></select></label>
      <button className="button-primary" type="submit">Search <Search size={16} /></button>
    </form>
    <section className="search-panel" aria-label="Search filters">
      <label className="search-field"><Search size={18} /><span className="sr-only">Address or area</span><input value={address} onChange={(event) => setAddress(event.target.value)} placeholder="Address or area" /></label>
      <label className="select-field"><SlidersHorizontal size={17} /><span className="sr-only">Minimum rating</span><select value={minStar} onChange={(event) => setMinStar(event.target.value)}><option value="">Any rating</option><option value="3">3 stars and up</option><option value="4">4 stars and up</option><option value="5">5 stars</option></select></label>
      <label className={`filter-chip ${kitchen ? 'selected' : ''}`}><input type="checkbox" checked={kitchen} onChange={(event) => setKitchen(event.target.checked)} /><Coffee size={16} /> Kitchen</label>
      <label className={`filter-chip ${metro ? 'selected' : ''}`}><input type="checkbox" checked={metro} onChange={(event) => setMetro(event.target.checked)} /><TrainFront size={16} /> Metro nearby</label><span className="filter-count"><MapPin size={15} /> Almaty</span>
    </section>
    <section className="results-section" id="spaces-results"><div className="section-heading"><div><p className="eyebrow">TODAY'S PICKS</p><h2>Spaces <span>{loading ? '...' : offices.length.toString().padStart(2, '0')}</span></h2></div><span className="sort-note"><Star size={15} fill="currentColor" /> Most popular first</span></div>
      {loading ? <div className="office-grid">{[0, 1, 2].map((item) => <div className="office-skeleton" key={item} />)}</div> : error ? <div className="empty-state"><strong>Couldn't load the list</strong><p>{error}</p><button className="text-button" onClick={() => setRetry((value) => value + 1)}>Retry <ArrowRight size={15} /></button></div> : offices.length ? <div className="office-grid">{offices.map((office, index) => <OfficeCard key={office.id} office={office} index={index} onOpen={(id) => onOpen(id, { date: bookingDate, time: bookingTime, duration })} />)}</div> : <div className="empty-state"><strong>Nothing found</strong><p>Adjust the filters to see other spaces.</p></div>}
    </section><div className="city-note"><span className="city-line" /><span>Explore the city. Work your way.</span><ArrowRight size={17} /></div>
  </>;
}

function OfficeCard({ office, index, onOpen }: { office: Office; index: number; onOpen: (id: number) => void }) {
  const photo = officePhotos[index % officePhotos.length];
  return <button className="office-card" onClick={() => onOpen(office.id)}><span className="office-photo-wrap"><img src={`https://images.unsplash.com/${photo}?auto=format&fit=crop&w=900&q=82`} alt={`${office.name}, coworking interior`} loading={index > 2 ? 'lazy' : 'eager'} />
    <span className="rating"><Star size={14} fill="currentColor" /> {office.star}.0</span><span className="photo-open" aria-hidden="true"><ArrowRight size={18} /></span></span>
    <span className="office-card-body"><span className="office-title-row"><strong>{office.name}</strong><ChevronRight size={17} /></span><span className="office-address"><MapPin size={14} /> {office.address}</span>
      <span className="office-tags">{office.metro_near && <span><TrainFront size={13} /> Metro</span>}{office.has_kitchen && <span><Coffee size={13} /> Kitchen</span>}{office.time_range_work && <span><Clock3 size={13} /> {office.time_range_work}</span>}</span></span>
  </button>;
}

function OfficePage({ officeId, token, onBack, onLogin, onToast, schedule }: { officeId: number; token: string; onBack: () => void; onLogin: () => void; onToast: (message: string) => void; schedule?: BookingSchedule }) {
  const [office, setOffice] = useState<Office | null>(null);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [roomType, setRoomType] = useState('');
  const [activeView, setActiveView] = useState<'rooms' | 'openspace'>('rooms');
  const [booking, setBooking] = useState<{ room: Room; deskId: number | null; deskLabel?: string } | null>(null);
  useEffect(() => {
    let active = true;
    setLoading(true);
    Promise.all([api.office(officeId), api.rooms({ office_id: officeId })])
      .then(([place, availableRooms]) => { if (active) { setOffice(place); setRooms(availableRooms); } })
      .catch((reason: unknown) => { if (active) setError(reason instanceof Error ? reason.message : 'Could not load the space'); })
      .finally(() => { if (active) setLoading(false); });
    return () => { active = false; };
  }, [officeId]);
  const filteredRooms = rooms.filter((room) => !roomType || room.type === roomType);
  if (loading) return <div className="detail-loading"><div className="office-skeleton" /><div className="office-skeleton" /></div>;
  if (error || !office) return <div className="empty-state"><strong>Couldn't open this space</strong><p>{error || 'Coworking space not found'}</p><button className="text-button" onClick={onBack}><ArrowLeft size={15} /> Back to spaces</button></div>;
  return <><button className="back-link" onClick={onBack}><ArrowLeft size={16} /> All spaces</button>
    <section className="detail-hero"><div className="detail-copy"><p className="eyebrow"><MapPin size={14} /> ALMATY · COWORKING</p><h1>{office.name}</h1><p className="detail-address">{office.address}</p>
      <div className="detail-meta"><span><Star size={15} fill="currentColor" /> {office.star}.0</span><span><Clock3 size={15} /> {office.time_range_work}</span>{office.metro_near && <span><TrainFront size={15} /> Metro nearby</span>}</div></div>
      <img src={`https://images.unsplash.com/${officePhotos[(office.id - 1) % officePhotos.length]}?auto=format&fit=crop&w=1500&q=85`} alt={`Interior of ${office.name}`} /></section>
    <section className="rooms-section"><div className="section-heading room-heading"><div><p className="eyebrow">CHOOSE A FORMAT</p><h2>Rooms and desks <span>{rooms.length.toString().padStart(2, '0')}</span></h2></div>
      <select className="room-select" value={roomType} onChange={(event) => setRoomType(event.target.value)} aria-label="Room type"><option value="">All formats</option>{[...new Set(rooms.map((room) => room.type))].map((type) => <option key={type} value={type}>{roomName(type)}</option>)}</select></div>
      <div className="space-tabs" role="tablist" aria-label="Space type"><button className={activeView === 'rooms' ? 'active' : ''} onClick={() => setActiveView('rooms')}><Building2 size={15} /> Rooms</button><button className={activeView === 'openspace' ? 'active' : ''} onClick={() => setActiveView('openspace')}><Armchair size={15} /> Open space map</button></div>
      {activeView === 'rooms' ? filteredRooms.length ? <div className="room-grid">{filteredRooms.map((room) => <RoomCard key={room.id} room={room} onBook={() => room.type === 'open' ? setActiveView('openspace') : token ? setBooking({ room, deskId: null }) : onLogin()} />)}</div> : <div className="empty-state"><strong>No rooms of this format</strong><p>Try choosing a different type.</p></div> : <OpenSpaceMap rooms={rooms.filter((room) => room.type === 'open')} schedule={schedule} onBook={(room, desk) => token ? setBooking({ room, deskId: desk.id, deskLabel: desk.label }) : onLogin()} />}
    </section>
    {booking && <BookingDialog room={booking.room} office={office} schedule={schedule} deskId={booking.deskId} deskLabel={booking.deskLabel} onClose={() => setBooking(null)} onToast={onToast} />}</>;
}

function RoomCard({ room, onBook }: { room: Room; onBook: () => void }) {
  return <article className="room-card"><div className="room-icon"><Monitor size={19} /></div><p className="room-type">{roomName(room.type)}</p>
    <div className="room-capacity"><span><Users size={15} /> up to {room.qty_person} people</span><span><Armchair size={15} /> {room.qty_desks} seats</span></div>
    <div className="room-features">{room.has_air_conditioner && <span><Wind size={14} /> Air conditioning</span>}{room.has_prayer_room && <span><Coffee size={14} /> Quiet room</span>}</div>
    <div className="room-card-footer"><p className="room-price">{money(room.price_hour)} <span>/ hour</span></p><button className="button-primary button-small" onClick={onBook}>{room.type === 'open' ? 'Desk map' : 'Select'} <ArrowRight size={15} /></button></div></article>;
}

function OpenSpaceMap({ rooms, schedule, onBook }: { rooms: Room[]; schedule?: BookingSchedule; onBook: (room: Room, desk: Desk) => void }) {
  const [roomId, setRoomId] = useState(rooms[0]?.id ?? 0);
  const [desks, setDesks] = useState<Desk[]>([]);
  const [selectedDesk, setSelectedDesk] = useState<Desk | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const room = rooms.find((item) => item.id === roomId);
  const interval = scheduleInterval(schedule);

  useEffect(() => {
    let active = true;
    if (!room) { setDesks([]); return; }
    setLoading(true);
    setError('');
    setSelectedDesk(null);
    api.desks(room.id, interval.start.toISOString(), interval.end.toISOString())
      .then((items) => { if (active) setDesks(items); })
      .catch((reason: unknown) => { if (active) setError(reason instanceof Error ? reason.message : 'Could not load the seating plan'); })
      .finally(() => { if (active) setLoading(false); });
    return () => { active = false; };
  }, [roomId, schedule?.date, schedule?.time, schedule?.duration]);

  if (!rooms.length) return <div className="empty-state"><strong>This space has no Open Space area</strong><p>Choose a meeting room instead.</p></div>;

  return <section className="desk-map-panel"><div className="desk-map-heading"><div><p className="eyebrow">SEATING PLAN</p><h3>Choose a free desk</h3></div>
      <label className="room-select-label">Room<select className="room-select" value={roomId} onChange={(event) => setRoomId(Number(event.target.value))}>{rooms.map((item) => <option key={item.id} value={item.id}>Room #{item.id} · {money(item.price_hour)}/hr</option>)}</select></label>
    </div>
    <div className="desk-legend"><span><i className="is-free" /> Free</span><span><i className="is-busy" /> Occupied</span><span><i className="is-selected" /> Selected</span></div>
    {loading ? <div className="office-skeleton desk-map-skeleton" /> : error ? <div className="empty-state"><strong>Plan unavailable</strong><p>{error}</p></div> : desks.length ? <div className="floor-plan"><div className="window-label">WINDOWS</div><div className="desk-grid">{desks.map((desk) => <button key={desk.id} className={`desk-seat ${!desk.available ? 'is-busy' : selectedDesk?.id === desk.id ? 'is-selected' : ''}`} disabled={!desk.available} aria-pressed={selectedDesk?.id === desk.id} onClick={() => setSelectedDesk(desk)}><span className="desk-top">{desk.label}</span><span className="desk-chair" /></button>)}</div></div> : <div className="empty-state"><strong>No desks available</strong><p>This room has no workstations set up.</p></div>}
    {selectedDesk && room && <div className="desk-checkout"><div><strong>Selected desk: {selectedDesk.label}</strong><span>{money(room.price_hour)} / hour · {schedule?.duration ?? 2} h</span></div><button className="button-primary button-small" onClick={() => onBook(room, selectedDesk)}>Book <ArrowRight size={15} /></button></div>}
  </section>;
}

function BookingDialog({ room, office, schedule, deskId, deskLabel, onClose, onToast }: { room: Room; office: Office; schedule?: BookingSchedule; deskId: number | null; deskLabel?: string; onClose: () => void; onToast: (message: string) => void }) {
  const interval = scheduleInterval(schedule);
  const [startAt, setStartAt] = useState(localDateTime(interval.start));
  const [endAt, setEndAt] = useState(localDateTime(interval.end));
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);
  const startDate = new Date(startAt);
  const endDate = new Date(endAt);
  const hours = (endDate.getTime() - startDate.getTime()) / 3_600_000;
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!Number.isFinite(hours) || hours <= 0) { setError('The end time must be after the start time.'); return; }
    setSaving(true); setError('');
    try {
      await api.createBooking({ room_id: room.id, desk_id: deskId, start_time: startDate.toISOString(), end_time: endDate.toISOString() });
      onClose(); onToast('Booking created'); navigate('/bookings');
    } catch (reason) { setError(reason instanceof Error ? reason.message : 'Could not create the booking'); setSaving(false); }
  }
  return <div className="modal-backdrop" onMouseDown={(event) => { if (event.target === event.currentTarget) onClose(); }}><section className="booking-modal" role="dialog" aria-modal="true" aria-labelledby="booking-title">
    <button className="modal-close" onClick={onClose} aria-label="Close"><X size={18} /></button><p className="eyebrow">NEW BOOKING</p><h2 id="booking-title">{roomName(room.type)}{deskLabel ? ` · ${deskLabel}` : ''}</h2><p className="modal-subtitle">{office.name} · {money(room.price_hour)} per hour</p>
    <form onSubmit={submit}><div className="datetime-fields"><label>Start<input required type="datetime-local" value={startAt} onChange={(event) => setStartAt(event.target.value)} /></label><label>End<input required type="datetime-local" value={endAt} onChange={(event) => setEndAt(event.target.value)} /></label></div>
      <div className="booking-total"><span><Clock3 size={16} /> {hours > 0 ? `${hours.toFixed(hours % 1 ? 1 : 0)} h` : 'Set an interval'}</span><strong>{hours > 0 ? money(Math.round(hours * room.price_hour)) : '—'}</strong></div>
      {error && <p className="form-error" role="alert">{error}</p>}<button className="button-primary button-wide" disabled={saving}>{saving ? 'Saving…' : 'Confirm booking'} <ArrowRight size={17} /></button></form>
  </section></div>;
}

function BookingsPage({ onBack, onToast }: { onBack: () => void; onToast: (message: string) => void }) {
  const [rows, setRows] = useState<Array<{ booking: Booking; room?: Room; office?: Office }>>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  async function load() {
    setLoading(true);
    try {
      const bookings = await api.bookings();
      const enriched = await Promise.all(bookings.map(async (booking) => {
        const room = await api.room(booking.room_id).catch(() => undefined);
        const office = room ? await api.office(room.office_id).catch(() => undefined) : undefined;
        return { booking, room, office };
      }));
      setRows(enriched.sort((a, b) => Date.parse(a.booking.start_time) - Date.parse(b.booking.start_time))); setError('');
    } catch (reason) { setError(reason instanceof Error ? reason.message : 'Could not load bookings'); }
    finally { setLoading(false); }
  }
  useEffect(() => { void load(); }, []);
  async function cancel(id: number) {
    if (!window.confirm('Cancel this booking?')) return;
    try { await api.cancelBooking(id); onToast('Booking cancelled'); await load(); }
    catch (reason) { onToast(reason instanceof Error ? reason.message : 'Could not cancel the booking'); }
  }
  const upcoming = rows.filter(({ booking }) => Date.parse(booking.end_time) > Date.now());
  const past = rows.filter(({ booking }) => Date.parse(booking.end_time) <= Date.now());
  return <section className="bookings-page"><button className="back-link" onClick={onBack}><ArrowLeft size={16} /> Spaces</button>
    <div className="section-heading"><div><p className="eyebrow">YOUR WORK SCHEDULE</p><h1>My bookings <span>{rows.length.toString().padStart(2, '0')}</span></h1></div><CalendarDays className="heading-mark" size={34} /></div>
    {loading ? <div className="booking-list"><div className="office-skeleton" /><div className="office-skeleton" /></div> : error ? <div className="empty-state"><strong>Bookings failed to load</strong><p>{error}</p><button className="text-button" onClick={() => void load()}>Retry <ArrowRight size={15} /></button></div> : rows.length ? <>
      {upcoming.length > 0 && <BookingGroup title="Upcoming" rows={upcoming} onCancel={cancel} />}{past.length > 0 && <BookingGroup title="Completed" rows={past} />}
    </> : <div className="empty-state"><span className="empty-icon"><CalendarDays size={22} /></span><strong>No bookings yet</strong><p>Find a space and book it for a time that suits you.</p><button className="button-primary button-small" onClick={onBack}>Find a space <ArrowRight size={15} /></button></div>}
  </section>;
}

function BookingGroup({ title, rows, onCancel }: { title: string; rows: Array<{ booking: Booking; room?: Room; office?: Office }>; onCancel?: (id: number) => void }) {
  return <section className="booking-group"><h2>{title} <span>{rows.length}</span></h2><div className="booking-list">{rows.map(({ booking, room, office }) => <article className="booking-row" key={booking.id}>
    <div className="booking-date"><span>{new Intl.DateTimeFormat('en-GB', { day: '2-digit' }).format(new Date(booking.start_time))}</span><small>{new Intl.DateTimeFormat('en-GB', { month: 'short' }).format(new Date(booking.start_time)).replace('.', '')}</small></div>
    <div className="booking-info"><strong>{office?.name ?? `Room #${booking.room_id}`}</strong><span>{room ? roomName(room.type) : ''} · {formatRange(booking.start_time, booking.end_time)}</span>
      {onCancel && booking.public_code && TELEGRAM_BOT_USERNAME && <a className="telegram-booking-link" href={`https://t.me/${TELEGRAM_BOT_USERNAME}?start=${encodeURIComponent(booking.public_code)}`} target="_blank" rel="noreferrer"><FaTelegramPlane size={14} /> Details on Telegram</a>}
    </div><span className="booking-code">{booking.public_code || 'Booking'}</span>
    {onCancel ? <button className="cancel-button" onClick={() => onCancel(booking.id)}>Cancel</button> : <span className="completed-mark"><Check size={15} /> Completed</span>}
  </article>)}</div></section>;
}

function AdminPage() {
  const [tab, setTab] = useState<'overview' | 'bookings' | 'offices' | 'users'>('overview');
  const [offices, setOffices] = useState<Office[]>([]);
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');

  useEffect(() => {
    let active = true;
    Promise.all([api.adminOffices(), api.adminBookings(), api.adminUsers()])
      .then(([officeRows, bookingRows, userRows]) => {
        if (!active) return;
        setOffices(officeRows);
        setBookings(bookingRows);
        setUsers(userRows);
      })
      .catch((reason: unknown) => { if (active) setError(reason instanceof Error ? reason.message : 'Could not load dashboard data'); })
      .finally(() => { if (active) setLoading(false); });
    return () => { active = false; };
  }, []);

  const visibleBookings = bookings.filter((booking) => `${booking.id} ${booking.user_id} ${booking.room_id} ${booking.public_code ?? ''}`.toLowerCase().includes(search.toLowerCase()));
  const visibleOffices = offices.filter((office) => `${office.name} ${office.address}`.toLowerCase().includes(search.toLowerCase()));
  const visibleUsers = users.filter((user) => `${user.name} ${user.email} ${user.role}`.toLowerCase().includes(search.toLowerCase()));
  const tabs = [
    { id: 'overview' as const, label: 'Overview', icon: ShieldCheck },
    { id: 'bookings' as const, label: 'Bookings', icon: CalendarDays },
    { id: 'offices' as const, label: 'Spaces', icon: Building2 },
    { id: 'users' as const, label: 'Users', icon: Users },
  ];

  async function cancelBooking(bookingId: number) {
    if (!window.confirm(`Cancel booking #${bookingId}?`)) return;
    try {
      await api.adminCancelBooking(bookingId);
      setBookings((items) => items.filter((item) => item.id !== bookingId));
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : 'Could not cancel the booking');
    }
  }

  return <section className="admin-page">
    <div className="admin-heading"><div><p className="eyebrow"><ShieldCheck size={14} /> ADMIN PANEL</p><h1>Administration</h1><p>CoworkGo overview and data lists.</p></div><span className="admin-role">Administrator</span></div>
    <nav className="admin-tabs" aria-label="Panel sections">{tabs.map(({ id, label, icon: Icon }) => <button key={id} className={tab === id ? 'active' : ''} onClick={() => { setTab(id); setSearch(''); }}><Icon size={16} />{label}</button>)}</nav>
    {loading ? <div className="admin-loading"><div className="office-skeleton" /><div className="office-skeleton" /></div> : error ? <div className="empty-state"><strong>Panel unavailable</strong><p>{error}</p></div> : <>
      {tab === 'overview' && <>
        <div className="admin-stats"><Stat label="Spaces" value={offices.length} icon={<Building2 size={18} />} /><Stat label="Bookings" value={bookings.length} icon={<CalendarDays size={18} />} /><Stat label="Users" value={users.length} icon={<Users size={18} />} /></div>
        <div className="admin-overview-grid"><section className="admin-panel"><h2>Latest bookings</h2><AdminBookingTable rows={[...bookings].sort((a, b) => Date.parse(b.start_time) - Date.parse(a.start_time)).slice(0, 6)} /></section><section className="admin-panel"><h2>Spaces</h2><div className="admin-office-list">{offices.slice(0, 6).map((office) => <div key={office.id}><span>{office.name}</span><small>{office.address}</small></div>)}</div></section></div>
      </>}
      {tab !== 'overview' && <>
        <div className="admin-list-toolbar"><div><p className="eyebrow">DATA MANAGEMENT</p><h2>{tabs.find((item) => item.id === tab)?.label}</h2></div><label className="admin-search"><Search size={16} /><input value={search} onChange={(event) => setSearch(event.target.value)} placeholder="Search" /></label></div>
        {tab === 'bookings' && <section className="admin-panel"><AdminBookingTable rows={visibleBookings} onCancel={cancelBooking} /></section>}
        {tab === 'offices' && <section className="admin-panel"><div className="admin-table-wrap"><table className="admin-table"><thead><tr><th>ID</th><th>Space</th><th>Address</th><th>Rating</th><th>Amenities</th></tr></thead><tbody>{visibleOffices.map((office) => <tr key={office.id}><td>#{office.id}</td><td>{office.name}</td><td>{office.address}</td><td>{office.star} / 5</td><td>{[office.has_kitchen && 'Kitchen', office.metro_near && 'Metro'].filter(Boolean).join(' · ') || '—'}</td></tr>)}</tbody></table></div></section>}
        {tab === 'users' && <section className="admin-panel"><div className="admin-table-wrap"><table className="admin-table"><thead><tr><th>ID</th><th>Name</th><th>Email</th><th>Role</th><th>Status</th></tr></thead><tbody>{visibleUsers.map((user) => <tr key={user.id}><td>#{user.id}</td><td>{user.name || '—'}</td><td>{user.email}</td><td><span className={`admin-role-pill ${user.role === 'admin' ? 'is-admin' : ''}`}>{user.role || 'Customer'}</span></td><td>{user.is_registered ? 'Registered' : 'Active'}</td></tr>)}</tbody></table></div></section>}
      </>}
    </>}
  </section>;
}

function Stat({ label, value, icon }: { label: string; value: number; icon: React.ReactNode }) {
  return <article className="admin-stat"><span>{icon}</span><strong>{value.toString().padStart(2, '0')}</strong><small>{label}</small></article>;
}

function AdminBookingTable({ rows, onCancel }: { rows: Booking[]; onCancel?: (bookingId: number) => void }) {
  return <div className="admin-table-wrap"><table className="admin-table"><thead><tr><th>Booking</th><th>User</th><th>Room</th><th>Desk</th><th>Time</th><th>Code</th>{onCancel && <th>Action</th>}</tr></thead><tbody>{rows.map((booking) => <tr key={booking.id}><td>#{booking.id}</td><td>#{booking.user_id}</td><td>#{booking.room_id}</td><td>{booking.desk_id ? `Desk #${booking.desk_id}` : 'Room'}</td><td>{formatRange(booking.start_time, booking.end_time)}</td><td><code>{booking.public_code || '—'}</code></td>{onCancel && <td><button className="admin-cancel-button" onClick={() => onCancel(booking.id)}>Cancel</button></td>}</tr>)}</tbody></table>{rows.length === 0 && <p className="admin-no-results">No records found.</p>}</div>;
}

function SupportWidget() {
  const [open, setOpen] = useState(false);
  const [messages, setMessages] = useState<string[]>(['Hello! Choose a topic or type a command. An operator will reply on Telegram.']);
  const [draft, setDraft] = useState('');
  const answers: Record<string, string> = {
    '/booking': 'A booking is created on the space page: choose a room or desk, then a date and time. The API will reject a conflicting time slot.',
    '/access': 'To get your booking details on Telegram, tap "Details on Telegram" on an active booking or send its code to the bot.',
    '/telegram': TELEGRAM_BOT_USERNAME ? 'Open Telegram using the button below and send /support to the operator.' : 'The Telegram bot has not been set up yet. Set VITE_TELEGRAM_BOT_USERNAME in .env.',
  };
  const telegramUrl = TELEGRAM_BOT_USERNAME ? `https://t.me/${TELEGRAM_BOT_USERNAME}?start=support` : '';

  function sendCommand(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const command = draft.trim().toLowerCase();
    if (!command) return;
    setMessages((items) => [...items, `You: ${draft.trim()}`, answers[command] ?? 'For an operator reply, continue the conversation on Telegram via /support.']);
    setDraft('');
  }

  return <div className="support-widget">
    {open && <section className="support-panel" aria-label="Support"><header><span><Headset size={17} /> CoworkGo Support</span><button onClick={() => setOpen(false)} aria-label="Close"><X size={16} /></button></header>
      <div className="support-messages" aria-live="polite">{messages.map((message, index) => <p key={`${index}-${message}`} className={message.startsWith('You:') ? 'from-user' : ''}>{message}</p>)}</div>
      <div className="support-commands"><button onClick={() => setDraft('/booking')}>/booking</button><button onClick={() => setDraft('/access')}>/access</button><button onClick={() => setDraft('/telegram')}>/telegram</button></div>
      <form className="support-compose" onSubmit={sendCommand}><input value={draft} onChange={(event) => setDraft(event.target.value)} placeholder="Command or question" aria-label="Support command" /><button aria-label="Send"><Send size={16} /></button></form>
      {telegramUrl ? <a className="support-telegram" href={telegramUrl} target="_blank" rel="noreferrer"><FaTelegramPlane size={17} /> Message an operator on Telegram <ArrowRight size={14} /></a> : <p className="support-config-note">Set the bot username in `.env` to enable Telegram.</p>}
    </section>}
    <button className="support-launcher" onClick={() => setOpen((value) => !value)} aria-label={open ? 'Close support' : 'Open support'}>{open ? <X size={21} /> : <CircleHelp size={22} />}<span>Support</span></button>
  </div>;
}

function AuthPage({ onSuccess }: { onSuccess: (token: string) => void }) {
  const [registering, setRegistering] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    const name = String(data.get('name') ?? '').trim();
    const email = String(data.get('email') ?? '').trim();
    const password = String(data.get('password') ?? '');
    setBusy(true); setError('');
    try { if (registering) await api.register(name, email, password); const result = await api.login(email, password); onSuccess(result.token); }
    catch (reason) { setError(reason instanceof ApiError ? reason.message : 'Could not complete the request'); }
    finally { setBusy(false); }
  }
  return <section className="auth-layout"><div className="auth-aside"><p className="eyebrow"><span className="live-dot" /> COWORKGO · ALMATY</p><h1>Your next<br /><em>working day</em><br />starts here.</h1>
    <p>Manage your bookings and come back to the spaces that work for you.</p><div className="auth-photo"><img src="https://images.unsplash.com/photo-1497366811353-6870744d04b2?auto=format&fit=crop&w=1200&q=82" alt="Bright coworking space" /></div></div>
    <div className="auth-form-wrap"><span className="auth-step">{registering ? 'NEW ACCOUNT' : 'WELCOME BACK'}</span><h2>{registering ? 'Create an account' : 'Sign in'}</h2><p>{registering ? 'Fill in your details to start booking.' : 'Enter your details to see your bookings.'}</p>
      <form onSubmit={submit}>{registering && <label>Name<input name="name" required minLength={2} autoComplete="name" placeholder="How should we address you" /></label>}
        <label>Email<input name="email" required type="email" autoComplete="email" placeholder="name@example.com" /></label><label>Password<input name="password" required minLength={6} type="password" autoComplete={registering ? 'new-password' : 'current-password'} placeholder="At least 6 characters" /></label>
        {error && <p className="form-error" role="alert">{error}</p>}
        {!registering && <p className="auth-inline-link"><button type="button" onClick={() => navigate('/forgot-password')}>Forgot your password?</button></p>}
        <button className="button-primary button-wide" disabled={busy}>{busy ? 'Please wait…' : registering ? 'Create account' : 'Sign in'} <ArrowRight size={17} /></button>
      </form><p className="auth-toggle">{registering ? 'Already have an account?' : 'First time here?'} <button onClick={() => { setRegistering(!registering); setError(''); }}>{registering ? 'Sign in' : 'Create account'}</button></p></div>
  </section>;
}

function ForgotPasswordPage() {
  const [email, setEmail] = useState('');
  const [sent, setSent] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setBusy(true);
    setError('');
    try {
      await api.forgotPassword(email.trim());
      setSent(true);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : 'Could not send the request');
    } finally {
      setBusy(false);
    }
  }

  return <section className="auth-layout">
    <div className="auth-aside"><p className="eyebrow"><span className="live-dot" /> COWORKGO · ALMATY</p><h1>Let's get you<br /><em>back in.</em></h1>
      <p>We'll send a secure link to set a new password.</p><div className="auth-photo"><img src="https://images.unsplash.com/photo-1497366811353-6870744d04b2?auto=format&fit=crop&w=1200&q=82" alt="Bright workspace" /></div></div>
    <div className="auth-form-wrap"><button className="text-button auth-back" onClick={() => navigate('/login')}><ArrowLeft size={15} /> Back to sign in</button>
      {sent ? <div className="reset-success"><span className="empty-icon"><Check size={22} /></span><h2>Check your inbox</h2><p>If an account with {email} exists, we've sent a link to reset the password.</p><button className="button-primary button-wide" onClick={() => navigate('/login')}>Back to sign in <ArrowRight size={16} /></button></div> : <>
        <span className="auth-step">ACCOUNT RECOVERY</span><h2>Forgot your password?</h2><p>Enter your account email and we'll send you a reset link.</p>
        <form onSubmit={submit}><label>Email<input required type="email" autoComplete="email" value={email} onChange={(event) => setEmail(event.target.value)} placeholder="name@example.com" /></label>
          {error && <p className="form-error" role="alert">{error}</p>}<button className="button-primary button-wide" disabled={busy}>{busy ? 'Sending…' : 'Send link'} <ArrowRight size={17} /></button></form>
      </>}</div>
  </section>;
}

function ResetPasswordPage({ token }: { token: string }) {
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [complete, setComplete] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!token) { setError('This link is missing a reset token. Request a new email.'); return; }
    if (password !== confirmPassword) { setError('Passwords do not match.'); return; }
    setBusy(true);
    setError('');
    try {
      await api.resetPassword(token, password);
      setComplete(true);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : 'This link is invalid or has expired');
    } finally {
      setBusy(false);
    }
  }

  return <section className="auth-layout">
    <div className="auth-aside"><p className="eyebrow"><span className="live-dot" /> COWORKGO · ALMATY</p><h1>New password.<br /><em>Fresh start.</em></h1>
      <p>Create the password you'll use to sign in to CoworkGo.</p><div className="auth-photo"><img src="https://images.unsplash.com/photo-1497366754035-f200968a6e72?auto=format&fit=crop&w=1200&q=82" alt="CoworkGo workspace" /></div></div>
    <div className="auth-form-wrap"><button className="text-button auth-back" onClick={() => navigate('/login')}><ArrowLeft size={15} /> Back to sign in</button>
      {complete ? <div className="reset-success"><span className="empty-icon"><Check size={22} /></span><h2>Password changed</h2><p>You can now sign in with your new password.</p><button className="button-primary button-wide" onClick={() => navigate('/login')}>Sign in <ArrowRight size={16} /></button></div> : <>
        <span className="auth-step">SECURE LINK</span><h2>Set a new password</h2><p>This link is valid for 30 minutes and can only be used once.</p>
        <form onSubmit={submit}><label>New password<input required minLength={6} type="password" autoComplete="new-password" value={password} onChange={(event) => setPassword(event.target.value)} placeholder="At least 6 characters" /></label>
          <label>Confirm password<input required minLength={6} type="password" autoComplete="new-password" value={confirmPassword} onChange={(event) => setConfirmPassword(event.target.value)} placeholder="Type the password again" /></label>
          {error && <p className="form-error" role="alert">{error}</p>}<button className="button-primary button-wide" disabled={busy}>{busy ? 'Saving…' : 'Change password'} <ArrowRight size={17} /></button></form>
      </>}</div>
  </section>;
}

export default App;
