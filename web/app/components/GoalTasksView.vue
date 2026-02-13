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
            <UButton icon="i-lucide-plus" color="primary" @click="createTaskOpen = true">
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
        <div class="flex flex-wrap gap-4 text-sm text-muted">
          <span>Владелец: {{ goal.ownerName || auth.displayName }}</span>
          <span>Задачи: {{ (goal.tasks || []).length }}</span>
        </div>

        <UAlert
          v-if="(goal.tasks || []).length === 0"
          color="neutral"
          variant="soft"
          icon="i-lucide-list-todo"
          title="Пока нет задач"
          description="Создайте первую задачу для этой цели."
        />

        <UCard v-for="task in goal.tasks" :key="task.id" variant="soft">
          <template #header>
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div class="space-y-1">
                <h3 class="font-semibold">{{ task.title }}</h3>
                <p class="text-sm text-muted">{{ task.description || 'Без описания' }}</p>
              </div>

              <div class="flex items-center gap-2">
                <UBadge :color="statusColor(task.status)" variant="soft">
                  {{ statusLabel(task.status) }}
                </UBadge>
                <UButton icon="i-lucide-pencil" color="neutral" variant="soft" @click="openEditTask(task)">
                  Редактировать
                </UButton>
                <UButton
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

          <UFormField label="Статус" name="status" required>
            <select v-model="editTaskState.status" class="w-full rounded-md border border-default bg-default p-2 text-sm">
              <option value="todo">К выполнению</option>
              <option value="in_progress">В работе</option>
              <option value="done">Готово</option>
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

const usersLookup = computed(() => tracker.usersLookup || [])

const createTaskSchema = v.object({
  title: v.pipe(v.string(), v.minLength(3, 'Название задачи должно быть не короче 3 символов')),
  description: v.pipe(v.string(), v.maxLength(2000, 'Описание должно быть не длиннее 2000 символов')),
  assigneeId: v.optional(v.string())
})

const updateTaskSchema = v.object({
  title: v.pipe(v.string(), v.minLength(3, 'Название задачи должно быть не короче 3 символов')),
  description: v.pipe(v.string(), v.maxLength(2000, 'Описание должно быть не длиннее 2000 символов')),
  status: v.pipe(v.string(), v.minLength(1, 'Статус обязателен')),
  assigneeId: v.optional(v.string())
})

type CreateTaskSchema = v.InferOutput<typeof createTaskSchema>
type UpdateTaskSchema = v.InferOutput<typeof updateTaskSchema>

const createTaskState = reactive<CreateTaskSchema>({
  title: '',
  description: '',
  assigneeId: ''
})

const editTaskState = reactive<{ id: number | null; title: string; description: string; status: string; assigneeId: string }>({
  id: null,
  title: '',
  description: '',
  status: 'todo',
  assigneeId: ''
})

const goalId = computed(() => Number(route.params.goalId))

function statusColor(status: string) {
  if (status === 'done') return 'success'
  if (status === 'in_progress') return 'warning'
  return 'neutral'
}

function statusLabel(status: string) {
  if (status === 'done') return 'Готово'
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
  creatingTask.value = true

  await withErrorToast(async () => {
    const assigneeId = parseOptionalPositiveInt(event.data.assigneeId)

    await tracker.createTask(
      goalId.value,
      {
        title: event.data.title.trim(),
        description: event.data.description.trim(),
        assigneeId
      },
      auth.authHeader()
    )

    createTaskState.title = ''
    createTaskState.description = ''
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
  editTaskState.id = task.id
  editTaskState.title = task.title
  editTaskState.description = task.description
  editTaskState.status = task.status || 'todo'
  editTaskState.assigneeId = task.assigneeId ? String(task.assigneeId) : ''
  editTaskOpen.value = true
}

async function onUpdateTask(event: FormSubmitEvent<UpdateTaskSchema>) {
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
        status: event.data.status,
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
