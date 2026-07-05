<template>
  <div>
    <el-card>
      <template #header><el-row justify="space-between" align="middle"><span>角色管理</span><el-button type="primary" @click="openDialog()">新增角色</el-button></el-row></template>
      <el-table :data="tableData" stripe v-loading="loading" empty-text="暂无数据">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="permissions" label="权限" min-width="200"><template #default="{ row }"><el-tag v-for="p in (row.permissions || [])" :key="p" size="small" style="margin: 2px">{{ p }}</el-tag></template></el-table-column>
        <el-table-column label="操作" width="180" fixed="right"><template #default="{ row }"><el-button size="small" @click="openDialog(row)">编辑</el-button><el-popconfirm title="确定删除？" @confirm="handleDelete(row.id)"><template #reference><el-button size="small" type="danger">删除</el-button></template></el-popconfirm></template></el-table-column>
      </el-table>
      <el-pagination style="margin-top: 16px; justify-content: flex-end" v-model:current-page="query.page" v-model:page-size="query.page_size" :total="total" layout="total, sizes, prev, pager, next" @change="fetchData" />
    </el-card>
    <el-dialog v-model="dialogVisible" :title="editId ? '编辑角色' : '新增角色'" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="80px">
        <el-form-item label="角色名" prop="name"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" /></el-form-item>
        <el-form-item label="权限"><el-checkbox-group v-model="form.permissions"><el-checkbox v-for="p in allPermissions" :key="p" :label="p" :value="p" /></el-checkbox-group></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button></template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { listRoles, createRole, updateRole, deleteRole } from '../../api/system'
const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20 })
const allPermissions = ['server:read', 'server:write', 'server:delete', 'bill:read', 'bill:write', 'bill:delete', 'user:read', 'user:write', 'user:delete', 'tenant:read', 'tenant:write', 'tenant:delete']
async function fetchData() { loading.value = true; try { const res: any = await listRoles(query); tableData.value = res.data?.items || res.data || []; total.value = res.data?.total || tableData.value.length } finally { loading.value = false } }
const dialogVisible = ref(false)
const editId = ref(0)
const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const form = reactive({ name: '', description: '', permissions: [] as string[] })
const formRules = { name: [{ required: true, message: '请输入角色名', trigger: 'blur' }] }
function openDialog(row?: any) { if (row) { editId.value = row.id; Object.assign(form, { name: row.name, description: row.description || '', permissions: row.permissions || [] }) } else { editId.value = 0; Object.assign(form, { name: '', description: '', permissions: [] }) } dialogVisible.value = true }
async function handleSubmit() { const valid = await formRef.value?.validate().catch(() => false); if (!valid) return; submitLoading.value = true; try { if (editId.value) { await updateRole(editId.value, { name: form.name, description: form.description, permissions: form.permissions }); ElMessage.success('更新成功') } else { await createRole(form); ElMessage.success('创建成功') } dialogVisible.value = false; fetchData() } finally { submitLoading.value = false } }
async function handleDelete(id: number) { await deleteRole(id); ElMessage.success('删除成功'); fetchData() }
onMounted(fetchData)
</script>
