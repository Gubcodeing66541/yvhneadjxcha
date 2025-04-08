<template>
    <div ref="warp" class="warp hidden">
        <div id="chatBox" ref="chatBox" class="chat-content scrollbar" :style="{ 'height': `calc(100% - ${inputHeight}px)` }" @scroll="divScroll">
            <div class="loadding tc ptb6" v-if="loading"><van-loading type="spinner" color="#7556fd" size="20" /></div>
            <notice :notice="notice" v-if="notice.is_show"></notice>
            <div class="item" :class="item.send_role == 'user' ? 'item-right' : ''" v-for="(item, index) in list" :key="index">
                <van-image :show-loading="false" width="44" height="44" radius="2" fit="cover"
                    :src="item.send_role == 'service' || item.send_role == 'hello' ? info.service.head : info.bot_head" class="avatar"
                    :class="item.send_role == 'user' ? 'ml4 fr' : 'mr4 fl'" v-show="item.send_role != 'user'" />
                <div class="box">
                    <p class="t-time lh12 f12">{{ item.create_time }}</p>
                    <div class="message-row">
                        <div class="message-box relative">
                            <!-- 文本信息 -->
                            <div class="nocopy ptb12 plr14 ico-arrow" v-if="item.type == 'text'" @touchstart="copMsg(item.content)" @touchmove="gtouchmove()" @touchend="showDeleteButton()">
                                <div style="white-space:pre-wrap" v-html="textReplaceEmoji(item.content)"></div>
                            </div>
                            <!-- 图片信息 -->
                            <div class="img-cont hidden br5" v-else-if="item.type == 'image'">
                                <img :src="item.content" height="200" @click="imagesZoom(item.content)" style="max-width: 100%;">
                            </div>
                            <!-- 视频信息 -->
                            <div class="video-cont" v-else-if="item.type == 'video'" @touchstart="copMsg(item.content)" @touchmove="gtouchmove()" @touchend="showDeleteButton()">
                                <video :src="item.content" :height="200" controls style="object-fit:fill;width: 100%;"></video>
                            </div>
                            <!-- 二维码 -->
                            <div class="code-cont pd5" v-else>
                                <div class="tc" style="height:200px"><vue-qr :logoSrc="info.service.head" :text="info.service.web" :backgroundColor="info.service.code_background" :colorLight="info.service.code_background" :colorDark="info.service.code_color" :logoMargin="4" :size="200" :margin="10"></vue-qr></div>
                                <p class="tc mtb5 f12">长按保存专属二维码,聊天记录不丢失</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <!-- 底部 -->
        <div ref="chatInput" class="bot-content">
            <bottomMenu :menu="menu" v-if="menu!=''"></bottomMenu>
            <div class="content phone-bottom">
                <div class="input-box plr12 ptb6 flex algin-center">
                    <van-field ref="textarea" v-model="messageValue" rows="1" autosize type="textarea" size="large" :border="false"
                        @focus="focusHland" @blur="blurHland" @update:modelValue="valueHland" />
                    <div class="flex align-bottom lh36">
                        <i class="iconfont icon-biaoqing f32 ml10" @click="abilHland('emojis')"></i>
                        <i class="iconfont icon-jia f32 ml14" @click="abilHland('upload')" v-show="!isBtn || messageValue == ''"></i>
                        <div class=" ml14" v-show="isBtn && messageValue != ''" @click="sendText">
                            <span class="plr12 ptb8 br4 btn-them nowrap">发送</span>
                        </div>
                    </div>
                </div>
                <div class="ability-box relative t-3s">
                    <div class="item" v-show="acitveType == 'emojis'">
                        <myEmo @emoHland="emoHland" class="pd12" v-if="acitveType == 'emojis'" />
                    </div>
                    <div class="item" v-show="acitveType == 'upload'">
                        <upload class="plr20 pt14 pb88" v-if="acitveType == 'upload'" />
                    </div>
                </div>
            </div>
        </div>
    </div>
    <report></report>
    <van-action-sheet v-model:show="copyShow" cancel-text="取消" close-on-click-action :actions="[{name:'复制'}]" @select="longPress" />
