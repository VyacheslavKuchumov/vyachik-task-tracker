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
            <label class="flex items-center gap-2 text-sm text-muted">
              <input v-model="hideCompletedTasks" type="checkbox" class="h-4 w-4 rounded border-default" />
              Скрыть выполненные
            </label>
            <UBadge color="primary" variant="subtle" size="lg">{{ visibleAssignedTasks.length }}</UBadge>
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
        v-else-if="visibleAssignedTasks.length === 0"
        icon="i-lucide-user-round-check"
        color="neutral"
        variant="soft"
        title="Нет назначенных задач"
        description="Попросите владельца цели назначить вам задачу."
      />

      <div v-else class="space-y-3">
        <UCard v-for="task in visibleAssignedTasks" :key="task.id" variant="soft">
          <div class="flex items-start justify-between gap-3">
            <div class="space-y-1">
              <p class="font-medium">{{ task.title }}</p>
              <p class="text-sm text-muted">{{ task.description || 'Без описания' }}</p>
            </div>

            <div class="flex items-center gap-2">
              <UBadge :color="completionColor(task.isCompleted)" variant="soft">
                {{ completionLabel(task.isCompleted) }}
              </UBadge>
              <UBadge :color="priorityColor(task.priority)" variant="soft">
                {{ priorityLabel(task.priority) }}
              </UBadge>
            </div>
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

          <div class="flex items-center gap-2">
            <label class="flex items-center gap-2 text-sm text-muted">
              <input v-model="hideAchievedGoals" type="checkbox" class="h-4 w-4 rounded border-default" />
              Скрыть достигнутые цели
            </label>
            <UBadge color="neutral" variant="subtle" size="lg">{{ tasksInGoals }}</UBadge>
          </div>
        </div>
      </template>

      <UProgress v-if="tracker.loadingGoals" />

      <UAlert
        v-else-if="visibleGoals.length === 0"
        icon="i-lucide-folder-open"
        color="neutral"
        variant="soft"
        title="Пока нет целей"
        description="Создайте первую цель, чтобы увидеть её задачи на главной."
      />

      <div v-else class="space-y-4">
        <UCard v-for="goal in visibleGoals" :key="goal.id" variant="soft">
          <template #header>
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div class="space-y-1">
                <h3 class="font-semibold">{{ goal.title }}</h3>
                <p class="text-sm text-muted">{{ goal.description || 'Без описания' }}</p>
              </div>

              <div class="flex items-center gap-2">
                <UBadge :color="goalStatusColor(goal.status)" variant="soft">
                  {{ goalStatusLabel(goal.status) }}
                </UBadge>
                <UBadge :color="priorityColor(goal.priority)" variant="soft">
                  {{ priorityLabel(goal.priority) }}
                </UBadge>
                <UBadge color="neutral" variant="soft">{{ goalVisibleTasks(goal).length }}</UBadge>
              </div>
            </div>
          </template>

          <UAlert
            v-if="goalVisibleTasks(goal).length === 0"
            icon="i-lucide-list-todo"
            color="neutral"
            variant="soft"
            title="У цели пока нет задач"
            description="Добавьте задачи на странице этой цели."
          />

          <div v-else class="space-y-2">
            <div
              v-for="task in goalVisibleTasks(goal)"
              :key="task.id"
              class="rounded-md border border-default bg-default px-3 py-2"
            >
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div class="space-y-1">
                  <p class="font-medium">{{ task.title }}</p>
                  <p class="text-sm text-muted">{{ task.description || 'Без описания' }}</p>
                </div>

                <div class="flex items-center gap-2">
                  <UBadge :color="completionColor(task.isCompleted)" variant="soft">
                    {{ completionLabel(task.isCompleted) }}
                  </UBadge>
                  <UBadge :color="priorityColor(task.priority)" variant="soft">
                    {{ priorityLabel(task.priority) }}
                  </UBadge>
                </div>
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

const hideCompletedTasks = ref(false)
const hideAchievedGoals = ref(false)

const visibleAssignedTasks = computed(() => {
  const tasks = tracker.assignedTasks || []
  if (!hideCompletedTasks.value) return tasks
  return tasks.filter((task) => !task.isCompleted)
})

const visibleGoals = computed(() => {
  const goals = tracker.goals || []
  if (!hideAchievedGoals.value) return goals
  return goals.filter((goal) => goal.status !== 'achieved')
})

const tasksInGoals = computed(() => {
  return visibleGoals.value.reduce((sum, goal) => sum + goalVisibleTasks(goal).length, 0)
})

function goalVisibleTasks(goal) {
  const tasks = goal?.tasks || []
  if (!hideCompletedTasks.value) return tasks
  return tasks.filter((task) => !task.isCompleted)
}

function priorityColor(priority) {
  if (priority === 'high') return 'error'
  if (priority === 'medium') return 'warning'
  return 'neutral'
}

function priorityLabel(priority) {
  if (priority === 'high') return 'Высокий'
  if (priority === 'medium') return 'Средний'
  return 'Низкий'
}

function completionColor(isCompleted) {
  return isCompleted ? 'success' : 'neutral'
}

function completionLabel(isCompleted) {
  return isCompleted ? 'Выполнена' : 'Не выполнена'
}

function goalStatusColor(status) {
  if (status === 'achieved') return 'success'
  if (status === 'in_progress') return 'warning'
  return 'neutral'
}

function goalStatusLabel(status) {
  if (status === 'achieved') return 'Достигнута'
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
