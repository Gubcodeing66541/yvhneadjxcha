<template>
    <div class="content">
        <div class="content-box flex dir-column">
            <div class="box flex tel-box">                
                <div class="sma-box bg-white t-3s mr20">
                    <a-statistic title="平台充值金额 (￥)" :precision="2" :value="data.count.renew_service_manager" />
                </div>
                <div class="sma-box bg-white t-3s mr20">
                    <a-statistic title="平台总余额 (￥)" :precision="2" :value="data.count.all_account" />
                </div>
                <div class="sma-box bg-white t-3s mr20">
                    <a-statistic title="平台总消费 (￥)" :precision="2" :value="data.count.all_pay" />
                </div>
                <div class="sma-box bg-white t-3s">
                    <a-statistic title="今日总消费 (￥)" :precision="2" :value="data.count.today_pay" />
                </div>                
            </div>
            <div class="flex h100 just-between mt20 tel-block">
                <div class="shrink w24 mr20 flex dir-column tel-flex tel-layout">                    
                    <div class="sma-box br40 flex dir-column just-center bg-white t-3s h50 mb20">
                        <p><i class="ico icon_service"></i></p>
                        <p class="f20 mt20">在线客服数</p>
                        <p class="f30 mt10">{{ data.socket.service_count }}</p>
                    </div>
                    <div class="sma-box flex dir-column just-center bg-white t-3s h50">
                        <p><i class="ico icon_users"></i></p>
                        <p class="f20 mt20">在线用户数</p>
                        <p class="f30 mt10">{{ data.socket.user_count }}</p>
                    </div>
                </div>
                <div id="tabBox" class="sma-box bg-white hidden">
                    <div id="title" class="f18 strong ptb30">时时在线客服</div>
                    <div class="hidden" :style="{ height: tabDetail.tabHeight + 50 + 'px' }">
                        <a-table :columns="tabDetail.columns" :data-source="data.onLineService" :scroll="{ y: tabDetail.tabHeight }" :pagination=false >
                            <template #bodyCell="{ column, record }">
                                <template v-if="column.dataIndex === 'status'">
                                    <span v-if="record.status" class="t-green">在线</span>
                                </template>
                            </template>
                        </a-table>
                    </div>
                    <div id="pagination" class="tr pt24 plr20">
                        <a-pagination size="small" :total="tabDetail.total" v-model:current="pagination.page" v-model:page-size="pagination.offset" show-quick-jumper @change="change" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, onBeforeUnmount } from "vue";
import { getElementData } from "@/utils";
import axios from "@/utils/axios";
const columns = [
    { title: '序号', dataIndex: 'id', align:'center', width: 100 },    
    { title: '客服名称', dataIndex: 'name', align: 'center', width: 180 },
    { title: '客服账号', dataIndex: 'username', align: 'center', width: 220 },
    { title: '管理账号', dataIndex: 'member', align: 'center', width: 220 },
    { title: '状态', dataIndex: 'status', align: 'center', width: 100 },
];
export default {
    name: "record",
    setup() {
        const state = reactive({       
            timer: null, //定时器     
            tabDetail: { //客服表格配置
                total: 0, //总条数
                columns: columns, //表格配置     
                tabHeight:0, //表格高度                
            },
            pagination: { page: 1, offset: 7 }, //客服列表分页
            data: {//数据
                count: {},
                socket:{},
                onLineService: [], //表格数据
            }, 
        });

        onMounted(() => {
            init();            
            state.tabDetail.tabHeight = getElementData('tabBox', 'h') - 220; //表格高度
        });

        onBeforeUnmount(() => {
            clearInterval(state.timer);
            state.timer = null;
        })

        //统计数据
        const init = async () => {
            getList();
            const res = await axios.post("/count", '', false);
            state.data = res.data;
            state.timer = setInterval(async () => {
                getList();
                const res = await axios.post("/count", '', false);
                state.data.count = res.data.count;
                state.data.socket = res.data.socket;
            }, 3000);
        };

        //获取在线客服列表
        const getList = async () => {
            let list = JSON.stringify(state.data.onLineService);
            const res = await axios.post("/count_service_list", state.pagination, false);
            state.tabDetail.total = res.data.count;            
            if (list !== JSON.stringify(res.data.services)){
                state.data.onLineService = res.data.services;
            }            
        }

        //分页
        const change = (e) => {
            state.pagination.page = e;
            getList(); //获取列表
        }

        return {
            ...toRefs(state),
            change, // 分页获取在线客服
        };
    },
};
</script>

<style lang="less" scoped>
.content-box {
    padding-bottom: 40px;
    .sma-box {
        text-align: center;
        border-radius: 20px;
        box-shadow: 0px 3px 6px 1px rgba(211, 229, 233, 0.2);
        &:hover {
            box-shadow: 0px 3px 6px 1px rgba(211, 229, 233, 0.5);
        }
    }
    .box{
        .sma-box{
            width: 25%;
            height: 100px;
            display: flex;
            align-items: center;
            justify-content: center;
        }
    }
}
</style>