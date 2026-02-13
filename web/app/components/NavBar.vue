<template>
  <header class="topbar">
    <div class="topbar-inner">
      <button class="brand" type="button" @click="goHome">
        <span class="brand-mark">TT</span>
        <span>Vyachik Task Tracker</span>
      </button>

      <div class="topbar-actions">
        <UBadge v-if="auth.isAuthenticated && auth.userId" color="neutral" variant="subtle" size="lg">
          user #{{ auth.userId }}
        </UBadge>

        <UButton
          v-if="auth.isAuthenticated"
          color="error"
          variant="soft"
          icon="i-lucide-log-out"
          @click="auth.logout()"
        >
          Logout
        </UButton>
      </div>
    </div>
  </header>
</template>

<script setup>
const router = useRouter()
const auth = useAuthStore()

function goHome() {
  if (auth.isAuthenticated) {
    router.push('/')
    return
  }

  router.push('/login')
}
</script>
