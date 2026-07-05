<template>
  <div class="login-container">
    <div class="login-card">
      <h2 class="login-title">星账系统</h2>
      <p class="login-subtitle">StarLedger · 财税管理平台</p>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="0" @submit.prevent="handleLogin">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="用户名" prefix-icon="User" size="large" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码" prefix-icon="Lock" size="large" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" :loading="loading" style="width: 100%" @click="handleLogin">
            登 录
          </el-button>
        </el-form-item>
      </el-form>
      <div class="login-footer">
        <el-button type="primary" link @click="showRegister = true">注册新租户</el-button>
      </div>
    </div>

    <!-- 注册弹窗 -->
    <el-dialog v-model="showRegister" title="注册新租户" width="520px" destroy-on-close>
      <el-form ref="regFormRef" :model="regForm" :rules="regRules" label-width="90px">
        <el-form-item label="租户类型" prop="tenant_type">
          <el-radio-group v-model="regForm.tenant_type" size="large">
            <el-radio-button value="personal">
              <el-icon><User /></el-icon> 个人
            </el-radio-button>
            <el-radio-button value="enterprise">
              <el-icon><OfficeBuilding /></el-icon> 企业
            </el-radio-button>
            <el-radio-button value="team">
              <el-icon><UserFilled /></el-icon> 团队
            </el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <div class="type-description">
            <template v-if="regForm.tenant_type === 'personal'">
              <el-tag type="info">个人版</el-tag>
              <span>适合个人用户，提供简单账单和财税统计功能</span>
            </template>
            <template v-else-if="regForm.tenant_type === 'enterprise'">
              <el-tag type="warning">企业版</el-tag>
              <span>适合企业用户，提供合同管理、多用户管理等功能</span>
            </template>
            <template v-else>
              <el-tag type="success">团队版</el-tag>
              <span>适合团队协作，提供任务分配、协作功能</span>
            </template>
          </div>
        </el-form-item>
        <el-form-item label="租户名称" prop="tenant_name">
          <el-input v-model="regForm.tenant_name" :placeholder="tenantNamePlaceholder" />
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="regForm.username" placeholder="登录用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="regForm.password" type="password" placeholder="至少6位" show-password />
        </el-form-item>
        <el-form-item label="真实姓名" prop="real_name">
          <el-input v-model="regForm.real_name" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="regForm.email" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="regForm.phone" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRegister = false">取消</el-button>
        <el-button type="primary" :loading="regLoading" @click="handleRegister">注册</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance } from 'element-plus'
import { User, OfficeBuilding, UserFilled } from '@element-plus/icons-vue'
import { useUserStore } from '../../stores/user'
import { login, register } from '../../api/auth'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const form = reactive({ username: '', password: '' })
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    const res: any = await login(form)
    userStore.setLogin(res.data)
    ElMessage.success('登录成功')
    router.push('/')
  } catch {
    // error handled by interceptor
  } finally {
    loading.value = false
  }
}

// 注册
const showRegister = ref(false)
const regFormRef = ref<FormInstance>()
const regLoading = ref(false)
const regForm = reactive({
  tenant_type: 'personal',
  tenant_name: '', username: '', password: '', real_name: '', email: '', phone: '',
})
const regRules = {
  tenant_type: [{ required: true, message: '请选择租户类型', trigger: 'change' }],
  tenant_name: [{ required: true, message: '请输入租户名称', trigger: 'blur' }],
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, min: 6, message: '密码至少6位', trigger: 'blur' }],
  real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
}

// 根据租户类型显示不同的占位符
const tenantNamePlaceholder = computed(() => {
  switch (regForm.tenant_type) {
    case 'personal': return '个人名称或昵称'
    case 'enterprise': return '企业名称'
    case 'team': return '团队名称'
    default: return '租户名称'
  }
})

async function handleRegister() {
  const valid = await regFormRef.value?.validate().catch(() => false)
  if (!valid) return
  regLoading.value = true
  try {
    await register(regForm)
    ElMessage.success('注册成功，请登录')
    showRegister.value = false
    form.username = regForm.username
  } catch {
    // error handled by interceptor
  } finally {
    regLoading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-card {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}
.login-title {
  text-align: center;
  font-size: 28px;
  color: #303133;
  margin-bottom: 4px;
}
.login-subtitle {
  text-align: center;
  color: #909399;
  margin-bottom: 30px;
  font-size: 14px;
}
.login-footer {
  text-align: center;
  margin-top: 10px;
}
.type-description {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 13px;
  color: #606266;
}
</style>
