import request from '../utils/request'

export interface Bill {
  id: number
  tenant_id: number
  bill_no: string
  bill_type: string
  related_resource_id: number
  amount: number
  bill_date: string
  due_date: string
  paid_date: string
  payment_status: string
  payment_method: string
  invoice_no: string
  remark: string
  created_at: string
  updated_at: string
}

export function listBills(params?: { page?: number; page_size?: number; payment_status?: string }) {
  return request.get('/bills', { params })
}

export function getBill(id: number) {
  return request.get(`/bills/${id}`)
}

export function createBill(data: Partial<Bill>) {
  return request.post('/bills', data)
}

export function updateBill(id: number, data: Partial<Bill>) {
  return request.put(`/bills/${id}`, data)
}

export function deleteBill(id: number) {
  return request.delete(`/bills/${id}`)
}

export function payBill(id: number, data: { payment_method: string }) {
  return request.put(`/bills/${id}/pay`, data)
}

export function cancelBill(id: number) {
  return request.put(`/bills/${id}/cancel`)
}

export function billSummary() {
  return request.get('/bills/summary')
}

export function batchPayBills(ids: number[], payment_method: string) {
  return request.post('/bills/batch-pay', { ids, payment_method })
}

export function batchDeleteBills(ids: number[]) {
  return request.post('/bills/batch-delete', { ids })
}
