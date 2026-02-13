<template>
  <UCard class="auth-card">
    <template #header>
      <div class="auth-header">
        <h1>Create Account</h1>
        <p>Register and immediately start planning goals and tasks.</p>
      </div>
    </template>

    <form class="auth-form" @submit.prevent="onSubmit">
      <div class="auth-grid">
        <UFormField label="First Name" required>
          <UInput v-model="form.firstName" placeholder="John" size="xl" />
        </UFormField>

        <UFormField label="Last Name" required>
          <UInput v-model="form.lastName" placeholder="Doe" size="xl" />
        </UFormField>
      </div>

      <UFormField label="Email" required>
        <UInput v-model="form.email" type="email" placeholder="you@example.com" size="xl" />
      </UFormField>

      <UFormField label="Password" required>
        <UInput v-model="form.password" type="password" placeholder="At least 3 characters" size="xl" />
      </UFormField>

      <UButton
        type="submit"
        color="primary"
        size="xl"
        block
        :loading="loading"
        :disabled="!isValid"
      >
        Register
      </UButton>
    </form>

    <template #footer>
      <div class="auth-footer">
        <span>Already registered?</span>
        <NuxtLink to="/login">Go to login</NuxtLink>
      </div>
    </template>
  </UCard>
</template>

<script setup>
const auth = useAuthStore()
const toast = useToast()

const loading = ref(false)
const form = reactive({
  firstName: '',
  lastName: '',
  email: '',
  password: ''
})

const isValid = computed(() => {
  return (
    form.firstName.trim().length > 0 &&
    form.lastName.trim().length > 0 &&
    form.email.trim().length > 0 &&
    form.password.length >= 3
  )
})

async function onSubmit() {
  if (!isValid.value) return

  loading.value = true

  try {
    await auth.signup({
      firstName: form.firstName,
      lastName: form.lastName,
      email: form.email,
      password: form.password
    })

    toast.add({
      title: 'Registered',
      description: 'Account created and logged in.',
      color: 'success'
    })

    await navigateTo('/')
  } catch (error) {
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
