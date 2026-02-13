import { readBody } from 'h3'
import { callBackend } from '~~/server/utils/backend'

export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return callBackend(event, 'POST', '/register', { body })
})
