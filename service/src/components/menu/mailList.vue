<template>
    <div class="drawer-list scrollbar" @scroll="listScroll">
        <div class="item flex align-center" v-for="(item, index) in list" :key="index" @click="checkHland(item)">
            <a-checkbox v-model:checked="item.checked" class="mr20" v-show="isCheckbox" @change="checkChange" />
            <a-avatar :src="item.user_head" shape="square" :size="50" style="background-color: #f56a00" class="shrink" />
            <p class="layout text1 ml10">{{ item.rename && item.rename != '' ? item.rename : item.user_name }}</p>            
        </div>

        <!-- 用户聊天记录 -->
        <a-modal v-model:visible="visibleChat" :centered="true" :mask="false" :footer="null"
            :title="`与 “${itemCheck.rename && itemCheck.rename != '' ? itemCheck.rename : itemCheck.user_name}” 的对话窗口`" width="60%" :bodyStyle="{ padding: 0, backgroundColor:'#f5f5f5'}"
            class="left-modal">
            <MyChat :user_id = itemCheck.user_id style="height: calc(100vh - 200px);" class="bt" v-if="visibleChat"></MyChat>
        </a-modal>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, watch } from "vue";
import MyChat from "../chat.vue";
import axios from "@/utils/axios";
import { message } from 'ant-design-vue';
export default {
    name: "mailList",
    components: { MyChat },
    props: ['isCheckbox', 'isAll'],
    setup(props,contact) {
        const state = reactive({
            visibleChat:false, //用户聊天记录
            list: [],
            listForm: { page: 1, offset: 20, type: "all" }, //分页
            page: 0, //总页数
            itemCheck:{}, //当前查看用户
        });

        // 监听是否全选
        watch(() => props.isAll, () => {
            checkboxHland(); //全选或取消全选
        })

        onMounted(()=>{
            getList(); //列表            
        })

        //列表
        const getList = async () => {
            const res = await axios.post("/users", state.listForm, false);
            state.list = state.listForm.page == 1 ? res.data.list : state.list.concat(res.data.list);
            state.page = res.data.page;
            if (props.isCheckbox) checkboxHland(); //全选或取消全选
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

        //用户点击查看记录
        const checkHland = (item) => {
            if (props.isCheckbox) return false;
            state.visibleChat = true
            state.itemCheck = item
        }

        //全选或取消全选
        const checkboxHland = () => {
            let checkList = [];
            state.list.map((i) => {
                i['checked'] = props.isAll;
                if (i['checked']) checkList.push(i.user_id);
            })
            contact.emit('checkList', checkList);
        }

        //用户选择监听
        const checkChange = (e) => {
            let checkList = [];
            state.list.map((i) => {
                if (i['checked']) checkList.push(i.user_id);
            })
            contact.emit('checkList', checkList);
        }

        return {
            ...toRefs(state),
            listScroll, //监听滚动
            checkHland, //用户点击查看记录
            checkChange, //用户选择监听
        };
    },
};
</script>

<style lang="less" scoped>

</style>