<template>
    <div class="content-box">
        <div id='main' class="main bg-white">
            <div id='top-box' class="top-box mb36">
                <div class="input-box flex flex-warp">
                    <div class="item mr20">
                        <span>敏感词：</span>
                        <a-input v-model:value="pagination.search" placeholder="请输入关键词" class="input" style="width: 300px" />
                    </div>
                    <div>
                        <a-button type="primary" class="br12" @click="searchBind">
                            <template #icon>
                                <SearchOutlined />
                            </template>搜索
                        </a-button>
                        <a-button type="primary" class="br12 bg-tipc bn mlr18" @click="reset()">重置</a-button>
                        <a-button type="primary" class="btn-add br12" @click="visibleHland(0)">
                            <template #icon>
                                <plus-outlined />
                            </template>添加
                        </a-button>
                    </div>
                </div>
            </div>
            <div class="tab-box">
                <div :style="{ height: tabDetail.tabHeight + 50 + 'px' }">
                    <a-table :columns="tabDetail.columns" :data-source="tabDetail.data"
                        :scroll="{ y: tabDetail.tabHeight }" :pagination=false>
                        <template #bodyCell="{ column, record }">
                            <template v-if="column.dataIndex === 'operation'">
                                <a-button type="primary" size="small" class="br6" @click="visibleHland(1, record)">修改
                                </a-button>
                                <a-popconfirm title="是否删除该条公告？" @confirm="deletHland(record.id)">
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
        <!-- 添加修改 -->
        <a-modal v-model:visible="visible" :centered="true" :closable="false" :footer="null"
            :title="visible_type ? '修改敏感词' : '添加敏感词'">
            <a-form :model="visibleForm" layout="vertical" autocomplete="off" @finish="onFinish">
                <a-form-item name="value" label="敏感词" :rules="[{ required: true, message: ' ' }]">
                    <a-textarea v-model:value="visibleForm.value" placeholder="请输入" :rows="6" />
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
import { message } from 'ant-design-vue';
import { getElementData } from '@/utils';
import axios from "@/utils/axios";
const columns = [
    { title: '序号', dataIndex: 'id', width: 120 },
    { title: '敏感词', dataIndex: 'value', ellipsis: true, width: 500 },
    { title: '创建时间', dataIndex: 'create_time', width: 200, ellipsis: true },
    { title: '操作', dataIndex: 'operation', width: 200 },
];
export default {
    name: "sensitiveWords",
    components: { PlusOutlined, SearchOutlined },
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
            state.pagination['type'] = 'keyword'; //敏感词
            const res = await axios.post("/setting/list", state.pagination, false);
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
        const onFinish = async (val) => {
            if (state.visible_type) {
                const res = await axios.post("/setting/update", state.visibleForm);
            } else {
                val.type = 'keyword'; //类型
                const res = await axios.post("/setting/create", val);
            }
            getList(); //获取列表
            state.visible = false;
        }
        //删除
        const deletHland = (id) => {
            return new Promise(async (resolve) => {
                setTimeout(() => resolve(true), 100);
                const res = await axios.post("/setting/delete", { id: id });
                getList(); //获取列表
            });

        }

        return {
            ...toRefs(state),
            change, //分页
            searchBind, //搜索
            reset, //重置
            visibleHland, //弹窗
            onFinish, //保存
            deletHland, //删除
        };
    },
};
</script>

<style lang="less" scoped>

</style>