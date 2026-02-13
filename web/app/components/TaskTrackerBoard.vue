<template>
  <section class="tracker-shell">
    <div class="tracker-summary">
      <UCard>
        <p class="metric-label">Goals</p>
        <p class="metric-value">{{ tracker.goals.length }}</p>
      </UCard>
      <UCard>
        <p class="metric-label">Tasks In Goals</p>
        <p class="metric-value">{{ tasksInGoals }}</p>
      </UCard>
      <UCard>
        <p class="metric-label">Assigned To You</p>
        <p class="metric-value">{{ tracker.assignedTasks.length }}</p>
      </UCard>
    </div>

    <div class="tracker-grid">
      <UCard class="panel panel--goals">
        <template #header>
          <div class="panel-header">
            <div>
              <h2>Goals Workspace</h2>
              <p>Create goals and nested tasks.</p>
            </div>
            <UButton
              icon="i-lucide-refresh-cw"
              color="neutral"
              variant="soft"
              :loading="tracker.loadingGoals || tracker.loadingAssigned"
              @click="loadDashboard"
            >
              Refresh
            </UButton>
          </div>
        </template>

        <form class="composer" @submit.prevent="onCreateGoal">
          <UFormField label="Goal title" required>
            <UInput v-model="goalDraft.title" size="lg" placeholder="Ship task tracker MVP" />
          </UFormField>

          <UFormField label="Description" required>
            <UTextarea v-model="goalDraft.description" :rows="3" placeholder="Clear objective and expected outcome" />
          </UFormField>

          <UButton type="submit" color="primary" :loading="creatingGoal" :disabled="!canCreateGoal">
            Add Goal
          </UButton>
        </form>

        <USeparator class="my-6" />

        <div v-if="tracker.loadingGoals" class="loading-row">
          <UIcon name="i-lucide-loader-circle" class="spin" />
          <span>Loading goals...</span>
        </div>

        <div v-else-if="tracker.goals.length === 0" class="empty-box">
          No goals yet. Create the first one.
        </div>

        <div v-else class="goal-list">
          <UCard v-for="goal in tracker.goals" :key="goal.id" class="goal-card" variant="soft">
            <template #header>
              <div class="goal-head">
                <div>
                  <h3>{{ goal.title }}</h3>
                  <p>{{ goal.description }}</p>
                </div>
                <UBadge color="neutral" variant="subtle">#{{ goal.id }}</UBadge>
              </div>
            </template>

            <div class="task-list">
              <div v-if="(goal.tasks || []).length === 0" class="empty-inline">No tasks in this goal.</div>

              <UCard
                v-for="task in goal.tasks"
                :key="task.id"
                class="task-card"
                variant="outline"
              >
                <div class="task-row">
                  <div>
                    <p class="task-title">{{ task.title }}</p>
                    <p class="task-desc">{{ task.description }}</p>
                  </div>
                  <UBadge :color="statusColor(task.status)" variant="soft">
                    {{ statusLabel(task.status) }}
                  </UBadge>
                </div>

                <div class="task-meta">
                  <span>Assignee: {{ task.assigneeName || (task.assigneeId ? `#${task.assigneeId}` : 'unassigned') }}</span>
                  <span>Creator: {{ task.createdByName || `#${task.createdBy}` }}</span>
                </div>

                <div class="assign-row">
                  <UInput
                    v-model="assignDrafts[task.id]"
                    type="number"
                    min="1"
                    placeholder="user id"
                    class="assign-input"
                  />
                  <UButton color="neutral" variant="soft" @click="setAssignToMe(task.id)">
                    Assign To Me
                  </UButton>
                  <UButton color="primary" variant="soft" :loading="assigningTaskId === task.id" @click="onAssignTask(task.id)">
                    Save Assignee
                  </UButton>
                </div>
              </UCard>
            </div>

            <template #footer>
              <form class="task-composer" @submit.prevent="onCreateTask(goal.id)">
                <UFormField label="Task title" required>
                  <UInput v-model="taskDraft(goal.id).title" placeholder="Define implementation" />
                </UFormField>
                <UFormField label="Task description" required>
                  <UInput v-model="taskDraft(goal.id).description" placeholder="API + UI + tests" />
                </UFormField>
                <UFormField label="Assignee id (optional)">
                  <UInput v-model="taskDraft(goal.id).assigneeId" type="number" min="1" placeholder="e.g. 2" />
                </UFormField>
                <UButton type="submit" :loading="creatingTaskGoalId === goal.id" :disabled="!canCreateTask(goal.id)">
                  Add Task
                </UButton>
              </form>
            </template>
          </UCard>
        </div>
      </UCard>

      <UCard class="panel panel--assigned">
        <template #header>
          <div class="panel-header">
            <div>
              <h2>Assigned To Me</h2>
              <p>Tasks where your user ID is assignee.</p>
            </div>
            <UBadge color="primary" variant="subtle" size="lg">
              {{ tracker.assignedTasks.length }}
            </UBadge>
          </div>
        </template>

        <div v-if="tracker.loadingAssigned" class="loading-row">
          <UIcon name="i-lucide-loader-circle" class="spin" />
          <span>Loading assigned tasks...</span>
        </div>

        <div v-else-if="tracker.assignedTasks.length === 0" class="empty-box">
          No tasks are assigned to your account.
        </div>

        <div v-else class="assigned-list">
          <UCard
            v-for="task in tracker.assignedTasks"
            :key="task.id"
            class="task-card"
            variant="soft"
          >
            <div class="task-row">
              <div>
                <p class="task-title">{{ task.title }}</p>
                <p class="task-desc">{{ task.description }}</p>
              </div>
              <UBadge :color="statusColor(task.status)" variant="soft">
                {{ statusLabel(task.status) }}
              </UBadge>
            </div>
            <div class="task-meta">
              <span>Goal: {{ task.goalTitle || `#${task.goalId}` }}</span>
              <span>By: {{ task.createdByName || `#${task.createdBy}` }}</span>
            </div>
          </UCard>
        </div>
      </UCard>
    </div>
  </section>
