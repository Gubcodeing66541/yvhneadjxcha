<template>
    <div class="user-container">
        <div class="top-box hh50 plr10 flex align-center">
            <a-input-search allowClear v-model:value="form.user_name" placeholder="请输入关键词" @search="onRefresh()" />
            <a-popover trigger="click">
                <template #content>
                    <ul class="popover-pointer pb0">
                        <li class="item" v-for="item in listCot" :key="item" @click="delUser(item)">
                            删除近{{ item }}天用户
                        </li>
                    </ul>
                </template>
                <unordered-list-outlined class="f20 t-tipc pointer ml4" title="删除会话" />
            </a-popover>
        </div>
        <div id="listBox" class="list-box scrollbar" @scroll="listScroll">
            <div class="item" :class="item.is_top == 1 ? 'active-important':''" v-for="(item, index) in list" :key="index" @click="userHland(item.user_id,index)">
                <div class="pd15 flex just-between" :class="tabId == item.user_id ? 'active':''">
                    <a-badge :count="item.service_no_read">
                        <a-badge :dot="item.is_online?true:false" status="success" :offset="[-2,42]">
                            <a-avatar :src="item.user_head" shape="square" :size="44" style="background-color: #f56a00" class="shrink" />
                        </a-badge>
                    </a-badge>
                    <div class="right hidden ml10">
                        <div class="flex align-center just-between">
                            <div class="text1 mb3">{{ item.rename && item.rename != '' ? item.rename : item.user_name }}</div>
                            <div class="t-tipc f12 shrink">{{ filterTime(item.update_time) }}</div>
                        </div>
                        <div class="text1 t-tipc f12">
                            <span v-if="item.late_type == 'image'">[图片]</span>
                            <span v-else-if="item.late_type == 'video'">[视频]</span>
                            <span v-else-if="item.late_type == 'code'">[二维码]</span>
                            <span v-else>{{ item.late_msg }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, onUnmounted, computed, watch, createVNode } from "vue";
import { UnorderedListOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue';
import axios from "@/utils/axios";
import { Modal } from 'ant-design-vue';
import bus from "@/utils/bus";
import { useStore } from "vuex";
export default {
    name: "user",
    components: { UnorderedListOutlined },
    setup(props, contact) {
        const store = useStore();
        const onMessage = computed(() => store.state.onMessage); //socket消息
        const state = reactive({
            listCot: [7, 15, 30], //用户删除天数
            tabId:'', //当前选中用户
            list:[], //列表
            form: { page: 1, offset: 20, user_name: '', is_clear_status:0, type:'' }, //分页
            pageCount: 0, //总页数
            scrollTop: 0, //滚动条据顶部距离            
        });

        onMounted (()=>{            
            getList();
            busOn(); //监听事件订阅
        })

        onUnmounted(()=>{
            bus.all.clear(); //关闭所有订阅事件
        })

        watch(() => onMessage, (v1) => { //监听消息
            if (v1.value.type == 'online') onRefresh(state.list.length);
            else if (v1.value.type == 'leave') onRefresh(state.list.length);
            else if (v1.value.type == 'message') onRefresh(state.list.length);
            else if (v1.value.type == 'clear') onRefresh(state.list.length);
        }, { deep: true })

        //监听事件订阅
        const busOn = () =>{
            bus.all.clear(); //关闭所有订阅事件
            bus.on("ListType", (type) => ListType(type)); //列表类型切换
            bus.on("onRefresh", () => onRefresh()); //刷新
            bus.on("blackUser", (black_id) => blackUser(black_id)); //当前聊天清空
        }

        //监听滚动
        const listScroll = (e) => {
            state.scrollTop = e.srcElement.scrollTop; //据顶部高度            
            if (e.srcElement.clientHeight + e.srcElement.scrollTop >= e.srcElement.scrollHeight) { //滚动到底部
                if (state.pageCount > state.list.length){
                    state.form.page += 1;
                    getList();
                }
            }
        }

        //获取用户列表
        const getList = async (is_clear_status = 0) => {
            state.form.is_clear_status = is_clear_status;
            const res = await axios.post("/rooms/list", state.form, false);
            state.list = state.form.page == 1 ? res.data.list : state.list.concat(res.data.list);
            state.pageCount = res.data.count;
        }

        //刷新
        const onRefresh = async (offset = 20) => {
            const form = JSON.parse(JSON.stringify(state.form))
            form.page = 1;
            form.offset = state.scrollTop >= 20 ? offset <= 20 ? 20 : offset : 20;
            const res = await axios.post("/rooms/list", form, false);
            state.list = res.data.list;
            if (form.offset <= 20) state.form.page = 1;
        }

        //列表类型切换
        const ListType = (type) => {
            state.form.type = type;
            state.form.offset = 20;
            state.tabId = '';
            onRefresh()
            contact.emit('user_id', null);
        }

        //删除会话用户
        const delUser = (day) => {
            Modal.confirm({
                title: '删除会话用户',
                icon: createVNode(ExclamationCircleOutlined),
                content: `请确定是否删除进${day}天会话用户`,
                centered:true,
                async onOk() {
                    const res = await axios.post('/rooms/delete_day', { day });
                    contact.emit('user_id', null);
                    onRefresh();
                }
            });            
        }

        //聊天切换
        const userHland = (id,index) => {
            state.list[index]['service_no_read'] = 0;
            state.tabId = id;
            contact.emit('user_id', id)
        }

        //用户拉黑
        const blackUser = (black_id) => {
            onRefresh();
            if (black_id == state.tabId){
                state.tabId = '';
                contact.emit('user_id', null);
            }
        }

        //时间过滤
        const filterTime = (date) => {
            var time = Math.floor(new Date().getTime() / 1000) - Math.floor(new Date(date).getTime() / 1000);
            if (time < 60) {
                return "刚刚";
            } else if (time / 60 < 60) {
                return Math.floor(time / 60) + "分钟前";
            } else if (time / 3600 < 24) {
                return Math.floor(time / 3600) + "小时前";
            }else{
                return Math.floor(time / 86400) + "天前";
            }
        }

        return {
            ...toRefs(state),
            onRefresh, //刷新
            delUser, //删除用户
            userHland, //用户切换
            listScroll, //监听滚动
            filterTime, //时间过滤
        };
    },
};
</script>

<style lang="less" scoped>
.user-container{
    width: 100%;
    min-width: 0;
    height: 100%;
    .top-box {
        background: #f5f5f5;
    }
    .list-box {
        cursor: context-menu;
        height: calc(100% - 50px);
        width: 100%;
        padding-bottom: 20px;
        .item {
            width: 100%;
            .right {
                width: 0;
                flex: 1;
            }

            &:hover {
                background: #dedbda;
            }
        }

        .active {
            background: #c5c3c4 !important;
        }

        .active-important {
            background: #e0dfdd !important;
        }
    }
}
</style>