const priorityRank = {
  high: 0,
  medium: 1,
  low: 2
}

function parseDateValue(value) {
  const parsed = Date.parse(value || '')
  return Number.isNaN(parsed) ? 0 : parsed
}

function sortGoals(goals = []) {
  return [...goals].sort((left, right) => {
    const leftAchieved = left?.status === 'achieved' ? 1 : 0
    const rightAchieved = right?.status === 'achieved' ? 1 : 0
    if (leftAchieved !== rightAchieved) return leftAchieved - rightAchieved

    const priorityDiff = (priorityRank[left?.priority] ?? 1) - (priorityRank[right?.priority] ?? 1)
    if (priorityDiff !== 0) return priorityDiff

    return parseDateValue(right?.createdAt) - parseDateValue(left?.createdAt)
  })
}

function sortTasks(tasks = []) {
  return [...tasks].sort((left, right) => {
    const leftCompleted = left?.isCompleted ? 1 : 0
    const rightCompleted = right?.isCompleted ? 1 : 0
    if (leftCompleted !== rightCompleted) return leftCompleted - rightCompleted

    const priorityDiff = (priorityRank[left?.priority] ?? 1) - (priorityRank[right?.priority] ?? 1)
    if (priorityDiff !== 0) return priorityDiff

    return parseDateValue(right?.createdAt) - parseDateValue(left?.createdAt)
  })
}

function sortGoalsWithNestedTasks(goals = []) {
  return sortGoals(goals).map((goal) => ({
    ...goal,
    tasks: sortTasks(goal?.tasks || [])
  }))
}

export const useTrackerStore = defineStore('tracker', {
  state: () => ({
    goals: [],
    assignedTasks: [],
    usersTaskBoard: [],
    usersLookup: [],
    loadingGoals: false,
    loadingAssigned: false,
    loadingUsersTaskBoard: false,
    loadingUsersLookup: false
  }),
  actions: {
    async fetchGoals(authHeader = {}) {
      this.loadingGoals = true
      try {
        const goals = await $fetch('/api/goals', {
          headers: authHeader
        })
        this.goals = sortGoalsWithNestedTasks(goals)
      } finally {
        this.loadingGoals = false
      }
    },

    async fetchAssignedTasks(authHeader = {}) {
      this.loadingAssigned = true
      try {
        const tasks = await $fetch('/api/tasks/assigned', {
          headers: authHeader
        })
        this.assignedTasks = sortTasks(tasks)
      } finally {
        this.loadingAssigned = false
      }
    },

    async fetchUsersTaskBoard(authHeader = {}) {
      this.loadingUsersTaskBoard = true
      try {
        const usersTaskBoard = await $fetch('/api/users/tasks', {
          headers: authHeader
        })

        this.usersTaskBoard = (usersTaskBoard || []).map((user) => ({
          ...user,
          tasks: sortTasks(user?.tasks || [])
        }))
      } finally {
        this.loadingUsersTaskBoard = false
      }
      return this.usersTaskBoard
    },

    async refresh(authHeader = {}) {
      await Promise.all([
        this.fetchGoals(authHeader),
        this.fetchAssignedTasks(authHeader)
      ])
    },

    async fetchGoalTasks(goalId, authHeader = {}) {
      const goal = await $fetch(`/api/goals/${goalId}/tasks`, {
        headers: authHeader
      })

      return {
        ...goal,
        tasks: sortTasks(goal?.tasks || [])
      }
    },

    async fetchUsersLookup(authHeader = {}) {
      this.loadingUsersLookup = true
      try {
        this.usersLookup = await $fetch('/api/users/lookup', {
          headers: authHeader
        })
      } finally {
        this.loadingUsersLookup = false
      }
      return this.usersLookup
    },

    upsertTaskInGoals(updatedTask) {
      for (const goal of this.goals) {
        goal.tasks = (goal.tasks || []).filter((task) => task.id !== updatedTask.id)
      }

      const targetGoal = this.goals.find((goal) => goal.id === updatedTask.goalId)
      if (!targetGoal) return

      targetGoal.tasks = sortTasks([...(targetGoal.tasks || []), updatedTask])
      this.goals = sortGoalsWithNestedTasks(this.goals)
    },

    removeTaskFromGoals(taskId) {
      for (const goal of this.goals) {
        goal.tasks = (goal.tasks || []).filter((task) => task.id !== taskId)
      }
      this.goals = sortGoalsWithNestedTasks(this.goals)
      this.assignedTasks = this.assignedTasks.filter((task) => task.id !== taskId)
    },

    async createGoal(payload, authHeader = {}) {
      const created = await $fetch('/api/goals', {
        method: 'POST',
        body: payload,
        headers: authHeader
      })

      this.goals = sortGoalsWithNestedTasks([
        { ...created, tasks: [] },
        ...this.goals
      ])

      return created
    },

    async updateGoal(goalId, payload, authHeader = {}) {
      const updated = await $fetch(`/api/goals/${goalId}`, {
        method: 'PUT',
        body: payload,
        headers: authHeader
      })

      this.goals = sortGoalsWithNestedTasks(this.goals.map((goal) => {
        if (goal.id !== goalId) return goal
        return {
          ...goal,
          ...updated,
          tasks: goal.tasks || []
        }
      }))

      return updated
    },

    async deleteGoal(goalId, authHeader = {}) {
      await $fetch(`/api/goals/${goalId}`, {
        method: 'DELETE',
        headers: authHeader
      })

      const deletedTaskIds = new Set(
        (this.goals.find((goal) => goal.id === goalId)?.tasks || []).map((task) => task.id)
      )

      this.goals = this.goals.filter((goal) => goal.id !== goalId)
      this.assignedTasks = this.assignedTasks.filter((task) => !deletedTaskIds.has(task.id))
    },

    async createTask(goalId, payload, authHeader = {}) {
      const created = await $fetch(`/api/goals/${goalId}/tasks`, {
        method: 'POST',
        body: payload,
        headers: authHeader
      })

      const goal = this.goals.find((item) => item.id === goalId)
      if (goal) {
        goal.tasks = sortTasks([...(goal.tasks || []), created])
      }

      return created
    },

    async updateTask(taskId, payload, authHeader = {}) {
      const updated = await $fetch(`/api/tasks/${taskId}`, {
        method: 'PUT',
        body: payload,
        headers: authHeader
      })

      this.upsertTaskInGoals(updated)
      return updated
    },

    async deleteTask(taskId, authHeader = {}) {
      await $fetch(`/api/tasks/${taskId}`, {
        method: 'DELETE',
        headers: authHeader
      })
      this.removeTaskFromGoals(taskId)
    },

    async assignTask(taskId, assigneeId, authHeader = {}) {
      const updated = await $fetch(`/api/tasks/${taskId}/assign`, {
        method: 'PUT',
        body: { assigneeId },
        headers: authHeader
      })

      this.upsertTaskInGoals(updated)
      this.assignedTasks = sortTasks(this.assignedTasks.map((task) => {
        if (task.id !== taskId) return task
        return { ...task, ...updated }
      }))

      return updated
    }
  }
})