</template>

<script setup>
const auth = useAuthStore()
const tracker = useTrackerStore()
const toast = useToast()

const creatingGoal = ref(false)
const creatingTaskGoalId = ref(null)
const assigningTaskId = ref(null)

const goalDraft = reactive({
  title: '',
  description: ''
})

const taskDrafts = reactive({})
const assignDrafts = reactive({})

const tasksInGoals = computed(() => {
  return tracker.goals.reduce((sum, goal) => sum + (goal.tasks?.length || 0), 0)
})

const canCreateGoal = computed(() => {
  return goalDraft.title.trim().length >= 3 && goalDraft.description.trim().length >= 3
})

function taskDraft(goalId) {
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

function canCreateTask(goalId) {
  const draft = taskDraft(goalId)
  return draft.title.trim().length >= 3 && draft.description.trim().length >= 3
}

function parseOptionalPositiveInt(value) {
  if (value === '' || value === null || value === undefined) return null

  const parsed = Number(value)
  if (!Number.isInteger(parsed) || parsed <= 0) {
    throw new Error('Assignee id must be a positive integer')
  }

  return parsed
}

function statusColor(status) {
  if (status === 'done') return 'success'
  if (status === 'in_progress') return 'warning'
  return 'neutral'
}

function statusLabel(status) {
  if (status === 'in_progress') return 'in progress'
  return status || 'todo'
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

async function onCreateGoal() {
  if (!canCreateGoal.value) return

  creatingGoal.value = true

  await withErrorToast(async () => {
    await tracker.createGoal(
      {
        title: goalDraft.title.trim(),
        description: goalDraft.description.trim()
      },
      auth.authHeader()
    )

    goalDraft.title = ''
    goalDraft.description = ''

    toast.add({
      title: 'Goal created',
      color: 'success'
    })
  })

  creatingGoal.value = false
}

async function onCreateTask(goalId) {
  const draft = taskDraft(goalId)
  if (!canCreateTask(goalId)) return

  creatingTaskGoalId.value = goalId

  await withErrorToast(async () => {
    const assigneeId = parseOptionalPositiveInt(draft.assigneeId)

    await tracker.createTask(
      goalId,
      {
        title: draft.title.trim(),
        description: draft.description.trim(),
        assigneeId
      },
      auth.authHeader()
    )

    draft.title = ''
    draft.description = ''
    draft.assigneeId = ''

    await tracker.fetchAssignedTasks(auth.authHeader())

    toast.add({
      title: 'Task created',
      color: 'success'
    })
  })

  creatingTaskGoalId.value = null
}

function setAssignToMe(taskId) {
  if (!auth.userId) return
  assignDrafts[taskId] = String(auth.userId)
}

async function onAssignTask(taskId) {
  assigningTaskId.value = taskId

  await withErrorToast(async () => {
    const assigneeId = parseOptionalPositiveInt(assignDrafts[taskId])

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
