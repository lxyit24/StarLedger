<template>
  <div>
    <el-card>
      <template #header><el-row justify="space-between" align="middle"><span>用户管理</span><el-row gap="12"><el-input v-model="searchKey" placeholder="搜索用户名/姓名" clearable style="width: 200px" @input="fetchData" /><el-button type="primary" @click="openDialog()">新增用户</el-button></el-row></el-row></template>
      <el-table :data="tableData" stripe v-loading="loading" empty-text="暂无数据">
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="real_name" label="真实姓名" width="120" />
        <el-table-column prop="email" label="邮箱" min-width="160"><template #default="{ row }">{{ row.email || '--' }}</template></el-table-column>
        <el-table-column prop="phone" label="手机号" width="130"><template #default="{ row }">{{ row.phone || '--' }}</template></el-table-column>
        <el-table-column prop="status" label="状态" width="80"><template #default="{ row }"><el-tag :type="userStatusType(row.status)" size="small">{{ formatUserStatus(row.status) }}</el-tag></template></el-table-column>
        <el-table-column label="操作" width="180" fixed="right"><template #default="{ row }"><el-button size="small" @click="openDialog(row)">编辑</el-button><el-popconfirm title="确定删除？" @confirm="handleDelete(row.id)"><template #reference><el-button size="small" type="danger">删除</el-button></template></el-popconfirm></template></el-table-column>
      </el-table>
      <el-pagination style="margin-top: 16px; justify-content: flex-end" v-model:current-page="query.page" v-model:page-size="query.page_size" :total="total" layout="total, sizes, prev, pager, next" @change="fetchData" />
    </el-card>
    <el-dialog v-model="dialogVisible" :title="editId ? '编辑用户' : '新增用户'" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="80px">
        <el-form-item label="用户名" prop="username"><el-input v-model="form.username" :disabled="!!editId" /></el-form-item>
        <el-form-item v-if="!editId" label="密码" prop="password"><el-input v-model="form.password" type="password" show-password /></el-form-item>
        <el-form-item label="真实姓名" prop="real_name"><el-input v-model="form.real_name" /></el-form-item>
        <el-form-item label="邮箱"><el-input v-model="form.email" /></el-form-item>
        <el-form-item label="手机号"><el-input v-model="form.phone" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button></template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { listUsers, createUser, updateUser, deleteUser } from '../../api/system'
import { formatUserStatus, userStatusType } from '../../utils/format'
const loading = ref(false)
const searchKey = ref('')
const tableData = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20 })
async function fetchData() { loading.value = true; try { const res: any = await listUsers(query); tableData.value = res.data?.items || res.data || []; total.value = res.data?.total || tableData.value.length } finally { loading.value = false } }
const dialogVisible = ref(false)
const editId = ref(0)
const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const form = reactive({ username: '', password: '', real_name: '', email: '', phone: '' })
const formRules = { username: [{ required: true, message: '请输入用户名', trigger: 'blur' }], password: [{ required: true, min: 6, message: '密码至少6位', trigger: 'blur' }], real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }] }
function openDialog(row?: any) { if (row) { editId.value = row.id; Object.assign(form, { username: row.username, password: '', real_name: row.real_name, email: row.email || '', phone: row.phone || '' }) } else { editId.value = 0; Object.assign(form, { username: '', password: '', real_name: '', email: '', phone: '' }) } dialogVisible.value = true }
async function handleSubmit() { const valid = await formRef.value?.validate().catch(() => false); if (!valid) return; submitLoading.value = true; try { if (editId.value) { await updateUser(editId.value, { real_name: form.real_name, email: form.email, phone: form.phone }); ElMessage.success('更新成功') } else { await createUser(form); ElMessage.success('创建成功') } dialogVisible.value = false; fetchData() } finally { submitLoading.value = false } }
async function handleDelete(id: number) { await deleteUser(id); ElMessage.success('删除成功'); fetchData() }
onMounted(fetchData)
</script>
