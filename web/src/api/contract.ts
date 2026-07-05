import request from '../utils/request'

export interface Contract {
  id: number
  tenant_id: number
  title: string
  party_a: string
  party_b: string
  amount: number
  start_date: string
  end_date: string
  status: 'draft' | 'active' | 'expired' | 'terminated'
  file_url: string
  remark: string
  created_at: string
  updated_at: string
}

export function listContracts(page = 1, page_size = 20) {
  return request.get('/contracts', { params: { page, page_size } })
}

export function getContract(id: number) {
  return request.get(`/contracts/${id}`)
}

export function createContract(data: Partial<Contract>) {
  return request.post('/contracts', data)
}

export function updateContract(id: number, data: Partial<Contract>) {
  return request.put(`/contracts/${id}`, data)
}

export function deleteContract(id: number) {
  return request.delete(`/contracts/${id}`)
}
