function parseTokenPayload(token) {
  try {
    const payload = token.split('.')[1]
    const normalized = payload.replace(/-/g, '+').replace(/_/g, '/')
    const decoded =
      typeof atob === 'function'
        ? atob(normalized)
        : Buffer.from(normalized, 'base64').toString('utf8')
    return JSON.parse(decoded)
  } catch {
    return null
  }
}

function getTokenExpiryMs(token) {
  const payload = parseTokenPayload(token)
  const expSeconds = Number(payload?.expiredAt)
  if (!Number.isFinite(expSeconds)) return 0
  return expSeconds * 1000
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null,
    userId: null,
    profile: null
  }),
  getters: {
    isAuthenticated: (state) => {
      if (!state.token) return false
      const expiryMs = getTokenExpiryMs(state.token)
      if (!expiryMs) return false
      return Date.now() < expiryMs
    },
    fullName: (state) => {
      const firstName = state.profile?.firstName?.trim() || ''
      const lastName = state.profile?.lastName?.trim() || ''
      return `${firstName} ${lastName}`.trim()
    },
    displayName() {
      return this.fullName || this.profile?.email || 'Пользователь'
    }
  },
  persist: true,
  actions: {
    hydrateFromToken() {
      if (!this.token) {
        this.userId = null
        this.profile = null
        return
      }

      const payload = parseTokenPayload(this.token)
      const parsedId = Number(payload?.userID)
      this.userId = Number.isFinite(parsedId) && parsedId > 0 ? parsedId : null

      if (!this.isAuthenticated) {
        this.logout(false)
      }
    },

    async login({ email, password }) {
      const response = await $fetch('/api/auth/login', {
        method: 'POST',
        body: { email, password }
      })

      this.token = response.token
      this.hydrateFromToken()
      await this.fetchProfile()
    },

    async signup({ firstName, lastName, email, password }) {
      await $fetch('/api/auth/register', {
        method: 'POST',
        body: { firstName, lastName, email, password }
      })

      await this.login({ email, password })
    },

    async fetchProfile() {
      if (!this.token) return null

      const profile = await $fetch('/api/profile', {
        headers: this.authHeader()
      })

      this.profile = profile
      const parsedId = Number(profile?.id)
      this.userId = Number.isFinite(parsedId) && parsedId > 0 ? parsedId : this.userId
      return profile
    },

    async updateProfile({ firstName, lastName }) {
      const profile = await $fetch('/api/profile', {
        method: 'PUT',
        headers: this.authHeader(),
        body: { firstName, lastName }
      })

      this.profile = profile
      return profile
    },

    async changePassword({ currentPassword, newPassword }) {
      await $fetch('/api/profile/password', {
        method: 'PUT',
        headers: this.authHeader(),
        body: { currentPassword, newPassword }
      })
    },

    authHeader() {
      if (!this.token) return {}
      return { Authorization: `Bearer ${this.token}` }
    },

    logout(redirect = true) {
      this.token = null
      this.userId = null
      this.profile = null

      if (redirect) {
        navigateTo('/login')
      }
    }
  }
})
