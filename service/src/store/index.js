import { createStore } from 'vuex';
import axios from "@/utils/axios";
import { localSet, localGet } from "@/utils";

const store = createStore({
    state: {
        token: localGet('token') || '',
        onMessage: {}, //socket消息更新
        info: localGet('info') || {}, //账号信息
        config: localGet('config') || {}, //后台配置信息
        messageReminder: localGet('messageReminder') == null ? true : localGet('messageReminder'), //消息提醒
    },
    // 获取状态(变量的值)
    getters: {
        getToken: state => state.token,
        getOnMessage: state => state.onMessage,
        getInfo: state => state.info,
        getConfig: state => state.config,
        getMessageReminder: state => state.messageReminder,
    },
    //修改state里面的变量(修改变量的值)
    mutations: {
        setToken(state, newValue) {
            state.token = newValue
        },
        setOnMessage(state, newValue) {
            state.onMessage = newValue
        },
        setInfo(state, newValue) {
            state.info = newValue
        },
        setConfig(state, newValue) {
            state.config = newValue
        },
        setMessageReminder(state, newValue) {
            state.messageReminder = newValue
        },
    },
    // Action 触发 mutation 函数,从而修改状态 通过dispatch()调用
    actions: {
        async asyncSetInfo(context) {
            const res = await axios.post("/info", '', false);
            res.data.service.web = res.data.service.web;
            localSet('info', res.data.service);
            context.commit("setInfo", res.data.service);
        }
    },
    // Module 当状态很多时,把状态分开来管理
    modules: {

    }
})

export default store