import request from '../utils/request'

export interface Task {
  id: number
  tenant_id: number
  title: string
  description: string
  assignee_id: number
  creator_id: number
  status: 'pending' | 'in_progress' | 'completed' | 'cancelled'
  priority: 'low' | 'medium' | 'high' | 'urgent'
  due_date: string
  created_at: string
  updated_at: string
  assignee?: { id: number; username: string; real_name: string }
  creator?: { id: number; username: string; real_name: string }
}

export function listTasks(page = 1, page_size = 20, assignee_id = 0) {
  return request.get('/tasks', { params: { page, page_size, assignee_id } })
}

export function getTask(id: number) {
  return request.get(`/tasks/${id}`)
}

export function createTask(data: Partial<Task>) {
  return request.post('/tasks', data)
}

export function updateTask(id: number, data: Partial<Task>) {
  return request.put(`/tasks/${id}`, data)
}

export function assignTask(id: number, assignee_id: number) {
  return request.put(`/tasks/${id}/assign`, { assignee_id })
}

export function deleteTask(id: number) {
  return request.delete(`/tasks/${id}`)
}

export function myTasks() {
  return request.get('/tasks/my')
}
