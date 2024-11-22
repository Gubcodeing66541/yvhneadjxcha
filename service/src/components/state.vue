<template>
    <div class="user-container flex align-center just-center" style="margin-top:-1px">
        <div class="list flex align-center just-center plr10">
            <div class="item pointer text1" :class="tabIndex == item.value ? 'active' : ''" v-for="item in tabList" :key="item" @click="tabHland(item.value)">
                {{ item.title }}
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs } from "vue";
import bus from "@/utils/bus";
export default {
    name: "state",
    setup() {
        const state = reactive({
            tabList: [{ title: "全部", value: 'all' }, { title: "已回复", value: 'server_read' }, { title: "未回复", value: 'server_no_read' }, { title: "置顶", value: 'top' }],
            tabIndex: 'all',
        });

        const tabHland = (value) => {
            state.tabIndex = value;
            bus.emit("ListType", value);  // 发布事件
        }

        return {
            ...toRefs(state),
            tabHland, //状态切换
        };
    },
};
</script>

<style lang="less" scoped>
.user-container{
    height: 100%;
    .list {
        .item {
            padding: 4px 15px;
            margin-left: -1px;
            border: 1px solid #e1e1e1;
            color: #606266;
        }

        .active {
            color: #fff;
            background-color: @primary-color;
            border-color: @primary-color;
        }
    }
}
</style>