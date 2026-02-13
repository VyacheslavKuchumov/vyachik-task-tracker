import { expect, test, type Page } from '@playwright/test'

const BACKEND_URL = process.env.PLAYWRIGHT_BACKEND_URL || 'http://127.0.0.1:8000'

type TestUser = {
  firstName: string
  lastName: string
  email: string
  password: string
}

function buildUser(prefix: string): TestUser {
  const suffix = `${Date.now()}_${Math.floor(Math.random() * 100000)}`
  return {
    firstName: `${prefix}Имя_${suffix}`,
    lastName: `${prefix}Фамилия_${suffix}`,
    email: `playwright_${prefix.toLowerCase()}_${suffix}@example.com`,
    password: `Pwd_${suffix}`
  }
}

async function waitForBackend() {
  const startedAt = Date.now()
  let lastError = ''

  while (Date.now() - startedAt < 60_000) {
    try {
      const response = await fetch(`${BACKEND_URL}/api/v1/login`, {
        method: 'POST',
        headers: { 'content-type': 'application/json' },
        body: '{}'
      })

      // Any non-5xx status means backend is reachable.
      if (response.status < 500) return

      lastError = `status ${response.status}`
    } catch (error: any) {
      lastError = error?.message || String(error)
    }

    await new Promise((resolve) => setTimeout(resolve, 1000))
  }

  throw new Error(
    `Backend is not reachable at ${BACKEND_URL}. Start postgres+server before e2e run. Last error: ${lastError}`
  )
}

async function signUp(page: Page, user: TestUser) {
  await page.goto('/signup')

  await page.getByLabel('Имя').fill(user.firstName)
  await page.getByLabel('Фамилия').fill(user.lastName)
  await page.getByLabel('Эл. почта').fill(user.email)
  await page.getByLabel('Пароль').fill(user.password)

  await Promise.all([
    page.waitForURL('**/'),
    page.getByRole('button', { name: 'Зарегистрироваться' }).click()
  ])
}

test.beforeAll(async () => {
  await waitForBackend()
})

test('user can register and complete goal/task flow', async ({ page }) => {
  const user = buildUser('Flow')
  const fullName = `${user.firstName} ${user.lastName}`
  const goalTitle = `Цель e2e ${Date.now()}`
  const goalDescription = `Описание цели ${Date.now()}`
  const taskTitle = `Задача e2e ${Date.now()}`
  const taskDescription = `Описание задачи ${Date.now()}`

  await signUp(page, user)
  await expect(page.getByRole('heading', { name: 'Главная' })).toBeVisible()

  await page.getByRole('link', { name: 'Цели', exact: true }).click()
  await expect(page).toHaveURL(/\/goals$/)

  await page.getByRole('button', { name: 'Новая цель' }).click()

  const createGoalModal = page.locator('[role="dialog"]').filter({ hasText: 'Создать цель' }).first()
  await expect(createGoalModal).toBeVisible()

  await createGoalModal.getByLabel('Название цели').fill(goalTitle)
  await createGoalModal.getByLabel('Описание').fill(goalDescription)
  await createGoalModal.getByRole('button', { name: 'Создать цель' }).click()

  await expect(createGoalModal).toBeHidden()
  await expect(page.getByText(goalTitle)).toBeVisible()

  await page.getByRole('link', { name: 'Перейти к задачам', exact: true }).click()
  await expect(page).toHaveURL(/\/tasks\/\d+$/)

  await page.getByRole('button', { name: 'Новая задача' }).click()

  const createTaskModal = page.locator('[role="dialog"]').filter({ hasText: 'Создать задачу' }).first()
  await expect(createTaskModal).toBeVisible()

  await createTaskModal.getByLabel('Название задачи').fill(taskTitle)
  await createTaskModal.getByLabel('Описание').fill(taskDescription)
  await createTaskModal.locator('select').first().selectOption({ label: fullName })
  await createTaskModal.getByRole('button', { name: 'Создать' }).click()

  await expect(createTaskModal).toBeHidden()
  await expect(page.getByText(taskTitle)).toBeVisible()
  await expect(page.getByText('К выполнению')).toBeVisible()

  await page.getByRole('button', { name: 'Редактировать' }).first().click()
  const editTaskModal = page.locator('[role="dialog"]').filter({ hasText: 'Редактировать задачу' }).first()
  await expect(editTaskModal).toBeVisible()

  await editTaskModal.locator('select').first().selectOption('in_progress')
  await editTaskModal.getByRole('button', { name: 'Сохранить' }).click()
  await expect(editTaskModal).toBeHidden()
  await expect(page.getByText('В работе')).toBeVisible()

  await page.getByRole('link', { name: 'Главная', exact: true }).click()
  await expect(page).toHaveURL(/\/$/)
  await expect(page.getByText(taskTitle).first()).toBeVisible()

  await page.getByRole('link', { name: 'Пользователи', exact: true }).click()
  await expect(page).toHaveURL(/\/users$/)
  await expect(page.getByText(fullName).first()).toBeVisible()
  await expect(page.getByText(taskTitle)).toBeVisible()
})

test('user can update profile and change password', async ({ page }) => {
  const user = buildUser('Profile')
  const updatedFirstName = `${user.firstName}_Обновлено`
  const updatedLastName = `${user.lastName}_Обновлено`
  const updatedPassword = `${user.password}_new`

  await signUp(page, user)
  await expect(page.getByRole('heading', { name: 'Главная' })).toBeVisible()

  await page.getByRole('link', { name: 'Профиль', exact: true }).click()
  await expect(page).toHaveURL(/\/profile$/)

  await page.getByLabel('Имя').fill(updatedFirstName)
  await page.getByLabel('Фамилия').fill(updatedLastName)
  await page.getByRole('button', { name: 'Сохранить профиль' }).click()

  await expect(page.getByText(`${updatedFirstName} ${updatedLastName}`)).toBeVisible()

  await page.getByLabel('Текущий пароль').fill(user.password)
  await page.getByLabel('Новый пароль').fill(updatedPassword)
  await page.getByRole('button', { name: 'Обновить пароль' }).click()

  await page.getByRole('button', { name: 'Выйти' }).click()
  await expect(page).toHaveURL(/\/login$/)

  await page.getByLabel('Эл. почта').fill(user.email)
  await page.getByLabel('Пароль').fill(updatedPassword)
  await page.getByRole('button', { name: 'Войти' }).click()

  await expect(page).toHaveURL(/\/$/)
  await expect(page.getByRole('heading', { name: 'Главная' })).toBeVisible()
})
