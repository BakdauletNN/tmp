// Формы данных ровно такие, какие отдаёт Go-бэкенд (см. handler/http).

export interface User {
  id: number;
  name: string;
  email: string;
}

export interface AuthResponse {
  token: string;
}

export interface Office {
  id: number;
  name: string;
  address: string;
  star: number;
  has_kitchen: boolean;
  time_range_work: string;
  metro_near: string;
}

export interface Room {
  id: number;
  office_id: number;
  type: string;
  price_hour: number;
  qty_desks: number;
  qty_person: number;
  access_code: string;
  has_air_conditioner: boolean;
  has_prayer_room: boolean;
}

export interface Booking {
  id: number;
  room_id: number;
  user_id: number;
  desk_id: number | null;
  start_time: string; // RFC3339
  end_time: string; // RFC3339
}

export interface CreateBookingPayload {
  room_id: number;
  desk_id: number | null;
  start_time: string;
  end_time: string;
}

export interface CreateBookingResponse {
  message: string;
  booking: Booking;
}

export interface OfficeFilters {
  address?: string;
  min_star?: number;
  has_kitchen?: boolean;
}

export interface RoomFilters {
  office_id?: number;
  type?: string;
  max_price_hour?: number;
  min_qty_person?: number;
  has_air_conditioner?: boolean;
}
