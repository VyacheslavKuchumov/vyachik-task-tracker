import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to) => {
  const auth = useAuthStore()
  auth.hydrateFromToken()

  const publicPages = ['/login', '/signup']
  const isPublicPage = publicPages.includes(to.path)

  if (!auth.isAuthenticated && !isPublicPage) {
    return navigateTo('/login')
  }

  if (auth.isAuthenticated && isPublicPage) {
    return navigateTo('/')
  }
})
