<template>
  <section class="space-y-6">
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
      <UCard>
        <div class="space-y-1">
          <p class="text-sm text-muted">Пользователи</p>
          <p class="text-2xl font-semibold">{{ usersWithTasks.length }}</p>
        </div>
      </UCard>

      <UCard>
        <div class="space-y-1">
          <p class="text-sm text-muted">Текущие задачи</p>
          <p class="text-2xl font-semibold">{{ currentTasksCount }}</p>
        </div>
      </UCard>
    </div>

    <UCard>
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 class="text-xl font-semibold">Пользователи и текущие задачи</h1>
            <p class="text-sm text-muted">Просмотр всех пользователей и назначенных им текущих задач.</p>
          </div>

          <UButton
            icon="i-lucide-refresh-cw"
            color="neutral"
            variant="soft"
            :loading="tracker.loadingUsersTaskBoard"
            @click="loadUsersTasks"
          >
            Обновить
          </UButton>
        </div>
      </template>

      <UProgress v-if="tracker.loadingUsersTaskBoard" />

      <UAlert
        v-else-if="usersWithTasks.length === 0"
        icon="i-lucide-users"
        color="neutral"
        variant="soft"
        title="Пользователи не найдены"
        description="Пока нет зарегистрированных пользователей."
      />

      <div v-else class="space-y-4">
        <UCard v-for="user in usersWithTasks" :key="user.id" variant="soft">
          <template #header>
            <div class="flex flex-wrap items-center justify-between gap-3">
              <div>
                <p class="font-semibold">{{ user.name || user.email }}</p>
                <p class="text-sm text-muted">{{ user.email }}</p>
              </div>

              <UBadge color="primary" variant="subtle">
                {{ (user.tasks || []).length }} активных
              </UBadge>
            </div>
          </template>

          <UAlert
            v-if="(user.tasks || []).length === 0"
            icon="i-lucide-list-checks"
            color="neutral"
            variant="soft"
            title="Нет текущих задач"
            description="У этого пользователя нет задач со статусами «к выполнению» или «в работе»."
          />

          <div v-else class="space-y-3">
            <UCard v-for="task in user.tasks" :key="task.id" class="border border-default bg-default/40">
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div class="space-y-1">
                  <p class="font-medium">{{ task.title }}</p>
                  <p class="text-sm text-muted">{{ task.description || 'Без описания' }}</p>
                </div>

                <UBadge :color="statusColor(task.status)" variant="soft">
                  {{ statusLabel(task.status) }}
                </UBadge>
              </div>

              <div class="mt-3 flex flex-wrap gap-3 text-xs text-muted">
                <span>Цель: {{ task.goalTitle || 'Неизвестная цель' }}</span>
                <span>Создал: {{ task.createdByName || 'Неизвестный пользователь' }}</span>
              </div>
            </UCard>
          </div>
        </UCard>
      </div>
    </UCard>
  </section>
</template>

<script setup>
const auth = useAuthStore()
const tracker = useTrackerStore()
const toast = useToast()

const usersWithTasks = computed(() => tracker.usersTaskBoard || [])
const currentTasksCount = computed(() => {
  return usersWithTasks.value.reduce((sum, user) => sum + (user.tasks?.length || 0), 0)
})

function statusColor(status) {
  if (status === 'done') return 'success'
  if (status === 'in_progress') return 'warning'
  return 'neutral'
}

function statusLabel(status) {
  if (status === 'done') return 'Готово'
  if (status === 'in_progress') return 'В работе'
  return 'К выполнению'
}

async function withErrorToast(action) {
  try {
    await action()
  } catch (error) {
    if (error?.statusCode === 401 || error?.statusCode === 403) {
      auth.logout()
      return
    }

    toast.add({
      title: 'Ошибка запроса',
      description: error?.data?.statusMessage || error?.statusMessage || error?.message || 'Непредвиденная ошибка.',
      color: 'error'
    })
  }
}

async function loadUsersTasks() {
  await withErrorToast(async () => {
    await tracker.fetchUsersTaskBoard(auth.authHeader())
  })
}

onMounted(async () => {
  await loadUsersTasks()
})
</script>
