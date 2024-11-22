import {
    createStore
} from 'vuex'

const store = createStore({
    state: {
        onMessage: {}, //socket消息更新
        socketError:0, //socket重连次数
    },
    // 获取状态(变量的值)
    getters: {
        getOnMessage: state => state.onMessage,
        getSocketError: state => state.socketError,
    },
    //修改state里面的变量(修改变量的值)
    mutations: {
        setOnMessage(state, newValue) {
            state.onMessage = newValue
        },
        setSocketError(state, newValue) {
            state.socketError = newValue
        },
    },
    // Action 触发 mutation 函数,从而修改状态
    actions: {

    },
    // Module 当状态很多时,把状态分开来管理
    modules: {

    }
})

export default store