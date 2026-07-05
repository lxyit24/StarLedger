<template>
  <div>
    <el-card style="margin-bottom: 16px">
      <el-row :gutter="16" align="middle">
        <el-col :span="4"><el-select v-model="query.payment_status" placeholder="支付状态" clearable @change="fetchData"><el-option label="待支付" value="pending" /><el-option label="已支付" value="paid" /><el-option label="已逾期" value="overdue" /><el-option label="已取消" value="cancelled" /></el-select></el-col>
        <el-col :span="8">
          <el-button type="primary" @click="openDialog()">新增账单</el-button>
          <el-button v-if="selectedIds.length > 0" type="success" @click="handleBatchPay">批量支付 ({{ selectedIds.length }})</el-button>
          <el-button v-if="selectedIds.length > 0" type="danger" @click="handleBatchDelete">批量删除 ({{ selectedIds.length }})</el-button>
        </el-col>
      </el-row>
    </el-card>
    <el-card>
      <el-table :data="tableData" stripe v-loading="loading" empty-text="暂无数据" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="45" />
        <el-table-column prop="bill_no" label="账单号" width="160" />
        <el-table-column prop="bill_type" label="类型" width="100"><template #default="{ row }">{{ formatBillType(row.bill_type) }}</template></el-table-column>
        <el-table-column prop="amount" label="金额" width="100" align="right"><template #default="{ row }">¥{{ row.amount }}</template></el-table-column>
        <el-table-column prop="bill_date" label="账单日" width="110"><template #default="{ row }">{{ formatDate(row.bill_date) }}</template></el-table-column>
        <el-table-column prop="due_date" label="到期日" width="110"><template #default="{ row }">{{ formatDate(row.due_date) }}</template></el-table-column>
        <el-table-column prop="paid_date" label="支付日" width="110"><template #default="{ row }">{{ formatDate(row.paid_date) }}</template></el-table-column>
        <el-table-column prop="payment_status" label="状态" width="80"><template #default="{ row }"><el-tag :type="paymentStatusType(row.payment_status)" size="small">{{ formatPaymentStatus(row.payment_status) }}</el-tag></template></el-table-column>
        <el-table-column prop="payment_method" label="支付方式" width="100"><template #default="{ row }">{{ formatPaymentMethod(row.payment_method) }}</template></el-table-column>
        <el-table-column prop="remark" label="备注" min-width="120" show-overflow-tooltip />
        <el-table-column label="操作" width="260" fixed="right"><template #default="{ row }"><el-button size="small" @click="openDialog(row)">编辑</el-button><el-button v-if="row.payment_status === 'pending' || row.payment_status === 'overdue'" size="small" type="success" @click="handlePay(row)">支付</el-button><el-button v-if="row.payment_status === 'pending'" size="small" type="warning" @click="handleCancel(row.id)">取消</el-button><el-popconfirm title="确定删除？" @confirm="handleDelete(row.id)"><template #reference><el-button size="small" type="danger">删除</el-button></template></el-popconfirm></template></el-table-column>
      </el-table>
      <el-pagination style="margin-top: 16px; justify-content: flex-end" v-model:current-page="query.page" v-model:page-size="query.page_size" :total="total" layout="total, sizes, prev, pager, next" @change="fetchData" />
    </el-card>
    <el-dialog v-model="dialogVisible" :title="editId ? '编辑账单' : '新增账单'" width="560px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="90px">
        <el-row :gutter="16"><el-col :span="12"><el-form-item label="账单类型" prop="bill_type"><el-select v-model="form.bill_type" style="width: 100%"><el-option label="服务器租赁" value="server_lease" /><el-option label="域名租赁" value="domain" /><el-option label="IDC交易" value="idc" /><el-option label="其他" value="other" /></el-select></el-form-item></el-col><el-col :span="12"><el-form-item label="金额" prop="amount"><el-input-number v-model="form.amount" :min="0" :precision="2" style="width: 100%" /></el-form-item></el-col></el-row>
        <el-row :gutter="16"><el-col :span="12"><el-form-item label="账单日" prop="bill_date"><el-date-picker v-model="form.bill_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" /></el-form-item></el-col><el-col :span="12"><el-form-item label="到期日" prop="due_date"><el-date-picker v-model="form.due_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" /></el-form-item></el-col></el-row>
        <el-form-item label="关联资源ID"><el-input-number v-model="form.related_resource_id" :min="0" style="width: 100%" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" type="textarea" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button></template>
    </el-dialog>
    <el-dialog v-model="payVisible" title="确认支付" width="400px">
      <el-form label-width="80px"><el-form-item label="支付方式"><el-select v-model="payMethod" style="width: 100%"><el-option label="银行转账" value="bank_transfer" /><el-option label="支付宝" value="alipay" /><el-option label="微信" value="wechat" /><el-option label="现金" value="cash" /></el-select></el-form-item></el-form>
      <template #footer><el-button @click="payVisible = false">取消</el-button><el-button type="primary" :loading="payLoading" @click="submitPay">确认支付</el-button></template>
    </el-dialog>
    <el-dialog v-model="batchPayVisible" title="批量支付确认" width="400px">
      <p>将批量支付 <strong>{{ selectedIds.length }}</strong> 个账单</p>
      <el-form label-width="80px"><el-form-item label="支付方式"><el-select v-model="batchPayMethod" style="width: 100%"><el-option label="银行转账" value="bank_transfer" /><el-option label="支付宝" value="alipay" /><el-option label="微信" value="wechat" /><el-option label="现金" value="cash" /></el-select></el-form-item></el-form>
      <template #footer><el-button @click="batchPayVisible = false">取消</el-button><el-button type="primary" :loading="batchPayLoading" @click="submitBatchPay">确认</el-button></template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { listBills, createBill, updateBill, deleteBill, payBill, cancelBill, batchPayBills, batchDeleteBills } from '../../api/bill'
