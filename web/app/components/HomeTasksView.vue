<template>
  <section class="space-y-6">
    <UCard>
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 class="text-xl font-semibold">Home</h1>
            <p class="text-sm text-muted">Tasks currently assigned to you.</p>
          </div>

          <div class="flex items-center gap-2">
            <UBadge color="primary" variant="subtle" size="lg">{{ tracker.assignedTasks.length }}</UBadge>
            <UButton
              icon="i-lucide-refresh-cw"
              color="neutral"
              variant="soft"
              :loading="tracker.loadingAssigned"
              @click="loadAssigned"
            >
              Refresh
            </UButton>
            <UButton to="/goals" icon="i-lucide-folder-kanban" color="primary">
              Open Goals
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
        title="No tasks assigned"
        description="Ask a goal owner to assign a task to your user id."
      />

      <div v-else class="space-y-3">
        <UCard v-for="task in tracker.assignedTasks" :key="task.id" variant="soft">
          <div class="flex items-start justify-between gap-3">
            <div class="space-y-1">
              <p class="font-medium">{{ task.title }}</p>
              <p class="text-sm text-muted">{{ task.description }}</p>
            </div>

            <UBadge :color="statusColor(task.status)" variant="soft">
              {{ task.status || 'todo' }}
            </UBadge>
          </div>

          <div class="mt-3 flex flex-wrap gap-3 text-xs text-muted">
            <span>Goal: {{ task.goalTitle || `#${task.goalId}` }}</span>
            <span>By: {{ task.createdByName || `#${task.createdBy}` }}</span>
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

function statusColor(status) {
  if (status === 'done') return 'success'
  if (status === 'in_progress') return 'warning'
  return 'neutral'
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

async function loadAssigned() {
  await withErrorToast(async () => {
    await tracker.fetchAssignedTasks(auth.authHeader())
  })
}

onMounted(async () => {
  await loadAssigned()
})
</script>
