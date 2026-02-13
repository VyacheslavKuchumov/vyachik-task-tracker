import { createError } from 'h3'
import { callBackend } from '~~/server/utils/backend'

export default defineEventHandler(async (event) => {
  const goalId = event.context.params?.goalId
  if (!goalId) {
    throw createError({ statusCode: 400, statusMessage: 'Missing goal id' })
  }

  return callBackend(event, 'GET', `/goals/${goalId}/tasks`, {
    requireAuth: true
  })
})
