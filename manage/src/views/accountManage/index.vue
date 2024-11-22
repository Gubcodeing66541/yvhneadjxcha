<template>
    <div class="content">
        <div class="content-box">
            <div id='main' class="main bg-white">
                <div id='top-box' class="top-box mb36">
                    <div class="input-box flex flex-warp">
                        <div class="item mr20">
                            <span>用户账号：</span>
                            <a-input v-model:value="pagination.search" placeholder="请输入" class="input mb10" style="width: 200px" />
                        </div>
                        <div class="item mr18">
                            <span>注册时间：</span>
                            <a-range-picker v-model:value="visit_time" format="YYYY-MM-DD" style="width:240px;border-radius: 8px;" @change="handleChange" />
                        </div>
                        <div>
                            <a-button type="primary" class="br12" @click="searchBind">
                                <template #icon>
                                    <SearchOutlined />
                                </template>
                                搜索
                            </a-button>
                            <a-button type="primary" class="br12 bg-tipc bn ml18" @click="reset()">重置</a-button>
                            <a-button type="primary" class="btn-add br12 ml18" @click="model.addVisible = true">
                                <template #icon>
                                    <PlusOutlined />
                                </template>
                                添加
                            </a-button>
                        </div>
                    </div>
                </div>
                <div class="tab-box">
                    <div :style="{ height: tabDetail.tabHeight + 50 + 'px' }">
                        <a-table :columns="tabDetail.columns" :data-source="tabDetail.data" :scroll="{ y: tabDetail.tabHeight }" :pagination=false>
                            <template #bodyCell="{ column, record }">
                                <template v-if="column.dataIndex === 'service_cnt'">
                                    <span class="link" @click="serviceList(record)" v-if="record.service_cnt" title="查看子客服">{{ record.service_cnt }}</span>
                                    <span v-else>{{ record.service_cnt }}</span>
                                </template>
                                <template v-else-if="column.dataIndex === 'account'">
                                    <span class="link" @click="accountDetail(record)" title="查看消费明细">{{ record.account }}</span>
                                </template>
                                <template v-else-if="column.dataIndex === 'status'">
                                    <span v-if="record.status =='no_use'" class="t-red">冻结中</span>
                                    <span v-else>正常</span>
                                </template>
                                <template v-else-if="column.dataIndex === 'operation'">
                                    <a-button type="primary" size="small" class="br6" @click="editVis(record.service_manager_id, record.name)">编辑</a-button>
                                    <a-button type="primary" size="small" class="br6 btn-red ml8" @click="deletHland(record.service_manager_id)">删除</a-button>
                                    <a-button type="primary" size="small" class="br6 btn-green ml8"  @click="renewVisibleForm.service_manager_id = record.service_manager_id, model.renewVisible=true">充值</a-button>
                                    <a-popconfirm placement="topRight" :title="record.status == 'no_use'?'是否取消冻结改账号':'冻结后该账号和子账号会立刻踢退？'" @confirm="frozen(record.service_manager_id, record.status == 'no_use' ? 'success' : 'no_use')">
                                        <a-button type="primary" size="small" class="btn-tipc br6 ml8" v-if="record.status == 'no_use'">解冻</a-button>
                                        <a-button type="primary" size="small" class="btn-add br6 ml8" v-else>冻结</a-button>
                                    </a-popconfirm>
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

        <!-- 新增账号 -->
        <a-modal v-model:visible="model.addVisible" :centered="true" :closable="false" :footer="null" title="新增账号">
            <a-form :model="addVisibleForm" layout="vertical" autocomplete="off" @finish="onFinish">
                <a-form-item label="账号数量" required>
                    <a-input type="number" v-model:value="addVisibleForm.count" placeholder="请输入账号数量" />
                </a-form-item>
                <a-form-item class="mb0">
                    <div class="tc">
                        <a-button class="br12" @click="model.addVisible = false">取消</a-button>
                        <a-button type="primary" html-type="submit" class="br12 ml26">确定</a-button>
                    </div>
                </a-form-item>
            </a-form>
        </a-modal>
        <!-- 新增客服成功弹窗 -->
        <a-modal v-model:visible="model.addSuccessVisible" :centered="true" :closable="false" :maskClosable="false" :footer="null" title="新增成功" :width="876">
            <div class="code-box flex just-center flex-warp">
                <div class="item br8 hh60 flex dir-column just-center t-them strong ml6 mr6 mb12 pl6 pr6" v-for="item in accountList" :key="item">
                    <p class="text1">账号：{{ item.username }}</p>
                    <p class="text1">密码：{{ item.password }}</p>
                </div>
            </div>
            <div class="tc mt40">
                <a-button type="primary" class="br12" @click="copy">一键复制</a-button>
            </div>
        </a-modal>
        <!-- 编辑账号 -->
        <a-modal v-model:visible="model.editVisible" :centered="true" :closable="false" :footer="null" title="编辑账号">
            <a-form :model="editVisibleForm" layout="vertical" autocomplete="off" @finish="editFinish">
                <a-form-item name="service_name" label="用户名称">
                    <a-input v-model:value="editVisibleForm.name" placeholder="请输入账号名称" />
                </a-form-item>
                <a-form-item name="password">
                    <template #label>
                        <div @click="editVisibleForm.password = 'abc123456'">
                            账号密码 <span class="them-link ml20 f12">重置密码</span>
                        </div>
                    </template>
                    <a-input v-model:value="editVisibleForm.password" readOnly placeholder="请输入账号密码" />
                </a-form-item>
                <a-form-item class="mb0">
                    <div class="tc">
                        <a-button class="br12" @click="model.editVisible = false">取消</a-button>
                        <a-button type="primary" html-type="submit" class="br12 ml26">保存</a-button>
                    </div>
                </a-form-item>
            </a-form>
        </a-modal>
        <!-- 充值 -->
        <a-modal v-model:visible="model.renewVisible" :centered="true" :closable="false" :footer="null" title="账号充值">
            <a-form :model="renewVisibleForm" layout="vertical" autocomplete="off" @finish="renewFinish">
                <a-form-item name="account" label="充值金额" :rules="[{ required: true, message: ' ' }]">
                    <a-input type="number" v-model:value="renewVisibleForm.account" placeholder="请输入充值金额" />
                </a-form-item>
                <a-form-item class="mb0">
                    <div class="tc">
                        <a-button class="br12" @click="model.renewVisible = false">取消</a-button>
                        <a-button type="primary" html-type="submit" class="br12 ml26">确定</a-button>
                    </div>
                </a-form-item>
            </a-form>
        </a-modal>        
        <!-- 子客服列表 -->
        <a-modal v-model:visible="model.serviceVisible" :width="1200" :centered="true" :closable="false" :footer="null" :title="itemDetail.name + ' 子客服列表'">
            <div class="tab-box">
                <div style="height:500px">
                    <a-table :columns="tabDetail.serviceColumns" :data-source="tabDetail.serviceData" :scroll="{ y: 500 }" :pagination=false>
                        <template #bodyCell="{ column, record }">
                            <template v-if="column.dataIndex === 'is_online'">
                                <span v-if="record.is_online" class="t-green">在线</span>
                                <span v-else class="t-ash">离线</span>
                            </template>
                            <template v-else-if="column.dataIndex === 'activate_time'">
                                <span v-if="record.is_activate">{{ record.activate_time }}</span>
                                <span v-else class="t-ash">暂未激活</span>
                            </template>
                            <template v-else-if="column.dataIndex === 'status'">
                                <span v-if="record.status =='time_out'" class="t-red">已过期</span>
                                <span v-else>正常</span>
                            </template>
                            <template v-else-if="column.dataIndex === 'time_out'">
                                <span v-if="record.is_activate">{{ record.time_out }}</span>
                                <span class="t-ash" v-else>-</span>
                            </template>
                        </template>
                    </a-table>
                </div>
                <div class="tr pt24">
                    <a-pagination size="small" :total="tabDetail.serviceTotal" v-model:current="servicePagination.page" :showSizeChanger="false" v-model:page-size="servicePagination.offset" :showTotal="(total) => `共有 ${total} 条数据`" show-quick-jumper @change="serviceChange" />
                </div>
            </div>
        </a-modal>
        <!-- 账单明细 -->
        <a-modal v-model:visible="model.accountVisible" :width="1200" :centered="true" :closable="false" :footer="null" :title="itemDetail.name + ' 账单明细'">
            <div class="td-bb">
                <div class="hidden" style="height: 550px;">
                    <a-table :columns="tabDetail.accountColumns" :data-source="tabDetail.accountData" :scroll="{ y: 500 }" :pagination=false>
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
                    <a-pagination size="small" :total="tabDetail.accountTotal" v-model:current="accountPagination.page" :showSizeChanger="false" v-model:page-size="accountPagination.offset" :showTotal="(total) => `共有 ${total} 条数据`" show-quick-jumper @change="accountChange" />
                </div>
            </div>
        </a-modal>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, createVNode } from "vue";
