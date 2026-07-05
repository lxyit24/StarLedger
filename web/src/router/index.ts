import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useAppStore } from '../stores/app'

// Map route paths to module names for access control
const routeModuleMap: Record<string, string> = {
  '/servers': 'server_lease',
  '/bills': 'billing',
}

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/index.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    component: () => import('../layouts/MainLayout.vue'),
    redirect: '/dashboard',
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
        path: 'market',
        name: 'Market',
        component: () => import('../views/market/index.vue'),
        meta: { title: '模块市场' },
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
    next('/market')
    return
  }

  next()
})

export default router
