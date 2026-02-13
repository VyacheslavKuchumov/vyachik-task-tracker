<template>
  <UCard class="rounded-none border-x-0 border-t-0">
    <div class="mx-auto w-full max-w-7xl space-y-3">
      <div class="flex items-center justify-between gap-3">
        <UButton
          color="neutral"
          variant="ghost"
          icon="i-lucide-check-square"
          size="lg"
          :label="'Трекер задач Vyachik'"
          @click="goHome"
        />

        <div class="flex items-center gap-2">
          <UBadge v-if="auth.isAuthenticated" color="neutral" variant="subtle" size="lg" class="hidden sm:inline-flex">
            {{ auth.displayName }}
          </UBadge>

          <div class="hidden items-center gap-2 sm:flex">
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
              <UButton color="error" variant="soft" icon="i-lucide-log-out" @click="handleLogout">
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

          <UButton
            class="sm:hidden"
            color="neutral"
            variant="soft"
            :icon="mobileMenuOpen ? 'i-lucide-x' : 'i-lucide-menu'"
            aria-label="Открыть меню"
            @click="mobileMenuOpen = !mobileMenuOpen"
          />
        </div>
      </div>

      <div v-if="mobileMenuOpen" class="space-y-2 border-t border-default pt-3 sm:hidden">
        <UBadge v-if="auth.isAuthenticated" color="neutral" variant="subtle" class="w-fit">
          {{ auth.displayName }}
        </UBadge>

        <template v-if="auth.isAuthenticated">
          <UButton block color="neutral" variant="soft" icon="i-lucide-house" @click="navigateToPath('/')">
            Главная
          </UButton>
          <UButton block color="neutral" variant="soft" icon="i-lucide-folder-kanban" @click="navigateToPath('/goals')">
            Цели
          </UButton>
          <UButton block color="neutral" variant="soft" icon="i-lucide-users" @click="navigateToPath('/users')">
            Пользователи
          </UButton>
          <UButton block color="neutral" variant="soft" icon="i-lucide-user-round" @click="navigateToPath('/profile')">
            Профиль
          </UButton>
          <UButton block color="error" variant="soft" icon="i-lucide-log-out" @click="handleLogout">
            Выйти
          </UButton>
        </template>

        <template v-else>
          <UButton block color="neutral" variant="soft" icon="i-lucide-log-in" @click="navigateToPath('/login')">
            Войти
          </UButton>
          <UButton block color="primary" icon="i-lucide-user-plus" @click="navigateToPath('/signup')">
            Регистрация
          </UButton>
        </template>
      </div>
    </div>
  </UCard>
</template>

<script setup>
const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const mobileMenuOpen = ref(false)

function goHome() {
  mobileMenuOpen.value = false

  if (auth.isAuthenticated) {
    router.push('/')
    return
  }

  router.push('/login')
}

function navigateToPath(path) {
  mobileMenuOpen.value = false
  router.push(path)
}

function handleLogout() {
  mobileMenuOpen.value = false
  auth.logout()
}

watch(
  () => route.fullPath,
  () => {
    mobileMenuOpen.value = false
  }
)

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
