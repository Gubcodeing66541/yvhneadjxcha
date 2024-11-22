<template>
    <div class="content">
        <div class="content-box">
            <div id='main' class="main bg-white">
                <div class="tab-box">
                    <div :style="{ height: tabDetail.tabHeight + 50 + 'px' }">
                        <a-table :columns="tabDetail.columns" :data-source="tabDetail.data" :scroll="{ y: tabDetail.tabHeight }" :pagination=false>
                            <template #bodyCell="{ column, record }">
                                <template v-if="column.dataIndex === 'user_head'">
                                    <a-avatar shape="square" :src="record.user_head" :size="38" style="background-color: #039B84;border-radius: 8px;" />
                                </template>
                                <template v-else-if="column.dataIndex === 'is_online'">
                                    <span v-if="record.is_online" class="t-green">在线</span>
                                    <span v-else class="t-ash">离线</span>
                                </template>
                            </template>
                        </a-table>
                    </div>
                    <div class="tr pt24">
                        <a-pagination size="small" :total="tabDetail.total" v-model:current="pagination.page" :showSizeChanger="false" v-model:page-size="pagination.offset" :showTotal="(total) => `共有 ${total} 条数据`" show-quick-jumper @change="change" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { FileSearchOutlined } from '@ant-design/icons-vue';
import { getElementData } from '@/utils';
import axios from "@/utils/axios";
const columns = [
    { title: '序号', dataIndex: 'user_id', width: 120 },
    { title: '游客昵称', dataIndex: 'user_name', width: 160, ellipsis: true },
    { title: '游客头像', dataIndex: 'user_head', width: 160 },
    { title: '游客IP', dataIndex: 'ip', width: 160 },
    { title: '地区', dataIndex: 'map', width: 150, ellipsis: true },
    // { title: '客服名称', dataIndex: 'name', width: 150 },
    { title: '注册时间', dataIndex: 'create_time', width: 150, ellipsis: true },
    { title: '状态', dataIndex: 'is_online', width: 120 },
];
export default {
    name: "userManage",
    components: { FileSearchOutlined },
    setup() {
        const state = reactive({
            pagination: { page: 1, offset: 8 }, //分页
            tabDetail: {
                total: 0, //总条数
                tabHeight: 0, //表格高度
                columns, //表格配置
                data: [], //表格数据
            },
        });

        onMounted(() => {
            getList();
            state.tabDetail.tabHeight = getElementData('main', 'h') - 160; //表格高度
        })

        //获取列表
        const getList = async () => {
            const res = await axios.post("/users/list", state.pagination, false);
            state.tabDetail.total = res.data.count;
            state.tabDetail.data = res.data.list;
        }
        //分页
        const change = (e) => {
            state.pagination.page = e;
            getList(); //获取列表
        }

        return {
            ...toRefs(state),
            change, //分页
        };
    },
};
</script>

<style lang="less" scoped>

</style>