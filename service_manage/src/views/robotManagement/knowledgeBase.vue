<template>
    <div class="content">
        <div class="content-box">
            <div id='main' class="main bg-white">
                <div id='top-box' class="top-box mb36">
                    <a-button type="primary" class="btn-add br12" @click="visibleHland(0)">
                        <template #icon>
                            <plus-outlined />
                        </template>
                        添加
                    </a-button>
                </div>
                <div class="tab-box">
                    <div :style="{ height: tabDetail.tabHeight + 50 + 'px' }">
                        <a-table :columns="tabDetail.columns" :data-source="tabDetail.data" :scroll="{ y: tabDetail.tabHeight }" :pagination=false>
                            <template #bodyCell="{ column, record }">
                                <template v-if="column.dataIndex === 'status'">
                                    <span v-if="record.status == 0" class="t-green">已通过</span>
                                    <span v-else-if="record.status == 1" class="t-orange">待审核</span>
                                    <span v-else class="t-red">未通过</span>
                                </template>
                                <template v-else-if="column.dataIndex === 'operation'">
                                    <a-button type="primary" size="small" class="br6" @click="visibleHland(1, record)">修改</a-button>                                    
                                    <a-popconfirm title="是否删除该条问题？" @confirm="deletHland(record.id)">
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
        <a-modal v-model:visible="visible" :centered="true" :closable="false" :footer="null" :title="visible_type ? '修改智能回复' : '添加智能回复'">
            <a-form :model="visibleForm" layout="vertical" name="basic" autocomplete="off" @finish="onFinish">
                <a-form-item name="problem" label="问题" :rules="[{ required: true, message: ' ' }]">
                    <a-input v-model:value="visibleForm.problem" placeholder="请输入问题" />
                </a-form-item>

                <a-form-item name="answer" label="话术内容" :rules="[{ required: true, message: ' ' }]">
                    <a-textarea v-model:value="visibleForm.answer" placeholder="请输入回复内容" :rows="6" />
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
import { PlusOutlined } from '@ant-design/icons-vue';
import { getElementData } from '@/utils';
import axios from "@/utils/axios";
const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '问题', dataIndex: 'problem', width: 180, ellipsis: true },
    { title: '回复内容', dataIndex: 'answer', width: 180, ellipsis: true },
    { title: '创建时间', dataIndex: 'create_time', width: 140 },
    { title: '操作', dataIndex: 'operation', width: 150 },
];
export default {
    name: "recharge",
    components: {
        PlusOutlined,
    },
    setup() {
        const state = reactive({
            pagination: { page: 1, offset: 9 }, //分页
            tabDetail: {
                total: 0, //总条数
                tabHeight: 0, //表格高度
                columns: columns, //表格配置
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

        //获取列表
        const getList = async () => {
            const res = await axios.post("/bot/list", state.pagination, false);
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
                const res = await axios.post("/bot/update", state.visibleForm);
            } else {
                const res = await axios.post("/bot/add", val);
            }
            getList(); //获取列表
            state.visible = false;
        }
        //删除
        const deletHland = (id) => {
            return new Promise(async (resolve) => {
                setTimeout(() => resolve(true), 100);
                const res = await axios.post("/bot/delete", { id: id });
                getList(); //获取列表
            });
        }
        return {
            ...toRefs(state),
            change, //分页
            visibleHland, //弹窗
            onFinish, //保存
            deletHland, //删除
        };
    },
};
</script>

<style lang="less" scoped>

</style>