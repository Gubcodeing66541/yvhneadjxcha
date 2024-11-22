/*
 * @详细参数: 封装socket方法
 */
import store from "@/store";
import { getTime,playMp3 } from "@/utils";
import axios from "@/utils/axios";

let socketUrl = 'ws://socket.51wqc.top' + `/api/websocket/conn?token=${token}`; // socket地址
let websocket = null; // websocket 实例
let heartTime = null; // 心跳定时器实例
let HeartTimeOut = 10000; // 心跳超时时间
let socketError = 0; // 错误次数
let isClose = false; //websocket 是否关闭连接

// 初始化socket
const initWebSocket = () => {
    websocket = new WebSocket(socketUrl);
    websocketonopen();  // socket 连接成功
    websocketonClose(); //断开连接
    websocketonmessage();   // socket 接收数据
};

// socket 连接成功
const websocketonopen = () => {    
    websocket.onopen = async(e) => {
        if (socketError == 0){
            let time = await getTime(new Date());
            await axios.post('/auth/send', { type: 'time', content: time + ' 扫码进入' }, false);
        }
        store.commit('setSocketError', socketError)
        resetHeart();  // socket 重置心跳
    }
};

// socket 断开链接
const websocketonClose = () => {
    websocket.onclose = function (e) {
        if (!isClose) reconnect();
    };
};

// socket 接收数据
const websocketonmessage = () => {
    websocket.onmessage = ((e) => {
        let socketMsg = JSON.parse(e.data);
        store.commit('setOnMessage', socketMsg)
        if (socketMsg.type == 'message' && socketMsg.content.send_role == 'service') playMp3() //消息音频播放
    })
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
    if (socketError <= 200) {
        setTimeout(() => {
            clearInterval(heartTime);
            initWebSocket();
            socketError += 1;
            console.log("socket重连", socketError);
        }, 1000);
    } else {
        // console.log("重试次数已用完的逻辑", socketError);
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

// socket链接状态
const readyState = () => {
    return websocket.readyState
};

export {
    initWebSocket, //websocket连接
    websocketonmessage, //接收数据
    sendMsg, //发送数据
    websocketClose, //关闭链接
    readyState, //socket链接状态
};