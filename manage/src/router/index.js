import store from "@/store";
import { createRouter, createWebHashHistory } from 'vue-router'

import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import ResetPassword from '../views/ResetPassword.vue'

import home from '../views/home/index.vue'
import accountManage from '../views/accountManage/index.vue'
import domainManage from '../views/domainManage/index.vue'
import independentDomainName from '../views/domainManage/independentDomainName.vue'
import transfer from '../views/domainManage/transfer.vue'
import LandedDomainName from '../views/domainManage/LandedDomainName.vue'
import authorizationManage from '../views/authorizationManage/index.vue'
import userManage from '../views/userManage/index.vue'
import orderManage from '../views/orderManage/index.vue'
import systemSettings from '../views/systemSettings/index.vue'
import service from '../views/systemSettings/service.vue'
import sensitiveWords from '../views/systemSettings/sensitiveWords.vue'
import notice from '../views/systemSettings/notice.vue'

const router = createRouter({
    history: createWebHashHistory(), // hash模式：createWebHashHistory，history模式：createWebHistory
    routes: [
        {
            path: '/',
            redirect: '/index/home'
        }, {
            path: '/login',
            name: 'login',
            component: Login,
            meta: { title: '登录' }
        }, {
            path: '/ResetPassword',
            name: 'ResetPassword',
            component: ResetPassword,
            meta: { title: '重置密码' }
        }, {
            path: '/index',
            name: 'index',
            component: Home,
            meta: { title: '总后台' },
            children: [
                {
                    path: '/index/home',
                    name: 'home',
                    component: home,
                    meta: { title: '首页' },
                }, {
                    path: '/index/accountManage',
                    name: 'accountManage',
                    component: accountManage,
                    meta: { title: '账号管理' },
                }, {
                    path: '/index/domainManage',
                    name: 'domainManage',
                    component: domainManage,
                    meta: { title: '域名管理' },
                    children: [
                        {
                            path: '/index/domainManage/independentDomainName',
                            name: 'independentDomainName',
                            component: independentDomainName,
                            meta: { title: '独立域名' },
                        }, {
                            path: '/index/domainManage/transfer',
                            name: 'transfer',
                            component: transfer,
                            meta: { title: '中转域名' },
                        }, {
                            path: '/index/domainManage/LandedDomainName',
                            name: 'LandedDomainName',
                            component: LandedDomainName,
                            meta: { title: '落地域名' },
                        }
                    ]
                }, {
                    path: '/index/authorizationManage',
                    name: 'authorizationManage',
                    component: authorizationManage,
                    meta: { title: '授权管理' },
                }, {
                    path: '/index/userManage',
                    name: 'userManage',
                    component: userManage,
                    meta: { title: '用户管理' },
                }, {
                    path: '/index/orderManage',
                    name: 'orderManage',
                    component: orderManage,
                    meta: { title: '订单管理' },
                }, {
                    path: '/index/systemSettings',
                    name: 'systemSettings',
                    component: systemSettings,
                    meta: { title: '系统配置' },
                    children: [
                        {
                            path: '/index/systemSettings/service',
                            name: 'service',
                            component: service,
                            meta: { title: '配置' },
                        }, {
                            path: '/index/systemSettings/notice',
                            name: 'notice',
                            component: notice,
                            meta: { title: '公告' },
                        }, {
                            path: '/index/systemSettings/sensitiveWords',
                            name: 'sensitiveWords',
                            component: sensitiveWords,
                            meta: { title: '敏感词' },
                        }
                    ]
                },
            ]
        },
    ]
})

router.beforeEach((to, from, next) => {
    if (store.state.config.system_name && store.state.config.system_name != '') document.title = store.state.config.system_name + '总后台-' + to.meta.title
    else document.title = '总后台-' + to.meta.title
    if (!window.localStorage.getItem("token")) {
        if (to.name == "login") {
            next()
        } else {
            next({
                name: 'login'
            })
        }
    } else {
        next()
    }

})

export default router