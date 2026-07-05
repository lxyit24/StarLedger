<template>
  <div>
    <el-card style="margin-bottom: 16px">
      <el-row :gutter="16" align="middle">
        <el-col :span="4"><el-select v-model="query.status" placeholder="状态筛选" clearable @change="fetchData"><el-option label="使用中" value="active" /><el-option label="已到期" value="expired" /><el-option label="已退租" value="cancelled" /></el-select></el-col>
        <el-col :span="4"><el-button type="primary" @click="openDialog()">新增服务器</el-button></el-col>
      </el-row>
    </el-card>
    <el-card>
      <el-table :data="tableData" stripe v-loading="loading" empty-text="暂无数据">
        <el-table-column prop="server_name" label="服务器名称" min-width="120" />
        <el-table-column prop="provider" label="服务商" width="100" />
        <el-table-column prop="ip_address" label="IP 地址" width="140" />
        <el-table-column prop="config" label="配置" min-width="120" show-overflow-tooltip />
        <el-table-column prop="monthly_cost" label="月费(元)" width="90" align="right"><template #default="{ row }">¥{{ row.monthly_cost }}</template></el-table-column>
        <el-table-column prop="start_date" label="开始日期" width="110" />
        <el-table-column prop="end_date" label="到期日期" width="110" />
        <el-table-column prop="status" label="状态" width="80"><template #default="{ row }"><el-tag :type="statusType(row.status)" size="small">{{ row.status }}</el-tag></template></el-table-column>
        <el-table-column label="操作" width="240" fixed="right"><template #default="{ row }"><el-button size="small" @click="openDialog(row)">编辑</el-button><el-button size="small" type="success" @click="handleRenew(row)">续租</el-button><el-popconfirm title="确定删除？" @confirm="handleDelete(row.id)"><template #reference><el-button size="small" type="danger">删除</el-button></template></el-popconfirm></template></el-table-column>
      </el-table>
      <el-pagination style="margin-top: 16px; justify-content: flex-end" v-model:current-page="query.page" v-model:page-size="query.page_size" :total="total" layout="total, sizes, prev, pager, next" @change="fetchData" />
    </el-card>
    <el-dialog v-model="dialogVisible" :title="editId ? '编辑服务器' : '新增服务器'" width="600px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="90px">
        <el-row :gutter="16"><el-col :span="12"><el-form-item label="名称" prop="server_name"><el-input v-model="form.server_name" /></el-form-item></el-col><el-col :span="12"><el-form-item label="服务商" prop="provider"><el-input v-model="form.provider" /></el-form-item></el-col></el-row>
        <el-row :gutter="16"><el-col :span="12"><el-form-item label="IP 地址"><el-input v-model="form.ip_address" /></el-form-item></el-col><el-col :span="12"><el-form-item label="机房位置"><el-input v-model="form.location" /></el-form-item></el-col></el-row>
        <el-form-item label="配置"><el-input v-model="form.config" placeholder="如: 8核16G 500G SSD" /></el-form-item>
        <el-row :gutter="16"><el-col :span="12"><el-form-item label="月费(元)" prop="monthly_cost"><el-input-number v-model="form.monthly_cost" :min="0" :precision="2" style="width: 100%" /></el-form-item></el-col><el-col :span="12"><el-form-item label="续租周期"><el-select v-model="form.renew_cycle" style="width: 100%"><el-option label="月付" value="monthly" /><el-option label="季付" value="quarterly" /><el-option label="年付" value="yearly" /></el-select></el-form-item></el-col></el-row>
        <el-row :gutter="16"><el-col :span="12"><el-form-item label="开始日期" prop="start_date"><el-date-picker v-model="form.start_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" /></el-form-item></el-col><el-col :span="12"><el-form-item label="到期日期" prop="end_date"><el-date-picker v-model="form.end_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" /></el-form-item></el-col></el-row>
        <el-form-item label="合同编号"><el-input v-model="form.contract_no" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" type="textarea" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button></template>
    </el-dialog>
    <el-dialog v-model="renewVisible" title="续租" width="400px">
      <el-form label-width="80px"><el-form-item label="续租月数"><el-input-number v-model="renewMonths" :min="1" :max="36" /></el-form-item></el-form>
      <template #footer><el-button @click="renewVisible = false">取消</el-button><el-button type="primary" :loading="renewLoading" @click="submitRenew">确定</el-button></template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { listServers, createServer, updateServer, deleteServer, renewServer } from '../../api/server'
const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, status: '' })
function statusType(s: string) { return s === 'active' ? 'success' : s === 'expired' ? 'warning' : 'info' }
async function fetchData() { loading.value = true; try { const res: any = await listServers(query); tableData.value = res.data?.items || res.data || []; total.value = res.data?.total || tableData.value.length } finally { loading.value = false } }
const dialogVisible = ref(false)
const editId = ref(0)
const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const form = reactive({ server_name: '', provider: '', ip_address: '', location: '', config: '', monthly_cost: 0, start_date: '', end_date: '', renew_cycle: 'monthly', contract_no: '', remark: '' })
const formRules = { server_name: [{ required: true, message: '请输入名称', trigger: 'blur' }], provider: [{ required: true, message: '请输入服务商', trigger: 'blur' }], monthly_cost: [{ required: true, message: '请输入月费', trigger: 'blur' }], start_date: [{ required: true, message: '请选择开始日期', trigger: 'change' }], end_date: [{ required: true, message: '请选择到期日期', trigger: 'change' }] }
function openDialog(row?: any) { if (row) { editId.value = row.id; Object.assign(form, { server_name: row.server_name, provider: row.provider, ip_address: row.ip_address, location: row.location, config: row.config, monthly_cost: row.monthly_cost, start_date: row.start_date, end_date: row.end_date, renew_cycle: row.renew_cycle, contract_no: row.contract_no, remark: row.remark }) } else { editId.value = 0; Object.assign(form, { server_name: '', provider: '', ip_address: '', location: '', config: '', monthly_cost: 0, start_date: '', end_date: '', renew_cycle: 'monthly', contract_no: '', remark: '' }) } dialogVisible.value = true }
async function handleSubmit() { const valid = await formRef.value?.validate().catch(() => false); if (!valid) return; submitLoading.value = true; try { if (editId.value) { await updateServer(editId.value, form); ElMessage.success('更新成功') } else { await createServer(form); ElMessage.success('创建成功') } dialogVisible.value = false; fetchData() } finally { submitLoading.value = false } }
async function handleDelete(id: number) { await deleteServer(id); ElMessage.success('删除成功'); fetchData() }
const renewVisible = ref(false)
const renewLoading = ref(false)
const renewMonths = ref(1)
const renewId = ref(0)
function handleRenew(row: any) { renewId.value = row.id; renewMonths.value = 1; renewVisible.value = true }
async function submitRenew() { renewLoading.value = true; try { await renewServer(renewId.value, { months: renewMonths.value }); ElMessage.success('续租成功'); renewVisible.value = false; fetchData() } finally { renewLoading.value = false } }
onMounted(fetchData)
</script>
