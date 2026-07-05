import request from '../utils/request'

export interface ServerLease {
  id: number
  tenant_id: number
  server_name: string
  provider: string
  config: string
  ip_address: string
  location: string
  monthly_cost: number
  start_date: string
  end_date: string
  renew_cycle: string
  status: string
  contract_no: string
  remark: string
  created_at: string
  updated_at: string
}

export function listServers(params?: { page?: number; page_size?: number; status?: string }) {
  return request.get('/servers', { params })
}

export function getServer(id: number) {
  return request.get(`/servers/${id}`)
}

export function createServer(data: Partial<ServerLease>) {
  return request.post('/servers', data)
}

export function updateServer(id: number, data: Partial<ServerLease>) {
  return request.put(`/servers/${id}`, data)
}

export function deleteServer(id: number) {
  return request.delete(`/servers/${id}`)
}

export function renewServer(id: number, data: { months: number }) {
  return request.put(`/servers/${id}/renew`, data)
}
