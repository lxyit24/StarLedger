import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userId = ref(0)
  const tenantId = ref(0)
  const username = ref('')

  function setLogin(data: { token: string; user_id: number; tenant_id: number; username: string }) {
    token.value = data.token
    userId.value = data.user_id
    tenantId.value = data.tenant_id
    username.value = data.username
    localStorage.setItem('token', data.token)
  }

  function logout() {
    token.value = ''
    userId.value = 0
    tenantId.value = 0
    username.value = ''
    localStorage.removeItem('token')
  }

  return { token, userId, tenantId, username, setLogin, logout }
})
