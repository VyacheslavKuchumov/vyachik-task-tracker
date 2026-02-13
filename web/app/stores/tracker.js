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
        this.goals = await $fetch('/api/goals', {
          headers: authHeader
        })
      } finally {
        this.loadingGoals = false
      }
    },

    async fetchAssignedTasks(authHeader = {}) {
      this.loadingAssigned = true
      try {
        this.assignedTasks = await $fetch('/api/tasks/assigned', {
          headers: authHeader
        })
      } finally {
        this.loadingAssigned = false
      }
    },

    async fetchUsersTaskBoard(authHeader = {}) {
      this.loadingUsersTaskBoard = true
      try {
        this.usersTaskBoard = await $fetch('/api/users/tasks', {
          headers: authHeader
        })
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
      return await $fetch(`/api/goals/${goalId}/tasks`, {
        headers: authHeader
      })
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

      targetGoal.tasks = [...(targetGoal.tasks || []), updatedTask]
    },

    removeTaskFromGoals(taskId) {
      for (const goal of this.goals) {
        goal.tasks = (goal.tasks || []).filter((task) => task.id !== taskId)
      }
      this.assignedTasks = this.assignedTasks.filter((task) => task.id !== taskId)
    },

    async createGoal(payload, authHeader = {}) {
      const created = await $fetch('/api/goals', {
        method: 'POST',
        body: payload,
        headers: authHeader
      })
      this.goals = [
        { ...created, tasks: [] },
        ...this.goals
      ]
      return created
    },

    async updateGoal(goalId, payload, authHeader = {}) {
      const updated = await $fetch(`/api/goals/${goalId}`, {
        method: 'PUT',
        body: payload,
        headers: authHeader
      })

      this.goals = this.goals.map((goal) => {
        if (goal.id !== goalId) return goal
        return {
          ...goal,
          ...updated,
          tasks: goal.tasks || []
        }
      })

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
        goal.tasks = [...(goal.tasks || []), created]
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
      this.assignedTasks = this.assignedTasks.map((task) => {
        if (task.id !== taskId) return task
        return { ...task, ...updated }
      })

      return updated
    }
  }
})
