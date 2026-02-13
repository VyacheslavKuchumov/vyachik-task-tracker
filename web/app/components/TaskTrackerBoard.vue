<template>
  <section class="space-y-6">
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
      <UCard>
        <div class="space-y-1">
          <p class="text-sm text-muted">Goals</p>
          <p class="text-2xl font-semibold">{{ tracker.goals.length }}</p>
        </div>
      </UCard>

      <UCard>
        <div class="space-y-1">
          <p class="text-sm text-muted">Tasks In Goals</p>
          <p class="text-2xl font-semibold">{{ tasksInGoals }}</p>
        </div>
      </UCard>
    </div>

    <UCard>
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold">Goals</h2>
            <p class="text-sm text-muted">Manage goals and tasks with full CRUD actions.</p>
          </div>

          <div class="flex items-center gap-2">
            <UButton
              icon="i-lucide-refresh-cw"
              color="neutral"
              variant="soft"
              :loading="tracker.loadingGoals"
              @click="loadGoals"
            >
              Refresh
            </UButton>

            <UButton icon="i-lucide-plus" color="primary" @click="createGoalOpen = true">
              New Goal
            </UButton>
          </div>
        </div>
      </template>

      <UProgress v-if="tracker.loadingGoals" />

      <UAlert
        v-else-if="tracker.goals.length === 0"
        icon="i-lucide-folder-open"
        color="neutral"
        variant="soft"
        title="No goals yet"
        description="Create your first goal to get started."
      />

      <div v-else class="space-y-4">
        <UCard v-for="goal in tracker.goals" :key="goal.id" variant="soft">
          <template #header>
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div class="space-y-1">
                <h3 class="font-semibold">{{ goal.title }}</h3>
                <p class="text-sm text-muted">{{ goal.description }}</p>
              </div>

              <div class="flex items-center gap-2">
                <UBadge color="neutral" variant="subtle">#{{ goal.id }}</UBadge>
                <UButton icon="i-lucide-pencil" color="neutral" variant="soft" @click="openEditGoal(goal)">
                  Edit
                </UButton>
                <UButton
                  icon="i-lucide-trash-2"
                  color="error"
                  variant="soft"
                  :loading="deletingGoalId === goal.id"
                  @click="onDeleteGoal(goal.id)"
                >
                  Delete
                </UButton>
              </div>
            </div>
          </template>

          <div class="space-y-3">
            <UAlert
              v-if="(goal.tasks || []).length === 0"
              color="neutral"
              variant="outline"
              icon="i-lucide-list-todo"
              title="No tasks in this goal"
            />

            <UCard v-for="task in goal.tasks" :key="task.id" variant="outline">
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div class="space-y-1">
                  <p class="font-medium">{{ task.title }}</p>
                  <p class="text-sm text-muted">{{ task.description }}</p>
                </div>

                <div class="flex items-center gap-2">
                  <UBadge :color="statusColor(task.status)" variant="soft">
                    {{ task.status || 'todo' }}
                  </UBadge>
                  <UButton icon="i-lucide-pencil" color="neutral" variant="soft" @click="openEditTask(goal.id, task)">
                    Edit
                  </UButton>
                  <UButton
                    icon="i-lucide-trash-2"
                    color="error"
                    variant="soft"
                    :loading="deletingTaskId === task.id"
                    @click="onDeleteTask(task.id)"
                  >
                    Delete
                  </UButton>
                </div>
              </div>

              <div class="mt-3 flex flex-wrap gap-3 text-xs text-muted">
                <span>Assignee: {{ task.assigneeName || (task.assigneeId ? `#${task.assigneeId}` : 'unassigned') }}</span>
                <span>Creator: {{ task.createdByName || `#${task.createdBy}` }}</span>
              </div>
            </UCard>
          </div>

          <template #footer>
            <UForm
              :schema="createTaskSchema"
              :state="taskDraft(goal.id)"
              class="grid grid-cols-1 gap-3 md:grid-cols-2"
              @submit="(event) => onCreateTask(goal.id, event)"
            >
              <UFormField label="Task title" name="title" required>
                <UInput v-model="taskDraft(goal.id).title" placeholder="Define implementation" class="w-full" />
              </UFormField>

              <UFormField label="Assignee ID (optional)" name="assigneeId">
                <UInput v-model="taskDraft(goal.id).assigneeId" type="number" min="1" placeholder="e.g. 2" class="w-full" />
              </UFormField>

              <UFormField label="Task description" name="description" required class="md:col-span-2">
                <UTextarea v-model="taskDraft(goal.id).description" :rows="2" placeholder="API + UI + tests" class="w-full" />
              </UFormField>

              <div class="md:col-span-2">
                <UButton type="submit" :loading="creatingTaskGoalId === goal.id">
                  Add Task
                </UButton>
              </div>
            </UForm>
          </template>
        </UCard>
      </div>
    </UCard>

    <UModal v-model:open="createGoalOpen" title="Create Goal">
      <template #body>
        <UForm :schema="goalSchema" :state="createGoalState" class="space-y-4" @submit="onCreateGoal">
          <UFormField label="Goal title" name="title" required>
            <UInput v-model="createGoalState.title" class="w-full" placeholder="Ship task tracker MVP" />
          </UFormField>

          <UFormField label="Description" name="description" required>
            <UTextarea
              v-model="createGoalState.description"
              :rows="4"
              class="w-full"
              placeholder="Clear objective and expected outcome"
            />
          </UFormField>

          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="createGoalOpen = false">
              Cancel
            </UButton>
            <UButton type="submit" color="primary" :loading="creatingGoal">
              Create Goal
            </UButton>
          </div>
        </UForm>
      </template>
    </UModal>

    <UModal v-model:open="editGoalOpen" title="Edit Goal">
      <template #body>
        <UForm :schema="goalSchema" :state="editGoalState" class="space-y-4" @submit="onUpdateGoal">
          <UFormField label="Goal title" name="title" required>
            <UInput v-model="editGoalState.title" class="w-full" />
          </UFormField>

          <UFormField label="Description" name="description" required>
            <UTextarea v-model="editGoalState.description" :rows="4" class="w-full" />
          </UFormField>

          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="editGoalOpen = false">
              Cancel
            </UButton>
            <UButton type="submit" color="primary" :loading="updatingGoal">
              Save
            </UButton>
          </div>
        </UForm>
      </template>
    </UModal>

    <UModal v-model:open="editTaskOpen" title="Edit Task">
      <template #body>
        <UForm :schema="updateTaskSchema" :state="editTaskState" class="space-y-4" @submit="onUpdateTask">
          <UFormField label="Task title" name="title" required>
            <UInput v-model="editTaskState.title" class="w-full" />
          </UFormField>

          <UFormField label="Description" name="description" required>
            <UTextarea v-model="editTaskState.description" :rows="3" class="w-full" />
          </UFormField>

          <UFormField label="Status" name="status" required>
            <UInput v-model="editTaskState.status" class="w-full" placeholder="todo | in_progress | done" />
          </UFormField>

          <UFormField label="Assignee ID (optional)" name="assigneeId">
            <UInput v-model="editTaskState.assigneeId" type="number" min="1" class="w-full" />
          </UFormField>

          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="editTaskOpen = false">
              Cancel
            </UButton>
            <UButton type="submit" color="primary" :loading="updatingTask">
              Save
            </UButton>
          </div>
        </UForm>
      </template>
    </UModal>
  </section>
