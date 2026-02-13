<template>
  <UCard>
    <template #header>
      <div class="space-y-1">
        <h1 class="text-xl font-semibold">Create Account</h1>
        <p class="text-sm text-muted">Register and start tracking your goals.</p>
      </div>
    </template>

    <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <UFormField label="First Name" name="firstName" required>
          <UInput v-model="state.firstName" size="xl" class="w-full" placeholder="John" />
        </UFormField>

        <UFormField label="Last Name" name="lastName" required>
          <UInput v-model="state.lastName" size="xl" class="w-full" placeholder="Doe" />
        </UFormField>
      </div>

      <UFormField label="Email" name="email" required>
        <UInput v-model="state.email" type="email" size="xl" class="w-full" placeholder="you@example.com" />
      </UFormField>

      <UFormField label="Password" name="password" required>
        <UInput v-model="state.password" type="password" size="xl" class="w-full" placeholder="At least 3 characters" />
      </UFormField>

      <UButton type="submit" color="primary" size="xl" block :loading="loading">
        Register
      </UButton>
    </UForm>

    <template #footer>
      <div class="flex items-center justify-center gap-1 text-sm text-muted">
        <span>Already registered?</span>
        <ULink to="/login">Go to login</ULink>
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
  firstName: v.pipe(v.string(), v.minLength(1, 'First name is required')),
  lastName: v.pipe(v.string(), v.minLength(1, 'Last name is required')),
  email: v.pipe(v.string(), v.nonEmpty('Email is required'), v.email('Invalid email address')),
  password: v.pipe(v.string(), v.minLength(3, 'Password should be at least 3 characters'))
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
      title: 'Registered',
      description: 'Account created and logged in.',
      color: 'success'
    })

    await navigateTo('/')
  } catch (error: any) {
    toast.add({
      title: 'Registration failed',
      description: error?.data?.statusMessage || error?.message || 'Try another email.',
      color: 'error'
    })
  } finally {
    loading.value = false
  }
}
</script>
