<template>
    <div class="content-box">
        <div id='main' class="main bg-white">
            <div id='top-box' class="top-box mb36">
                <div class="input-box flex flex-warp">
                    <div class="item mr20">
                        <span>访客名称：</span>
                        <a-input v-model:value="pagination.user_name" placeholder="请输入" class="input" style="width: 170px" />
                    </div>
                    <div class="item mr20">
                        <span>拉黑日期：</span>
                        <a-range-picker v-model:value="visit_time" format="YYYY-MM-DD" style="width:240px" @change="handleChange" />
                    </div>                    
                    <div>
                        <a-button type="primary" class="br12" @click="searchBind">搜索</a-button>
                        <a-button type="primary" class="br12 bg-tipc bn ml12" @click="reset()">重置</a-button>
                        <a-button type="primary" class="btn-add br12 ml18" @click="addVisible = true">
                            <template #icon>
                                <PlusOutlined />
                            </template>
                            新增黑名单
                        </a-button>
                    </div>
                </div>
            </div>
            <div class="tab-box">
                <div :style="{ height: tabDetail.tabHeight + 50 + 'px' }">
                    <a-table :columns="tabDetail.columns" :data-source="tabDetail.data" :scroll="{ y: tabDetail.tabHeight }" :pagination=false>
                        <template #bodyCell="{ column,record }">
                            <template v-if="column.dataIndex === 'user_head'">
                                <a-avatar shape="square" :src="record.user_head" :size="38" style="background-color: #039B84;border-radius: 8px;" v-if="record.type == 'user'" />
                                <a-avatar shape="square" :size="38" style="background-color: #039B84;border-radius: 8px;" v-else>IP</a-avatar>
                            </template>
                            <template v-else-if="column.dataIndex === 'type'">
                                <span>{{ record.type=='ip'?'ip拉黑':'用户拉黑'}}</span>
                            </template>
                            <template v-else-if="column.dataIndex === 'user_name'">
                                <span>{{ record.type == 'user' ? record.user_name : '-' }}</span>
                            </template>
                            <template v-else-if="column.dataIndex === 'operation'">
                                <a-button type="primary" danger ghost size="small" class="br6" @click="cancelHland(record.Id)">
                                    <template #icon>
                                        <CloseSquareOutlined />
                                    </template>
                                    取消黑名单
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

        <!-- 新增黑名单 -->
        <a-modal v-model:visible="addVisible" :centered="true" :closable="false" :footer="null" title="新建黑名单">
            <a-form :model="addVisibleForm" layout="vertical" autocomplete="off" @finish="onFinish">
                <a-form-item label="新建黑名单" required>
                    <a-radio-group v-model:value="addVisibleForm.type" name="radioGroup">
                        <a-radio value="user">访客</a-radio>
                        <a-radio value="ip">IP地址</a-radio>
                    </a-radio-group>
                </a-form-item>        
                <a-form-item label="访客" required v-if="addVisibleForm.type=='user'">
                    <a-select v-model:value="addVisibleForm.user_id" show-search placeholder="请输入后选择"
                        :default-active-first-option="false" :show-arrow="false" :filter-option="false" :not-found-content="null" :field-names="{ label: 'user_name', value: 'user_id' }"
                        :options="visitOptions" @search="handleSearch" style="border-radius: 8px;"></a-select>
                </a-form-item>
                <a-form-item label="IP地址" required v-else>
                    <a-input v-model:value="addVisibleForm.ip" placeholder="请输入IP地址" />
                </a-form-item>
                <a-form-item label="时间" required>
                    <a-checkbox v-model:checked="forever">无期限</a-checkbox>
                    <a-input type="number" v-model:value="addVisibleForm.day" :disabled="forever" placeholder="请输入天数" class="mt10" />
                </a-form-item>        
                <a-form-item class="mb0">
                    <div class="tc">
                        <a-button class="br12" @click="addVisible = false">取消</a-button>
                        <a-button type="primary" html-type="submit" class="br12 ml26">确定</a-button>
                    </div>
                </a-form-item>
            </a-form>
        </a-modal>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, createVNode } from "vue";
import { CloseSquareOutlined, InfoCircleOutlined, PlusOutlined } from '@ant-design/icons-vue';
import { getElementData } from '@/utils';
import { message, Modal } from 'ant-design-vue';
import axios from "@/utils/axios";
const columns = [
    { title: '序号', dataIndex: 'Id', width: 80 },
    { title: '访客头像', dataIndex: 'user_head', width: 160 },
    { title: '访客名称', dataIndex: 'user_name', width: 160 },
    { title: 'IP', dataIndex: 'ip', width: 140 },
    { title: '类型', dataIndex: 'type', width: 140 },
    { title: '拉黑时间', dataIndex: 'create_time', width: 160, ellipsis: true },
    { title: '操作', dataIndex: 'operation', width: 100 },
];
export default {
    name: "blacklist",
    components: { CloseSquareOutlined, InfoCircleOutlined, PlusOutlined },
    setup() {
        const state = reactive({
            visit_time: [],//日期选择储存
            pagination: { page: 1, offset: 9 }, //分页
            tabDetail: {
                total: 0, //总条数
                tabHeight: 0, //表格高度
                columns: columns, //表格配置
                data: [], //表格数据
            },
            addVisible:false, //添加黑名单
            addVisibleForm: { //添加黑名单form
                type: 'user',
            },
            forever: true,
            visitOptions:[],//访客
        });

        onMounted(() => {            
            state.tabDetail.tabHeight = getElementData('main', 'h') - getElementData('top-box', 'h') - 200; //表格高度
            getList();
        })

        //获取列表
        const getList = async () => {
            const res = await axios.post("/users/black", state.pagination, false);
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

        //取消拉黑
        const cancelHland = (id) => {
            Modal.confirm({
                title: '温馨提示',
                icon: createVNode(InfoCircleOutlined),
                centered: true,
                content: '确定取消该用户的黑名单吗？',
                async onOk() {
                    const res = await axios.post("/black/delete", { id });
                    getList(); //获取列表
                },
                onCancel() { },
            });
        }

        //访客搜索
        const handleSearch = async(val) => {
            if(val=='') return
            const res = await axios.post("/users/list", { page: 1, offset: 30, user_name:val }, false);
            state.visitOptions = res.data.list;
        };

        //新增客服
        const onFinish = async() => {
            if (state.addVisibleForm.type == 'user' && !state.addVisibleForm.user_id || state.addVisibleForm.type=='ip' && !state.addVisibleForm.ip) return message.error("请填写完整后提交")
            if (state.forever) state.addVisibleForm.day = 0;
            const res = await axios.post("/black/add", state.addVisibleForm);
            getList(); //获取列表
            state.addVisible = false;
            state.addVisibleForm = { type: 'user'};
            state.forever = true;
        }

        return {
            ...toRefs(state),
            handleChange, //日期选择监听
            searchBind, //搜索
            reset, //重置
            cancelHland, //取消拉黑
            handleSearch, //访客搜索
            onFinish, //新增黑名单提交
            change, //分页
        };
    },
};
</script>

<style lang="less" scoped>
:deep(.ant-select:not(.ant-select-customize-input) .ant-select-selector){
    border-radius: 8px;
}
</style>