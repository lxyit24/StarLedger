<template>
  <div>
    <el-card>
      <template #header><el-row justify="space-between" align="middle"><span>租户管理</span><el-button type="primary" @click="openDialog()">新增租户</el-button></el-row></template>
      <el-table :data="tableData" stripe v-loading="loading" empty-text="暂无数据">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="租户名称" min-width="150" />
        <el-table-column prop="contact" label="联系人" width="120" />
        <el-table-column prop="phone" label="电话" width="140" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column prop="status" label="状态" width="80"><template #default="{ row }"><el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">{{ row.status }}</el-tag></template></el-table-column>
        <el-table-column label="操作" width="180" fixed="right"><template #default="{ row }"><el-button size="small" @click="openDialog(row)">编辑</el-button><el-popconfirm title="确定删除？" @confirm="handleDelete(row.id)"><template #reference><el-button size="small" type="danger">删除</el-button></template></el-popconfirm></template></el-table-column>
      </el-table>
      <el-pagination style="margin-top: 16px; justify-content: flex-end" v-model:current-page="query.page" v-model:page-size="query.page_size" :total="total" layout="total, sizes, prev, pager, next" @change="fetchData" />
    </el-card>
    <el-dialog v-model="dialogVisible" :title="editId ? '编辑租户' : '新增租户'" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="80px">
        <el-form-item label="名称" prop="name"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="联系人" prop="contact"><el-input v-model="form.contact" /></el-form-item>
        <el-form-item label="电话"><el-input v-model="form.phone" /></el-form-item>
        <el-form-item label="邮箱"><el-input v-model="form.email" /></el-form-item>
        <el-form-item v-if="editId" label="状态"><el-select v-model="form.status" style="width: 100%"><el-option label="正常" value="active" /><el-option label="停用" value="suspended" /></el-select></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button></template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { listTenants, createTenant, updateTenant, deleteTenant } from '../../api/system'
const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20 })
async function fetchData() { loading.value = true; try { const res: any = await listTenants(query); tableData.value = res.data?.items || res.data || []; total.value = res.data?.total || tableData.value.length } finally { loading.value = false } }
const dialogVisible = ref(false)
const editId = ref(0)
const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const form = reactive({ name: '', contact: '', phone: '', email: '', status: 'active' })
const formRules = { name: [{ required: true, message: '请输入名称', trigger: 'blur' }], contact: [{ required: true, message: '请输入联系人', trigger: 'blur' }] }
function openDialog(row?: any) { if (row) { editId.value = row.id; Object.assign(form, { name: row.name, contact: row.contact, phone: row.phone || '', email: row.email || '', status: row.status }) } else { editId.value = 0; Object.assign(form, { name: '', contact: '', phone: '', email: '', status: 'active' }) } dialogVisible.value = true }
async function handleSubmit() { const valid = await formRef.value?.validate().catch(() => false); if (!valid) return; submitLoading.value = true; try { if (editId.value) { await updateTenant(editId.value, { name: form.name, contact: form.contact, phone: form.phone, email: form.email, status: form.status }); ElMessage.success('更新成功') } else { await createTenant(form); ElMessage.success('创建成功') } dialogVisible.value = false; fetchData() } finally { submitLoading.value = false } }
async function handleDelete(id: number) { await deleteTenant(id); ElMessage.success('删除成功'); fetchData() }
onMounted(fetchData)
</script>
