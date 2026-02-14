<template>
  <section class="space-y-6">
    <UCard>
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 class="text-xl font-semibold">{{ goal?.title || 'Задачи цели' }}</h1>
            <p class="text-sm text-muted">{{ goal?.description || 'Задачи выбранной цели.' }}</p>
          </div>

          <div class="flex items-center gap-2">
            <UButton to="/goals" color="neutral" variant="soft" icon="i-lucide-arrow-left">
              К целям
            </UButton>
            <UButton
              icon="i-lucide-refresh-cw"
              color="neutral"
              variant="soft"
              :loading="loadingGoal"
              @click="loadGoalTasks"
            >
              Обновить
            </UButton>
            <UButton v-if="canManageGoal" icon="i-lucide-plus" color="primary" @click="createTaskOpen = true">
              Новая задача
            </UButton>
          </div>
        </div>
      </template>

      <UProgress v-if="loadingGoal" />

      <UAlert
        v-else-if="!goal"
        color="warning"
        variant="soft"
        icon="i-lucide-triangle-alert"
        title="Цель не найдена"
        description="Возможно, у вас нет доступа к этой цели."
      />

      <div v-else class="space-y-3">
        <div class="flex flex-wrap items-center justify-between gap-3 text-sm text-muted">
          <div class="flex flex-wrap gap-4">
            <span>Владелец: {{ goal.ownerName || auth.displayName }}</span>
            <span>Задачи: {{ visibleTasks.length }}</span>
          </div>

          <div class="flex items-center gap-2">
            <UBadge :color="goalStatusColor(goal.status)" variant="soft">{{ goalStatusLabel(goal.status) }}</UBadge>
            <UBadge :color="priorityColor(goal.priority)" variant="soft">{{ priorityLabel(goal.priority) }}</UBadge>
            <label class="flex items-center gap-2 text-sm text-muted">
              <input v-model="hideCompletedTasks" type="checkbox" class="h-4 w-4 rounded border-default" />
              Скрыть выполненные
            </label>
          </div>
        </div>

        <UAlert
          v-if="visibleTasks.length === 0"
          color="neutral"
          variant="soft"
          icon="i-lucide-list-todo"
          title="Пока нет задач"
          description="Создайте первую задачу для этой цели."
        />

        <UCard v-for="task in visibleTasks" :key="task.id" variant="soft">
          <template #header>
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div class="space-y-1">
                <h3 class="font-semibold">{{ task.title }}</h3>
                <p class="text-sm text-muted">{{ task.description || 'Без описания' }}</p>
              </div>

              <div class="flex items-center gap-2">
                <UBadge :color="completionColor(task.isCompleted)" variant="soft">
                  {{ completionLabel(task.isCompleted) }}
                </UBadge>
                <UBadge :color="priorityColor(task.priority)" variant="soft">
                  {{ priorityLabel(task.priority) }}
                </UBadge>
                <UButton
                  v-if="canManageGoal"
                  icon="i-lucide-pencil"
                  color="neutral"
                  variant="soft"
                  @click="openEditTask(task)"
                >
                  Редактировать
                </UButton>
                <UButton
                  v-if="canManageGoal"
                  icon="i-lucide-trash-2"
                  color="error"
                  variant="soft"
                  :loading="deletingTaskId === task.id"
                  @click="onDeleteTask(task.id)"
                >
                  Удалить
                </UButton>
              </div>
            </div>
          </template>

          <div class="flex flex-wrap gap-4 text-xs text-muted">
            <span>Исполнитель: {{ task.assigneeName || 'Не назначен' }}</span>
            <span>Создал: {{ task.createdByName || 'Неизвестный пользователь' }}</span>
          </div>
        </UCard>
      </div>
    </UCard>

    <UModal v-model:open="createTaskOpen" title="Создать задачу">
      <template #body>
        <UForm :schema="createTaskSchema" :state="createTaskState" class="space-y-4" @submit="onCreateTask">
          <UFormField label="Название задачи" name="title" required>
            <UInput v-model="createTaskState.title" class="w-full" />
          </UFormField>

          <UFormField label="Описание" name="description">
            <UTextarea v-model="createTaskState.description" :rows="3" class="w-full" placeholder="Необязательно" />
          </UFormField>

          <UFormField label="Приоритет" name="priority" required>
            <select v-model="createTaskState.priority" class="w-full rounded-md border border-default bg-default p-2 text-sm">
              <option value="high">Высокий</option>
              <option value="medium">Средний</option>
              <option value="low">Низкий</option>
            </select>
          </UFormField>

          <UFormField label="Исполнитель" name="assigneeId">
            <select v-model="createTaskState.assigneeId" class="w-full rounded-md border border-default bg-default p-2 text-sm">
              <option value="">Не назначен</option>
              <option v-for="user in usersLookup" :key="user.id" :value="String(user.id)">
                {{ user.name }}
              </option>
            </select>
          </UFormField>

          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="createTaskOpen = false">Отмена</UButton>
            <UButton type="submit" color="primary" :loading="creatingTask">Создать</UButton>
          </div>
        </UForm>
      </template>
    </UModal>

    <UModal v-model:open="editTaskOpen" title="Редактировать задачу">
      <template #body>
        <UForm :schema="updateTaskSchema" :state="editTaskState" class="space-y-4" @submit="onUpdateTask">
          <UFormField label="Название задачи" name="title" required>
            <UInput v-model="editTaskState.title" class="w-full" />
          </UFormField>

          <UFormField label="Описание" name="description">
            <UTextarea v-model="editTaskState.description" :rows="3" class="w-full" placeholder="Необязательно" />
          </UFormField>

          <UFormField label="Приоритет" name="priority" required>
            <select v-model="editTaskState.priority" class="w-full rounded-md border border-default bg-default p-2 text-sm">
              <option value="high">Высокий</option>
              <option value="medium">Средний</option>
              <option value="low">Низкий</option>
            </select>
          </UFormField>

          <UFormField label="Исполнитель" name="assigneeId">
            <select v-model="editTaskState.assigneeId" class="w-full rounded-md border border-default bg-default p-2 text-sm">
              <option value="">Не назначен</option>
              <option v-for="user in usersLookup" :key="user.id" :value="String(user.id)">
                {{ user.name }}
              </option>
            </select>
          </UFormField>

          <label class="flex items-center gap-2 text-sm text-muted">
            <input v-model="editTaskState.isCompleted" type="checkbox" class="h-4 w-4 rounded border-default" />
            Задача выполнена
          </label>

          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="editTaskOpen = false">Отмена</UButton>
            <UButton type="submit" color="primary" :loading="updatingTask">Сохранить</UButton>
          </div>
        </UForm>
      </template>
    </UModal>
  </section>
