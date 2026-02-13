export const useTrackerStore = defineStore('tracker', {
  state: () => ({
    goals: [],
    assignedTasks: [],
    loadingGoals: false,
    loadingAssigned: false
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

    async refresh(authHeader = {}) {
      await Promise.all([
        this.fetchGoals(authHeader),
        this.fetchAssignedTasks(authHeader)
      ])
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

    async assignTask(taskId, assigneeId, authHeader = {}) {
      const updated = await $fetch(`/api/tasks/${taskId}/assign`, {
        method: 'PUT',
        body: { assigneeId },
        headers: authHeader
      })

      for (const goal of this.goals) {
        const task = (goal.tasks || []).find((item) => item.id === taskId)
        if (task) {
          task.assigneeId = updated.assigneeId
          task.assigneeName = updated.assigneeName || task.assigneeName || ''
        }
      }

      return updated
    }
  }
})
