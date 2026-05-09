export type ApiResponse<T> = {
  code: number
  data: T
}

export type PaginatedResponse<T> = {
  items: T[]
  page: number
  limit: number
  total: number
  total_pages: number
}
