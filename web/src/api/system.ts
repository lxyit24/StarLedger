import request from '../utils/request'

// ===== 用户管理 =====
export function listUsers(params?: { page?: number; page_size?: number }) {
  return request.get('/users', { params })
}

export function createUser(data: { username: string; password: string; real_name: string; email: string; phone: string; role_ids?: number[] }) {
  return request.post('/users', data)
}

export function updateUser(id: number, data: { real_name?: string; email?: string; phone?: string; status?: string; role_ids?: number[] }) {
  return request.put(`/users/${id}`, data)
}

export function deleteUser(id: number) {
  return request.delete(`/users/${id}`)
}

// ===== 角色管理 =====
export function listRoles(params?: { page?: number; page_size?: number }) {
  return request.get('/roles', { params })
}

export function createRole(data: { name: string; description: string; permissions: string[] }) {
  return request.post('/roles', data)
}

export function updateRole(id: number, data: { name?: string; description?: string; permissions?: string[] }) {
  return request.put(`/roles/${id}`, data)
}

export function deleteRole(id: number) {
  return request.delete(`/roles/${id}`)
}

// ===== 租户管理 =====
export function listTenants(params?: { page?: number; page_size?: number }) {
  return request.get('/tenants', { params })
}

export function createTenant(data: { name: string; contact: string; phone: string; email: string }) {
  return request.post('/tenants', data)
}

export function updateTenant(id: number, data: { name?: string; contact?: string; phone?: string; email?: string; status?: string }) {
  return request.put(`/tenants/${id}`, data)
}

export function deleteTenant(id: number) {
  return request.delete(`/tenants/${id}`)
}
