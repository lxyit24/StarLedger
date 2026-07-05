import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'

// Map route paths to module names for access control
const routeModuleMap: Record<string, string> = {
  '/app/servers': 'server_lease',
  '/app/bills': 'billing',
  '/app/contracts': 'contract',
  '/app/tasks': 'task',
  '/app/reports': 'report',
  '/app/invoices': 'invoice',
}

const routes = [
  {
    path: '/',
    name: 'Landing',
    component: () => import('../views/landing/index.vue'),
    meta: { public: true },
  },
  {
    path: '/terms',
    name: 'Terms',
    component: () => import('../views/landing/terms.vue'),
    meta: { public: true },
  },
  {
    path: '/privacy',
    name: 'Privacy',
    component: () => import('../views/landing/privacy.vue'),
    meta: { public: true },
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/index.vue'),
    meta: { public: true },
  },
  {
    path: '/app',
    component: () => import('../layouts/MainLayout.vue'),
    redirect: '/app/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/dashboard/index.vue'),
        meta: { title: '仪表盘' },
      },
      {
        path: 'servers',
        name: 'Servers',
        component: () => import('../views/server/index.vue'),
        meta: { title: '服务器租赁' },
      },
      {
        path: 'bills',
        name: 'Bills',
        component: () => import('../views/billing/index.vue'),
        meta: { title: '账单管理' },
      },
      {
        path: 'contracts',
        name: 'Contracts',
        component: () => import('../views/contract/index.vue'),
        meta: { title: '合同管理' },
      },
      {
        path: 'tasks',
        name: 'Tasks',
        component: () => import('../views/task/index.vue'),
        meta: { title: '任务协作' },
      },
      {
        path: 'reports',
        name: 'Reports',
        component: () => import('../views/report/index.vue'),
        meta: { title: '数据报表' },
      },
      {
        path: 'market',
        name: 'Market',
        component: () => import('../views/market/index.vue'),
        meta: { title: '模块市场' },
      },
      {
        path: 'invoices',
        name: 'Invoices',
        component: () => import('../views/invoice/index.vue'),
        meta: { title: '发票管理' },
      },
      {
        path: 'system/users',
        name: 'Users',
        component: () => import('../views/system/users.vue'),
        meta: { title: '用户管理' },
      },
      {
        path: 'system/roles',
        name: 'Roles',
        component: () => import('../views/system/roles.vue'),
        meta: { title: '角色管理' },
      },
      {
        path: 'system/tenants',
        name: 'Tenants',
        component: () => import('../views/system/tenants.vue'),
        meta: { title: '租户管理' },
      },
      {
        path: 'system/audit',
        name: 'Audit',
        component: () => import('../views/system/audit.vue'),
        meta: { title: '审计日志' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore()
  if (!to.meta.public && !userStore.token) {
    next('/login')
    return
  }

  // Fetch modules if logged in and not loaded yet
  const appStore = useAppStore()
  if (userStore.token && !appStore.loaded) {
    await appStore.fetchModules()
  }

  // Check module access for protected routes
  const moduleName = routeModuleMap[to.path]
  if (moduleName && !appStore.isModuleEnabled(moduleName)) {
    next('/app/market')
    return
  }

  next()
})

export default router
