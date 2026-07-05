<template>
  <div>
    <!-- 搜索栏 -->
    <el-card style="margin-bottom: 16px">
      <el-row :gutter="16" align="middle">
        <el-col :span="4">
          <el-select v-model="query.status" placeholder="状态筛选" clearable @change="fetchData">
            <el-option label="草稿" value="draft" />
            <el-option label="生效中" value="active" />
            <el-option label="已过期" value="expired" />
            <el-option label="已终止" value="terminated" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="openDialog()">新增合同</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 数据表格 -->
    <el-card>
      <el-table :data="tableData" stripe v-loading="loading" empty-text="暂无数据">
        <el-table-column prop="title" label="合同标题" min-width="150" show-overflow-tooltip />
        <el-table-column prop="party_a" label="甲方" width="120" show-overflow-tooltip />
        <el-table-column prop="party_b" label="乙方" width="120" show-overflow-tooltip />
        <el-table-column prop="amount" label="金额(元)" width="100" align="right">
          <template #default="{ row }">¥{{ row.amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="start_date" label="开始日期" width="110" />
        <el-table-column prop="end_date" label="结束日期" width="110" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openDialog(row)">编辑</el-button>
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
    <el-dialog v-model="dialogVisible" :title="editId ? '编辑合同' : '新增合同'" width="600px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="90px">
        <el-form-item label="合同标题" prop="title"><el-input v-model="form.title" /></el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="甲方"><el-input v-model="form.party_a" /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="乙方"><el-input v-model="form.party_b" /></el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="金额(元)"><el-input-number v-model="form.amount" :min="0" :precision="2" style="width: 100%" /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" v-if="editId">
              <el-select v-model="form.status" style="width: 100%">
                <el-option label="草稿" value="draft" />
                <el-option label="生效中" value="active" />
                <el-option label="已过期" value="expired" />
                <el-option label="已终止" value="terminated" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="开始日期"><el-date-picker v-model="form.start_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束日期"><el-date-picker v-model="form.end_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" /></el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注"><el-input v-model="form.remark" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { listContracts, createContract, updateContract, deleteContract } from '../../api/contract'

const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, status: '' })

function statusType(s: string) {
  return s === 'active' ? 'success' : s === 'expired' ? 'warning' : s === 'terminated' ? 'danger' : 'info'
}

function statusText(s: string) {
  return s === 'draft' ? '草稿' : s === 'active' ? '生效中' : s === 'expired' ? '已过期' : '已终止'
}

async function fetchData() {
  loading.value = true
  try {
    const res: any = await listContracts(query.page, query.page_size)
    tableData.value = res.data?.items || res.data || []
    total.value = res.data?.total || tableData.value.length
  } finally { loading.value = false }
}

// 新增/编辑
const dialogVisible = ref(false)
const editId = ref(0)
const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const form = reactive({
  title: '', party_a: '', party_b: '', amount: 0, start_date: '', end_date: '', status: '', remark: '',
})
const formRules = {
  title: [{ required: true, message: '请输入合同标题', trigger: 'blur' }],
}

function openDialog(row?: any) {
  if (row) {
    editId.value = row.id
    Object.assign(form, { title: row.title, party_a: row.party_a, party_b: row.party_b, amount: row.amount, start_date: row.start_date, end_date: row.end_date, status: row.status, remark: row.remark })
  } else {
    editId.value = 0
    Object.assign(form, { title: '', party_a: '', party_b: '', amount: 0, start_date: '', end_date: '', status: '', remark: '' })
  }
  dialogVisible.value = true
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitLoading.value = true
  try {
    if (editId.value) {
      await updateContract(editId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createContract(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } finally { submitLoading.value = false }
}

async function handleDelete(id: number) {
  await deleteContract(id)
  ElMessage.success('删除成功')
  fetchData()
}

onMounted(fetchData)
</script>