</template>

<script setup lang="ts">
import * as v from 'valibot'
import type { FormSubmitEvent } from '@nuxt/ui'

type GoalEntity = {
  id: number
  title: string
  description: string
  tasks?: TaskEntity[]
}

type TaskEntity = {
  id: number
  goalId: number
  title: string
  description: string
  status: string
  assigneeId?: number | null
  assigneeName?: string
  createdBy: number
  createdByName?: string
}

const auth = useAuthStore()
const tracker = useTrackerStore()
const toast = useToast()

const createGoalOpen = ref(false)
const editGoalOpen = ref(false)
const editTaskOpen = ref(false)

const creatingGoal = ref(false)
const updatingGoal = ref(false)
const deletingGoalId = ref<number | null>(null)

const creatingTaskGoalId = ref<number | null>(null)
const updatingTask = ref(false)
const deletingTaskId = ref<number | null>(null)

const goalSchema = v.object({
  title: v.pipe(v.string(), v.minLength(3, 'Goal title should be at least 3 characters')),
  description: v.pipe(v.string(), v.minLength(3, 'Description should be at least 3 characters'))
})

const createTaskSchema = v.object({
  title: v.pipe(v.string(), v.minLength(3, 'Task title should be at least 3 characters')),
  description: v.pipe(v.string(), v.minLength(3, 'Description should be at least 3 characters')),
  assigneeId: v.optional(v.string())
})

const updateTaskSchema = v.object({
  title: v.pipe(v.string(), v.minLength(3, 'Task title should be at least 3 characters')),
  description: v.pipe(v.string(), v.minLength(3, 'Description should be at least 3 characters')),
  status: v.pipe(v.string(), v.minLength(1, 'Status is required')),
  assigneeId: v.optional(v.string())
})

type GoalSchema = v.InferOutput<typeof goalSchema>
type CreateTaskSchema = v.InferOutput<typeof createTaskSchema>
type UpdateTaskSchema = v.InferOutput<typeof updateTaskSchema>

const createGoalState = reactive<GoalSchema>({
  title: '',
  description: ''
})

const editGoalState = reactive<{ id: number | null; title: string; description: string }>({
  id: null,
  title: '',
  description: ''
})

const editTaskState = reactive<{ id: number | null; goalId: number | null; title: string; description: string; status: string; assigneeId: string }>({
  id: null,
  goalId: null,
  title: '',
  description: '',
  status: 'todo',
  assigneeId: ''
})

const taskDrafts = reactive<Record<string, CreateTaskSchema>>({})

const tasksInGoals = computed(() => {
  return tracker.goals.reduce((sum: number, goal: GoalEntity) => sum + ((goal.tasks || []).length || 0), 0)
})

