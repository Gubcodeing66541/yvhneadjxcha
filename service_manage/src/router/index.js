import store from "@/store";
import { createRouter, createWebHashHistory } from 'vue-router'

import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import ResetPassword from '../views/ResetPassword.vue'
import LargeDataScreen from '../views/LargeDataScreen.vue'

import home from '../views/home/index.vue'
import record from '../views/home/record.vue'
import certification from '../views/home/certification.vue'

import account from '../views/account/index.vue'
import recharge from '../views/account/recharge.vue'
import rechargeRecord from '../views/account/rechargeRecord.vue'

import quickReply from '../views/quickReply/index.vue'

import robotManage from '../views/robotManagement/index.vue'
import robot from '../views/robotManagement/robot.vue'
import knowledgeBase from '../views/robotManagement/knowledgeBase.vue'

import service from '../views/serviceManage/index.vue'

import userManage from '../views/userManage/index.vue'
import serviceRecord from '../views/userManage/serviceRecord.vue'
import visitorRecords from '../views/userManage/visitorRecords.vue'
import blacklist from '../views/userManage/blacklist.vue'

const router = createRouter({
    history: createWebHashHistory(), // hash模式：createWebHashHistory，history模式：createWebHistory
    routes: [
        {
            path: '/',
            redirect: '/index/home/record'
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
            path: '/LargeDataScreen',
            name: 'LargeDataScreen',
            component: LargeDataScreen,
            meta: { title: '数据大屏' }
        }, {
            path: '/index',
            name: 'index',
            component: Home,
            meta: { title: '框架页' },
            children: [{
                path: '/index/home',
                name: 'home',
                redirect: '/index/home/record',
                component: home,
                meta: { title: '首页' },
                children: [{
                    path: 'record',
                    name: 'record',
                    component: record,
                    meta: { title: "首页数据" }
                }, {
                    path: 'certification',
                    name: 'certification',
                    component: certification,
                        meta: { title: "实名认证" }
                    }]
            }, {
                    path: '/index/account',
                    name: 'account',
                    redirect: '/index/account/recharge',
                    component: account,
                    meta: { title: '账户' },
                    children: [{
                        path: 'recharge',
                        name: 'recharge',
                        component: recharge,
                        meta: { title: "充值" }
                    }, {
                            path: 'rechargeRecord',
                            name: 'rechargeRecord',
                            component: rechargeRecord,
                            meta: { title: "充值记录" }
                        }]
                }, {
                    path: '/index/quickReply',
                    name: 'quickReply',
                    component: quickReply,
                    meta: { title: '快捷回复' }
                }, {
                    path: '/index/robotManage',
                    name: 'robotManage',
                    redirect: '/index/robotManage/robot',
                    component: robotManage,
                    meta: { title: '机器人管理' },
                    children: [{
                        path: 'robot',
                        name: 'robot',
                        component: robot,
                        meta: { title: "机器人" }
                    }, {
                        path: 'knowledgeBase',
                        name: 'knowledgeBase',
                        component: knowledgeBase,
                        meta: { title: "机器人知识库" }
                    }]
                }, {
                    path: '/index/service',
                    name: 'service',
                    component: service,
                    meta: { title: '客服管理' }
                }, {
                    path: '/index/userManage',
                    name: 'userManage',
                    redirect: '/index/userManage/serviceRecord',
                    component: userManage,
                    meta: { title: '用户管理' },
                    children: [{
                        path: 'serviceRecord',
                        name: 'serviceRecord',
                        component: serviceRecord,
                        meta: { title: "服务记录" }
                    }, {
                        path: 'visitorRecords',
                        name: 'visitorRecords',
                        component: visitorRecords,
                        meta: { title: "游客记录" }
                    }, {
                        path: 'blacklist',
                        name: 'blacklist',
                        component: blacklist,
                        meta: { title: "黑名单" }
                    }]
                }]
        }
    ]
})

router.beforeEach((to, from, next) => {
    if (store.state.config.SystemName && store.state.config.SystemName!='') document.title = store.state.config.SystemName + '管理后台-' + to.meta.title
    else document.title = '管理后台-' + to.meta.title
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