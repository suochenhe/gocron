import Vue from 'vue'
import Router from 'vue-router'
import store from '../store/index'
import NotFound from '../components/common/notFound'

import TaskList from '../pages/task/list'
import TaskEdit from '../pages/task/edit'
import TaskLog from '../pages/taskLog/list'

import HostList from '../pages/host/list'
import HostEdit from '../pages/host/edit'

import UserList from '../pages/user/list'
import UserEdit from '../pages/user/edit'
import UserLogin from '../pages/user/login'
import UserEditPassword from '../pages/user/editPassword'
import UserEditMyPassword from '../pages/user/editMyPassword'

import NotificationDing from '../pages/system/notification/ding'
import NotificationEmail from '../pages/system/notification/email'
import NotificationSlack from '../pages/system/notification/slack'
import NotificationWebhook from '../pages/system/notification/webhook'

import Install from '../pages/install/index'
import LoginLog from '../pages/system/loginLog'

Vue.use(Router)
// authRoles -1未登录 0游客 1管理员 2开发者
const router = new Router({
  routes: [
    {
      path: '*',
      component: NotFound,
      meta: {
        authRoles: [-1]
      }
    },
    {
      path: '/',
      redirect: '/task'
    },
    {
      path: '/install',
      name: 'install',
      component: Install,
      meta: {
        // -1 未登录
        authRoles: [-1]
      }
    },
    {
      path: '/task',
      name: 'task-list',
      component: TaskList,
      meta: {
        authRoles: [-1, 0, 1, 2],
        keepAlive: true
      }
    },
    {
      path: '/task/create',
      name: 'task-create',
      component: TaskEdit,
      meta: {
        authRoles: [1, 2]
      }
    },
    {
      path: '/task/edit/:id',
      name: 'task-edit',
      component: TaskEdit,
      meta: {
        authRoles: [1, 2]
      }
    },
    {
      path: '/task/log',
      name: 'task-log',
      component: TaskLog,
      meta: {
        authRoles: [0, 1, 2]
      }
    },
    {
      path: '/host',
      name: 'host-list',
      component: HostList,
      meta: {
        authRoles: [0, 1, 2]
      }
    },
    {
      path: '/host/create',
      name: 'host-create',
      component: HostEdit,
      meta: {
      }
    },
    {
      path: '/host/edit/:id',
      name: 'host-edit',
      component: HostEdit,
      meta: {
      }
    },
    {
      path: '/user',
      name: 'user-list',
      component: UserList,
      meta: {
      }
    },
    {
      path: '/user/create',
      name: 'user-create',
      component: UserEdit,
      meta: {
      }
    },
    {
      path: '/user/edit/:id',
      name: 'user-edit',
      component: UserEdit,
      meta: {
      }
    },
    {
      path: '/user/login',
      name: 'user-login',
      component: UserLogin,
      meta: {
        authRoles: [-1]
      }
    },
    {
      path: '/user/edit-password/:id',
      name: 'user-edit-password',
      component: UserEditPassword,
      meta: {
      }
    },
    {
      path: '/user/edit-my-password',
      name: 'user-edit-my-password',
      component: UserEditMyPassword,
      meta: {
        authRoles: [0, 1, 2]
      }
    },
    {
      path: '/system',
      redirect: '/system/notification/ding',
      meta: {}
    },
    {
      path: '/system/notification/ding',
      name: 'system-notification-ding',
      component: NotificationDing,
      meta: {}
    },
    {
      path: '/system/notification/email',
      name: 'system-notification-email',
      component: NotificationEmail,
      meta: {}
    },
    {
      path: '/system/notification/slack',
      name: 'system-notification-slack',
      component: NotificationSlack,
      meta: {}
    },
    {
      path: '/system/notification/webhook',
      name: 'system-notification-webhook',
      component: NotificationWebhook,
      meta: {}
    },
    {
      path: '/system/login-log',
      name: 'login-log',
      component: LoginLog,
      meta: {}
    }
  ]
})

router.beforeEach((to, from, next) => {
  let authRoles = to.meta.authRoles || []
  console.log(to.name + ' path:' + to.path + ' ' + authRoles)
  console.log('role:' + store.getters.user.role)
  if (authRoles.includes(-1)) {
    next()
    return
  }

  if (store.getters.user.token) {
    let hasAuth = authRoles.includes(store.getters.user.role)
    if (store.getters.user.isAdmin || hasAuth) {
      next()
      return
    }
    if (!hasAuth) {
      next(
        {
          path: '/404.html'
        }
      )
      return
    }
  }

  next({
    path: '/user/login',
    query: {redirect: to.fullPath}
  })
})

export default router
