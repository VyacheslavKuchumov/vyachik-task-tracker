<template>
  <UCard>
    <template #header>
      <div class="space-y-1">
        <h1 class="text-xl font-semibold">С возвращением</h1>
        <p class="text-sm text-muted">Войдите, чтобы управлять целями и задачами.</p>
      </div>
    </template>

    <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
      <UFormField label="Эл. почта" name="email" required>
        <UInput v-model="state.email" type="email" size="xl" class="w-full" placeholder="you@example.com" />
      </UFormField>

      <UFormField label="Пароль" name="password" required>
        <UInput v-model="state.password" type="password" size="xl" class="w-full" placeholder="Ваш пароль" />
      </UFormField>

      <UButton type="submit" color="primary" size="xl" block :loading="loading">
        Войти
      </UButton>
    </UForm>

    <template #footer>
      <div class="flex items-center justify-center gap-1 text-sm text-muted">
        <span>Нет аккаунта?</span>
        <ULink to="/signup">Зарегистрироваться</ULink>
      </div>
    </template>
  </UCard>
</template>

<script setup lang="ts">
import * as v from 'valibot'
import type { FormSubmitEvent } from '@nuxt/ui'

const auth = useAuthStore()
const toast = useToast()

const loading = ref(false)

const schema = v.object({
  email: v.pipe(v.string(), v.nonEmpty('Эл. почта обязательна'), v.email('Некорректная эл. почта')),
  password: v.pipe(v.string(), v.nonEmpty('Пароль обязателен'))
})

type LoginSchema = v.InferOutput<typeof schema>

const state = reactive<LoginSchema>({
  email: '',
  password: ''
})

async function onSubmit(event: FormSubmitEvent<LoginSchema>) {
  loading.value = true

  try {
    await auth.login({
      email: event.data.email,
      password: event.data.password
    })

    toast.add({
      title: 'Вход выполнен',
      description: 'Сессия успешно запущена.',
      color: 'success'
    })

    await navigateTo('/')
  } catch (error: any) {
    toast.add({
      title: 'Ошибка входа',
      description: error?.data?.statusMessage || error?.message || 'Проверьте эл. почту и пароль.',
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}
</script>
