<template>
    <div class="drawer-list scrollbar" @scroll="listScroll">
        <div class="item flex align-center" v-for="(item, index) in list" :key="index" @contextmenu.prevent="openMenu($event, item.id, index)">
            <div class="content layout text1 mr10" v-if="item.msg_type == 'text'">
                {{ item.msg_info }}
            </div>
            <div class="content layout text1 mr10" v-else>
                <a-image :src="item.msg_info" :height="60" class="img"></a-image>
            </div>
            <a-button type="dashed" size="small" @click="sendHland(item)" class="ml10 bc-them t-them">发送</a-button>
        </div>
    </div>

    <!-- 通讯录 -->
    <a-drawer v-model:visible="visible.mailList" :width="drawer.width" :closable="false" class="custom-class" title="选择群发用户" placement="right" :bodyStyle="drawer.bodyStyle">
        <mailList :isCheckbox="true" :isAll="isAll" @checkList="listHland" v-if="visible.mailList"></mailList>
        <template #footer>
            <div class="flex just-between algin-center">
                <a-checkbox v-model:checked="isAll" class="ml24 lh32">全选</a-checkbox>
                <a-button type="primary" @click="sendMsg">发送</a-button>
            </div>
        </template>
    </a-drawer>

    <!-- 新建 -->
    <a-modal v-model:visible="addVisible" :centered="true" :mask="false" :width="600" title="新建群发" @ok="handleOk" @cancel="addClose">
        <newlyAdded :detailForm="detailForm" @visibleForm="visibleForm" v-if="addVisible"></newlyAdded>
    </a-modal>

    <!-- 聊天记录右键菜单 -->
    <ul v-show="menu.visible" :style="{ left: menu.left + 'px', top: menu.top + 'px' }" class="contextmenu">
        <li @click="MenuHland('edit')">编辑</li>
        <li @click="MenuHland('delet')">删除</li>
    </ul>
    
</template>

<script>
import { reactive, toRefs, onMounted, onUnmounted } from "vue";
import { message } from 'ant-design-vue';
import mailList from './mailList.vue';
import newlyAdded from './newlyAdded.vue';
import axios from "@/utils/axios";
export default {
    name: "mass_dispatch",
    components: { mailList, newlyAdded },
    props: ['addVisible'],
    emits: ['addClose'],
    setup(props, contact) {
        const state = reactive({
            list: [],
            listForm: { page: 1, offset: 20, type: "group" }, //分页
            page: 0, //总页数
            visible:{}, //弹窗
            detailForm: {}, //修改表单
            drawer: { bodyStyle: { backgroundColor: '#f5f5f5' }, width: 500 },
            isAll:true, //是否全选
            checkList:[], //选中的用户
            form:{}, //新建表单
            menu: { top: '', left: '', visible: '', id: '', index: '' }, //右键菜单
        });

        onMounted(() => {
            getList(); //列表
            window.addEventListener('click', closeMenu); //监听页面点击关闭右键菜单
        })

        onUnmounted(() => {
            window.removeEventListener('click', closeMenu); //销毁页面监听
        })

        //列表
        const getList = async () => {
            const res = await axios.post("/service_message/list", state.listForm, false);
            state.list = state.listForm.page == 1 ? res.data.list : state.list.concat(res.data.list);
            state.page = res.data.page;
        }

        //监听滚动
        const listScroll = (e) => {
            if (e.srcElement.clientHeight + e.srcElement.scrollTop >= e.srcElement.scrollHeight) { //滚动到底部
                if (state.page > state.listForm.page) {
                    state.listForm.page += 1;
                    getList(); //列表
                }
            }
        }

        //新建；修改
        const handleOk = async () => {
            state.form['type'] = state.listForm.type;
            if (state.form['id'] && state.form['id'] != '') {
                const res = await axios.post("/service_message/update", state.form); //修改
            } else {
                const res = await axios.post("/service_message/create", state.form); //新增
            }
            state.listForm.page = 1;
            getList(); //列表
            addClose(); //关闭弹窗
        }

        //删除
        const deletHland = async (id,index) => {
            state.list.splice(state.menu.index, 1);
            const res = await axios.post("/service_message/delete", { id: state.menu.id });
        }

        //点击列表发送显示用户
        const sendHland = (item) => {
            state.visible['mailList'] = true
            state.detailForm = item;
        }

        //接收选择的用户
        const listHland = (list) => {
            state.checkList = list;            
        }

        //发送信息
        const sendMsg = async() => {
            let list = JSON.parse(JSON.stringify(state.checkList));
            if(list == '') return message.warning("请选择群发用户");
            const res = await axios.post('/message/send_all', { user_id: list, type: state.detailForm['msg_type'], content: state.detailForm['msg_info'] });
            state.detailForm = {};
            state.visible['mailList'] = false;
        }

        //新建表单
        const visibleForm = (e) => {
            state.form = e;
        }

        //关闭弹窗
        const addClose = () => {
            state.form = {};
            contact.emit('addClose',false)
        }

        //鼠标右键显示菜单
        const openMenu = (e, id, index) => {
            state.menu.top = e.pageY;
            state.menu.left = e.pageX;
            state.menu.visible = true;
            state.menu.id = id;
            state.menu.index = index;
        }

        //关闭右键菜单
        const closeMenu = () => {
            state.menu.visible = false;
            state.menu.id = '';
            state.menu.index = '';
        }

        //菜单点击共用
        const MenuHland = async (type) => {
            switch (type) {
                case 'edit':
                    state.detailForm = state.list[state.menu.index];
                    contact.emit('addClose', true)
                    break;
                case 'delet':
                    deletHland()
                    break;
            }
        }

        return {
            ...toRefs(state),
            listScroll, //监听滚动
            sendHland, //发送
            sendMsg, //推送信息
            listHland, //接收选择的用户
            visibleForm, //新建表单
            handleOk, //新建
            addClose, //关闭弹窗
            openMenu, //鼠标右键
            MenuHland, //右键菜单点击事件
        };
    },
};
</script>

<style lang="less" scoped>

</style>