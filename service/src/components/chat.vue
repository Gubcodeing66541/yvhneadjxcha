<template>
    <div class="chat-container flex">
        <div class="left-content">
            <!-- 顶部 -->
            <div class="top-box hh50 flex align-center just-between plr20">
                <div class="f16 text1">{{ detail.rename && detail.rename != '' ? detail.rename : detail.user_name }}</div>
                <div class="flex algin-center">
                    <div class="t-tipc lh40 pointer text1" @click="topping">
                        <i class="iconfont icon-zhiding"></i>
                        <span>{{ detail.is_top == 1 ? '取消置顶' : '置顶' }}</span>
                    </div>
                    <div class="mlr20 t-tipc lh40 pointer text1" @click="endSession">
                        <i class="iconfont icon-shanchu"></i>
                        <span>结束会话</span>
                    </div>
                    <div class="t-tipc lh40 pointer" @click="showinfo = !showinfo">
                        <i class="iconfont icon-gengduo f26"></i>
                    </div>
                </div>
            </div>
            <!-- 聊天记录 -->
            <div ref="chatBox" class="chat-content ptb12 relative scrollbar" @scroll="divScroll">
                <div class="item hidden" :class="item.send_role != 'user' ?'item-right':''" v-for="(item, index) in list" :key="index">
                    <div class="flex just-center" v-if="item.type == 'time'">
                        <p class="bg-tipc tc" style="width:200px;color: #91d382;">{{ item.content }}</p>
                    </div>
                    <div class="flex just-center" v-else-if="item.type =='remind'">
                        <p class="bg-tipc tc" style="width:270px;color: #f7b5b5;">{{ item.content }}</p>
                    </div>
                    <div v-else>
                        <a-avatar :src="item.send_role == 'service' || item.send_role == 'hello' ? info.head : item.send_role == 'user' ? detail.user_head : info.bot_head" shape="square" :size="40" style="background-color: #f56a00" class="avatar" />
                        <div class="box">
                            <p class="t-time plr10 lh12 f12">{{ item.create_time }}</p>
                            <div class="message-row pt6 plr10 hidden">
                                <div class="message-box relative"
                                    :style="{ 'background': item.send_role == 'user' ? '#ffffff' : '#e9e9fb' }"
                                    @contextmenu.prevent="openMenu($event, item, index)">
                                    <!-- 文本信息 -->
                                    <div class="ptb10 plr16 ico-arrow relative" v-if="item.type =='text'">
                                        <div style="white-space:pre-wrap" v-html="textReplaceEmoji(item.content)"></div>
                                    </div>
                                    <!-- 图片信息 -->
                                    <div class="img-cont hidden br5" v-else-if="item.type == 'image'">
                                        <a-image :src="item.content" :height="200" fallback="/src/assets/image/erro.png" :placeholder="true" />
                                    </div>
                                    <!-- 视频信息 -->
                                    <div class="video-cont" v-else-if="item.type == 'video'">
                                        <video :src="item.content" :height="200" controls style="object-fit:fill;width: 100%;"></video>
                                    </div>
                                    <!-- 二维码信息 -->
                                    <div class="code-cont pd10" v-else>
                                        <div class="tc" style="height:240px"><vue-qr :logoSrc="info.head" :text="info.web" :backgroundColor="info.code_background" :colorLight="info.code_background" :colorDark="info.code_color" :logoMargin="4" :size="240" :margin="10"></vue-qr></div>
                                        <p class="tc mt10">长按保存专属二维码，聊天记录不丢失</p>
                                    </div>
                                    <div class="read h100 mr10 f12 flex align-center" :style="{ 'color': item.is_read ? '#D6D6D6' : '#b1b1fd' }">
                                        <span>{{ item.is_read ? '已读' :'未读'}}</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    
                </div>                
            </div>
            <!-- 聊天编辑 -->
            <div class="chat-edit ptb10 plr20">
                <editBox :user_id="user_id" @funHland="funHland"></editBox>
            </div>
        </div>
        <!-- 详细信息 -->
        <div class="t-3s" :class="showinfo && detail.user_id ? 'right-content' :'hide-box'">
            <userInfo :userInfo="detail" :login_log="login_log" @changeDetail="getDetail" v-if="showinfo && detail.user_id"></userInfo>
        </div>
        <!-- 聊天记录右键菜单 -->
        <ul v-show="menu.visible" :style="{ left: menu.left + 'px', top: menu.top + 'px'}" class="contextmenu">
            <li @click="MenuHland('group')" v-show="menu.item.type !== 'video' && menu.item.type !== 'code'">加入群发</li>
            <li @click="MenuHland('hello')" v-show="menu.item.type !== 'video' && menu.item.type !== 'code'">加入打招呼</li>
            <li @click="MenuHland('quick_reply')" v-show="menu.item.type !== 'video' && menu.item.type !== 'code'">加入快捷回复</li>
            <li @click="MenuHland('withdraw')" v-show="menu.item.send_role == 'service'">撤回</li>
        </ul>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, onUnmounted, nextTick, createVNode, ref, watch, computed } from "vue";
