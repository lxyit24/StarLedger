import request from '../utils/request'

export function getInvoices(params: any) {
  return request.get('/invoices', { params })
}

export function getInvoice(id: number) {
  return request.get(`/invoices/${id}`)
}

export function createInvoice(data: any) {
  return request.post('/invoices', data)
}

export function updateInvoice(id: number, data: any) {
  return request.put(`/invoices/${id}`, data)
}

export function issueInvoice(id: number) {
  return request.put(`/invoices/${id}/issue`)
}

export function cancelInvoice(id: number) {
  return request.put(`/invoices/${id}/cancel`)
}

export function deleteInvoice(id: number) {
  return request.delete(`/invoices/${id}`)
}

export function getInvoiceSummary() {
  return request.get('/invoices/summary')
}