import { SearchOutlined, PlusOutlined, InfoCircleOutlined } from '@ant-design/icons-vue';
import { getElementData } from '@/utils';
import useClipboard from "vue-clipboard3";
import { message, Modal } from 'ant-design-vue';
import axios from "@/utils/axios";
const columns = [
    { title: '序号', dataIndex: 'service_manager_id', width: 100 },
    { title: '用户名称', dataIndex: 'name', width: 150, ellipsis: true },
    { title: '用户账号', dataIndex: 'member', width: 230, ellipsis: true },
    { title: '余额/￥', dataIndex: 'account', width: 120 },    
    { title: '子客服数', dataIndex: 'service_cnt', width: 100 },
    { title: '登录IP', dataIndex: 'ip', width: 130, ellipsis: true },
    { title: '账号状态', dataIndex: 'status', width: 100 },
    { title: '注册时间', dataIndex: 'create_time', ellipsis: true, width: 160 },
    { title: '操作', dataIndex: 'operation', width: 250 },
];
const serviceColumns = [
    { title: '序号', dataIndex: 'id', width: 60 },
    { title: '客服名称', dataIndex: 'name', width: 110, ellipsis: true },
    { title: '创建时间', dataIndex: 'create_time', width: 160, ellipsis: true },
    { title: '激活时间', dataIndex: 'activate_time', width: 160, ellipsis: true },
    { title: '过期时间', dataIndex: 'time_out', width: 160, ellipsis: true },
    { title: '通讯录', dataIndex: 'user_cnt', width: 80 },
    { title: '在线状态', dataIndex: 'is_online', width: 90 },
    { title: '账号状态', dataIndex: 'status', width: 90 },
    { title: '激活码', dataIndex: 'username', width: 220 },
];
const accountColumns = [
    { title: '序号', dataIndex: 'id', width: 80 },
    { title: '订单编号', dataIndex: 'order_id', width: 150, ellipsis: true },
    { title: '订单路径', dataIndex: 'member', width: 220 },
    { title: '订单前金额/￥', dataIndex: 'old_account', width: 120 },
    { title: '订单后金额/￥', dataIndex: 'account', width: 120 },
    { title: '金额/￥', dataIndex: 'renew', width: 90 },
    { title: '订单类型', dataIndex: 'reason', width: 92 },
    { title: '付款方式', dataIndex: 'pay_type', width: 90 },
    { title: '创建时间', dataIndex: 'create_time', width: 170, ellipsis: true },
]
export default {
    name: "accountManage",
    components: { SearchOutlined, PlusOutlined, InfoCircleOutlined },
    setup() {
        const { toClipboard } = useClipboard();
        const state = reactive({
            visit_time: [], //时间
            pagination: { page: 1, offset: 9 }, //分页
            servicePagination: { page: 1, offset: 8 }, //子客服分页
            accountPagination: { page: 1, offset: 8 }, //余额明细分页
            itemDetail:{}, //当前点击行
            tabDetail: {
                total: 0, //总条数
                tabHeight: 0, //表格高度
                columns, //表格配置
                data: [], //表格数据
                serviceTotal: 0, //子客服总条数
                serviceColumns, //子客服表格配置
                serviceData: [], //子客服表格数据
                accountTotal: 0, //消费明细总条数
                accountColumns, //消费明细表格配置
                accountData: [], //消费表格数据
                pay_count: {}, //消费总计
            },            
            model: {}, //弹窗控制显示
            addVisibleForm: { count: 1, }, //新增账号表单
            accountList: [], //新增的账号
            editVisibleForm: {}, //编辑账号表单
            renewVisibleForm: {}, //充值表单
        });

        onMounted(() => {
            getList();
            state.tabDetail.tabHeight = getElementData('main', 'h') - getElementData('top-box', 'h') - 200; //表格高度
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
            state.pagination = { page: 1, offset: 9 };
            state.visit_time = [];
            getList(); //获取列表
        }

        //获取列表
        const getList = async () => {
            const res = await axios.post("/service_manager/list", state.pagination, false);
            state.tabDetail.total = res.data.count;
            state.tabDetail.data = res.data.list;
        }
        //分页
        const change = (e) => {
            state.pagination.page = e;
            getList(); //获取列表
        }

        //新增账号
        const onFinish = async () => {            
            if (state.addVisibleForm.count == '' || state.addVisibleForm.count <= 0) return message.error("至少生成一个账号");
            state.addVisibleForm.count = parseInt(state.addVisibleForm.count);
            const res = await axios.post("/service_manager/create", state.addVisibleForm);
            state.accountList = res.data.list;
            state.model['addVisible'] = false;
            state.model['addSuccessVisible'] = true;
            getList(); //获取列表
        }

        //编辑账号
        const editVis = (service_manager_id, name) => {
            state.editVisibleForm.name = name;
            state.editVisibleForm.service_manager_id = service_manager_id;
            state.editVisibleForm.password = "";
            state.model['editVisible'] = true;
        }

        //编辑提交
        const editFinish = async () => {
            const res = await axios.post("/service_manager/reset_password", state.editVisibleForm);
            getList(); //获取列表
            state.model['editVisible'] = false;            
        }

        //充值
        const renewFinish = async () => {
            state.renewVisibleForm['account'] = parseInt(state.renewVisibleForm['account'])
            const res = await axios.post("/service_manager/renew", state.renewVisibleForm);
            getList(); //获取列表
            state.model['renewVisible'] = false;
            state.renewVisibleForm['account'] = '';
        }

        //删除
        const deletHland = (service_manager_id) => {
            Modal.confirm({
                title: '温馨提示',
                icon: createVNode(InfoCircleOutlined),
                centered: true,
                content: '确定删除该账号吗？',
                async onOk() {
                    const res = await axios.post("/service_manager/delete", { service_manager_id });
                    getList(); //获取列表
                },
            });
        }

        //账号冻结或解冻
        const frozen = async (service_manager_id,status) => {
            const res = await axios.post("/service_manager/ban", { service_manager_id,status });
            getList(); //获取列表
        }

        //子客服列表
        const serviceList = (item) => {
            state.itemDetail = item;            
            state.model['serviceVisible'] = true;
            serviceChange(1) //子客服分页
        }

        //子客服分页
        const serviceChange = async(e) => {
            state.servicePagination.page = e;
            state.servicePagination['service_manager_id'] = state.itemDetail['service_manager_id'];
            const res = await axios.post("/service_manager/get_service_list", state.servicePagination, false);
            state.tabDetail.serviceTotal = res.data.count;
            state.tabDetail.serviceData = res.data.list;
        }

        //余额明细
        const accountDetail = (item) => {
            state.itemDetail = item;
            state.model['accountVisible'] = true;
            accountChange(1) //子客服分页
        }

        //余额明细分页
        const accountChange = async(e) => {
            state.tabDetail.accountData = []
            state.accountPagination.page = e;            
            state.accountPagination['service_manager_id'] = state.itemDetail['service_manager_id'];
            const res = await axios.post("/pay/recorder", state.accountPagination, false);
            state.tabDetail.accountTotal = res.data.count;
            state.tabDetail.accountData = res.data.list;
            state.tabDetail.pay_count = res.data.pay_count;
        }

        //复制
        const copy = async () => {
            const accountList = state.accountList.map(item => '账号：' + item.username + '  ' + '密码：' +  item.password);
            try {
                await toClipboard(accountList.join('\r'));
                state.model['addSuccessVisible'] = false;
                message.success("复制成功")
            } catch (e) {
                message.error("复制失败")
            }
        }

        return {
            ...toRefs(state),
            change, //分页
            serviceChange, //子客服分页
            accountChange, //余额明细分页
            accountDetail, //查看余额明细
            handleChange, //客服名称选择
            searchBind, //搜索
            reset, //重置
            onFinish, //新增客服
            editVis, //编辑
            editFinish, //编辑提交
            renewFinish, //续费
            deletHland, //删除
            frozen, //账号冻结或解冻
            serviceList, //子客服列表
            copy, //复制
        };
    },
};
</script>

<style lang="less" scoped>
.code-box {
    .item {
        width: 180px;
        border: 1px solid #6971F8;
    }
}
.td-bb{
    .ant-table-wrapper {
        border: 1px solid #EEEEEE;
    }
}
</style>