<template>
    <div class="warp">
        <div class="container">
            <h2 class="tc ptb40 f24">客服近七天数据统计</h2>
            <div class="tc">
                <a-input-search v-model:value="searValue" :allowClear="true" placeholder="输入激活码查询" size="large" enter-button @search="onSearch" class="w40" />
            </div>
            <div class="tab-box mt50">
                <a-table :dataSource="list" :columns="columns" :pagination=false />
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import {useRoute} from "vue-router";
import axios from "@/utils/axios";
const columns = [
    { title: "时间", dataIndex: "dates", align: "center" },
    { title: "新增用户", dataIndex: "user_cnt", align: "center" },
    { title: "回复率%", dataIndex: "rate", align: "center" },
];
export default {
    name: "statistics",    
    setup() {
        const route = useRoute();
        const state = reactive({
            searValue:'', //搜索值
            columns, //表格配置
            list: [],
        });

        onMounted(()=>{
            document.title = '客服统计'
            state.searValue = route.query.username || "";
            getList();
        })

        const getList = async () => {            
            if (state.searValue == '') return false;
            const res = await axios.post("/common/count", { username: state.searValue }, false);
            state.list = res.data.count;
        }

        const onSearch = (e) => {
            state.searValue = e;
            getList();
        }

        return {
            ...toRefs(state),
            onSearch, //搜索
        };
    },
};
</script>

<style lang="less" scoped>
.warp{
    .container{
        width: 1200px;
        margin: 0 auto;
    }
}
</style>