</template>
<script>
import { reactive, toRefs, ref, onMounted, nextTick, onUnmounted, computed, watch } from "vue";
import myEmo from "@/components/emo.vue";
import upload from "@/components/upload.vue";
import notice from "@/components/notice.vue";
import bottomMenu from "@/components/menu.vue";
import report from "@/components/report.vue";
import vueQr from 'vue-qr/src/packages/vue-qr.vue';
import useClipboard from "vue-clipboard3";
import { textReplaceEmoji } from "@/utils/emojis";
import { playMp3 } from "@/utils";
import { initWebSocket, websocketClose, readyState } from "@/socket";
import { useStore } from "vuex";
import axios from "@/utils/axios";
import { ImagePreview, Toast } from 'vant';
export default {
    name: "Home",
    components: { myEmo, upload, notice, report, bottomMenu, vueQr },
    setup() {
        const { toClipboard } = useClipboard();
        const store = useStore();
        const onMessage = computed(() => store.state.onMessage); //socket消息
        const socketError = computed(() => store.state.socketError); //socket重连次数
        const warp = ref(0);
        const chatBox = ref(0);
        const chatInput = ref(0);
        const textarea = ref(0);
        const state = reactive({
            info: {}, //基本信息
            list: [], //聊天记录
            notice:{}, //公告
            menu: [], //快捷菜单
            form: { page: 1, offset: 20 }, //分页
            pageCount: 0, //聊天记录总页数
            messageValue: "", //输入框值
            blurIndex: '',   //光标下标记录
            isBtn: false, //发送按钮展示隐藏
            inputHeight: 0, //底部高度
            scrollHeight: 0, //滚动条总高度
            acitveType: '', //底部功能展示
            loading: false,
            timeOutEvent:null, //长按定时
            copyShow:false, //复制面板
            copyCont:"", //复制内容
        });

        onMounted(() => {
            if (navigator.onLine) init();
            else Toast.fail('没有网络，请先连接网络后刷新尝试');
        })

        onUnmounted(() => {
            websocketClose(); //关闭websocket
            // document.getElementById('chatBox').removeEventListener('touchstart', touchstart, { passive: true }); //销毁页面监听
        })

        watch(() => onMessage, (v1) => { //监听消息
            if (v1.value.type == 'message') onmessage(v1.value);
            else if (v1.value.type == 'remove') remove(v1.value);
            else if (v1.value.type == 'clear') state.list = [];
        }, { deep: true })

        watch(() => state.messageValue, (v1) => { //监听消息
            if (v1 != '') state.isBtn = true
        }, { deep: true })

        watch(() => socketError, (v1) => { //socket重连次数
            if (v1.value >= 1){
                console.log("socket重连成功");
                tipcSocket();
            }
        }, { deep: true })

        //初始化
        const init = () => {
            getChatHeight(); //获取输入框高度
            login(); //获取用户信息
            initWebSocket(); //连接websocket
            document.getElementById('chatBox').addEventListener('touchstart', touchstart, { passive: true }); //监听手指触摸屏幕
        }
        //获取用户信息
        const login = async () => {
            const res = await axios.post("/auth/info", '', false);
            state.info['bot_head'] = res.data.bot_head;
            state.info['service'] = res.data.service;
            state.info['users'] = res.data.users;
            state.menu = res.data.menu;
            state.notice = res.data.notice;
            document.title = state.info['service']['name']; //标题
            getChatHeight(); //获取输入框高度
            getMessageList(); //获取聊天记录
        }

        //监听滚动
        const divScroll = (e) => {
            nextTick(() => {
                if (e.srcElement.scrollTop == 0 && state.pageCount > state.form.page) {
                    state.form.page += 1;
                    getMessageList(); //获取聊天记录
                    state.scrollHeight = e.srcElement.scrollHeight; //储存旧全文高度
                }
            });
        }

        //获取聊天记录
        const getMessageList = async () => {
            state.loading = true;
            const res = await axios.post('/message/list', state.form, false);
            state.list = state.form.page == 1 ? res.data.list : res.data.list.concat(state.list);
            state.pageCount = res.data.page;
            state.loading = false;
            nextTick(() => {
                if (state.form.page == 1) scrollBottom()
                else chatBox.value['scrollTop'] = chatBox.value['scrollHeight'] - state.scrollHeight; //重置滚动高度
            });
        }

        //发送文本消息
        const sendText = () => {
            msgPush('text', state.messageValue);
            state.messageValue = ''; //清空输入框内容
            if (state.acitveType=='') textarea.value.focus(); //清空后重新聚焦        
        }

        //消息统一发送
        const msgPush = async (type, content) => {
            const socketState = readyState();
            if (socketState){
                await axios.post('/auth/send', { type, content }, false);
                await playMp3(); //音频播放
                abilHland(); //底部功能隐藏
            }else{
                tipcSocket();
            }
        }

        //发送消息时检测用户是否连接websocket
        const tipcSocket = () => {
            Toast.loading({
                message: '连接超时正在重新接入...',
                forbidClick: true,
            });
            setTimeout(() => {
                window.location.reload();
            }, 1500);
        }

        // 预览图片
        const imagesZoom = (img) => {
            ImagePreview([img]);
        }

        //表情插入
        const emoHland = (val) => {
            state.messageValue = state.messageValue.slice(0, state.blurIndex) + val + state.messageValue.slice(state.blurIndex);
        }

        //焦点获取监听
        const focusHland = (e) => {
            state.isBtn = true;
            abilHland();
        }

        //输入框焦点失去
        const blurHland = (e) => {
            state.blurIndex = e.srcElement.selectionStart;
            setTimeout(() => {
                state.isBtn = false
            }, 200);
        }

        //输入框内容变化
        const valueHland = () => {
            state.isBtn = true;
            nextTick(()=>{
                getChatHeight();
            })
        }

        //手指触摸监听
        const touchstart = () => {
            if (state.acitveType != '') abilHland(); //底部功能隐藏
            textarea.value.blur();
        }        

        //底部功能显示隐藏
        const abilHland = (type = '') => {
            let scrollBom = true;
            if (state.acitveType != '') scrollBom = false;
            state.acitveType = type;
            nextTick(() => {
                getChatHeight(scrollBom);
            });
        }

        //获取聊天页面高度
        const getChatHeight = (scrollBom = true) => {
            nextTick(() => {
                state.inputHeight = chatInput.value['offsetHeight'];
                if (scrollBom) scrollBottom();
            });
        }

        //消息监听
        const onmessage = (data) => {
            if (data.content.type == 'time' || data.content.type == 'remind') return false
            state.list.push(data.content);
            scrollBottom(); //滚动到底部
        }

        //滚动到底部
        const scrollBottom = () => {
            nextTick(() => {
                chatBox.value['scrollTop'] = chatBox.value['scrollHeight']; //滚动到底
            });
        }

        //消息撤回
        const remove = (data) => {
            for (let i in state.list) {
                if (state.list[i]["id"] == data.content.id) {
                    state.list.splice(i, 1);
                }
            }            
        }

        //长按复制
        const copMsg = (content) => {
            state.timeOutEvent = setTimeout(() => {
                state.copyShow = true
                state.copyCont = content
            }, 500);
            return false;
        }

        //手释放，如果在500毫秒内就释放，则取消长按事件，此时可以执行onclick应该执行的事件
        const showDeleteButton = () => {
            clearTimeout(state.timeOutEvent); //清除定时器
            return false;
        }

        //如果手指有移动，则取消所有事件，此时说明用户只是要移动而不是长按
        const gtouchmove = () => {
            clearTimeout(state.timeOutEvent); //清除定时器
            state.timeOutEvent = 0;
        }

        //长按事件
        const longPress = async() => {
            try {
                if (state.copyCont.startsWith('http') && state.copyCont.match(/\.(mp4|webm|ogg)$/i)) {
                    // 如果是视频链接，则下载视频
                    const link = document.createElement('a');
                    link.href = state.copyCont;
                    link.download = 'video_' + new Date().getTime() + '.mp4';
                    document.body.appendChild(link);
                    link.click();
                    document.body.removeChild(link);
                    Toast("视频保存中...");
                } else {
                    // 如果是文本或普通链接，则复制
                    await toClipboard(state.copyCont);
                    Toast("复制成功");
                }
            } catch (e) {
                Toast("操作失败");
            }
        }

        return {
            ...toRefs(state),
            sendText, //发送文本消息
            emoHland, //表情
            warp, //整个页面ref
            chatBox, //聊天盒子ref
            chatInput,//聊天输入框盒子ref
            textarea, //输入框ref
            focusHland, //焦点获取监听
            blurHland, //输入框焦点失去
            divScroll, //监听滚动
            scrollBottom, //滚动到底部
            textReplaceEmoji, //文字表情处理
            imagesZoom, //图片预览
            abilHland, //显示底部功能
            valueHland, //输入框内容变化
            copMsg, //长按复制
            showDeleteButton, //手释放，如果在500毫秒内就释放，则取消长按事件，此时可以执行onclick应该执行的事件
            gtouchmove, //如果手指有移动，则取消所有事件，此时说明用户只是要移动而不是长按
            longPress, //长按显示复制面板
        };
    },
};
</script>

