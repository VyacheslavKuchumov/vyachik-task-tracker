import { callBackend } from '~~/server/utils/backend'

export default defineEventHandler(async (event) => {
  return callBackend(event, 'GET', '/users/lookup', { requireAuth: true })
})
