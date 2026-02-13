<template>
  <UCard class="auth-card">
    <template #header>
      <div class="auth-header">
        <h1>Welcome Back</h1>
        <p>Log in to manage your goals and assigned tasks.</p>
      </div>
    </template>

    <form class="auth-form" @submit.prevent="onSubmit">
      <UFormField label="Email" required>
        <UInput v-model="form.email" type="email" placeholder="you@example.com" size="xl" />
      </UFormField>

      <UFormField label="Password" required>
        <UInput v-model="form.password" type="password" placeholder="Your password" size="xl" />
      </UFormField>

      <UButton
        type="submit"
        color="primary"
        size="xl"
        block
        :loading="loading"
        :disabled="!form.email || !form.password"
      >
        Login
      </UButton>
    </form>

    <template #footer>
      <div class="auth-footer">
        <span>Need an account?</span>
        <NuxtLink to="/signup">Create one</NuxtLink>
      </div>
    </template>
  </UCard>
</template>

<script setup>
const auth = useAuthStore()
const toast = useToast()

const loading = ref(false)
const form = reactive({
  email: '',
  password: ''
})

async function onSubmit() {
  loading.value = true

  try {
    await auth.login({
      email: form.email,
      password: form.password
    })

    toast.add({
      title: 'Logged in',
      description: 'JWT session started successfully.',
      color: 'success'
    })

    await navigateTo('/')
  } catch (error) {
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
