<template>
    <div class="drawer-list scrollbar" @scroll="listScroll">
        <div class="item flex align-center" v-for="(item, index) in list" :key="index">
            <a-avatar :src="item.user_head" shape="square" :size="50" style="background-color: #f56a00" class="shrink" />
            <div class="layout hidden ml10">
                <p class="text1">{{ item.user_name }}</p>
                <p class="f12 mt6" :class="item.type == 'ip' ? 't-them' :'t-tipc'">{{ item.type == 'ip' ? 'IP 拉黑：' + item.ip : '用户拉黑' }}</p>
            </div>
            <a-button type="dashed" size="small" @click="cancelHland(item.user_id, item.ip, item.type)">取消拉黑</a-button>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { message } from 'ant-design-vue';
import axios from "@/utils/axios";
import bus from "@/utils/bus";
export default {
    name: "blacklist",
    setup() {
        const state = reactive({
            list: [],
            form: { page: 1, offset: 10 }, //分页
            page: 0, //总页数
        });

        onMounted(() => {
            getList(); //列表
        })

        //监听滚动
        const listScroll = (e) => {    
            if (e.srcElement.clientHeight + e.srcElement.scrollTop >= e.srcElement.scrollHeight) { //滚动到底部
                if (state.page > state.form.page) {
                    state.form.page += 1;
                    getList(); //列表
                }
            }
        }

        //获取用户列表
        const getList = async () => {
            const res = await axios.post("/rooms/black_list", state.form, false);
            state.list = state.form.page == 1 ? res.data.list : state.list.concat(res.data.list);
            state.page = res.data.page;
        }

        //取消拉黑
        const cancelHland = async (user_id,ip,type) => {
            const res = await axios.post("/rooms/black", { user_id,ip, type }, false);
            state.form.page = 1;
            getList();
            bus.emit("onRefresh", '');
        }

        return {
            ...toRefs(state),
            listScroll, //监听滚动
            cancelHland, //取消拉黑
        };
    },
};
</script>

<style lang="less" scoped>

</style>