<template>
  <section class="space-y-6">
    <UCard>
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 class="text-xl font-semibold">Главная</h1>
            <p class="text-sm text-muted">Задачи, назначенные вам.</p>
          </div>

          <div class="flex items-center gap-2">
            <UBadge color="primary" variant="subtle" size="lg">{{ tracker.assignedTasks.length }}</UBadge>
            <UButton
              icon="i-lucide-refresh-cw"
              color="neutral"
              variant="soft"
              :loading="tracker.loadingAssigned || tracker.loadingGoals"
              @click="loadHomeData"
            >
              Обновить
            </UButton>
            <UButton to="/goals" icon="i-lucide-folder-kanban" color="primary">
              Открыть цели
            </UButton>
          </div>
        </div>
      </template>

      <UProgress v-if="tracker.loadingAssigned" />

      <UAlert
        v-else-if="tracker.assignedTasks.length === 0"
        icon="i-lucide-user-round-check"
        color="neutral"
        variant="soft"
        title="Нет назначенных задач"
        description="Попросите владельца цели назначить вам задачу."
      />

      <div v-else class="space-y-3">
        <UCard v-for="task in tracker.assignedTasks" :key="task.id" variant="soft">
          <div class="flex items-start justify-between gap-3">
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
            <span>Автор: {{ task.createdByName || 'Неизвестный пользователь' }}</span>
          </div>
        </UCard>
      </div>
    </UCard>

    <UCard>
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold">Задачи по целям</h2>
            <p class="text-sm text-muted">Список всех целей и их задач.</p>
          </div>

          <UBadge color="neutral" variant="subtle" size="lg">{{ tasksInGoals }}</UBadge>
        </div>
      </template>

      <UProgress v-if="tracker.loadingGoals" />

      <UAlert
        v-else-if="tracker.goals.length === 0"
        icon="i-lucide-folder-open"
        color="neutral"
        variant="soft"
        title="Пока нет целей"
        description="Создайте первую цель, чтобы увидеть её задачи на главной."
      />

      <div v-else class="space-y-4">
        <UCard v-for="goal in tracker.goals" :key="goal.id" variant="soft">
          <template #header>
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div class="space-y-1">
                <h3 class="font-semibold">{{ goal.title }}</h3>
                <p class="text-sm text-muted">{{ goal.description || 'Без описания' }}</p>
              </div>

              <UBadge color="neutral" variant="soft">{{ (goal.tasks || []).length }}</UBadge>
            </div>
          </template>

          <UAlert
            v-if="(goal.tasks || []).length === 0"
            icon="i-lucide-list-todo"
            color="neutral"
            variant="soft"
            title="У цели пока нет задач"
            description="Добавьте задачи на странице этой цели."
          />

          <div v-else class="space-y-2">
            <div
              v-for="task in goal.tasks"
              :key="task.id"
              class="rounded-md border border-default bg-default px-3 py-2"
            >
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div class="space-y-1">
                  <p class="font-medium">{{ task.title }}</p>
                  <p class="text-sm text-muted">{{ task.description || 'Без описания' }}</p>
                </div>

                <UBadge :color="statusColor(task.status)" variant="soft">
                  {{ statusLabel(task.status) }}
                </UBadge>
              </div>

              <div class="mt-2 flex flex-wrap gap-3 text-xs text-muted">
                <span>Исполнитель: {{ task.assigneeName || 'Не назначен' }}</span>
                <span>Создал: {{ task.createdByName || 'Неизвестный пользователь' }}</span>
              </div>
            </div>
          </div>

          <template #footer>
            <div class="flex justify-end">
              <UButton :to="`/tasks/${goal.id}`" icon="i-lucide-list-checks" color="neutral" variant="soft">
                Открыть задачи цели
              </UButton>
            </div>
          </template>
        </UCard>
      </div>
    </UCard>
  </section>
</template>

<script setup>
const auth = useAuthStore()
const tracker = useTrackerStore()
const toast = useToast()
const tasksInGoals = computed(() => {
  return tracker.goals.reduce((sum, goal) => sum + ((goal.tasks || []).length || 0), 0)
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

async function loadHomeData() {
  await withErrorToast(async () => {
    await Promise.all([
      tracker.fetchAssignedTasks(auth.authHeader()),
      tracker.fetchGoals(auth.authHeader())
    ])
  })
}

onMounted(async () => {
  await loadHomeData()
})
</script>
