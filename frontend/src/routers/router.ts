import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import type { UserInfo } from '@/api/auth'

// 静态路由
export const staticRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/upload-record'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/register.vue'),
    meta: { title: '注册' }
  }
]

// 动态路由
export const asyncRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    children: [
      {
        path: 'upload-record',
        name: 'UploadRecordDashboard',
        component: () => import('@/views/upload-record/dashboard.vue'),
        meta: { title: '数据概览' }
      },
      {
        path: 'upload-record/list',
        name: 'UploadRecordList',
        component: () => import('@/views/upload-record/list.vue'),
        meta: { title: '上传记录' }
      },
      {
        path: 'field-config',
        name: 'FieldConfig',
        component: () => import('@/views/field-config/list.vue'),
        meta: { title: '字段配置' }
      },
      {
        path: 'projects',
        name: 'ProjectList',
        component: () => import('@/views/project/list.vue'),
        meta: { title: '项目管理' }
      },
      {
        path: 'projects/network',
        name: 'ProjectNetwork',
        component: () => import('@/views/project/network.vue'),
        meta: { title: '关系网络' }
      },
      {
        path: 'personnel',
        name: 'PersonnelList',
        component: () => import('@/views/personnel/list.vue'),
        meta: { title: '人员管理' }
      },
      {
        path: 'users',
        name: 'UserList',
        component: () => import('@/views/user/list.vue'),
        meta: { title: '用户管理' }
      },
      {
        path: 'roles',
        name: 'RoleList',
        component: () => import('@/views/role/list.vue'),
        meta: { title: '角色管理' }
      },
      {
        path: 'user-groups',
        name: 'UserGroupList',
        component: () => import('@/views/user-group/list.vue'),
        meta: { title: '用户组' }
      },
      {
        path: 'audit/operation-log',
        name: 'OperationLog',
        component: () => import('@/views/audit/operation-log.vue'),
        meta: { title: '操作日志' }
      },
      {
        path: 'audit/login-log',
        name: 'LoginLog',
        component: () => import('@/views/audit/login-log.vue'),
        meta: { title: '登录日志' }
      },
      {
        path: 'system/ad-settings',
        name: 'ADSettings',
        component: () => import('@/views/system/ad-settings.vue'),
        meta: { title: 'AD域配置' }
      },
      {
        path: 'system/security-settings',
        name: 'SecuritySettings',
        component: () => import('@/views/system/security-settings.vue'),
        meta: { title: '安全设置' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人设置' }
      },
      {
        path: 'sync',
        name: 'Sync',
        component: () => import('@/views/sync/list.vue'),
        meta: { title: '数据同步' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory('/'),
  routes: [...staticRoutes, ...asyncRoutes],
  scrollBehavior: () => ({ left: 0, top: 0 })
})

export default router
