<template>
  <el-container class="main-layout">
    <el-aside :width="isCollapse ? '64px' : '220px'" class="sidebar">
      <div class="logo">
        <h2 v-if="!isCollapse">星账系统</h2>
        <h2 v-else>星</h2>
      </div>
      <el-menu
        :default-active="$route.path"
        :collapse="isCollapse"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
      >
        <el-menu-item index="/app/dashboard">
          <el-icon><Odometer /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item v-if="appStore.isModuleEnabled('server_lease')" index="/app/servers">
          <el-icon><Monitor /></el-icon>
          <span>服务器租赁</span>
        </el-menu-item>
        <el-menu-item v-if="appStore.isModuleEnabled('billing')" index="/app/bills">
          <el-icon><Document /></el-icon>
          <span>账单管理</span>
        </el-menu-item>
        <el-menu-item v-if="appStore.isModuleEnabled('contract')" index="/app/contracts">
          <el-icon><Tickets /></el-icon>
          <span>合同管理</span>
        </el-menu-item>
        <el-menu-item v-if="appStore.isModuleEnabled('task')" index="/app/tasks">
          <el-icon><List /></el-icon>
          <span>任务协作</span>
        </el-menu-item>
        <el-menu-item v-if="appStore.isModuleEnabled('report')" index="/app/reports">
          <el-icon><DataAnalysis /></el-icon>
          <span>数据报表</span>
        </el-menu-item>
        <el-menu-item index="/app/market">
          <el-icon><ShoppingCart /></el-icon>
          <span>模块市场</span>
        </el-menu-item>
        <el-sub-menu index="system">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/app/system/users">用户管理</el-menu-item>
          <el-menu-item index="/app/system/roles">角色管理</el-menu-item>
          <el-menu-item index="/app/system/tenants">租户管理</el-menu-item>
          <el-menu-item index="/app/system/audit">审计日志</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <el-icon class="collapse-btn" @click="isCollapse = !isCollapse">
          <Fold v-if="!isCollapse" />
          <Expand v-else />
        </el-icon>
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/app/dashboard' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item v-if="$route.meta.title">{{ $route.meta.title }}</el-breadcrumb-item>
        </el-breadcrumb>
        <div class="header-right">
          <span>{{ userStore.username }}</span>
          <el-button type="danger" text @click="handleLogout">退出</el-button>
        </div>
      </el-header>
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'
import { Odometer, Monitor, Document, Setting, Fold, Expand, ShoppingCart, Tickets, List, DataAnalysis } from '@element-plus/icons-vue'

const isCollapse = ref(false)
const userStore = useUserStore()
const appStore = useAppStore()
const router = useRouter()

onMounted(() => {
  if (userStore.token) {
    appStore.fetchModules()
  }
})

function handleLogout() {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
}
.sidebar {
  background-color: #304156;
  transition: width 0.3s;
  overflow: hidden;
}
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}
.logo h2 {
  margin: 0;
  font-size: 18px;
}
.header {
  display: flex;
  align-items: center;
  border-bottom: 1px solid #eee;
  padding: 0 20px;
}
.collapse-btn {
  font-size: 20px;
  cursor: pointer;
  margin-right: 16px;
}
.header-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 12px;
}
.main-content {
  background-color: #f5f7fa;
}
</style>
