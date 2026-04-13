export interface ApiError {
  code: string;
  message: string;
  details: string | null;
  request_id: string;
}
