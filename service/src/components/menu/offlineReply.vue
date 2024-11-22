<template>
    <div class="drawer-list scrollbar" @scroll="listScroll">
        <div class="item draggable flex align-center" v-for="(item, index) in list" :key="index" :id="index"
            :draggable="true" @dragenter="dragenter($event, index)" @dragover="dragover($event, index)"
            @dragstart="dragstart(index)" title="鼠标长按移动排序" @contextmenu.prevent="openMenu($event, item.id, index)">
            <div class="content layout text1 mr10" v-if="item.msg_type == 'text'">
                {{ item.msg_info }}
            </div>
            <div class="content layout text1 mr10" v-else>
                <a-image :src="item.msg_info" :height="60" class="img"></a-image>
            </div>
            <a-switch :checked="item.status == 'enable' ? true:false" checked-children="开" un-checked-children="关" @change="helloChange($event, item)" />
        </div>
    </div>

    <!-- 新建 -->
    <a-modal v-model:visible="addVisible" :centered="true" :mask="false" :width="600" title="新建离线回复" @ok="handleOk" @cancel="addClose">
        <newlyAdded :detailForm="detailForm" @visibleForm="visibleForm" v-if="addVisible"></newlyAdded>
    </a-modal>

    <!-- 聊天记录右键菜单 -->
    <ul v-show="menu.visible" :style="{ left: menu.left + 'px', top: menu.top + 'px'}" class="contextmenu">
        <li @click="MenuHland('edit')">编辑</li>
        <li @click="MenuHland('delet')">删除</li>
    </ul>

</template>

<script>
import { reactive, toRefs, onMounted, onUnmounted } from "vue";
import { message } from 'ant-design-vue';
import newlyAdded from './newlyAdded.vue';
import axios from "@/utils/axios";
export default {
    name: "offlineReply",
    components: { newlyAdded },
    props: ['addVisible'],
    emits: ['addClose'],
    setup(props, contact) {
        const state = reactive({
            menu: { top: '', left: '', visible: '', id: '', index: '' }, //右键菜单
            list: [],
            listForm: { page: 1, offset: 20, type:"leave" }, //分页
            page:0, //总页数
            form: {}, //新建表单
            detailForm:{}, //修改表单
            timeout: null,  //拖拽防抖定时器
            swapForm:{}, //交换位置
            dragIndex: '', //源对象的下标
            enterIndex: '',//目标对象的下标
        });

        onMounted(() => {
            getList(); //列表
            window.addEventListener('click', closeMenu); //监听页面点击关闭右键菜单
        })

        onUnmounted(() => {
            window.removeEventListener('click', closeMenu); //销毁页面监听
        })

        //列表
        const getList = async() => {
            const res = await axios.post("/service_message/list", state.listForm, false);
            state.list = state.listForm.page == 1 ? res.data.list : state.list.concat(res.data.list);
            state.page = res.data.page;
        }

        //监听滚动
        const listScroll = (e) => {
            if (e.srcElement.clientHeight + e.srcElement.scrollTop >= e.srcElement.scrollHeight) { //滚动到底部
                if (state.page > state.listForm.page){
                    state.listForm.page += 1;
                    getList(); //列表
                }
            }
        }

        //新建；修改
        const handleOk = async () => {
            state.form['type'] = state.listForm.type;
            if (state.form['id'] && state.form['id']!= '' ){
                const res = await axios.post("/service_message/update", state.form); //修改
            }else{
                const res = await axios.post("/service_message/create", state.form); //新增
            }
            state.listForm.page = 1;
            getList(); //列表
            addClose(); //关闭弹窗
        }

        //启动或关闭
        const helloChange = async(e, item) => {
            item.status = e ? 'enable' : 'un_enable'
            const res = await axios.post("/service_message/update", item, false);
            state.listForm.page = 1;
            getList(); //列表
        }

        //删除
        const deletHland = async() => {
            state.list.splice(state.menu.index, 1);
            const res = await axios.post("/service_message/delete", { id: state.menu.id });
        }

        //位置交换排序
        const swap = async() => {
            const res = await axios.post("/service_message/swap", state.swapForm,false);
            state.listForm.page = 1;
            getList(); //列表
        }

        //源对象开始被拖动时触发
        const dragstart = (index) => {
            state.dragIndex = index;
            state.swapForm['from'] = state.list[index].id;
        }

        //源对象开始进入目标对象范围内触发
        const dragenter = (e, index) => {
            e.preventDefault();
            state.enterIndex = index;            
            if (state.timeout !== null) { clearTimeout(state.timeout) }
            // 拖拽事件的防抖
            state.timeout = setTimeout(() => {
                if (state.dragIndex !== index) {
                    state.swapForm['to'] = state.list[index].id;
                    swap();            
                }
            }, 200);
        }

        //源对象在目标对象范围内移动时触发
        const dragover = (e, index) => {
            e.preventDefault();
        }

        //表单值接收
        const visibleForm = (e) => {
            state.form = e;
        }

        //关闭弹窗
        const addClose = () => {
            state.form = {};
            state.detailForm = {};
            contact.emit('addClose', false)
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
        const MenuHland = async(type) => {
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
            helloChange, //启动或关闭
            visibleForm, //新建表单
            handleOk, //新建
            addClose, //关闭弹窗
            dragstart, //源对象开始被拖动时触发
            dragenter, //源对象开始进入目标对象范围内触发
            dragover, //源对象在目标对象范围内移动时触发
            openMenu, //鼠标右键
            MenuHland, //右键菜单点击事件
        };
    },
};
</script>

<style lang="less" scoped>

</style>