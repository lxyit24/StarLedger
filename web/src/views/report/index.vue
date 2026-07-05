<template>
  <div class="report-page">
    <!-- Overview Cards -->
    <el-row :gutter="16" class="stat-row">
      <el-col :xs="12" :sm="6" :md="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #667eea, #764ba2)">
            <el-icon :size="24"><Money /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">¥{{ formatNum(overview.total_revenue) }}</div>
            <div class="stat-label">累计收入</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb, #f5576c)">
            <el-icon :size="24"><Warning /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">¥{{ formatNum(overview.pending_amount) }}</div>
            <div class="stat-label">待收金额</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe, #00f2fe)">
            <el-icon :size="24"><Monitor /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ overview.server_count || 0 }}</div>
            <div class="stat-label">活跃服务器</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #43e97b, #38f9d7)">
            <el-icon :size="24"><Tickets /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ overview.active_contracts || 0 }}</div>
            <div class="stat-label">生效合同</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Charts Row 1: Monthly Trend -->
    <el-card class="chart-card">
      <template #header>
        <div class="chart-header">
          <span>收支趋势（近12个月）</span>
        </div>
      </template>
      <v-chart :option="trendOption" autoresize style="height: 350px" />
    </el-card>

    <!-- Charts Row 2: Bill Type + Bill Status -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :md="12">
        <el-card class="chart-card">
          <template #header><span>账单类型分布</span></template>
          <v-chart :option="billTypeOption" autoresize style="height: 300px" />
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card class="chart-card">
          <template #header><span>账单状态统计</span></template>
          <v-chart :option="billStatusOption" autoresize style="height: 300px" />
        </el-card>
      </el-col>
    </el-row>

    <!-- Charts Row 3: Server Cost + Task Stats -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :md="12">
        <el-card class="chart-card">
          <template #header><span>服务器成本分析（按服务商）</span></template>
          <v-chart :option="serverCostOption" autoresize style="height: 300px" />
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card class="chart-card">
          <template #header><span>任务完成情况</span></template>
          <v-chart :option="taskOption" autoresize style="height: 300px" />
        </el-card>
      </el-col>
    </el-row>

    <!-- Overdue Details -->
    <el-card style="margin-top: 16px">
      <template #header><span>逾期账单明细</span></template>
      <el-table :data="overdueData" stripe size="small" empty-text="暂无逾期账单">
        <el-table-column prop="bill_no" label="账单号" width="160" />
        <el-table-column prop="bill_type" label="类型" width="100" />
        <el-table-column prop="amount" label="金额" width="100" align="right">
          <template #default="{ row }">¥{{ row.amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="due_date" label="到期日" width="120"><template #default="{ row }">{{ formatDate(row.due_date) }}</template></el-table-column>
        <el-table-column prop="remark" label="备注" min-width="120" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { Money, Warning, Monitor, Tickets } from '@element-plus/icons-vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, PieChart, BarChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components'
import { reportOverview, monthlyTrend, billTypeDistribution, billStatusSummary, serverCostAnalysis, taskStatistics } from '../../api/report'
import { listBills } from '../../api/bill'
import { formatDate } from '../../utils/format'

use([CanvasRenderer, LineChart, PieChart, BarChart, TitleComponent, TooltipComponent, LegendComponent, GridComponent])

const overview = reactive({
  total_revenue: 0, pending_amount: 0, overdue_count: 0, overdue_amount: 0,
  server_count: 0, server_month_cost: 0, active_contracts: 0, contract_amount: 0,
  pending_tasks: 0, month_revenue: 0,
})
const overdueData = ref<any[]>([])

function formatNum(n: number) {
  if (!n) return '0.00'
  return Number(n).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// Monthly trend chart
const trendOption = computed(() => {
  const months = trendData.value.map((d: any) => d.month)
  const totals = trendData.value.map((d: any) => Number(d.total).toFixed(2))
  const paid = trendData.value.map((d: any) => Number(d.paid).toFixed(2))
  const unpaid = trendData.value.map((d: any) => Number(d.unpaid).toFixed(2))
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['总金额', '已收款', '未收款'] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: months, axisLabel: { rotate: 30 } },
    yAxis: { type: 'value', axisLabel: { formatter: '¥{value}' } },
    series: [
      { name: '总金额', type: 'line', data: totals, smooth: true, areaStyle: { opacity: 0.1 }, itemStyle: { color: '#409eff' } },
      { name: '已收款', type: 'line', data: paid, smooth: true, itemStyle: { color: '#67c23a' } },
      { name: '未收款', type: 'line', data: unpaid, smooth: true, itemStyle: { color: '#e6a23c' } },
    ],
  }
})

