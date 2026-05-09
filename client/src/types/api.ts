export type ApiResponse<T> = {
  code: number
  data: T
}

export type PaginatedResponse<T> = {
  items: T[]
  pagination: {
    page: number
    limit: number
    total: number
  }
}
