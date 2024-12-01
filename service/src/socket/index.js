/*
 * @详细参数: 封装socket方法
 */
import { localRemove } from '@/utils'
import { message, notification } from 'ant-design-vue';
import store from "@/store";
let websocket = null; // websocket 实例
let heartTime = null; // 心跳定时器实例
let HeartTimeOut = 10000; // 心跳超时时间
let socketError = 0; // 错误次数
let isClose = false; //websocket 是否关闭连接
let mp3 = new Audio("http://"+import.meta.env.VITE_BASE_URL + '/common/audio/mesg.mp3');

// 初始化socket
const initWebSocket = () => {
    let socketUrl = 'ws://'+import.meta.env.VITE_BASE_URL + `/api/websocket/conn?token=${store.state.token}`; // socket地址
    websocket = new WebSocket(socketUrl);
    websocketonopen();  // socket 连接成功
    websocketclose(); //断开连接   
    websocketonmessage();   // socket 接收数据
};

// socket 连接成功
const websocketonopen = () => {
    websocket.onopen = (e) => {
        resetHeart();  // socket 重置心跳
    }
};

// socket 断开链接
const websocketclose = () => {
    websocket.onclose = function (e) {
        if (!isClose) reconnect();
    };
};

// socket 接收数据
const websocketonmessage = () => {
    websocket.onmessage = ((e) => {
        let socketMsg = JSON.parse(e.data);        
        store.commit('setOnMessage', socketMsg);
        switch (socketMsg.type) {
            case "message":
                if (store.state.messageReminder && socketMsg.content.send_role == 'user' && socketMsg.content.type != 'time') mp3.play();  // 播放
                if (socketMsg.content.send_role == 'user') msgTipc(socketMsg.content.content)
            break;
            case "ban":
                store.dispatch("asyncSetInfo");
                notification['warning']({
                    message: '域名拦截提醒',
                    description: '二维码已重置请重新保存二维码！',
                    duration:null,
                    placement:"bottomRight",
                    onClick: () => {
                        console.log('Notification Clicked!');
                    },
                });
            break;
            case "out_login": //踢出登录
                websocketClose();
                localRemove();
                message.error(socketMsg.content.message);
                window.location.reload();
            break;
        }
    })
};

//接收消息提示
const msgTipc = (content) => {
    if (window.Notification && window.Notification.permission == 'granted'){
        const n = new Notification('你有一条新消息！', {
            body: content,
            icon: store.state.info.head,
        });
    // setTimeout(() => {
    //     n.close();
    // }, 6000);
    }
};

// socket 发送数据
const sendMsg = (data) => {
    websocket.send(JSON.stringify(data));
};

// socket 重置心跳
const resetHeart = () => {
    socketError = 0;
    clearInterval(heartTime);
    sendSocketHeart();
};

// socket心跳发送
const sendSocketHeart = () => {
    heartTime = setInterval(() => {
        sendMsg({ type: "ping" });
    }, HeartTimeOut);
};

// socket重连
const reconnect = () => {
    if (socketError <= 100) {
        setTimeout(() => {
            clearInterval(heartTime);
            initWebSocket();
            socketError += 1;
            console.log("socket重连", socketError);
        }, 500);
    } else {
        console.log("重试次数已用完的逻辑", socketError);
        clearInterval(heartTime);
    }
};

// 关闭websocket
const websocketClose = () => {
    if (!websocket.readyState) return false
    isClose = true;
    websocket.close();
    clearInterval(heartTime);
    websocket = null;
    socketError = 0;
};

export {
    initWebSocket, //websocket连接
    websocketonmessage, //接收数据
    sendMsg, //发送数据
    websocketClose, //关闭链接
};