import { useStore } from "vuex";
import editBox from "@/components/editBox.vue"; //编辑框
import userInfo from "@/components/userInfo.vue"; //用户信息
import { textReplaceEmoji } from "@/utils/emojis";
import { message, Modal } from 'ant-design-vue';
import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
import vueQr from 'vue-qr/src/packages/vue-qr.vue';
import axios from "@/utils/axios";
import bus from "@/utils/bus";
export default {
    name: "MyChat",
    components: { userInfo, editBox, ExclamationCircleOutlined, vueQr },
    props: ['user_id'],
    setup(props, contact) {        
        const store = useStore();
        const chatBox = ref(0);
        const onMessage = computed(() => store.state.onMessage); //socket消息
        const state = reactive({
            info: computed(() => store.state.info), //客服信息
            list:[], //聊天记录
            form: { page: 1, offset: 20 }, //分页
            pageCount: 0, //聊天记录总页数
            scrollHeight: 0, //滚动条总高度
            detail:{}, //详情
            login_log:[], //用户登录信息
            showinfo:true, //控制详细信息                      
            menu:{ top:'', left:'', visible:'', item:{}, index:'' }, //右键菜单
        });

        watch(() => props.user_id, (v1) => {
            init(); //监听user_id切换
        }, { deep: true })

        watch(() => onMessage, (v1) => { //监听消息
            if (v1.value.type == 'message') onmessage(v1.value); //消息
            else if (v1.value.type == 'online') online(v1.value); //上线
        }, { deep: true })

        onMounted(() => {
            init(); //初始化
            window.addEventListener('click', closeMenu); //监听页面点击关闭右键菜单
        })

        onUnmounted(() => {
            window.removeEventListener('click', closeMenu); //销毁页面监听
        })

        //房间初始化
        const init = () => {
            state.form.page = 1;
            getMessageList(); //获取聊天记录
            getDetail(); //房间详情
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
        const getMessageList = async() => {
            state.form['user_id'] = props.user_id;
            const res = await axios.post('/message/list',state.form,false);
            state.list = state.form.page == 1 ? res.data.list : res.data.list.concat(state.list);
            state.pageCount = res.data.page;
            nextTick(() => {
                if (state.form.page == 1) chatBox.value['scrollTop'] = chatBox.value['scrollHeight']; //滚动到底
                else chatBox.value['scrollTop'] = chatBox.value['scrollHeight'] - state.scrollHeight; //重置滚动高度
            });
        }

        //获取房间详细信息
        const getDetail = async() => {
            const res = await axios.post("/rooms/detail", {user_id:props.user_id}, false);
            state.detail = res.data.user;
            state.login_log = res.data.login_log;
        }

        //置顶或取消置顶
        const topping = async () => {
            const res = await axios.post("/rooms/top", { user_id: props.user_id, top: state.detail['is_top'] == 1 ? 0 : 1});
            getDetail();
            bus.emit("onRefresh", ''); //列表刷新
        }

        //结束会话
        const endSession = () => {
            Modal.confirm ({
                title: '是否结束当前会话？',
                icon: createVNode(ExclamationCircleOutlined),
                content: "结束会话后聊天列表中将会清楚该用户和聊天记录",
                centered: true,
                async onOk() {
                    const res = await axios.post("/rooms/end", { user_id: props.user_id }, false);
                    bus.emit("blackUser", state.detail['user_id']);
                }
            });
        }

        //鼠标右键显示菜单
        const openMenu = (e,item,index) => {
            state.menu.top = e.pageY;
            state.menu.left = e.pageX;
            state.menu.visible = true;
            state.menu.item = item;
            state.menu.index = index;
        }       
        //关闭右键菜单
        const closeMenu = () => {
            state.menu.visible = false;
            state.menu.item = {};
            state.menu.index = '';
        }
        //菜单点击共用
        const MenuHland = (type) => {
            switch (type) {
                case 'group':       //加入群发
                    createMesg(type)
                    break;
                case 'hello':       //加入打招呼
                    createMesg(type)
                    break;
                case 'quick_reply': //快捷回复
                    createMesg(type)
                    break;
                case 'withdraw':    //撤回
                    withdraw()
                    state.list.splice(state.menu.index, 1);                    
                    break;
            }
        }
        //加入打招呼和快捷回复
        const createMesg = async (type) => {
            const res = axios.post('/service_message/create', { type, msg_type: state.menu.item['type'], msg_info: state.menu.item['content'] })
            closeMenu()
        }

        //撤回单个消息
        const withdraw = async() => {
            const res = await axios.post("/message/remove_msg", { user_id: props.user_id, id: state.menu.item['id']});
            bus.emit("onRefresh", ''); //列表刷新
        }

        //输入框回调事件
        const funHland = (item) => {
            if (item.type == 'withdraw'){ //一键撤回
                state.list = [];
                bus.emit("onRefresh", ''); //列表刷新
            }
        }

        //消息监听
        const onmessage = (data) => {
            if (data.content.user_id != props.user_id) return false
            state.list.push(data.content);
            nextTick(() => {
                chatBox.value['scrollTop'] = chatBox.value['scrollHeight']; //滚动到底
            })
        }

        //上线消息
        const online = (data) => {
            if (data.content.user_id != props.user_id) return false;            
            state.list.forEach(item => {
                if (!item.is_read) item.is_read = 1
            })
            getDetail(); //房间详情
        }

        return {
            ...toRefs(state),
            getDetail, //房间详情
            topping, //置顶或取消置顶
            endSession, //结束会话
            divScroll, //监听滚动
            textReplaceEmoji, //文字表情处理
            openMenu, //鼠标右键
            MenuHland, //右键菜单点击事件
            chatBox,
            funHland,            
        };
    },
};
</script>

<style lang="less" scoped>
.chat-container{
    height: 100%;
    .left-content{
        width: 100%;
        height: 100%;
        min-width: 300px;
        .top-box {
            border-bottom: 1px solid #eee;
        }
        .chat-content {
            height: calc(100% - 250px);
            .item {
                padding: 10px 20px;
                .avatar {
                    margin-right: 6px;
                    float: left;
                }
                .box{
                    float: left;
                    width: 70%;                    
                    .t-time {
                        color: #AEAEAE;
                    }
                    .message-row{
                        max-width: 100%;
                    }
                    .message-box {
                        border-radius: 5px;
                        float: left;
                        max-width: 100%;                        
                        .ico-arrow {
                            word-break: break-all;
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
                        .img-cont {
                            :deep(.ant-image-img) {
                                width: auto;
                                height: 100%;
                            }
                            :deep(.ant-image-img-placeholder) {
                                background: none;
                            }
                        }                
                        .video-cont {
                            height: 200px;
                            overflow: hidden;
                            border-radius: 5px;
                        }
                        .code-cont{
                            height: 296px;
                        }
                    }
                    .read{
                        position: absolute;
                        height: 100%;
                        top: 0;
                        left: -30px;    
                        width: 26px;
                        opacity: 0;
                    }
                }
            }        
            .item-right {
                .avatar {
                    margin-left: 6px;
                    float: right !important;
                }
                .box{
                    float: right;
                    .t-time { text-align: right; }                    
                    .message-box {
                        float: right;
                        background: #e9e9fb;
                        margin-left: 30px;
                        .ico-arrow {
                            &::after {
                                left: auto;
                                right: -5px;
                                background: #e9e9fb;
                            }
                        }
                    }
                    .read{ right: -30px; opacity: 1;}
                }
            }
        }
    
        .chat-edit {
            height: 200px;
            border-top: 1px solid #eee;
        }
    }
    .right-content{
        width: 25%;
        border-left: 1px solid #eee;
    }
    .hide-box{
        width: 0;
    }
}
</style>