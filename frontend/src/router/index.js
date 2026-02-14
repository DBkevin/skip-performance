import { createRouter, createWebHistory } from 'vue-router'
import { ElMessage } from 'element-plus'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('../views/Layout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: '首页', icon: 'HomeFilled' }
      },
      {
        path: 'customers',
        name: 'Customers',
        component: () => import('../views/Customers.vue'),
        meta: { title: '顾客管理', icon: 'UserFilled' }
      },
      {
        path: 'employees',
        name: 'Employees',
        component: () => import('../views/Employees.vue'),
        meta: { title: '员工管理', icon: 'Avatar', admin: true }
      },
      {
        path: 'projects',
        name: 'Projects',
        component: () => import('../views/Projects.vue'),
        meta: { title: '项目管理', icon: 'FirstAidKit', admin: true }
      },
      {
        path: 'visits',
        name: 'Visits',
        component: () => import('../views/Visits.vue'),
        meta: { title: '就诊管理', icon: 'DocumentChecked' }
      },
      {
        path: 'revisit-records',
        name: 'RevisitRecords',
        component: () => import('../views/RevisitRecords.vue'),
        meta: { title: '回访记录', icon: 'PhoneFilled' }
      },
      {
        path: 'reports',
        name: 'Reports',
        component: () => import('../views/Reports.vue'),
        meta: { title: '业绩报表', icon: 'TrendCharts' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const userRole = localStorage.getItem('userRole')

  // 公开页面直接放行
  if (to.meta.public) {
    next()
    return
  }

  // 未登录跳转到登录页
  if (!token) {
    ElMessage.warning('请先登录')
    next('/login')
    return
  }

  // 检查是否需要管理员权限
  if (to.meta.admin && userRole !== '管理员') {
    ElMessage.error('需要管理员权限')
    next('/dashboard')
    return
  }

  next()
})

export default router
