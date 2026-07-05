import request from '../utils/request'

export function listAuditLogs(params?: { page?: number; page_size?: number; action?: string; resource_type?: string }) {
  return request.get('/audit/logs', { params })
}

export function createAuditLog(data: { action: string; resource_type?: string; resource_id?: number; detail?: string }) {
  return request.post('/audit/logs', data)
}
