import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userId = ref(0)
  const tenantId = ref(0)
  const tenantType = ref('')
  const username = ref('')

  function setLogin(data: { token: string; user_id: number; tenant_id: number; tenant_type: string; username: string }) {
    token.value = data.token
    userId.value = data.user_id
    tenantId.value = data.tenant_id
    tenantType.value = data.tenant_type
    username.value = data.username
    localStorage.setItem('token', data.token)
  }

  function logout() {
    token.value = ''
    userId.value = 0
    tenantId.value = 0
    tenantType.value = ''
    username.value = ''
    localStorage.removeItem('token')
  }

  return { token, userId, tenantId, tenantType, username, setLogin, logout }
})
