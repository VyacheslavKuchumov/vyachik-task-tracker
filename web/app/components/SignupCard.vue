<template>
  <UCard>
    <template #header>
      <div class="space-y-1">
        <h1 class="text-xl font-semibold">Создать аккаунт</h1>
        <p class="text-sm text-muted">Зарегистрируйтесь и начните отслеживать цели.</p>
      </div>
    </template>

    <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <UFormField label="Имя" name="firstName" required>
          <UInput v-model="state.firstName" size="xl" class="w-full" placeholder="Иван" />
        </UFormField>

        <UFormField label="Фамилия" name="lastName" required>
          <UInput v-model="state.lastName" size="xl" class="w-full" placeholder="Иванов" />
        </UFormField>
      </div>

      <UFormField label="Эл. почта" name="email" required>
        <UInput v-model="state.email" type="email" size="xl" class="w-full" placeholder="you@example.com" />
      </UFormField>

      <UFormField label="Пароль" name="password" required>
        <UInput v-model="state.password" type="password" size="xl" class="w-full" placeholder="Минимум 3 символа" />
      </UFormField>

      <UButton type="submit" color="primary" size="xl" block :loading="loading">
        Зарегистрироваться
      </UButton>
    </UForm>

    <template #footer>
      <div class="flex items-center justify-center gap-1 text-sm text-muted">
        <span>Уже есть аккаунт?</span>
        <ULink to="/login">Войти</ULink>
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
  firstName: v.pipe(v.string(), v.minLength(1, 'Имя обязательно')),
  lastName: v.pipe(v.string(), v.minLength(1, 'Фамилия обязательна')),
  email: v.pipe(v.string(), v.nonEmpty('Эл. почта обязательна'), v.email('Некорректная эл. почта')),
  password: v.pipe(v.string(), v.minLength(3, 'Пароль должен быть не короче 3 символов'))
})

type SignupSchema = v.InferOutput<typeof schema>

const state = reactive<SignupSchema>({
  firstName: '',
  lastName: '',
  email: '',
  password: ''
})

async function onSubmit(event: FormSubmitEvent<SignupSchema>) {
  loading.value = true

  try {
    await auth.signup({
      firstName: event.data.firstName,
      lastName: event.data.lastName,
      email: event.data.email,
      password: event.data.password
    })

    toast.add({
      title: 'Регистрация выполнена',
      description: 'Аккаунт создан, вы вошли в систему.',
      color: 'success'
    })

    await navigateTo('/')
  } catch (error: any) {
    toast.add({
      title: 'Ошибка регистрации',
      description: error?.data?.statusMessage || error?.message || 'Попробуйте другой адрес эл. почты.',
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}
</script>
