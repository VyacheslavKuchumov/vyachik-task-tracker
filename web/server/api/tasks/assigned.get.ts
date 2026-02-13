import { callBackend } from '~~/server/utils/backend'

export default defineEventHandler(async (event) => {
  return callBackend(event, 'GET', '/tasks/assigned', { requireAuth: true })
})
