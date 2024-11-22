<template>
    <div class="content-box">
        <div id='main' class="main bg-white">
            <div id='top-box' class="top-box mb36">
                <div class="input-box flex flex-warp">
                    <div class="item mr16">
                        <span>访客名称：</span>
                        <a-input v-model:value="pagination.user_name" placeholder="请输入" style="width: 170px" />
                    </div>
                    <div class="item mr14">
                        <span>客服名称：</span>
                        <a-input v-model:value="pagination.service_name" placeholder="请输入"
                            style="width: 170px;border-radius: 8px;" />
                    </div>
                    <div class="item mr14">
                        <span>客服账号：</span>
                        <a-input v-model:value="pagination.service_member" placeholder="请输入"
                            style="width: 170px;border-radius: 8px;" />
                    </div>
                    <div class="item mr14">
                        <span>访问时间：</span>
                        <a-range-picker v-model:value="visit_time" format="YYYY-MM-DD" style="width:240px" @change="handleChange" />
                    </div>                    
                    <div>
                        <a-button type="primary" class="br12" @click="searchBind">搜索</a-button>
                        <a-button type="primary" class="br12 bg-tipc bn ml12" @click="reset()">重置</a-button>
                    </div>
                </div>
            </div>
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
                            <template v-else-if="column.dataIndex === 'operation'">
                                <a-button type="primary" ghost size="small" class="br6" @click="visible = true, formMsg = record">
                                    <template #icon>
                                        <file-search-outlined />
                                    </template>
                                    查看记录
                                </a-button>
                            </template>
                        </template>
                    </a-table>
                </div>
                <div class="tr pt24">
                    <a-pagination size="small" :total="tabDetail.total" v-model:current="pagination.page" v-model:page-size="pagination.offset" :showTotal="(total) => `共有 ${total} 条数据`" show-quick-jumper @change="change" />
                </div>
            </div>
        </div>
    </div>
    
    <!-- 聊天记录 -->
    <a-modal v-model:visible="visible" :centered="true" :footer="null" wrapClassName="chat-modal" :width="1034" :bodyStyle="{padding:0}">
        <chatRecord :formMsg="formMsg" v-if="visible"></chatRecord>
    </a-modal>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { FileSearchOutlined } from '@ant-design/icons-vue';
import chatRecord from '../../components/chatRecord.vue'
import { getElementData } from '@/utils';
import axios from "@/utils/axios";
import { message } from 'ant-design-vue';
const columns = [
    { title: '序号', dataIndex: 'user_id', width: 120 },
    { title: '访客头像', dataIndex: 'user_head', width: 160 },
    { title: '访客昵称', dataIndex: 'user_name', width: 160 },    
    { title: '访客IP', dataIndex: 'ip', width: 160 },
    { title: '地区', dataIndex: 'map', width: 160 },
    { title: '客服名称', dataIndex: 'name', width: 160 },
    { title: '状态', dataIndex: 'is_online', width: 120 },
    { title: '操作', dataIndex: 'operation', width: 120 },
];
export default {
    name: "visitorRecords",
    components: { FileSearchOutlined, chatRecord },
    setup() {
        const state = reactive({
            visit_time:[],
            pagination: { page: 1, offset: 7 }, //分页
            tabDetail: {
                total: 0, //总条数
                tabHeight: 0, //表格高度
                columns: columns, //表格配置
                data: [], //表格数据
            },
            visible: false, //聊天记录弹窗
            formMsg: {}//当前查看记录传送数据
        });

        onMounted(() => {
            state.tabDetail.tabHeight = getElementData('main', 'h') - getElementData('top-box', 'h') - 200; //表格高度
            getList();
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

        //日期选择监听
        const handleChange = (date, dateString) => {
            state.pagination['start_time'] = dateString[0];
            state.pagination['end_time'] = dateString[1];
        }

        //搜索
        const searchBind = () => {
            state.pagination.page = 1;
            getList(); //获取列表
        }

        //重置
        const reset = () => {
            state.visit_time = [];
            state.pagination = { page: 1, offset: 9 };
            getList();
            message.success("重置成功");
        }

        return {
            ...toRefs(state),
            handleChange, //日期选择监听
            searchBind, //搜索
            reset, //重置
            change, //分页
        };
    },
};
</script>

<style lang="less" scoped>

</style>