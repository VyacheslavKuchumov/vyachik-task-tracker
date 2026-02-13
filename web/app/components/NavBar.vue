<template>
  <UCard class="rounded-none border-x-0 border-t-0">
    <div class="mx-auto flex w-full max-w-7xl flex-wrap items-center justify-between gap-3">
      <UButton
        color="neutral"
        variant="ghost"
        icon="i-lucide-check-square"
        size="lg"
        :label="'Vyachik Task Tracker'"
        @click="goHome"
      />

      <div class="flex items-center gap-2">
        <UBadge v-if="auth.isAuthenticated && auth.userId" color="neutral" variant="subtle" size="lg">
          User #{{ auth.userId }}
        </UBadge>

        <template v-if="auth.isAuthenticated">
          <UButton to="/" color="neutral" variant="soft" icon="i-lucide-layout-dashboard">
            Dashboard
          </UButton>
          <UButton color="error" variant="soft" icon="i-lucide-log-out" @click="auth.logout()">
            Logout
          </UButton>
        </template>

        <template v-else>
          <UButton to="/login" color="neutral" variant="soft" icon="i-lucide-log-in">
            Login
          </UButton>
          <UButton to="/signup" color="primary" icon="i-lucide-user-plus">
            Sign up
          </UButton>
        </template>
      </div>
    </div>
  </UCard>
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
