<template>
  <div class="invoice-page">
    <!-- Summary Cards -->
    <el-row :gutter="16" class="summary-row">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="summary-item">
            <div class="summary-label">本月开票数</div>
            <div class="summary-value">{{ summary.monthly_issued_count || 0 }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="summary-item">
            <div class="summary-label">本月开票金额</div>
            <div class="summary-value">¥{{ (summary.monthly_total_amount || 0).toFixed(2) }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="summary-item">
            <div class="summary-label">草稿数</div>
            <div class="summary-value warning">{{ summary.draft_count || 0 }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="summary-item">
            <div class="summary-label">总发票数</div>
            <div class="summary-value">{{ summary.total_count || 0 }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Toolbar -->
    <el-card style="margin-top: 16px">
      <div class="toolbar">
        <div class="toolbar-left">
          <el-select v-model="query.status" placeholder="发票状态" clearable style="width: 140px" @change="fetchData">
            <el-option label="草稿" value="draft" />
            <el-option label="已开具" value="issued" />
            <el-option label="已作废" value="cancelled" />
            <el-option label="红冲" value="red" />
          </el-select>
          <el-select v-model="query.invoice_type" placeholder="发票类型" clearable style="width: 160px" @change="fetchData">
            <el-option label="增值税普通发票" value="vat_normal" />
            <el-option label="增值税专用发票" value="vat_special" />
          </el-select>
        </div>
        <el-button type="primary" @click="openCreateDialog">
          <el-icon><Plus /></el-icon> 新建发票
        </el-button>
      </div>

      <!-- Table -->
      <el-table :data="list" stripe v-loading="loading" style="margin-top: 16px">
        <el-table-column prop="invoice_no" label="发票号码" width="160" />
        <el-table-column prop="invoice_type" label="类型" width="130">
          <template #default="{ row }">
            <el-tag :type="row.invoice_type === 'vat_special' ? 'danger' : 'info'" size="small">
              {{ row.invoice_type === 'vat_special' ? '专用发票' : '普通发票' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="buyer_name" label="购买方" width="160" show-overflow-tooltip />
        <el-table-column prop="total_amount" label="价税合计" width="120" align="right">
          <template #default="{ row }">¥{{ row.total_amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="tax_amount" label="税额" width="100" align="right">
          <template #default="{ row }">¥{{ row.tax_amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="invoice_date" label="开票日期" width="120">
          <template #default="{ row }">{{ formatDate(row.invoice_date) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 'draft'" type="primary" link size="small" @click="handleIssue(row)">开具</el-button>
            <el-button v-if="row.status === 'draft'" type="warning" link size="small" @click="openEditDialog(row)">编辑</el-button>
            <el-button v-if="row.status === 'issued'" type="danger" link size="small" @click="handleCancel(row)">作废</el-button>
            <el-button v-if="row.status === 'draft'" type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
            <el-button type="info" link size="small" @click="viewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- Create/Edit Dialog -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑发票' : '新建发票'" width="720px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="发票类型" prop="invoice_type">
              <el-select v-model="form.invoice_type" style="width: 100%">
                <el-option label="增值税普通发票" value="vat_normal" />
                <el-option label="增值税专用发票" value="vat_special" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="发票号码" prop="invoice_no">
              <el-input v-model="form.invoice_no" placeholder="留空则自动生成" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">购买方信息</el-divider>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="名称" prop="buyer_name">
              <el-input v-model="form.buyer_name" placeholder="购买方名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="税号" prop="buyer_tax_no">
              <el-input v-model="form.buyer_tax_no" placeholder="纳税人识别号" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="地址电话">
              <el-input v-model="form.buyer_address" placeholder="地址+电话" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="开户行账号">
              <el-input v-model="form.buyer_bank" placeholder="开户行及账号" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">销售方信息</el-divider>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="名称">
              <el-input v-model="form.seller_name" placeholder="销售方名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="税号">
              <el-input v-model="form.seller_tax_no" placeholder="纳税人识别号" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="地址电话">
              <el-input v-model="form.seller_address" placeholder="地址+电话" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="开户行账号">
              <el-input v-model="form.seller_bank" placeholder="开户行及账号" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">金额信息</el-divider>
        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="不含税金额" prop="amount">
              <el-input-number v-model="form.amount" :min="0" :precision="2" style="width: 100%" @change="calcTotal" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="税率(%)" prop="tax_rate">
              <el-input-number v-model="form.tax_rate" :min="0" :max="100" :precision="2" style="width: 100%" @change="calcTax" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="税额">
              <el-input-number v-model="form.tax_amount" :min="0" :precision="2" disabled style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="价税合计">
              <el-input :model-value="'¥' + (form.amount + form.tax_amount).toFixed(2)" disabled />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- Detail Dialog -->
    <el-dialog v-model="detailVisible" title="发票详情" width="600px">
      <el-descriptions :column="2" border v-if="currentInvoice">
        <el-descriptions-item label="发票号码">{{ currentInvoice.invoice_no }}</el-descriptions-item>
        <el-descriptions-item label="发票代码">{{ currentInvoice.invoice_code }}</el-descriptions-item>
        <el-descriptions-item label="发票类型">{{ currentInvoice.invoice_type === 'vat_special' ? '增值税专用发票' : '增值税普通发票' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusTagType(currentInvoice.status)" size="small">{{ statusLabel(currentInvoice.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="购买方">{{ currentInvoice.buyer_name }}</el-descriptions-item>
        <el-descriptions-item label="购买方税号">{{ currentInvoice.buyer_tax_no }}</el-descriptions-item>
        <el-descriptions-item label="销售方">{{ currentInvoice.seller_name }}</el-descriptions-item>
        <el-descriptions-item label="销售方税号">{{ currentInvoice.seller_tax_no }}</el-descriptions-item>
        <el-descriptions-item label="不含税金额">¥{{ currentInvoice.amount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="税额">¥{{ currentInvoice.tax_amount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="价税合计" :span="2">
          <span style="font-size: 18px; font-weight: bold; color: #e6a23c">¥{{ currentInvoice.total_amount?.toFixed(2) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="税率">{{ currentInvoice.tax_rate }}%</el-descriptions-item>
        <el-descriptions-item label="开票日期">{{ formatDate(currentInvoice.invoice_date) }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ currentInvoice.remark }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getInvoices, createInvoice, updateInvoice, issueInvoice, cancelInvoice, deleteInvoice, getInvoiceSummary } from '../../api/invoice'

const loading = ref(false)
const submitLoading = ref(false)
const list = ref<any[]>([])
const total = ref(0)
const summary = ref<any>({})
const dialogVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const formRef = ref<FormInstance>()
const currentInvoice = ref<any>(null)

const query = reactive({ page: 1, pageSize: 20, status: '', invoice_type: '' })

const defaultForm = {
  invoice_no: '', invoice_code: '', invoice_type: 'vat_normal', bill_id: 0,
  amount: 0, tax_amount: 0, tax_rate: 0,
  buyer_name: '', buyer_tax_no: '', buyer_address: '', buyer_bank: '',
  seller_name: '', seller_tax_no: '', seller_address: '', seller_bank: '',
  items_detail: '', remark: '',
}
const form = reactive({ ...defaultForm })

const rules = {
  buyer_name: [{ required: true, message: '请输入购买方名称', trigger: 'blur' }],
  amount: [{ required: true, message: '请输入金额', trigger: 'blur' }],
}

onMounted(() => {
  fetchData()
  fetchSummary()
})

async function fetchData() {
  loading.value = true
  try {
    const res = await getInvoices(query)
    list.value = res.data?.items || []
    total.value = res.data?.total || 0
  } finally {
    loading.value = false
  }
}

async function fetchSummary() {
  try {
    const res = await getInvoiceSummary()
    summary.value = res.data || {}
  } catch { /* ignore */ }
}

function openCreateDialog() {
  isEdit.value = false
  Object.assign(form, defaultForm)
  dialogVisible.value = true
}

function openEditDialog(row: any) {
  isEdit.value = true
  editId.value = row.id
  Object.assign(form, {
    invoice_no: row.invoice_no,
    invoice_code: row.invoice_code || '',
    invoice_type: row.invoice_type,
    amount: row.amount,
    tax_amount: row.tax_amount,
    tax_rate: row.tax_rate,
    buyer_name: row.buyer_name,
    buyer_tax_no: row.buyer_tax_no,
    buyer_address: row.buyer_address || '',
    buyer_bank: row.buyer_bank || '',
    seller_name: row.seller_name || '',
    seller_tax_no: row.seller_tax_no || '',
    seller_address: row.seller_address || '',
    seller_bank: row.seller_bank || '',
    items_detail: row.items_detail || '',
    remark: row.remark || '',
  })
  dialogVisible.value = true
}

function calcTotal() {
  form.tax_amount = Math.round(form.amount * form.tax_rate) / 100
}

function calcTax() {
  calcTotal()
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitLoading.value = true
  try {
    if (isEdit.value) {
      await updateInvoice(editId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createInvoice(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
    fetchSummary()
  } finally {
    submitLoading.value = false
  }
}

async function handleIssue(row: any) {
  await ElMessageBox.confirm('确认开具此发票？开具后将无法编辑。', '开具确认')
  await issueInvoice(row.id)
  ElMessage.success('开具成功')
  fetchData()
  fetchSummary()
}

async function handleCancel(row: any) {
  await ElMessageBox.confirm('确认作废此发票？', '作废确认')
  await cancelInvoice(row.id)
  ElMessage.success('作废成功')
  fetchData()
  fetchSummary()
}

async function handleDelete(row: any) {
  await ElMessageBox.confirm('确认删除此发票？删除后不可恢复。', '删除确认', { type: 'warning' })
  await deleteInvoice(row.id)
  ElMessage.success('删除成功')
  fetchData()
  fetchSummary()
}

function viewDetail(row: any) {
  currentInvoice.value = row
  detailVisible.value = true
}

function statusLabel(status: string) {
  const map: Record<string, string> = { draft: '草稿', issued: '已开具', cancelled: '已作废', red: '红冲' }
  return map[status] || status
}

function statusTagType(status: string) {
  const map: Record<string, string> = { draft: 'info', issued: 'success', cancelled: 'danger', red: 'warning' }
  return map[status] || 'info'
}

function formatDate(d: string) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.summary-row .summary-item {
  text-align: center;
}
.summary-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 8px;
}
.summary-value {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
}
.summary-value.warning {
  color: #e6a23c;
}
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.toolbar-left {
  display: flex;
  gap: 12px;
}
.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
