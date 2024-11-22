<template>
    <div class="drawer-list">
        <div class="item flex align-center" v-for="(item, index) in list" :key="index">
            <div class="td"> {{ item.dates }} </div>
            <div class="td"> 新增客户：{{ item.user_cnt }} </div>
            <div class="td"> 有效回复率：{{ item.rate }}% </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import axios from "@/utils/axios";
export default {
    name: "statistics",
    setup() {
        const state = reactive({
            list: []
        });

        onMounted(() => {
            getList();
        })

        const getList = async () => {
            const res = await axios.post("/count",'',false);
            state.list = res.data.count;
        }

        return {
            ...toRefs(state),
        };
    },
};
</script>

<style lang="less" scoped>
.drawer-list{
    .item{
        .td{
            width: 33.33%;
        }
    }
}
</style>