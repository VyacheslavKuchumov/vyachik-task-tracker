import { createError } from 'h3'
import { callBackend } from '~~/server/utils/backend'

export default defineEventHandler(async (event) => {
  const taskId = event.context.params?.taskId
  if (!taskId) {
    throw createError({ statusCode: 400, statusMessage: 'Missing task id' })
  }

  return callBackend(event, 'DELETE', `/tasks/${taskId}`, {
    requireAuth: true
  })
})
