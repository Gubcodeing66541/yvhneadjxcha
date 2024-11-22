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
                    <div :style="{ height: tabDetail.tabHeight + 50 + 'px'}">
                        <a-table :columns="tabDetail.columns" :data-source="tabDetail.data" :scroll="{ y: tabDetail.tabHeight }" :pagination=false>
                            <template #bodyCell="{ column, record }">
                                <template v-if="column.dataIndex === 'content'">
                                    <a-image :src="record.content" v-if="record.type == 'image'" style="height:50px" />
                                    <span v-else>{{ record.content }}</span>
                                </template>
                                <template v-else-if="column.dataIndex === 'operation'">
                                    <a-button type="primary" size="small" class="br6" @click="visibleHland(1, { id: record.id, type: record.type, content: record.content})">修改</a-button>
                                    <a-popconfirm title="是否删除该条话术？" @confirm="deletHland(record.id)">
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
        <a-modal v-model:visible="visible" :centered="true" :closable="false" :footer="null" :title="visible_type ? '修改快捷回复' :'添加快捷回复'">
            <a-form :model="visibleForm" layout="vertical" name="basic" autocomplete="off" @finish="onFinish">
                <a-form-item name="content" label="话术内容" :rules="[{ required: true, message: ' ' }]">
                    <div class="bd br8">
                        <a-form-item class="ml11 mr11 pt11 pb4 bb-tipc mb0">
                            <div class="flex align-center just-between pl6 pr6">
                                <a-upload accept=".png, .jpg, .jpeg" :show-upload-list="false" :before-upload="beforeUpload" :customRequest="handleChange">
                                    <i class="ico icon_img ml2 pointer"></i>
                                </a-upload>
                                <span class="f12 t-red pointer" @click="visibleForm.type = 'text', visibleForm.content=''" v-if="visibleForm.type == 'image'">删除图片</span>
                            </div>
                        </a-form-item>
                        <a-textarea v-model:value="visibleForm.content" :bordered="false" placeholder="请输入话术内容" :rows="4" v-if="visibleForm.type=='text'" />
                        <div class="mg10">
                            <a-image :src="visibleForm.content" v-if="visibleForm.type=='image'" />
                        </div>
                    </div>
                </a-form-item>
            
                <a-form-item class="mb0">
                    <div class="tc">
                        <a-button class="br12" @click="visible=false">取消</a-button>
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
import { getElementData, beforeUpload, uploadFile } from '@/utils';
import axios from "@/utils/axios";
const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '话术内容', dataIndex: 'content',width: 180 },
    { title: '创建时间', dataIndex: 'create_time', width: 140, ellipsis: true },
    { title: '更新时间', dataIndex: 'update_time', width: 140, ellipsis: true },
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
            visible:false, //弹窗
            visible_type:0,//0添加 1修改
            visibleForm:{ type:"text", content:"" },
        });

        onMounted(() => {
            getList();
            state.tabDetail.tabHeight = getElementData('main', 'h') - getElementData('top-box', 'h') - 200; //表格高度
        })

        //获取列表
        const getList = async () => {
            const res = await axios.post("/message/list", state.pagination, false);
            state.tabDetail.total = res.data.count;
            state.tabDetail.data = res.data.list;
        }

        //分页
        const change = (e) => {
            state.pagination.page = e;
            getList(); //获取列表
        }

        //弹窗
        const visibleHland = (type, record = { title: ' ', type: "text", content :''}) => {
            state.visible = true;
            state.visible_type = type;
            state.visibleForm = record;
        }
        //上传图片
        const handleChange = async(files) => {
            const res = await uploadFile(files.file);
            state.visibleForm.content = res;
            state.visibleForm.type = 'image';
        };
        //保存
        const onFinish = async () => {
            if (state.visible_type) {
                const res = await axios.post("/message/update", state.visibleForm);
            } else {
                const res = await axios.post("/message/add", state.visibleForm);
            }
            getList(); //获取列表
            state.visible = false;
        }
        //删除
        const deletHland = (id) => {
            return new Promise(async (resolve) => {
                setTimeout(() => resolve(true), 100);
                const res = await axios.post("/message/delete", { id });
                getList(); //获取列表
            });
        }
        return {
            ...toRefs(state),
            change, //分页
            visibleHland, //弹窗
            onFinish, //保存
            deletHland, //删除
            beforeUpload, //上传图片前
            handleChange, //上传图片
        };
    },
};
</script>

<style lang="less" scoped>

</style>