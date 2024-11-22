<template>
    <div class="user-menu layout flex pd20">
        <div class="left shrink pc">           
            <examplePhone :type="'userMenu'"></examplePhone>
        </div>
        <div class="right ml20">
            <div class="right-top mb20">
                <a-button type="primary" @click="form = {}, formType = 'add', formVisible = true">+ 添加</a-button>
            </div>
            <div class="tab-box" :style="{ height: tabDetail.tabHeight + 50 + 'px'}">
                <a-table :columns="tabDetail.columns" :data-source="tabDetail.data" :scroll="{ y: tabDetail.tabHeight }" :pagination=false>
                    <template #bodyCell="{ column,record }">
                        <template v-if="column.dataIndex === 'action'">
                            <a-tag color="cyan" v-if="record.action == 'http'">链接跳转</a-tag>
                            <a-tag color="green" v-else-if="record.action == 'copy'">复制</a-tag>
                            <a-tag color="blue" v-else-if="record.action == 'call'">打电话</a-tag>
                            <a-tag color="purple" v-else>回复</a-tag>
                        </template>
                        <template v-else-if="column.dataIndex === 'operation'">
                            <a-button type="primary" size="small" class="br6" @click="form = record, formVisible = true, formType = 'edit'">修改</a-button>
                            <a-popconfirm title="是否删除该快捷菜单？" @confirm="delet(record.id)">
                                <a-button type="primary" size="small" class="btn-red br6 ml10">删除</a-button>
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

    <!-- 用户快捷菜单 -->
    <a-modal v-model:visible="formVisible" :width="550" :centered="true" :footer="null" :mask='false' :title="`${formType == 'add'?'添加':'编辑'}快捷菜单`">
        <div v-if="formVisible">
            <a-form :model="form" name="basic" :label-col="{ span: 3 }" :wrapper-col="{ span: 21 }" autocomplete="off" @finish="onFinish">
                <a-form-item label="标题" name="title" :rules="[{ required: true, message: ' ' }]">
                    <a-input v-model:value="form.title" placeholder="标题最多支持8个字" show-count :maxlength="8" />
                </a-form-item>                
                <a-form-item label="功能" name="action" :rules="[{ required: true, message: ' ' }]">
                    <a-radio-group v-model:value="form.action" button-style="solid">
                        <a-radio-button value="copy">点击复制</a-radio-button>
                        <a-radio-button value="banswer">点击自动回复</a-radio-button>
                        <a-radio-button value="http">点击链接跳转</a-radio-button>
                        <a-radio-button value="call">点击打电话</a-radio-button>
                    </a-radio-group>
                </a-form-item>
                <a-form-item label="内容" name="content" :rules="[{ required: true, message: ' ' }]">
                    <a-input v-model:value="form.content" placeholder="请输入内容" />
                    <p class="t-tipc f12 mt4" v-show="form.action == 'http'">链接要加http://或者https://</p>
                </a-form-item>
                <a-form-item label="排序" name="sort">
                    <a-input type="number" v-model:value="form.sort" placeholder="数字越小越排前" />
                </a-form-item>
                <a-form-item label="提示语" name="tag" v-if="form.action == 'copy'">
                    <a-input v-model:value="form.tag" placeholder="复制成功后的提示，如：复制成功、请妥善保存等" />
                </a-form-item>
                <a-form-item :wrapper-col="{ span: 24 }">
                    <div class="tc">
                        <a-button class="br12" @click="formVisible =false">取消</a-button>
                        <a-button type="primary" html-type="submit" class="br12 ml26">保存</a-button>
                    </div>
                </a-form-item>
            </a-form>
        </div>
    </a-modal>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import examplePhone from "./examplePhone.vue";
import axios from "@/utils/axios";
const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '标题', dataIndex: 'title', width: 160, ellipsis: true },
    { title: '内容', dataIndex: 'content', width: 160, ellipsis: true },
    { title: '功能', dataIndex: 'action', width: 80 },
    { title: '排序', dataIndex: 'sort', width: 80 },
    { title: '操作', dataIndex: 'operation', width: 140 },
];
export default {
    name: "userMenu",
    components: { examplePhone },
    setup() {
        const state = reactive({
            formVisible:false,
            form:{},
            formType:'', //添加或编辑
            pagination: { page: 1, offset: 8 }, //分页
            tabDetail: {
                total: 0, //总条数
                columns: columns, //表格配置
                tabHeight: 500,
                data: [], //表格数据
            },
        });

        onMounted(() => {
            getList();
        })

        //获取列表
        const getList = async () => {
            const res = await axios.post("/menu_setting/list", state.pagination, false);
            state.tabDetail.total = res.data.count;
            state.tabDetail.data = res.data.list;
        }

        //分页
        const change = (e) => {
            state.pagination.page = e;
            getList(); //获取列表
        }

        //提交
        const onFinish = async() => {
            state.form['sort'] = parseInt(state.form['sort']);
            if (state.formType=='add') {
                const res = await axios.post("/menu_setting/create", state.form);
            } else {
                const res = await axios.post("/menu_setting/update", state.form);
            }
            getList(); //获取列表
            state.formVisible = false;
        }

        //删除
        const delet = async(id) => {
            const res = await axios.post("/menu_setting/delete", { id });
            getList(); //获取列表
        }

        return {
            ...toRefs(state),
            change, //分页
            onFinish, //提交表单
            delet, //删除
        };
    },
};
</script>

<style lang="less" scoped>
.left{width: 370px;}
.right{max-width: calc(100% - 390px);}
</style>