</template>

<script setup lang="ts">
import * as v from 'valibot'
import type { FormSubmitEvent } from '@nuxt/ui'

const route = useRoute()
const auth = useAuthStore()
const tracker = useTrackerStore()
const toast = useToast()

const goal = ref<any>(null)
const loadingGoal = ref(false)
const createTaskOpen = ref(false)
const editTaskOpen = ref(false)
const creatingTask = ref(false)
const updatingTask = ref(false)
const deletingTaskId = ref<number | null>(null)
const hideCompletedTasks = ref(false)

const usersLookup = computed(() => tracker.usersLookup || [])
const canManageGoal = computed(() => Number(goal.value?.ownerId) > 0 && Number(goal.value?.ownerId) === Number(auth.userId))
const visibleTasks = computed(() => {
  const tasks = goal.value?.tasks || []
  if (!hideCompletedTasks.value) return tasks
  return tasks.filter((task: any) => !task.isCompleted)
})

const createTaskSchema = v.object({
  title: v.pipe(v.string(), v.minLength(3, 'Название задачи должно быть не короче 3 символов')),
  description: v.pipe(v.string(), v.maxLength(2000, 'Описание должно быть не длиннее 2000 символов')),
  priority: v.pipe(v.string(), v.minLength(1, 'Приоритет обязателен')),
  assigneeId: v.optional(v.string())
})

const updateTaskSchema = v.object({
  title: v.pipe(v.string(), v.minLength(3, 'Название задачи должно быть не короче 3 символов')),
  description: v.pipe(v.string(), v.maxLength(2000, 'Описание должно быть не длиннее 2000 символов')),
  priority: v.pipe(v.string(), v.minLength(1, 'Приоритет обязателен')),
  assigneeId: v.optional(v.string()),
  isCompleted: v.boolean()
})

type CreateTaskSchema = v.InferOutput<typeof createTaskSchema>
type UpdateTaskSchema = v.InferOutput<typeof updateTaskSchema>

const createTaskState = reactive<CreateTaskSchema>({
  title: '',
  description: '',
  priority: 'medium',
  assigneeId: ''
})

const editTaskState = reactive<{ id: number | null; title: string; description: string; priority: string; assigneeId: string; isCompleted: boolean }>({
  id: null,
  title: '',
  description: '',
  priority: 'medium',
  assigneeId: '',
  isCompleted: false
})

const goalId = computed(() => Number(route.params.goalId))

function priorityColor(priority: string) {
  if (priority === 'high') return 'error'
  if (priority === 'medium') return 'warning'
  return 'neutral'
}