import { formatDate, formatPaymentStatus, paymentStatusType, formatBillType, formatPaymentMethod } from '../../utils/format'
const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, payment_status: '' })
const selectedIds = ref<number[]>([])
function payStatusType(s: string) { return s === 'paid' ? 'success' : s === 'pending' ? 'warning' : s === 'overdue' ? 'danger' : 'info' }
function handleSelectionChange(rows: any[]) { selectedIds.value = rows.map(r => r.id) }
async function fetchData() { loading.value = true; try { const res: any = await listBills(query); tableData.value = res.data?.items || res.data || []; total.value = res.data?.total || tableData.value.length } finally { loading.value = false } }
const dialogVisible = ref(false)
const editId = ref(0)
const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const form = reactive({ bill_type: 'server_lease', amount: 0, bill_date: '', due_date: '', related_resource_id: 0, remark: '' })
const formRules = { bill_type: [{ required: true, message: '请选择类型', trigger: 'change' }], amount: [{ required: true, message: '请输入金额', trigger: 'blur' }], bill_date: [{ required: true, message: '请选择账单日', trigger: 'change' }], due_date: [{ required: true, message: '请选择到期日', trigger: 'change' }] }
function openDialog(row?: any) { if (row) { editId.value = row.id; Object.assign(form, { bill_type: row.bill_type, amount: row.amount, bill_date: row.bill_date, due_date: row.due_date, related_resource_id: row.related_resource_id, remark: row.remark }) } else { editId.value = 0; Object.assign(form, { bill_type: 'server_lease', amount: 0, bill_date: '', due_date: '', related_resource_id: 0, remark: '' }) } dialogVisible.value = true }
async function handleSubmit() { const valid = await formRef.value?.validate().catch(() => false); if (!valid) return; submitLoading.value = true; try { if (editId.value) { await updateBill(editId.value, form); ElMessage.success('更新成功') } else { await createBill(form); ElMessage.success('创建成功') } dialogVisible.value = false; fetchData() } finally { submitLoading.value = false } }
async function handleDelete(id: number) { await deleteBill(id); ElMessage.success('删除成功'); fetchData() }
const payVisible = ref(false)
const payLoading = ref(false)
const payMethod = ref('bank_transfer')
const payId = ref(0)
function handlePay(row: any) { payId.value = row.id; payMethod.value = 'bank_transfer'; payVisible.value = true }
async function submitPay() { payLoading.value = true; try { await payBill(payId.value, { payment_method: payMethod.value }); ElMessage.success('支付成功'); payVisible.value = false; fetchData() } finally { payLoading.value = false } }
async function handleCancel(id: number) { await cancelBill(id); ElMessage.success('已取消'); fetchData() }
// Batch operations
const batchPayVisible = ref(false)
const batchPayLoading = ref(false)
const batchPayMethod = ref('bank_transfer')
function handleBatchPay() { batchPayMethod.value = 'bank_transfer'; batchPayVisible.value = true }
async function submitBatchPay() { batchPayLoading.value = true; try { const res: any = await batchPayBills(selectedIds.value, batchPayMethod.value); ElMessage.success(`成功支付 ${res.data?.paid_count || selectedIds.value.length} 个账单`); batchPayVisible.value = false; fetchData() } finally { batchPayLoading.value = false } }
async function handleBatchDelete() { try { await ElMessageBox.confirm(`确定批量删除 ${selectedIds.value.length} 个账单？仅待支付/已取消的账单可删除`, '批量删除确认'); const res: any = await batchDeleteBills(selectedIds.value); ElMessage.success(`成功删除 ${res.data?.deleted_count || 0} 个账单`); fetchData() } catch { /* cancelled */ } }
onMounted(fetchData)
</script>
