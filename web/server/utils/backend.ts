import { createError, getHeader, type H3Event } from 'h3'

type Method = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

type RequestOptions = {
  body?: unknown
  query?: Record<string, unknown>
  requireAuth?: boolean
}

export async function callBackend(
  event: H3Event,
  method: Method,
  path: string,
  options: RequestOptions = {}
) {
  const config = useRuntimeConfig(event)
  const headers: Record<string, string> = {}

  if (options.requireAuth) {
    const authHeader = getHeader(event, 'authorization')
    if (!authHeader) {
      throw createError({ statusCode: 401, statusMessage: 'Missing Authorization header' })
    }
    headers.Authorization = authHeader
  }

  try {
    return await $fetch(`${config.backendUrl}/api/v1${path}`, {
      method,
      body: options.body,
      query: options.query,
      headers
    })
  } catch (error: any) {
    const statusCode = Number(error?.statusCode || error?.status || error?.response?.status) || 500
    const payload = error?.data || error?.response?._data || {}
    const statusMessage =
      payload?.error ||
      payload?.message ||
      error?.statusMessage ||
      error?.message ||
      'Backend request failed'

    throw createError({
      statusCode,
      statusMessage
    })
  }
}
