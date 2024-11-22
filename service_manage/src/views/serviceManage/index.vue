<template>
    <div class="content">
        <div class="content-box">
            <div id='main' class="main bg-white">
                <div id='top-box' class="top-box mb36">
                    <div class="input-box flex flex-warp">
                        <div class="item mr20">
                            <span>客服名称：</span>
                            <a-input v-model:value="pagination.search" placeholder="请输入" style="width: 170px" />
                        </div>
                        <div class="item mr20">
                            <span>客服账号：</span>
                            <a-input v-model:value="pagination.username" placeholder="请输入" style="width: 170px" />
                        </div>
                        <div class="item mr18">
                            <span>创建时间：</span>
                            <a-range-picker v-model:value="visit_time" format="YYYY-MM-DD" style="width:240px;border-radius: 8px;" @change="handleChange" />
                        </div>
                        <div class="btn-cont">
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
                            <a-button type="primary" class="btn-green br12 ml18" @click="model.batchRenewalVisible = true">
                                <template #icon>
                                    <FileAddOutlined />
                                </template>
                                批量续费
                            </a-button>
                        </div>
                    </div>
                </div>
                <div class="tab-box">
                    <div :style="{ height: tabDetail.tabHeight + 50 + 'px' }">
                        <a-table :columns="tabDetail.columns" :data-source="tabDetail.data" :scroll="{ y: tabDetail.tabHeight }" :pagination=false>
                            <template #bodyCell="{ column, record }">
                                <template v-if="column.dataIndex === 'is_online'">
                                    <span v-if="record.is_online" class="t-green">在线</span>
                                    <span v-else class="t-ash">离线</span>
                                </template>
                                <template v-else-if="column.dataIndex === 'activate_time'">
                                    <span v-if="record.is_activate">{{ record.activate_time }}</span>
                                    <span v-else class="t-ash">暂未激活</span>
                                </template>
                                <template v-else-if="column.dataIndex === 'time_out'">
                                    <span v-if="record.is_activate">{{ record.time_out }}</span>
                                    <span v-else class="t-ash">-</span>
                                </template>
                                <template v-else-if="column.dataIndex === 'status'">
                                    <span v-if="record.status =='time_out'" class="t-red">已过期</span>
                                    <span v-else>正常</span>
                                </template>
                                <template v-else-if="column.dataIndex === 'operation'">
                                    <a-button type="primary" size="small" class="br6" @click="editVis(record.service_id, record.name)">编辑</a-button>
                                    <a-popconfirm title="是否删除该激活码？" @confirm="deletHland(record.service_id)">
                                        <a-button type="primary" size="small" class="br6 btn-red ml8">删除</a-button>
                                    </a-popconfirm>                                    
                                    <a-button type="primary" size="small" class="br6 btn-green ml8" @click="renewVisibleForm.username = record.username, model.renewVisible=true">续费</a-button>
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

        <!-- 新增客服 -->
        <a-modal v-model:visible="model.addVisible" :centered="true" :closable="false" :footer="null" title="新增客服">
            <a-form :model="addVisibleForm" layout="vertical" autocomplete="off" @finish="onFinish">
                <a-form-item label="新增类型" required>
                    <a-radio-group v-model:value="addVisibleForm.type" name="radioGroup">
                        <a-radio :value="0">单个新增</a-radio>
                        <a-radio :value="1">批量新增</a-radio>
                    </a-radio-group>
                </a-form-item>        
                <a-form-item label="客服名称" required v-if="!addVisibleForm.type">
                    <a-input v-model:value="addVisibleForm.service_name" placeholder="请输入客服名称" />
                </a-form-item>
                <a-form-item label="批量生成客服数" required v-else>
                    <a-input v-model:value="addVisibleForm.service_number" placeholder="请输入新增客服个数" />
                </a-form-item>
                <a-form-item label="续费天数" required>
                    <a-input type="number" v-model:value="addVisibleForm.day" placeholder="请输入天数" />
                </a-form-item>        
                <a-form-item class="mb0">
                    <div class="tc">
                        <a-button class="br12" @click="model.addVisible = false">取消</a-button>
                        <a-button type="primary" html-type="submit" class="br12 ml26">保存</a-button>
                    </div>
                </a-form-item>
            </a-form>
        </a-modal>
        <!-- 编辑客服 -->
        <a-modal v-model:visible="model.editVisible" :centered="true" :closable="false" :footer="null" title="编辑客服">
            <a-form :model="editVisibleForm" layout="vertical" autocomplete="off" @finish="editFinish">
                <a-form-item name="name" label="客服名称" :rules="[{ required: true, message: ' ' }]">
                    <a-input v-model:value="editVisibleForm.name" placeholder="请输入客服名称" />
                </a-form-item>
                <a-form-item class="mb0">
                    <div class="tc">
                        <a-button class="br12" @click="model.editVisible = false">取消</a-button>
                        <a-button type="primary" html-type="submit" class="br12 ml26">保存</a-button>
                    </div>
                </a-form-item>
            </a-form>
        </a-modal>
        <!-- 续费 -->
        <a-modal v-model:visible="model.renewVisible" :centered="true" :closable="false" :footer="null" title="客服续费">
            <a-form :model="renewVisibleForm" layout="vertical" autocomplete="off" @finish="renewFinish">
                <a-form-item name="day" label="续费天数" :rules="[{ required: true, message: ' ' }]">
                    <a-input type="number" v-model:value="renewVisibleForm.day" placeholder="请输入续费天数" />
                </a-form-item>
                <a-form-item class="mb0">
                    <div class="tc">
                        <a-button class="br12" @click="model.renewVisible = false">取消</a-button>
                        <a-button type="primary" html-type="submit" class="br12 ml26">保存</a-button>
                    </div>
                </a-form-item>
            </a-form>
        </a-modal>
        <!-- 批量续费 -->
        <a-modal v-model:visible="model.batchRenewalVisible" :centered="true" :closable="false" :footer="null" title="批量续费">
            <a-form :model="batchRenewalVisibleForm" layout="vertical" name="basic" autocomplete="off" @finish="batchRenewalFinish">
                <a-form-item name="day" label="续费天数" :rules="[{ required: true, message: ' ' }]">
                    <a-input type="number" v-model:value="batchRenewalVisibleForm.day" placeholder="请输入续费天数" />
                </a-form-item>
                <a-form-item name="username_list" label="授权码" :rules="[{ required: true, message: ' ' }]">
                    <a-textarea v-model:value="batchRenewalVisibleForm.username_list" placeholder="多个授权码请用回车分割" :rows="6" />
                </a-form-item>
                <a-form-item class="mb0">
                    <div class="tc">
                        <a-button class="br12" @click="model.batchRenewalVisible = false">取消</a-button>
                        <a-button type="primary" html-type="submit" class="br12 ml26">保存</a-button>
                    </div>
                </a-form-item>
            </a-form>
        </a-modal>
        <!-- 新增客服成功弹窗 -->
        <a-modal v-model:visible="model.addSuccessVisible" :centered="true" :closable="false" :maskClosable="false" :footer="null" title="新增成功" :width="876">
            <div class="code-box flex just-center flex-warp">
                <p class="item br8 hh40 lh40 t-them tc strong ml6 mr6 mb12 text1 pl6 pr6" v-for="item in addService" :key="item" >{{item}}</p>
            </div>
            <div class="tc mt40">
                <a-button type="primary" class="br12" @click="copy">一键复制</a-button>
            </div>
        </a-modal>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { SearchOutlined, PlusOutlined, FileAddOutlined, InfoCircleOutlined } from '@ant-design/icons-vue';
