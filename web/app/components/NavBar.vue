<template>
  <UCard class="rounded-none border-x-0 border-t-0">
    <div class="mx-auto flex w-full max-w-7xl flex-wrap items-center justify-between gap-3">
      <UButton
        color="neutral"
        variant="ghost"
        icon="i-lucide-check-square"
        size="lg"
        :label="'Трекер задач Vyachik'"
        @click="goHome"
      />

      <div class="flex items-center gap-2">
        <UBadge v-if="auth.isAuthenticated" color="neutral" variant="subtle" size="lg">
          {{ auth.displayName }}
        </UBadge>

        <template v-if="auth.isAuthenticated">
          <UButton to="/" color="neutral" variant="soft" icon="i-lucide-house">
            Главная
          </UButton>
          <UButton to="/goals" color="neutral" variant="soft" icon="i-lucide-folder-kanban">
            Цели
          </UButton>
          <UButton to="/users" color="neutral" variant="soft" icon="i-lucide-users">
            Пользователи
          </UButton>
          <UButton to="/profile" color="neutral" variant="soft" icon="i-lucide-user-round">
            Профиль
          </UButton>
          <UButton color="error" variant="soft" icon="i-lucide-log-out" @click="auth.logout()">
            Выйти
          </UButton>
        </template>

        <template v-else>
          <UButton to="/login" color="neutral" variant="soft" icon="i-lucide-log-in">
            Войти
          </UButton>
          <UButton to="/signup" color="primary" icon="i-lucide-user-plus">
            Регистрация
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

onMounted(async () => {
  if (auth.isAuthenticated && !auth.profile) {
    try {
      await auth.fetchProfile()
    } catch {
      auth.logout(false)
    }
  }
})
</script>
