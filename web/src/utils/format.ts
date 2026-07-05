/**
 * Format a date string to YYYY-MM-DD format.
 */
export function formatDate(dateStr: string | null | undefined): string {
  if (!dateStr) return '--'
  try {
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    // Go zero-value check: year 0 or 1
    if (d.getFullYear() <= 1) return '--'
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    return `${y}-${m}-${day}`
  } catch {
    return dateStr
  }
}

/**
 * Format a date string to YYYY-MM-DD HH:mm format.
 */
export function formatDateTime(dateStr: string | null | undefined): string {
  if (!dateStr) return '--'
  try {
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    const h = String(d.getHours()).padStart(2, '0')
    const min = String(d.getMinutes()).padStart(2, '0')
    return `${y}-${m}-${day} ${h}:${min}`
  } catch {
    return dateStr
  }
}

/**
 * Map server status to Chinese label.
 */
export function formatServerStatus(status: string): string {
  const map: Record<string, string> = {
    active: '运行中',
    expired: '已到期',
    disabled: '已停用',
  }
  return map[status] || status
}

/**
 * Map server status to Element Plus tag type.
 */
export function serverStatusType(status: string): string {
  const map: Record<string, string> = {
    active: 'success',
    expired: 'warning',
    disabled: 'info',
  }
  return map[status] || 'info'
}

/**
 * Map payment status to Chinese label.
 */
export function formatPaymentStatus(status: string): string {
  const map: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    overdue: '已逾期',
    cancelled: '已取消',
  }
  return map[status] || status
}

/**
 * Map payment status to Element Plus tag type.
 */
export function paymentStatusType(status: string): string {
  const map: Record<string, string> = {
    pending: 'warning',
    paid: 'success',
    overdue: 'danger',
    cancelled: 'info',
  }
  return map[status] || 'info'
}

/**
 * Map user status to Chinese label.
 */
export function formatUserStatus(status: string): string {
  const map: Record<string, string> = {
    active: '启用',
    disabled: '禁用',
  }
  return map[status] || status
}

/**
 * Map user status to Element Plus tag type.
 */
export function userStatusType(status: string): string {
  const map: Record<string, string> = {
    active: 'success',
    disabled: 'danger',
  }
  return map[status] || 'info'
}

/**
 * Map invoice status to Chinese label.
 */
export function formatInvoiceStatus(status: string): string {
  const map: Record<string, string> = {
    draft: '草稿',
    issued: '已开具',
    cancelled: '已作废',
    red: '红冲',
  }
  return map[status] || status
}

/**
 * Map invoice status to Element Plus tag type.
 */
export function invoiceStatusType(status: string): string {
  const map: Record<string, string> = {
    draft: 'info',
    issued: 'success',
    cancelled: 'danger',
    red: 'warning',
  }
  return map[status] || 'info'
}

/**
 * Map contract status to Chinese label.
 */
export function formatContractStatus(status: string): string {
  const map: Record<string, string> = {
    draft: '草稿',
    active: '生效中',
    expired: '已到期',
    terminated: '已终止',
  }
  return map[status] || status
}

/**
 * Map task status to Chinese label.
 */
export function formatTaskStatus(status: string): string {
  const map: Record<string, string> = {
    pending: '待处理',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消',
  }
  return map[status] || status
}

/**
 * Map task priority to Chinese label.
 */
export function formatTaskPriority(priority: string): string {
  const map: Record<string, string> = {
    low: '低',
    medium: '中',
    high: '高',
    urgent: '紧急',
  }
  return map[priority] || priority
}

/**
 * Map bill type to Chinese label.
 */
export function formatBillType(type: string): string {
  const map: Record<string, string> = {
    server_lease: '服务器租赁',
    lease: '服务器租赁',
    domain: '域名租赁',
    idc: 'IDC交易',
    other: '其他',
  }
  return map[type] || type
}

/**
 * Map payment method to Chinese label.
 */
export function formatPaymentMethod(method: string): string {
  const map: Record<string, string> = {
    bank_transfer: '银行转账',
    alipay: '支付宝',
    wechat: '微信',
    cash: '现金',
  }
  return map[method] || method || '--'
}
