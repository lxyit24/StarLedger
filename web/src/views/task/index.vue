<template>
  <div>
    <!-- 搜索栏 -->
    <el-card style="margin-bottom: 16px">
      <el-row :gutter="16" align="middle">
        <el-col :span="4">
          <el-select v-model="query.status" placeholder="状态筛选" clearable @change="fetchData">
            <el-option label="待处理" value="pending" />
            <el-option label="进行中" value="in_progress" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="query.priority" placeholder="优先级筛选" clearable @change="fetchData">
            <el-option label="低" value="low" />
            <el-option label="中" value="medium" />
            <el-option label="高" value="high" />
            <el-option label="紧急" value="urgent" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="openDialog()">新增任务</el-button>
        </el-col>
        <el-col :span="4">
          <el-button @click="fetchMyTasks">我的任务</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 数据表格 -->
    <el-card>
      <el-table :data="tableData" stripe v-loading="loading" empty-text="暂无数据">
        <el-table-column prop="title" label="任务标题" min-width="150" show-overflow-tooltip />
        <el-table-column prop="priority" label="优先级" width="80">
          <template #default="{ row }">
            <el-tag :type="priorityType(row.priority)" size="small">{{ priorityText(row.priority) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="负责人" width="100">
          <template #default="{ row }">{{ row.assignee?.real_name || row.assignee?.username || '未分配' }}</template>
        </el-table-column>
        <el-table-column label="创建人" width="100">
          <template #default="{ row }">{{ row.creator?.real_name || row.creator?.username || '-' }}</template>
        </el-table-column>
        <el-table-column prop="due_date" label="截止日期" width="110" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openDialog(row)">编辑</el-button>
            <el-button size="small" type="success" @click="openAssignDialog(row)">分配</el-button>
            <el-popconfirm title="确定删除？" @confirm="handleDelete(row.id)">
              <template #reference><el-button size="small" type="danger">删除</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        style="margin-top: 16px; justify-content: flex-end"
        v-model:current-page="query.page"
        v-model:page-size="query.page_size"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @change="fetchData"
      />
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="editId ? '编辑任务' : '新增任务'" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="90px">
        <el-form-item label="任务标题" prop="title"><el-input v-model="form.title" /></el-form-item>
        <el-form-item label="任务描述"><el-input v-model="form.description" type="textarea" :rows="3" /></el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="优先级">
              <el-select v-model="form.priority" style="width: 100%">
                <el-option label="低" value="low" />
                <el-option label="中" value="medium" />
                <el-option label="高" value="high" />
                <el-option label="紧急" value="urgent" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="截止日期"><el-date-picker v-model="form.due_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" /></el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="状态" v-if="editId">
          <el-select v-model="form.status" style="width: 100%">
            <el-option label="待处理" value="pending" />
            <el-option label="进行中" value="in_progress" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 分配弹窗 -->
    <el-dialog v-model="assignVisible" title="分配任务" width="400px">
      <el-form label-width="80px">
        <el-form-item label="负责人">
          <el-select v-model="assignUserId" placeholder="选择负责人" style="width: 100%">
            <el-option v-for="u in users" :key="u.id" :label="u.real_name || u.username" :value="u.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignVisible = false">取消</el-button>
        <el-button type="primary" :loading="assignLoading" @click="submitAssign">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { listTasks, createTask, updateTask, deleteTask, assignTask, myTasks } from '../../api/task'
import { listUsers } from '../../api/system'

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, status: '', priority: '' })

function priorityType(p: string) {
  return p === 'urgent' ? 'danger' : p === 'high' ? 'warning' : p === 'medium' ? '' : 'info'
}

function priorityText(p: string) {
  return p === 'low' ? '低' : p === 'medium' ? '中' : p === 'high' ? '高' : '紧急'
}

function statusType(s: string) {
  return s === 'completed' ? 'success' : s === 'in_progress' ? '' : s === 'cancelled' ? 'info' : 'warning'
}

function statusText(s: string) {
  return s === 'pending' ? '待处理' : s === 'in_progress' ? '进行中' : s === 'completed' ? '已完成' : '已取消'
}

async function fetchData() {
  loading.value = true
  try {
    const res: any = await listTasks(query.page, query.page_size)
    tableData.value = res.data?.items || res.data || []
    total.value = res.data?.total || tableData.value.length
  } finally { loading.value = false }
}

async function fetchMyTasks() {
  loading.value = true
  try {
    const res: any = await myTasks()
    tableData.value = res.data || []
    total.value = tableData.value.length
  } finally { loading.value = false }
}

// 新增/编辑
const dialogVisible = ref(false)
const editId = ref(0)
const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const form = reactive({
  title: '', description: '', priority: 'medium', due_date: '', status: '',
})
const formRules = {
  title: [{ required: true, message: '请输入任务标题', trigger: 'blur' }],
}

function openDialog(row?: any) {
  if (row) {
    editId.value = row.id
    Object.assign(form, { title: row.title, description: row.description, priority: row.priority, due_date: row.due_date, status: row.status })
  } else {
    editId.value = 0
    Object.assign(form, { title: '', description: '', priority: 'medium', due_date: '', status: '' })
  }
  dialogVisible.value = true
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitLoading.value = true
  try {
    if (editId.value) {
      await updateTask(editId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createTask(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } finally { submitLoading.value = false }
}

async function handleDelete(id: number) {
  await deleteTask(id)
  ElMessage.success('删除成功')
  fetchData()
}

// 分配
const assignVisible = ref(false)
const assignLoading = ref(false)
const assignTaskId = ref(0)
const assignUserId = ref(0)
const users = ref<any[]>([])

async function openAssignDialog(row: any) {
  assignTaskId.value = row.id
  assignUserId.value = row.assignee_id || 0
  // Load users
  try {
    const res: any = await listUsers(1, 100)
    users.value = res.data?.items || res.data || []
  } catch { /* ignore */ }
  assignVisible.value = true
}

async function submitAssign() {
  if (!assignUserId.value) {
    ElMessage.warning('请选择负责人')
    return
  }
  assignLoading.value = true
  try {
    await assignTask(assignTaskId.value, assignUserId.value)
    ElMessage.success('分配成功')
    assignVisible.value = false
    fetchData()
  } finally { assignLoading.value = false }
}

onMounted(fetchData)
</script>
