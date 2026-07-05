<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stat-row">
      <el-col :span="6"><el-card shadow="hover" class="stat-card"><div class="stat-icon" style="background: #409eff"><el-icon :size="28"><Monitor /></el-icon></div><div class="stat-info"><div class="stat-value">{{ stats.serverCount }}</div><div class="stat-label">服务器总数</div></div></el-card></el-col>
      <el-col :span="6"><el-card shadow="hover" class="stat-card"><div class="stat-icon" style="background: #e6a23c"><el-icon :size="28"><Warning /></el-icon></div><div class="stat-info"><div class="stat-value">{{ stats.expiringCount }}</div><div class="stat-label">即将到期</div></div></el-card></el-col>
      <el-col :span="6"><el-card shadow="hover" class="stat-card"><div class="stat-icon" style="background: #f56c6c"><el-icon :size="28"><Document /></el-icon></div><div class="stat-info"><div class="stat-value">{{ stats.overdueBills }}</div><div class="stat-label">逾期账单</div></div></el-card></el-col>
      <el-col :span="6"><el-card shadow="hover" class="stat-card"><div class="stat-icon" style="background: #67c23a"><el-icon :size="28"><Money /></el-icon></div><div class="stat-info"><div class="stat-value">¥{{ stats.totalAmount?.toFixed(2) }}</div><div class="stat-label">本月账单总额</div></div></el-card></el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12"><el-card><template #header><span>即将到期服务器</span></template><el-table :data="expiringServers" stripe size="small" empty-text="暂无数据"><el-table-column prop="server_name" label="名称" /><el-table-column prop="provider" label="服务商" /><el-table-column prop="end_date" label="到期日期" width="120"><template #default="{ row }">{{ formatDate(row.end_date) }}</template></el-table-column><el-table-column prop="status" label="状态" width="80"><template #default="{ row }"><el-tag :type="serverStatusType(row.status)" size="small">{{ formatServerStatus(row.status) }}</el-tag></template></el-table-column></el-table></el-card></el-col>
      <el-col :span="12"><el-card><template #header><span>待处理账单</span></template><el-table :data="pendingBills" stripe size="small" empty-text="暂无数据"><el-table-column prop="bill_no" label="账单号" /><el-table-column prop="amount" label="金额" width="100"><template #default="{ row }">¥{{ row.amount }}</template></el-table-column><el-table-column prop="due_date" label="到期日" width="120"><template #default="{ row }">{{ formatDate(row.due_date) }}</template></el-table-column><el-table-column prop="payment_status" label="状态" width="80"><template #default="{ row }"><el-tag :type="paymentStatusType(row.payment_status)" size="small">{{ formatPaymentStatus(row.payment_status) }}</el-tag></template></el-table-column></el-table></el-card></el-col>
    </el-row>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Monitor, Warning, Document, Money } from '@element-plus/icons-vue'
import { listServers } from '../../api/server'
import { listBills, billSummary } from '../../api/bill'
import { formatDate, formatServerStatus, serverStatusType, formatPaymentStatus, paymentStatusType } from '../../utils/format'
const stats = ref({ serverCount: 0, expiringCount: 0, overdueBills: 0, totalAmount: 0 })
const expiringServers = ref<any[]>([])
const pendingBills = ref<any[]>([])
onMounted(async () => {
  try {
    const [serverRes, billRes, summaryRes]: any[] = await Promise.all([listServers({ page: 1, page_size: 5, status: 'active' }), listBills({ page: 1, page_size: 5, payment_status: 'pending' }), billSummary()])
    const servers = serverRes.data?.items || serverRes.data || []
    const bills = billRes.data?.items || billRes.data || []
    const summary = summaryRes.data || {}
    stats.value = { serverCount: serverRes.data?.total || servers.length, expiringCount: servers.filter((s: any) => { const diff = new Date(s.end_date).getTime() - Date.now(); return diff > 0 && diff < 30 * 86400000 }).length, overdueBills: summary.overdue_count || 0, totalAmount: summary.month_total || 0 }
    expiringServers.value = servers.filter((s: any) => { const diff = new Date(s.end_date).getTime() - Date.now(); return diff > 0 && diff < 30 * 86400000 }).slice(0, 5)
    pendingBills.value = bills.slice(0, 5)
  } catch { }
})
</script>
<style scoped>
.stat-row { margin-bottom: 0; }
.stat-card { display: flex; align-items: center; padding: 10px; }
.stat-card :deep(.el-card__body) { display: flex; align-items: center; gap: 16px; width: 100%; }
.stat-icon { width: 56px; height: 56px; border-radius: 12px; display: flex; align-items: center; justify-content: center; color: #fff; flex-shrink: 0; }
.stat-value { font-size: 24px; font-weight: 600; color: #303133; }
.stat-label { font-size: 13px; color: #909399; margin-top: 4px; }
</style>
