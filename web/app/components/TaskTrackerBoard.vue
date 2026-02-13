<template>
  <section class="space-y-6">
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
      <UCard>
        <div class="space-y-1">
          <p class="text-sm text-muted">Цели</p>
          <p class="text-2xl font-semibold">{{ tracker.goals.length }}</p>
        </div>
      </UCard>

      <UCard>
        <div class="space-y-1">
          <p class="text-sm text-muted">Задачи по всем целям</p>
          <p class="text-2xl font-semibold">{{ tasksInGoals }}</p>
        </div>
      </UCard>
    </div>

    <UCard>
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold">Цели</h2>
            <p class="text-sm text-muted">Создавайте цели и открывайте отдельную страницу задач для каждой цели.</p>
          </div>

          <div class="flex items-center gap-2">
            <UButton
              icon="i-lucide-refresh-cw"
              color="neutral"
              variant="soft"
              :loading="tracker.loadingGoals"
              @click="loadGoals"
            >
              Обновить
            </UButton>

            <UButton icon="i-lucide-plus" color="primary" @click="createGoalOpen = true">
              Новая цель
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
        title="Пока нет целей"
        description="Создайте первую цель, чтобы начать."
      />

      <div v-else class="space-y-4">
        <UCard v-for="goal in tracker.goals" :key="goal.id" variant="soft">
          <template #header>
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div class="space-y-1">
                <h3 class="font-semibold">{{ goal.title }}</h3>
                <p class="text-sm text-muted">{{ goal.description || 'Без описания' }}</p>
              </div>

              <div class="flex items-center gap-2">
                <UButton icon="i-lucide-pencil" color="neutral" variant="soft" @click="openEditGoal(goal)">
                  Редактировать
                </UButton>
                <UButton
                  icon="i-lucide-trash-2"
                  color="error"
                  variant="soft"
                  :loading="deletingGoalId === goal.id"
                  @click="onDeleteGoal(goal.id)"
                >
                  Удалить
                </UButton>
              </div>
            </div>
          </template>

          <div class="flex flex-wrap items-center justify-between gap-3 text-sm text-muted">
            <span>Владелец: {{ goal.ownerName || auth.displayName }}</span>
            <span>Задачи: {{ (goal.tasks || []).length }}</span>
          </div>

          <template #footer>
            <div class="flex justify-end">
              <UButton :to="`/tasks/${goal.id}`" icon="i-lucide-list-checks" color="primary">
                Перейти к задачам
              </UButton>
            </div>
          </template>
        </UCard>
      </div>
    </UCard>

    <UModal v-model:open="createGoalOpen" title="Создать цель">
      <template #body>
        <UForm :schema="goalSchema" :state="createGoalState" class="space-y-4" @submit="onCreateGoal">
          <UFormField label="Название цели" name="title" required>
            <UInput v-model="createGoalState.title" class="w-full" placeholder="Запустить MVP трекера задач" />
          </UFormField>

          <UFormField label="Описание" name="description">
            <UTextarea
              v-model="createGoalState.description"
              :rows="4"
              class="w-full"
              placeholder="Четкая цель и ожидаемый результат (необязательно)"
            />
          </UFormField>

          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="createGoalOpen = false">
              Отмена
            </UButton>
            <UButton type="submit" color="primary" :loading="creatingGoal">
              Создать цель
            </UButton>
          </div>
        </UForm>
      </template>
    </UModal>

    <UModal v-model:open="editGoalOpen" title="Редактировать цель">
      <template #body>
        <UForm :schema="goalSchema" :state="editGoalState" class="space-y-4" @submit="onUpdateGoal">
          <UFormField label="Название цели" name="title" required>
            <UInput v-model="editGoalState.title" class="w-full" />
          </UFormField>

          <UFormField label="Описание" name="description">
            <UTextarea v-model="editGoalState.description" :rows="4" class="w-full" placeholder="Необязательно" />
          </UFormField>

          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="soft" @click="editGoalOpen = false">
              Отмена
            </UButton>
            <UButton type="submit" color="primary" :loading="updatingGoal">
              Сохранить
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
  ownerName?: string
  tasks?: Array<{ id: number }>
}

const auth = useAuthStore()
const tracker = useTrackerStore()
const toast = useToast()

const createGoalOpen = ref(false)
const editGoalOpen = ref(false)

const creatingGoal = ref(false)
const updatingGoal = ref(false)
const deletingGoalId = ref<number | null>(null)

const goalSchema = v.object({
  title: v.pipe(v.string(), v.minLength(3, 'Название цели должно быть не короче 3 символов')),
  description: v.pipe(v.string(), v.maxLength(2000, 'Описание должно быть не длиннее 2000 символов'))
})

type GoalSchema = v.InferOutput<typeof goalSchema>

const createGoalState = reactive<GoalSchema>({
  title: '',
  description: ''
})

const editGoalState = reactive<{ id: number | null; title: string; description: string }>({
  id: null,
  title: '',
  description: ''
})

const tasksInGoals = computed(() => {
  return tracker.goals.reduce((sum: number, goal: GoalEntity) => sum + (goal.tasks?.length || 0), 0)
})

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

    await tracker.fetchGoals(auth.authHeader())

    toast.add({ title: 'Цель создана', color: 'success' })
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
      editGoalState.id,
      {
        title: event.data.title.trim(),
        description: event.data.description.trim()
      },
      auth.authHeader()
    )

    editGoalOpen.value = false
    await tracker.fetchGoals(auth.authHeader())
    toast.add({ title: 'Цель обновлена', color: 'success' })
  })

  updatingGoal.value = false
}

async function onDeleteGoal(goalId: number) {
  if (!confirmAction('Удалить эту цель и все вложенные задачи?')) return

  deletingGoalId.value = goalId

  await withErrorToast(async () => {
    await tracker.deleteGoal(goalId, auth.authHeader())
    toast.add({ title: 'Цель удалена', color: 'success' })
  })

  deletingGoalId.value = null
}

onMounted(async () => {
  await loadGoals()
})
</script>
