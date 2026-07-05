<template>
  <div>
    <el-card style="margin-bottom: 16px">
      <el-row :gutter="16" align="middle">
        <el-col :span="4">
          <el-select v-model="query.action" placeholder="操作类型" clearable @change="fetchData">
            <el-option label="创建" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
            <el-option label="登录" value="login" />
            <el-option label="登出" value="logout" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="query.resource_type" placeholder="资源类型" clearable @change="fetchData">
            <el-option label="账单" value="bill" />
            <el-option label="合同" value="contract" />
            <el-option label="任务" value="task" />
            <el-option label="服务器" value="server" />
            <el-option label="用户" value="user" />
          </el-select>
        </el-col>
      </el-row>
    </el-card>
    <el-card>
      <el-table :data="tableData" stripe v-loading="loading" empty-text="暂无日志记录">
        <el-table-column prop="created_at" label="时间" width="170">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column prop="username" label="操作人" width="100" />
        <el-table-column prop="action" label="操作" width="80">
          <template #default="{ row }">
            <el-tag :type="actionType(row.action)" size="small">{{ actionText(row.action) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="90" />
        <el-table-column prop="resource_id" label="资源ID" width="80" />
        <el-table-column prop="detail" label="详情" min-width="200" show-overflow-tooltip />
        <el-table-column prop="ip_address" label="IP地址" width="130" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">{{ row.status === 'success' ? '成功' : '失败' }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination style="margin-top: 16px; justify-content: flex-end" v-model:current-page="query.page" v-model:page-size="query.page_size" :total="total" layout="total, sizes, prev, pager, next" @change="fetchData" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { listAuditLogs } from '../../api/audit'

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, action: '', resource_type: '' })

function actionType(a: string) { return a === 'create' ? 'success' : a === 'delete' ? 'danger' : a === 'update' ? 'warning' : 'info' }
function actionText(a: string) { return a === 'create' ? '创建' : a === 'delete' ? '删除' : a === 'update' ? '更新' : a === 'login' ? '登录' : a === 'logout' ? '登出' : a }
function formatTime(t: string) { if (!t) return '-'; const d = new Date(t); return d.toLocaleString('zh-CN') }

async function fetchData() {
  loading.value = true
  try {
    const res: any = await listAuditLogs(query)
    tableData.value = res.data?.items || res.data || []
    total.value = res.data?.total || tableData.value.length
  } finally { loading.value = false }
}

onMounted(fetchData)
</script>
