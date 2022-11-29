export interface BaseAPIResponse<T = any> {
  error: boolean;
  data?: T;
  message?: string;
}
