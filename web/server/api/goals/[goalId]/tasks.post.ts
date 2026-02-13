import { createError, readBody } from 'h3'
import { callBackend } from '~~/server/utils/backend'

export default defineEventHandler(async (event) => {
  const goalId = event.context.params?.goalId
  if (!goalId) {
    throw createError({ statusCode: 400, statusMessage: 'Missing goal id' })
  }

  const body = await readBody(event)
  return callBackend(event, 'POST', `/goals/${goalId}/tasks`, {
    body,
    requireAuth: true
  })
})
