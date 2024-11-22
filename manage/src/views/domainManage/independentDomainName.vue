<template>
    <div class="content">
        <div class="content-box">
            <div id='main' class="main bg-white">
                <div id='top-box' class="top-box mb36">
                    <div class="input-box flex flex-warp">
                        <div class="item mr20">
                            <span>域名地址：</span>
                            <a-input v-model:value="pagination.domain" placeholder="请输入" class="input mb10" style="width: 170px" />
                        </div>
                        <div class="item mr20">
                            <span>客服账号：</span>
                            <a-input v-model:value="pagination.username" placeholder="请输入" class="input mb10" style="width: 170px" />
                        </div>
                        <div class="item mr20">
                            <span>是否绑定：</span>
                            <a-select ref="select" v-model:value="pagination.is_bind_service" class="select" style="width: 120px">
                                <a-select-option :value="1">绑定</a-select-option>
                                <a-select-option :value="2">未绑定</a-select-option>
                            </a-select>
                        </div>
                        <div>
                            <a-button type="primary" class="br12" @click="searchBind">
                                <template #icon>
                                    <SearchOutlined />
                                </template>
                                搜索
                            </a-button>
                            <a-button type="primary" class="br12 bg-tipc bn mlr18" @click="reset()">重置</a-button>
                            <a-button type="primary" class="btn-add br12" @click="visibleHland(0)">
                                <template #icon>
                                    <plus-outlined />
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
                                <template v-if="column.dataIndex === 'bind_service_id'">
                                    <div v-if="record.bind_service_id">
                                        <span>管理账号：{{ record.service_manager_member }}</span>
                                        <span class="t-tipc mlr10">>></span>
                                        <span>激活码：{{ record.username }}</span>
                                    </div>
                                    <span class="t-ash" v-else>暂未绑定</span>
                                </template>
                                <template v-if="column.dataIndex === 'status'">
                                    <a-tag color="processing" v-if="record.status == 'enable'">启用</a-tag>
                                    <a-tag color="warning" v-else-if="record.status == 'un_enable'">禁用</a-tag>
                                    <a-tag color="default" v-else>下架</a-tag>
                                </template>
                                <template v-else-if="column.dataIndex === 'operation'">
                                    <a-popconfirm title="是否解绑该条域名？" @confirm="relieve(record.id)" v-if="record.bind_service_id">
                                        <a-button type="primary" size="small" class="br6 btn-green mr14">解绑</a-button>
                                    </a-popconfirm>
                                    <a-button type="primary" size="small" :disabled="true" class="br6 btn-green mr14" v-else>解绑</a-button>
                                    <a-button type="primary" size="small" class="br6" @click="visibleHland(1, record)">修改</a-button>
                                    <a-popconfirm title="是否删除该条域名？" @confirm="deletHland(record.id)">
                                        <a-button type="primary" size="small" class="btn-red br6 ml14">删除</a-button>
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
        <!-- 添加修改 -->
        <a-modal v-model:visible="visible" :centered="true" :closable="false" :footer="null" :title="visible_type ? '修改' : '添加'">
            <a-form :model="visibleForm" layout="vertical" name="basic" autocomplete="off" @finish="onFinish">
                <a-form-item name="status" :rules="[{ required: true, message: ' ' }]" v-if="visible_type">
                    <a-radio-group v-model:value="visibleForm.status">
                        <a-radio value="enable">启用</a-radio>
                        <a-radio value="un_enable">禁用</a-radio>
                    </a-radio-group>
                </a-form-item>
                <a-form-item name="domain" :rules="[{ required: true, message: ' ' }]">
                    <a-textarea v-model:value="visibleForm.domain" :placeholder="visible_type ? '请输入域名' :'批量添加域名用回车分隔'" :rows="6" />
                </a-form-item>
                <a-form-item class="mb0">
                    <div class="tc">
                        <a-button class="br12" @click="visible = false">取消</a-button>
                        <a-button type="primary" html-type="submit" class="br12 ml26">保存</a-button>
                    </div>
                </a-form-item>
            </a-form>
        </a-modal>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { PlusOutlined, SearchOutlined } from '@ant-design/icons-vue';
import { getElementData } from '@/utils';
import { message } from 'ant-design-vue';
import axios from "@/utils/axios";
const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '域名', dataIndex: 'domain', width: 180 },
    { title: '绑定指向', dataIndex: 'bind_service_id', width: 300 },
    { title: '域名状态', dataIndex: 'status', width: 150 },
    { title: '创建时间', dataIndex: 'create_time', width: 180, ellipsis: true },
    { title: '操作', dataIndex: 'operation', width: 220 },
];
export default {
    name: "independentDomainName",
    components: {
        PlusOutlined, SearchOutlined
    },
    setup() {
        const state = reactive({
            pagination: { page: 1, offset: 9 }, //分页
            tabDetail:{
                total: 0, //总条数
                tabHeight: 0, //表格高度
                columns, //表格配置
                data: [], //表格数据
            },
            visible: false, //弹窗
            visible_type: 0,//0添加 1修改
            visibleForm: {},
        });

        onMounted(() => {
            getList(); //获取列表
            state.tabDetail.tabHeight = getElementData('main', 'h') - getElementData('top-box', 'h') - 200; //表格高度
        })

        //搜索
        const searchBind = () => {
            state.pagination.page = 1;
            getList(); //获取列表
        }
        //重置
        const reset = () => {
            state.pagination = { page: 1, offset: 9 };
            getList(); //获取列表
        }
        //获取列表
        const getList = async () => {
            state.pagination['type'] =  'private'; //独立域名
            const res = await axios.post("/domain/list", state.pagination, false);
            state.tabDetail.total = res.data.count;
            state.tabDetail.data = res.data.list;
        }
        //分页
        const change = (e) => {
            state.pagination.page = e;
            getList(); //获取列表
        }
        //弹窗
        const visibleHland = (type, record = {}) => {
            state.visible = true;
            state.visible_type = type;
            state.visibleForm = JSON.parse(JSON.stringify(record));
        }
        //保存
        const onFinish = async(val) => {
            if (state.visible_type) {
                const res = await axios.post("/domain/update", state.visibleForm);
            } else {
                val.type = 'private'; //域名类型
                const res = await axios.post("/domain/create", val);
            }
            getList(); //获取列表
            state.visible = false;
        }
        //删除
        const deletHland = (id) => {
            return new Promise(async (resolve) => {
                setTimeout(() => resolve(true), 100);
                const res = await axios.post("/domain/delete", { domain_id: id });
                getList(); //获取列表
            });
            
        }
        //解绑
        const relieve = async(id) => {
            const res = await axios.post("/domain/un_bind", { domain_id: id });
            getList(); //获取列表
        }
        return {
            ...toRefs(state),
            searchBind, //搜索
            reset, //重置
            visibleHland, //弹窗
            onFinish, //保存
            deletHland, //删除
            relieve, //解绑
            change, //分页
        };
    },
};
</script>

<style lang="less" scoped>

</style>