<style lang="less" scoped>
.warp {
    width: 100vw;
    height: 100vh;
    background: #ededed;
}
.loadding{
    position: absolute;
    width: 100%;
    top: 0;
    left: 0;
    z-index: 999;
}
.chat-content {
    .item {
        overflow: hidden;
        padding: 10px 12px;
        .box {
            float: left;
            max-width: 76%;
            .message-row {
                padding: 0 10px;
                margin-top: 6px;
                overflow: hidden;
            }
            .t-time {
                color: #AEAEAE;
                padding: 0 10px;
            }
            .message-box {
                max-width: 100%;
                background: white;
                border-radius: 5px;
                float: left;
                .ico-arrow {
                    position: relative;
                    word-break: break-all;

                    img {
                        float: left;
                    }

                    &::after {
                        content: '';
                        position: absolute;
                        left: -5px;
                        top: 14px;
                        width: 10px;
                        height: 10px;
                        background: white;
                        transform: rotate(45deg);
                    }
                }
                .video-cont {
                    height: 200px;
                    overflow: hidden;
                    border-radius: 5px;
                }
            }
        }
        .avatar {
            background: #7556fd;
        }
    }
    .item-right {
        .box {
            float: right;

            .t-time {
                text-align: right;
            }

            .message-box {
                float: right;
                background: #e1e1f3;

                .ico-arrow {
                    &::after {
                        left: auto;
                        right: -5px;
                        background: #e1e1f3;
                    }
                }
            }
        }
    }
}
.nocopy{
    -webkit-touch-callout: none;
    -webkit-user-select: none;
    -khtml-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
}
.bot-content {
    .content{
        position: relative;
        background: #f6f6f6;        
        &::after{
            content: '';
            position: absolute;
            top: 0;
            background: #e8e8e8;
            width: 100%;
            height: 1px;
            -webkit-transform: scaleY(0.5);
            transform: scaleY(0.5);
            -webkit-transform-origin: 0 0;
            transform-origin: 0 0;
        }
        .input-box {
            .van-cell {
                overflow: overlay;
                max-height: 200px;
                padding: 8px 10px;
                border-radius: 6px;
            }
        }
        .ability-box {            
            .item {
                height: 224px;
                position: relative;
                &::after {
                    content: '';
                    position: absolute;
                    top: 0;
                    background: #e8e8e8;
                    width: 100%;
                    height: 1px;
                    -webkit-transform: scaleY(0.5);
                    transform: scaleY(0.5);
                    -webkit-transform-origin: 0 0;
                    transform-origin: 0 0;
                }
            }
        }
    }
}
</style>