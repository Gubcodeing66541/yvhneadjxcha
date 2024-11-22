import store from "@/store";
import { createRouter, createWebHashHistory } from 'vue-router'

import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Statistics from '../views/Statistics.vue'

const router = createRouter({
    history: createWebHashHistory(), // hash模式：createWebHashHistory，history模式：createWebHistory
    routes: [
        {
            path: '/',
            redirect: '/index'
        }, {
            path: '/login',
            name: 'login',
            component: Login,
            meta: { title: '云客服-登录' }
        }, {
            path: '/statistics',
            name: 'statistics',
            component: Statistics,
            meta: { title: '云客服-统计' }
        }, {
            path: '/index',
            name: 'index',
            component: Home,
            meta: { title: '云客服' },
        },
    ]
})

router.beforeEach((to, from, next) => {
    if (!window.localStorage.getItem("token")) {
        if (to.name == "login" || to.name == "statistics") {
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