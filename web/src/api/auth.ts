import request from '../utils/request'

export function login(data: { username: string; password: string }) {
  return request.post('/auth/login', data)
}

export function register(data: {
  tenant_name: string
  tenant_type: string
  username: string
  password: string
  real_name: string
  email: string
  phone: string
}) {
  return request.post('/auth/register', data)
}

export function changePassword(data: { old_password: string; new_password: string }) {
  return request.put('/auth/password', data)
}

export function getProfile() {
  return request.get('/auth/profile')
}
