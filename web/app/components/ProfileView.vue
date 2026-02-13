<template>
  <section class="grid grid-cols-1 gap-4 lg:grid-cols-2">
    <UCard>
      <template #header>
        <div class="space-y-1">
          <h1 class="text-xl font-semibold">Профиль</h1>
          <p class="text-sm text-muted">Обновите персональные данные.</p>
        </div>
      </template>

      <UForm :schema="profileSchema" :state="profileState" class="space-y-4" @submit="onUpdateProfile">
        <UFormField label="Имя" name="firstName" required>
          <UInput v-model="profileState.firstName" class="w-full" />
        </UFormField>

        <UFormField label="Фамилия" name="lastName" required>
          <UInput v-model="profileState.lastName" class="w-full" />
        </UFormField>

        <UFormField label="Эл. почта">
          <UInput :model-value="auth.profile?.email || ''" class="w-full" disabled />
        </UFormField>

        <UButton type="submit" color="primary" :loading="savingProfile">
          Сохранить профиль
        </UButton>
      </UForm>
    </UCard>

    <UCard>
      <template #header>
        <div class="space-y-1">
          <h2 class="text-lg font-semibold">Пароль</h2>
          <p class="text-sm text-muted">Измените пароль аккаунта.</p>
        </div>
      </template>

      <UForm :schema="passwordSchema" :state="passwordState" class="space-y-4" @submit="onChangePassword">
        <UFormField label="Текущий пароль" name="currentPassword" required>
          <UInput v-model="passwordState.currentPassword" type="password" class="w-full" />
        </UFormField>

        <UFormField label="Новый пароль" name="newPassword" required>
          <UInput v-model="passwordState.newPassword" type="password" class="w-full" />
        </UFormField>

        <UFormField label="Повторите новый пароль" name="confirmNewPassword" required>
          <UInput v-model="passwordState.confirmNewPassword" type="password" class="w-full" />
        </UFormField>

        <UButton type="submit" color="primary" :loading="changingPassword">
          Обновить пароль
        </UButton>
      </UForm>
    </UCard>
  </section>
</template>

<script setup lang="ts">
import * as v from 'valibot'
import type { FormSubmitEvent } from '@nuxt/ui'

const auth = useAuthStore()
const toast = useToast()

const savingProfile = ref(false)
const changingPassword = ref(false)

const profileSchema = v.object({
  firstName: v.pipe(v.string(), v.minLength(1, 'Имя обязательно')),
  lastName: v.pipe(v.string(), v.minLength(1, 'Фамилия обязательна'))
})

const passwordSchema = v.object({
  currentPassword: v.pipe(v.string(), v.minLength(3, 'Текущий пароль обязателен')),
  newPassword: v.pipe(v.string(), v.minLength(3, 'Новый пароль должен быть не короче 3 символов')),
  confirmNewPassword: v.pipe(v.string(), v.minLength(3, 'Подтвердите новый пароль'))
})

type ProfileSchema = v.InferOutput<typeof profileSchema>
type PasswordSchema = v.InferOutput<typeof passwordSchema>

const profileState = reactive<ProfileSchema>({
  firstName: '',
  lastName: ''
})

const passwordState = reactive<PasswordSchema>({
  currentPassword: '',
  newPassword: '',
  confirmNewPassword: ''
})

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

async function loadProfile() {
  await withErrorToast(async () => {
    const profile = await auth.fetchProfile()
    profileState.firstName = profile?.firstName || ''
    profileState.lastName = profile?.lastName || ''
  })
}

async function onUpdateProfile(event: FormSubmitEvent<ProfileSchema>) {
  savingProfile.value = true

  await withErrorToast(async () => {
    await auth.updateProfile({
      firstName: event.data.firstName.trim(),
      lastName: event.data.lastName.trim()
    })

    toast.add({ title: 'Профиль обновлен', color: 'success' })
  })

  savingProfile.value = false
}

async function onChangePassword(event: FormSubmitEvent<PasswordSchema>) {
  changingPassword.value = true

  await withErrorToast(async () => {
    if (event.data.newPassword !== event.data.confirmNewPassword) {
      throw new Error('Новые пароли не совпадают')
    }

    await auth.changePassword({
      currentPassword: event.data.currentPassword,
      newPassword: event.data.newPassword
    })

    passwordState.currentPassword = ''
    passwordState.newPassword = ''
    passwordState.confirmNewPassword = ''

    toast.add({ title: 'Пароль обновлен', color: 'success' })
  })

  changingPassword.value = false
}

onMounted(async () => {
  await loadProfile()
})
</script>
