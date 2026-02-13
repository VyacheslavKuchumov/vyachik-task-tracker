import { createError, readBody } from 'h3'
import { callBackend } from '~~/server/utils/backend'

export default defineEventHandler(async (event) => {
  const taskId = event.context.params?.taskId
  if (!taskId) {
    throw createError({ statusCode: 400, statusMessage: 'Missing task id' })
  }

  const body = await readBody(event)
  return callBackend(event, 'PUT', `/tasks/${taskId}/assign`, {
    body,
    requireAuth: true
  })
})