// Bill type pie chart
const billTypeOption = computed(() => {
  const data = billTypeData.value.map((d: any) => ({ name: d.type, value: Number(d.amount).toFixed(2) }))
  return {
    tooltip: { trigger: 'item', formatter: '{b}: ¥{c} ({d}%)' },
    legend: { orient: 'vertical', left: 'left' },
    series: [{
      type: 'pie', radius: ['40%', '70%'], center: ['60%', '50%'],
      label: { formatter: '{b}\n{d}%' },
      data,
    }],
  }
})

// Bill status bar chart
const billStatusOption = computed(() => {
  const statusNames: Record<string, string> = { pending: '待支付', paid: '已支付', overdue: '已逾期', cancelled: '已取消' }
  const statusColors: Record<string, string> = { pending: '#e6a23c', paid: '#67c23a', overdue: '#f56c6c', cancelled: '#909399' }
  const labels: string[] = []
  const counts: number[] = []
  const amounts: number[] = []
  const colors: string[] = []
  for (const d of billStatusData.value) {
    labels.push(statusNames[d.status] || d.status)
    counts.push(d.count)
    amounts.push(Number(d.amount).toFixed(2))
    colors.push(statusColors[d.status] || '#409eff')
  }
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['笔数', '金额'] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: labels },
    yAxis: [
      { type: 'value', name: '笔数' },
      { type: 'value', name: '金额(¥)', axisLabel: { formatter: '¥{value}' } },
    ],
    series: [
      { name: '笔数', type: 'bar', data: counts.map((v, i) => ({ value: v, itemStyle: { color: colors[i] } })), barWidth: '30%' },
      { name: '金额', type: 'bar', yAxisIndex: 1, data: amounts.map((v, i) => ({ value: v, itemStyle: { color: colors[i], opacity: 0.6 } })), barWidth: '30%' },
    ],
  }
})

// Server cost pie chart
const serverCostOption = computed(() => {
  const data = serverCostData.value.map((d: any) => ({ name: d.provider, value: Number(d.monthly_cost).toFixed(2) }))
  return {
    tooltip: { trigger: 'item', formatter: '{b}: ¥{c}/月 ({d}%)' },
    legend: { orient: 'vertical', left: 'left' },
    series: [{
      type: 'pie', radius: '65%', center: ['60%', '50%'],
      label: { formatter: '{b}\n¥{c}/月' },
      data,
    }],
  }
})

// Task stats donut chart
const taskOption = computed(() => {
  const d = taskStatsData.value
  return {
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { orient: 'vertical', left: 'left' },
    series: [{
      type: 'pie', radius: ['40%', '70%'], center: ['60%', '50%'],
      label: { formatter: '{b}\n{c}个' },
      data: [
        { name: '待处理', value: d.pending || 0, itemStyle: { color: '#e6a23c' } },
        { name: '进行中', value: d.in_progress || 0, itemStyle: { color: '#409eff' } },
        { name: '已完成', value: d.completed || 0, itemStyle: { color: '#67c23a' } },
        { name: '已取消', value: d.cancelled || 0, itemStyle: { color: '#909399' } },
      ],
    }],
  }
})

const trendData = ref<any[]>([])
const billTypeData = ref<any[]>([])
const billStatusData = ref<any[]>([])
const serverCostData = ref<any[]>([])
const taskStatsData = ref<any>({})

onMounted(async () => {
  try {
    const [ovRes, trendRes, btRes, bsRes, scRes, tsRes, overdueRes]: any[] = await Promise.all([
      reportOverview(),
      monthlyTrend(),
      billTypeDistribution(),
      billStatusSummary(),
      serverCostAnalysis(),
      taskStatistics(),
      listBills({ page: 1, page_size: 20, payment_status: 'overdue' }),
    ])
    Object.assign(overview, ovRes.data || {})
    trendData.value = trendRes.data || []
    billTypeData.value = btRes.data || []
    billStatusData.value = bsRes.data || []
    serverCostData.value = scRes.data || []
    taskStatsData.value = tsRes.data || {}
    overdueData.value = overdueRes.data?.items || overdueRes.data || []
  } catch { /* ignore */ }
})
</script>

<style scoped>
.report-page { padding: 20px; }
.stat-row { margin-bottom: 0; }
.stat-card { display: flex; align-items: center; padding: 8px; }
.stat-card :deep(.el-card__body) { display: flex; align-items: center; gap: 14px; width: 100%; }
.stat-icon { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; color: #fff; flex-shrink: 0; }
.stat-value { font-size: 20px; font-weight: 600; color: #303133; }
.stat-label { font-size: 12px; color: #909399; margin-top: 2px; }
.chart-card { margin-bottom: 0; }
.chart-header { font-weight: 500; }
</style>