function taskDraft(goalId: number) {
  const key = String(goalId)

  if (!taskDrafts[key]) {
    taskDrafts[key] = {
      title: '',
      description: '',
      assigneeId: ''
    }
  }

  return taskDrafts[key]
}

function parseOptionalPositiveInt(value?: string) {
  if (value === '' || value === null || value === undefined) return null

  const parsed = Number(value)
  if (!Number.isInteger(parsed) || parsed <= 0) {
    throw new Error('Assignee id must be a positive integer')
  }

  return parsed
}

function parseTaskStatus(status: string) {
  const normalized = (status || '').trim()
  const allowed = ['todo', 'in_progress', 'done']
  if (!allowed.includes(normalized)) {
    throw new Error('Status must be one of: todo, in_progress, done')
  }
  return normalized
}

function statusColor(status: string) {
  if (status === 'done') return 'success'
  if (status === 'in_progress') return 'warning'
  return 'neutral'
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
      title: 'Request failed',
      description: error?.data?.statusMessage || error?.statusMessage || error?.message || 'Unexpected error.',
      color: 'error'
    })
  }
}

async function loadGoals() {
  await withErrorToast(async () => {
    await tracker.fetchGoals(auth.authHeader())
  })
}

async function onCreateGoal(event: FormSubmitEvent<GoalSchema>) {
  creatingGoal.value = true

  await withErrorToast(async () => {
    await tracker.createGoal(
      {
        title: event.data.title.trim(),
        description: event.data.description.trim()
      },
      auth.authHeader()
    )

    createGoalState.title = ''
    createGoalState.description = ''
    createGoalOpen.value = false

    toast.add({ title: 'Goal created', color: 'success' })
  })

  creatingGoal.value = false
}

function openEditGoal(goal: GoalEntity) {
  editGoalState.id = goal.id
  editGoalState.title = goal.title
  editGoalState.description = goal.description
  editGoalOpen.value = true
}

async function onUpdateGoal(event: FormSubmitEvent<GoalSchema>) {
  if (!editGoalState.id) return

  updatingGoal.value = true

  await withErrorToast(async () => {
    await tracker.updateGoal(
      editGoalState.id as number,
      {
        title: event.data.title.trim(),
        description: event.data.description.trim()
      },
      auth.authHeader()
    )

    editGoalOpen.value = false
    toast.add({ title: 'Goal updated', color: 'success' })
  })

  updatingGoal.value = false
}

async function onDeleteGoal(goalId: number) {
  if (!confirmAction('Delete this goal and all nested tasks?')) return

  deletingGoalId.value = goalId

  await withErrorToast(async () => {
    await tracker.deleteGoal(goalId, auth.authHeader())
    toast.add({ title: 'Goal deleted', color: 'success' })
  })

  deletingGoalId.value = null
}

async function onCreateTask(goalId: number, event: FormSubmitEvent<CreateTaskSchema>) {
  creatingTaskGoalId.value = goalId

  await withErrorToast(async () => {
    const assigneeId = parseOptionalPositiveInt(event.data.assigneeId)

    await tracker.createTask(
      goalId,
      {
        title: event.data.title.trim(),
        description: event.data.description.trim(),
        assigneeId
      },
      auth.authHeader()
    )

    taskDraft(goalId).title = ''
    taskDraft(goalId).description = ''
    taskDraft(goalId).assigneeId = ''

    toast.add({ title: 'Task created', color: 'success' })
  })

  creatingTaskGoalId.value = null
}

function openEditTask(goalId: number, task: TaskEntity) {
  editTaskState.id = task.id
  editTaskState.goalId = goalId
  editTaskState.title = task.title
  editTaskState.description = task.description
  editTaskState.status = task.status || 'todo'
  editTaskState.assigneeId = task.assigneeId ? String(task.assigneeId) : ''
  editTaskOpen.value = true
}

async function onUpdateTask(event: FormSubmitEvent<UpdateTaskSchema>) {
  if (!editTaskState.id || !editTaskState.goalId) return

  updatingTask.value = true

  await withErrorToast(async () => {
    const assigneeId = parseOptionalPositiveInt(event.data.assigneeId)
    const status = parseTaskStatus(event.data.status)

    await tracker.updateTask(
      editTaskState.id as number,
      {
        goalId: editTaskState.goalId,
        title: event.data.title.trim(),
        description: event.data.description.trim(),
        status,
        assigneeId
      },
      auth.authHeader()
    )

    editTaskOpen.value = false
    await tracker.fetchAssignedTasks(auth.authHeader())
    toast.add({ title: 'Task updated', color: 'success' })
  })

  updatingTask.value = false
}

async function onDeleteTask(taskId: number) {
  if (!confirmAction('Delete this task?')) return

  deletingTaskId.value = taskId

  await withErrorToast(async () => {
    await tracker.deleteTask(taskId, auth.authHeader())
    toast.add({ title: 'Task deleted', color: 'success' })
  })

  deletingTaskId.value = null
}

onMounted(async () => {
  await loadGoals()
})
</script>
