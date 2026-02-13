<template>
  <section class="space-y-6">
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
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

      <UCard>
        <div class="space-y-1">
          <p class="text-sm text-muted">Assigned To You</p>
          <p class="text-2xl font-semibold">{{ tracker.assignedTasks.length }}</p>
        </div>
      </UCard>
    </div>

    <div class="grid grid-cols-1 gap-4 xl:grid-cols-3">
      <UCard class="xl:col-span-2">
        <template #header>
          <div class="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 class="text-lg font-semibold">Goals Workspace</h2>
              <p class="text-sm text-muted">Create goals and nested tasks.</p>
            </div>

            <div class="flex items-center gap-2">
              <UButton
                icon="i-lucide-refresh-cw"
                color="neutral"
                variant="soft"
                :loading="tracker.loadingGoals || tracker.loadingAssigned"
                @click="loadDashboard"
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
              <div class="flex flex-wrap items-start justify-between gap-2">
                <div class="space-y-1">
                  <h3 class="font-semibold">{{ goal.title }}</h3>
                  <p class="text-sm text-muted">{{ goal.description }}</p>
                </div>
                <UBadge color="neutral" variant="subtle">#{{ goal.id }}</UBadge>
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

                  <UBadge :color="statusColor(task.status)" variant="soft">
                    {{ statusLabel(task.status) }}
                  </UBadge>
                </div>

                <div class="mt-3 flex flex-wrap gap-3 text-xs text-muted">
                  <span>Assignee: {{ task.assigneeName || (task.assigneeId ? `#${task.assigneeId}` : 'unassigned') }}</span>
                  <span>Creator: {{ task.createdByName || `#${task.createdBy}` }}</span>
                </div>

                <UForm
                  :state="assignState(task.id)"
                  class="mt-4 flex flex-wrap items-end gap-2"
                  @submit="(event) => onAssignTask(task.id, event)"
                >
                  <UFormField label="Assignee ID" name="assigneeId" class="min-w-[140px]">
                    <UInput v-model="assignState(task.id).assigneeId" type="number" min="1" placeholder="e.g. 2" />
                  </UFormField>

                  <UButton type="button" color="neutral" variant="soft" @click="setAssignToMe(task.id)">
                    Assign To Me
                  </UButton>

                  <UButton type="submit" color="primary" variant="soft" :loading="assigningTaskId === task.id">
                    Save Assignee
                  </UButton>
                </UForm>
              </UCard>
            </div>

            <template #footer>
              <UForm
                :schema="taskSchema"
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

      <UCard>
        <template #header>
          <div class="flex items-center justify-between gap-3">
            <div>
              <h2 class="text-lg font-semibold">Assigned To Me</h2>
              <p class="text-sm text-muted">Tasks where your user ID is assignee.</p>
            </div>
            <UBadge color="primary" variant="subtle" size="lg">{{ tracker.assignedTasks.length }}</UBadge>
          </div>
        </template>

        <UProgress v-if="tracker.loadingAssigned" />

        <UAlert
          v-else-if="tracker.assignedTasks.length === 0"
          icon="i-lucide-user-round-check"
          color="neutral"
          variant="soft"
          title="Nothing assigned yet"
        />

        <div v-else class="space-y-3">
          <UCard v-for="task in tracker.assignedTasks" :key="task.id" variant="soft">
            <div class="flex items-start justify-between gap-3">
              <div class="space-y-1">
                <p class="font-medium">{{ task.title }}</p>
                <p class="text-sm text-muted">{{ task.description }}</p>
              </div>

              <UBadge :color="statusColor(task.status)" variant="soft">
                {{ statusLabel(task.status) }}
              </UBadge>
            </div>

            <div class="mt-3 flex flex-wrap gap-3 text-xs text-muted">
              <span>Goal: {{ task.goalTitle || `#${task.goalId}` }}</span>
              <span>By: {{ task.createdByName || `#${task.createdBy}` }}</span>
            </div>
          </UCard>
        </div>
      </UCard>
    </div>

    <UModal v-model:open="createGoalOpen" title="Create Goal">
      <template #body>
        <UForm :schema="goalSchema" :state="goalState" class="space-y-4" @submit="onCreateGoal">
          <UFormField label="Goal title" name="title" required>
            <UInput v-model="goalState.title" class="w-full" placeholder="Ship task tracker MVP" />
          </UFormField>

          <UFormField label="Description" name="description" required>
            <UTextarea
              v-model="goalState.description"
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
  </section>
