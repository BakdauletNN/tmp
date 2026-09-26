export interface Office {
  id: number;
  name: string;
  address: string;
  star: number;
  has_kitchen: boolean;
  time_range_work: string;
  metro_near: boolean;
}

export interface Room {
  id: number;
  office_id: number;
  type: string;
  price_hour: number;
  qty_desks: number;
  qty_person: number;
  access_code: number | null;
  has_air_conditioner: boolean;
  has_prayer_room: boolean;
}

export interface Booking {
  id: number;
  room_id: number;
  user_id: number;
  desk_id?: number | null;
  start_time: string;
  end_time: string;
  public_code?: string;
}

export interface Desk {
  id: number;
  room_id: number;
  label: string;
  available: boolean;
}

export interface AdminUser {
  id: number;
  name: string;
  email: string;
  role: string;
  is_registered: boolean;
}

export interface OfficeFilter { address?: string; min_star?: number; has_kitchen?: boolean; metro_near?: boolean }
export interface RoomFilter { office_id?: number; type?: string; max_price_hour?: number; has_air_conditioner?: boolean }
export interface CreateBookingBody { room_id: number; desk_id: number | null; start_time: string; end_time: string }

export interface RouteCtx {
  root: HTMLElement;
  params: Record<string, string>;
  query: URLSearchParams;
  /** true if the user has already navigated away — the request result should be discarded */
  stale: () => boolean;
}