function priorityLabel(priority: string) {
  if (priority === 'high') return 'Высокий'
  if (priority === 'medium') return 'Средний'
  return 'Низкий'
}

function completionColor(isCompleted: boolean) {
  return isCompleted ? 'success' : 'neutral'
}

function completionLabel(isCompleted: boolean) {
  return isCompleted ? 'Выполнена' : 'Не выполнена'
}

function goalStatusColor(status: string) {
  if (status === 'achieved') return 'success'
  if (status === 'in_progress') return 'warning'
  return 'neutral'
}

function goalStatusLabel(status: string) {
  if (status === 'achieved') return 'Достигнута'
  if (status === 'in_progress') return 'В работе'
  return 'К выполнению'
}

function parseOptionalPositiveInt(value?: string) {
  if (value === '' || value === null || value === undefined) return null

  const parsed = Number(value)
  if (!Number.isInteger(parsed) || parsed <= 0) {
    throw new Error('Исполнитель должен быть корректным пользователем')
  }

  return parsed
}

function confirmAction(message: string) {
  if (typeof window === 'undefined') return false
  return window.confirm(message)
}

async function withErrorToast(action: () => Promise<void>) {
  try {
    await action()
  } catch (error: any) {
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

async function loadGoalTasks() {
  if (!Number.isInteger(goalId.value) || goalId.value <= 0) {
    goal.value = null
    return
  }

  loadingGoal.value = true

  await withErrorToast(async () => {
    const response = await tracker.fetchGoalTasks(goalId.value, auth.authHeader())
    goal.value = response
  })

  loadingGoal.value = false
}

async function ensureLookupsLoaded() {
  await withErrorToast(async () => {
    await tracker.fetchUsersLookup(auth.authHeader())
  })
}

async function onCreateTask(event: FormSubmitEvent<CreateTaskSchema>) {
  if (!canManageGoal.value) return

  creatingTask.value = true

  await withErrorToast(async () => {
    const assigneeId = parseOptionalPositiveInt(event.data.assigneeId)

    await tracker.createTask(
      goalId.value,
      {
        title: event.data.title.trim(),
        description: event.data.description.trim(),
        priority: event.data.priority,
        assigneeId
      },
      auth.authHeader()
    )

    createTaskState.title = ''
    createTaskState.description = ''
    createTaskState.priority = 'medium'
    createTaskState.assigneeId = ''
    createTaskOpen.value = false

    await Promise.all([
      loadGoalTasks(),
      tracker.fetchGoals(auth.authHeader()),
      tracker.fetchAssignedTasks(auth.authHeader())
    ])

    toast.add({ title: 'Задача создана', color: 'success' })
  })

  creatingTask.value = false
}

function openEditTask(task: any) {
  if (!canManageGoal.value) return

  editTaskState.id = task.id
  editTaskState.title = task.title
  editTaskState.description = task.description
  editTaskState.priority = task.priority || 'medium'
  editTaskState.assigneeId = task.assigneeId ? String(task.assigneeId) : ''
  editTaskState.isCompleted = Boolean(task.isCompleted)
  editTaskOpen.value = true
}

async function onUpdateTask(event: FormSubmitEvent<UpdateTaskSchema>) {
  if (!canManageGoal.value) return
  if (!editTaskState.id) return

  updatingTask.value = true

  await withErrorToast(async () => {
    const assigneeId = parseOptionalPositiveInt(event.data.assigneeId)

    await tracker.updateTask(
      editTaskState.id,
      {
        goalId: goalId.value,
        title: event.data.title.trim(),
        description: event.data.description.trim(),
        priority: event.data.priority,
        isCompleted: event.data.isCompleted,
        assigneeId
      },
      auth.authHeader()
    )

    editTaskOpen.value = false

    await Promise.all([
      loadGoalTasks(),
      tracker.fetchGoals(auth.authHeader()),
      tracker.fetchAssignedTasks(auth.authHeader())
    ])

    toast.add({ title: 'Задача обновлена', color: 'success' })
  })

  updatingTask.value = false
}

async function onDeleteTask(taskId: number) {
  if (!canManageGoal.value) return
  if (!confirmAction('Удалить эту задачу?')) return

  deletingTaskId.value = taskId

  await withErrorToast(async () => {
    await tracker.deleteTask(taskId, auth.authHeader())

    await Promise.all([
      loadGoalTasks(),
      tracker.fetchGoals(auth.authHeader()),
      tracker.fetchAssignedTasks(auth.authHeader())
    ])

    toast.add({ title: 'Задача удалена', color: 'success' })
  })

  deletingTaskId.value = null
}

onMounted(async () => {
  await Promise.all([
    ensureLookupsLoaded(),
    loadGoalTasks()
  ])
})
</script>