import { getElementData } from '@/utils';
import useClipboard from "vue-clipboard3";
import { message } from 'ant-design-vue';
import axios from "@/utils/axios";
import { useStore } from "vuex";
const columns = [
    { title: '序号', dataIndex: 'id', width: 70 },
    { title: '客服名称', dataIndex: 'name', width: 110, ellipsis: true },
    { title: '创建时间', dataIndex: 'create_time', width: 160, ellipsis: true },
    { title: '激活时间', dataIndex: 'activate_time', width: 160, ellipsis: true },
    { title: '过期时间', dataIndex: 'time_out', width: 160, ellipsis: true },
    { title: '通讯录', dataIndex: 'user_cnt', width: 90 },
    { title: '在线状态', dataIndex: 'is_online', width: 90 },
    { title: '账号状态', dataIndex: 'status', width: 90 },
    { title: '激活码', dataIndex: 'username', width: 220 },
    { title: '操作', dataIndex: 'operation', width: 200 },
];
export default {
    name: "serviceManage",
    components: { SearchOutlined, PlusOutlined, FileAddOutlined, InfoCircleOutlined },
    setup() {
        const { toClipboard } = useClipboard();
        const store = useStore();
        const state = reactive({
            visit_time:[], //时间
            pagination: { page: 1, offset: 9 }, //分页
            tabDetail: {
                total: 0, //总条数
                tabHeight: 0, //表格高度
                columns:columns, //表格配置
                data: [], //表格数据
            },
            model:{}, //弹窗
            addVisibleForm: { //新增客服表单
                type: 0,
                service_name:'',
                service_number:'',
                day: 1,
            },
            addService:[], //新增的激活码
            renewVisibleForm: {}, //续费
            batchRenewalVisibleForm: {}, //批量续费
            editVisibleForm: {}, //编辑客服表单
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
            const res = await axios.post("/member/list", state.pagination, false);
            state.tabDetail.total = res.data.count;
            state.tabDetail.data = res.data.list;
        }

        //分页
        const change = (e) => {
            state.pagination.page = e;
            getList(); //获取列表
        }

        //新增客服
        const onFinish = async() => {
            if (!state.addVisibleForm.type && state.addVisibleForm.service_name == '' || state.addVisibleForm.type && state.addVisibleForm.service_number == '') return message.error("请填写完整后提交")
            if (state.addVisibleForm.type){
                const res = await axios.post("/member/create_list", { service_number: parseInt(state.addVisibleForm.service_number), day: parseInt(state.addVisibleForm.day) });
                state.addService = res.data.member_list;
            }else{
                const res = await axios.post("/member/create", { service_name: state.addVisibleForm.service_name, day: parseInt(state.addVisibleForm.day) });
                state.addService = [res.data.member];
            }            
            state.addVisibleForm = { type: 0, service_name: '', service_number: '',  day: 1 }
            state.model['addVisible'] = false;
            state.model['addSuccessVisible'] = true;
            getList(); //获取列表
            store.dispatch("asyncSetInfo");
        }

        //续费
        const renewFinish = async () => {
            state.renewVisibleForm['day'] = parseInt(state.renewVisibleForm['day']);
            const res = await axios.post("/member/renewal", state.renewVisibleForm);
            getList(); //获取列表
            store.dispatch("asyncSetInfo");
            state.model['renewVisible'] = false;
            state.renewVisibleForm['day'] = ''
        }

        //批量续费 
        const batchRenewalFinish = async () => {
            state.batchRenewalVisibleForm['day'] = parseInt(state.batchRenewalVisibleForm['day']);
            const res = await axios.post("/member/renewal_all", state.batchRenewalVisibleForm);
            getList(); //获取列表
            store.dispatch("asyncSetInfo");
            state.model['batchRenewalVisible'] = false;
            state.batchRenewalVisibleForm['day'] = ''
        }
        
        //编辑客服
        const editVis = (service_id, name) => {
            state.editVisibleForm['name'] = name;
            state.editVisibleForm['service_id'] = service_id;           
            state.model['editVisible'] = true;
        }

        //编辑提交
        const editFinish = async () => {
            const res = await axios.post("/member/update", state.editVisibleForm);
            getList(); //获取列表
            state.model['editVisible'] = false;
        }

        //删除
        const deletHland = (service_id) => {
            return new Promise(async (resolve) => {
                setTimeout(() => resolve(true), 100);
                const res = await axios.post("/member/delete", { service_id });
                getList(); //获取列表
            });
        }

        //复制
        const copy = async() => {
            let addService = state.addService.map(item => item);
            try {
                await toClipboard(addService.join('\r'));
                state.model['addSuccessVisible'] = false;
                message.success("复制成功")
            } catch (e) {
                message.error("复制失败")
            }
        }

        return {
            ...toRefs(state),
            change, //分页
            handleChange, //客服名称选择
            searchBind, //搜索
            reset, //重置
            onFinish, //新增客服
            editVis, //编辑
            editFinish, //编辑提交
            renewFinish, //续费
            batchRenewalFinish, //批量续费
            deletHland, //删除
            copy, //复制
        };
    },
};
</script>

<style lang="less" scoped>
.code-box{
    .item{
        width: 126px;

        border: 1px solid #6971F8;
    }
}
</style>