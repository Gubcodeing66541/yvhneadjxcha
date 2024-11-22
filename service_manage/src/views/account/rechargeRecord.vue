<template>
    <div class="content">
        <div class="content-box">
            <div id='main' class="main bg-white">
                <div id='top-box' class="top-box mb36">
                    <div class="input-box flex flex-warp">
                        <div class="item mr20">
                            <span>订单类型：</span>
                            <a-select ref="select" v-model:value="pagination.reason" class="select" style="width: 120px">
                                <a-select-option value="renew_service_manager">充值</a-select-option>
                                <a-select-option value="create_service">新增子账号</a-select-option>
                                <a-select-option value="renew_service">续费子账号</a-select-option>
                            </a-select>
                        </div>
                        <div class="item mr20">
                            <span>订单路径：</span>
                            <a-input v-model:value="pagination.search" placeholder="请输入" style="width: 170px" />
                        </div>
                        <div class="item mr20">
                            <span>创建时间：</span>
                            <a-range-picker v-model:value="visit_time" format="YYYY-MM-DD" style="width:240px;border-radius: 8px;" @change="handleChange" />
                        </div>
                        <div class="item ml20">
                            <a-button type="primary" class="br12" @click="searchBind">
                                <template #icon>
                                    <SearchOutlined />
                                </template>
                                搜索
                            </a-button>
                            <a-button type="primary" class="br12 bg-tipc bn ml18" @click="reset()">重置</a-button>
                        </div>
                    </div>
                </div>
                <div class="tab-box">
                    <div class="hidden" :style="{ height: tabDetail.tabHeight + 50 + 'px' }">
                        <a-table :columns="tabDetail.columns" :data-source="tabDetail.data" :scroll="{ y: tabDetail.tabHeight - 66 }" :pagination=false>
                            <template #bodyCell="{ column, record }">
                                <template v-if="column.dataIndex === 'renew'">                                    
                                    <span :class="record.reason == 'renew_service_manager' ? 't-orange' : ''">{{ record.renew }}</span>
                                </template>
                                <template v-else-if="column.dataIndex === 'reason'">
                                    <!-- renew_service 续费客服 create_service 创建客服 renew_service_manager 账号充值 -->
                                    <a-tag color="orange" v-if="record.reason == 'renew_service_manager'">充值</a-tag>
                                    <a-tag color="blue" v-else-if="record.reason == 'create_service'">新增子账号</a-tag>
                                    <a-tag color="green" v-else>续费子账号</a-tag>
                                </template>
                                <template v-else-if="column.dataIndex === 'pay_type'">
                                    <span v-if="record.reason == 'renew_service_manager'">其它</span>
                                    <span v-else>余额</span>
                                </template>
                            </template>
                            <template #footer>
                                <div class="flex">
                                    <p>
                                        <span class="t-tipc">总计充值金额：</span>
                                        <span>￥{{ tabDetail.pay_count.renew_service_manager }}.00</span>
                                    </p>
                                    <p class="mlr30">
                                        <span class="t-tipc">总计新增子账号：</span>
                                        <span>{{ tabDetail.pay_count.create_service }}.00</span>
                                    </p>
                                    <p>
                                        <span class="t-tipc">总计续费子账号：</span>
                                        <span>{{ tabDetail.pay_count.renew_service }}.00</span>
                                    </p>
                                </div>
                            </template>
                        </a-table>
                    </div>
                    <div class="tr pt24">
                        <a-pagination size="small" :total="tabDetail.total" v-model:current="pagination.page" v-model:page-size="pagination.offset" :showTotal="(total) => `共有 ${total} 条数据`"
                            show-quick-jumper @change="change" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { SearchOutlined } from '@ant-design/icons-vue';
import { getElementData } from '@/utils';
import axios from "@/utils/axios";
const columns = [
    { title: '序号', dataIndex: 'id', width: 80 },
    { title: '订单编号', dataIndex: 'order_id', width: 170, ellipsis: true },
    { title: '订单账户', dataIndex: 'service_manager_member', width: 220 },
    { title: '订单路径', dataIndex: 'member', width: 220 },
    { title: '订单前金额/￥', dataIndex: 'old_account', width: 120 },
    { title: '订单后金额/￥', dataIndex: 'account', width: 120 },
    { title: '金额/￥', dataIndex: 'renew', width: 90 },
    { title: '订单类型', dataIndex: 'reason', width: 92 },
    { title: '付款方式', dataIndex: 'pay_type', width: 90 },
    { title: '创建时间', dataIndex: 'create_time', width: 170, ellipsis: true },
];
export default {
    name: "recharge",
    components: {
        SearchOutlined,
    },
    setup() {
        const state = reactive({
            visit_time: [], //时间
            pagination: { page: 1, offset: 8 }, //分页
            tabDetail: {
                total: 0, //总条数
                tabHeight: 0, //表格高度
                columns, //表格配置
                data: [], //表格数据
                pay_count:{}, //统计
            },
        });

        onMounted(() => {
            state.tabDetail.tabHeight = getElementData('main', 'h') - getElementData('top-box', 'h') - 200; //表格高度
            getList(); //获取列表
        })

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
            state.pagination = { page: 1, offset: 8 };
            state.visit_time = [];
            getList(); //获取列表
        }

        //获取列表
        const getList = async () => {
            const res = await axios.post("/pay/recorder", state.pagination, false);
            state.tabDetail.total = res.data.count;
            state.tabDetail.data = res.data.list;
            state.tabDetail.pay_count = res.data.pay_count;
        }
        //分页
        const change = (e) => {
            state.pagination.page = e;
            getList(); //获取列表
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
:deep(.tab-box .ant-table-wrapper tr:last-child td) {
    border-bottom: 1px solid #f0f0f0;
}
</style>