</template>

<script setup lang="ts">
import * as v from 'valibot'
import type { FormSubmitEvent } from '@nuxt/ui'

const auth = useAuthStore()
const tracker = useTrackerStore()
const toast = useToast()

const createGoalOpen = ref(false)

const creatingGoal = ref(false)
const creatingTaskGoalId = ref<number | null>(null)
const assigningTaskId = ref<number | null>(null)

const goalSchema = v.object({
  title: v.pipe(v.string(), v.minLength(3, 'Goal title should be at least 3 characters')),
  description: v.pipe(v.string(), v.minLength(3, 'Description should be at least 3 characters'))
})

type GoalSchema = v.InferOutput<typeof goalSchema>

const taskSchema = v.object({
  title: v.pipe(v.string(), v.minLength(3, 'Task title should be at least 3 characters')),
  description: v.pipe(v.string(), v.minLength(3, 'Description should be at least 3 characters')),
  assigneeId: v.optional(v.string())
})

type TaskSchema = v.InferOutput<typeof taskSchema>

type AssignSchema = {
  assigneeId?: string
}

const goalState = reactive<GoalSchema>({
  title: '',
  description: ''
})

const taskDrafts = reactive<Record<string, TaskSchema>>({})
const assignDrafts = reactive<Record<number, AssignSchema>>({})

const tasksInGoals = computed(() => {
  return tracker.goals.reduce((sum, goal) => sum + (goal.tasks?.length || 0), 0)
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

function assignState(taskId: number) {
  if (!assignDrafts[taskId]) {
    assignDrafts[taskId] = { assigneeId: '' }
  }

  return assignDrafts[taskId]
}

function parseOptionalPositiveInt(value?: string) {
  if (value === '' || value === null || value === undefined) return null

  const parsed = Number(value)
  if (!Number.isInteger(parsed) || parsed <= 0) {
    throw new Error('Assignee id must be a positive integer')
  }

  return parsed
}

function statusColor(status: string) {
  if (status === 'done') return 'success'
  if (status === 'in_progress') return 'warning'
  return 'neutral'
}

function statusLabel(status: string) {
  if (status === 'in_progress') return 'in progress'
  return status || 'todo'
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

async function loadDashboard() {
  await withErrorToast(async () => {
    await tracker.refresh(auth.authHeader())
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

    goalState.title = ''
    goalState.description = ''
    createGoalOpen.value = false

    toast.add({
      title: 'Goal created',
      color: 'success'
    })
  })

  creatingGoal.value = false
}

async function onCreateTask(goalId: number, event: FormSubmitEvent<TaskSchema>) {
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

    await tracker.fetchAssignedTasks(auth.authHeader())

    toast.add({
      title: 'Task created',
      color: 'success'
    })
  })

  creatingTaskGoalId.value = null
}

function setAssignToMe(taskId: number) {
  if (!auth.userId) return
  assignState(taskId).assigneeId = String(auth.userId)
}

async function onAssignTask(taskId: number, event: FormSubmitEvent<AssignSchema>) {
  assigningTaskId.value = taskId

  await withErrorToast(async () => {
    const assigneeId = parseOptionalPositiveInt(event.data.assigneeId)

    await tracker.assignTask(taskId, assigneeId, auth.authHeader())
    await tracker.refresh(auth.authHeader())

    toast.add({
      title: 'Assignee updated',
      color: 'success'
    })
  })

  assigningTaskId.value = null
}

onMounted(async () => {
  await loadDashboard()
})
</script>
