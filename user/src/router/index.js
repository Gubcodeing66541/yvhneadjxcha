import { createRouter, createWebHashHistory } from 'vue-router'

import Home from '../views/Home.vue'
import Report from '../views/Report.vue'

const router = createRouter({
    history: createWebHashHistory(), // hash模式：createWebHashHistory，history模式：createWebHistory
    routes: [
        {
            path: '/',
            redirect: '/index'
        }, {
            path: '/index',
            name: 'index',
            component: Home,
            meta: { title: '在线客服' },
        },
        {
            path: '/report',
            name: 'report',
            component: Report,
            meta: { title: '举报成功' },
        },
    ]
})

export default router