<template>
  <UCard>
    <template #header>
      <div class="space-y-1">
        <h1 class="text-xl font-semibold">Welcome Back</h1>
        <p class="text-sm text-muted">Log in to manage your goals and tasks.</p>
      </div>
    </template>

    <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
      <UFormField label="Email" name="email" required>
        <UInput v-model="state.email" type="email" size="xl" class="w-full" placeholder="you@example.com" />
      </UFormField>

      <UFormField label="Password" name="password" required>
        <UInput v-model="state.password" type="password" size="xl" class="w-full" placeholder="Your password" />
      </UFormField>

      <UButton type="submit" color="primary" size="xl" block :loading="loading">
        Login
      </UButton>
    </UForm>

    <template #footer>
      <div class="flex items-center justify-center gap-1 text-sm text-muted">
        <span>Need an account?</span>
        <ULink to="/signup">Create one</ULink>
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
  email: v.pipe(v.string(), v.nonEmpty('Email is required'), v.email('Invalid email address')),
  password: v.pipe(v.string(), v.nonEmpty('Password is required'))
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
      title: 'Logged in',
      description: 'JWT session started successfully.',
      color: 'success'
    })

    await navigateTo('/')
  } catch (error: any) {
    toast.add({
      title: 'Login failed',
      description: error?.data?.statusMessage || error?.message || 'Please verify your credentials.',
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}
</script>
