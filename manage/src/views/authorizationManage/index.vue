<template>
    <div class="content">
        <div class="content-box">
            <div id='main' class="main bg-white">
                <div id='top-box' class="top-box mb36">
                    <div class="input-box flex flex-warp">
                        <div class="item mr20">
                            <span>授权类型：</span>
                            <a-select ref="select" v-model:value="pagination.type" class="select" style="width: 120px">
                                <a-select-option value="auth">登录</a-select-option>
                                <a-select-option value="push">推送</a-select-option>
                            </a-select>
                        </div>
                        <div class="item mr20">
                            <span>公众号名称：</span>
                            <a-input v-model:value="pagination.search" class="input mb10 tel-w74" style="width: 170px" />
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
                                <template v-if="column.dataIndex === 'type'">
                                    <a-tag color="green" v-if="record.type == 'auth'">登录</a-tag>
                                    <a-tag color="blue" v-else>推送</a-tag>
                                </template>
                                <template v-if="column.dataIndex === 'status'">
                                    <a-switch :checked="record.status =='enable'?true:false" checked-children="开" un-checked-children="关" @change="switchCheck($event, record.id)" />
                                </template>
                                <template v-else-if="column.dataIndex === 'operation'">
                                    <a-button type="primary" size="small" class="br6 btn-green" @click="switchUrl(record.id, record.url, record.url_spare)">切换URL</a-button>
                                    <a-button type="primary" size="small" class="br6 ml14" @click="visibleHland(1, record)">修改</a-button>
                                    <a-popconfirm title="是否删除该公众号？" @confirm="deletHland(record.id)">
                                        <a-button type="primary" size="small" class="btn-red br6 ml14">删除</a-button>
                                    </a-popconfirm>
                                </template>
                            </template>
                        </a-table>
                    </div>
                    <div class="tr pt24">
                        <a-pagination size="small" :total="tabDetail.total" v-model:current="pagination.page" v-model:page-size="pagination.offset" show-quick-jumper @change="change" />
                    </div>
                </div>
            </div>
        </div>
        <!-- 添加修改 -->
        <a-modal v-model:visible="visible" :centered="true" :closable="false" :footer="null" :title="visible_type ? '修改公众号' : '添加公众号'">
            <a-form :model="visibleForm" autocomplete="off" :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }" @finish="onFinish">
                <a-form-item label="授权类型" name="type" :rules="[{ required: true, message: ' ' }]">
                    <a-radio-group v-model:value="visibleForm.type">                        
                        <a-radio value="push">推送</a-radio>
                        <a-radio value="auth">登录</a-radio>
                    </a-radio-group>
                </a-form-item>
                <a-form-item label="公众号名称" name="name" :rules="[{ required: true, message: ' ' }]">
                    <a-input v-model:value="visibleForm.name" />
                </a-form-item>
                <a-form-item label="APP_ID" name="app_id" :rules="[{ required: true, message: ' ' }]">
                    <a-input v-model:value="visibleForm.app_id" />
                </a-form-item>
                <a-form-item label="APP_SECRET" name="app_secret" :rules="[{ required: true, message: ' ' }]">
                    <a-input v-model:value="visibleForm.app_secret" />
                </a-form-item>
                <a-form-item label="URL链接" name="url" :rules="[{ required: true, message: ' ' }]">
                    <a-input v-model:value="visibleForm.url" />
                </a-form-item>
                <a-form-item label="URL链接2" name="url_spare">
                    <a-input v-model:value="visibleForm.url_spare" />
                </a-form-item>
                <a-form-item label="消息模板ID" name="message_id" :rules="[{ required: true, message: ' ' }]" v-if="visibleForm.type =='push'">
                    <a-input v-model:value="visibleForm.message_id" />
                </a-form-item>
                <a-form-item label="文件名称" name="file_name" :rules="[{ required: true, message: ' ' }]">
                    <a-input v-model:value="visibleForm.file_name" />
                </a-form-item>
                <a-form-item label="文件内容" name="file_value" :rules="[{ required: true, message: ' ' }]">
                    <a-input v-model:value="visibleForm.file_value" />
                </a-form-item>
                <div class="tc mt30">
                    <a-button class="br12" @click="visible = false">取消</a-button>
                    <a-button type="primary" html-type="submit" class="br12 ml26">确定</a-button>
                </div>
            </a-form>
        </a-modal>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { PlusOutlined, SearchOutlined } from '@ant-design/icons-vue';
import { getElementData } from '@/utils';
import axios from "@/utils/axios";
const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '公众号名称', dataIndex: 'name', width: 200 },
    { title: 'APP_ID', dataIndex: 'app_id', width: 200 },
    { title: '当前使用URL', dataIndex: 'url', width: 200 },
    { title: '公众号类型', dataIndex: 'type', width: 150 },
    { title: '状态', dataIndex: 'status', width: 120 },
    { title: '创建时间', dataIndex: 'create_time', width: 160, ellipsis: true },
    { title: '操作', dataIndex: 'operation', width: 250 },
];
export default {
    name: "authorizationManage",
    components: {
        PlusOutlined, SearchOutlined
    },
    setup() {
        const state = reactive({
            pagination: { page: 1, offset: 9 }, //分页
            tabDetail: {
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
            getList();
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
            const res = await axios.post("/wechat_auths/list", state.pagination, false);
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
            state.visibleForm = record;
        }
        //保存
        const onFinish = async (val) => {
            if (state.visible_type) {
                const res = await axios.post("/wechat_auths/update", state.visibleForm);
            } else {
                const res = await axios.post("/wechat_auths/create", val);
            }
            getList(); //获取列表
            state.visible = false;
        }
        //修改状态
        const switchCheck = async (status,id) => {
            const res = await axios.post("/wechat_auths/enable_disable", { id: id, status: status ? 'enable' :'un_enable'});
            getList(); //获取列表
        }
        //切换URL
        const switchUrl = async (id, url, url_spare) => {
            const res = await axios.post("/wechat_auths/switch", { id, url_spare: url, url: url_spare });
            getList(); //获取列表
        }
        //删除
        const deletHland = (id) => {
            return new Promise(async (resolve) => {
                setTimeout(() => resolve(true), 100);
                const res = await axios.post("/wechat_auths/delete", { id: id });
                getList(); //获取列表
            });
        }
        return {
            ...toRefs(state),
            searchBind, //搜索
            reset, //重置
            visibleHland, //弹窗
            onFinish, //保存
            deletHland, //删除
            switchUrl, //切换URL            
            switchCheck, //修改状态
            change, //分页
        };
    },
};
</script>

<style lang="less" scoped>

</style>