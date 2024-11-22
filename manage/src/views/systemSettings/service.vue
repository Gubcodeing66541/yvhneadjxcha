<template>
    <div class="content-box">
        <div id='main' class="main bg-white">
            
            <a-form :model="visibleForm" layout="vertical" autocomplete="off" :wrapper-col="{ span: 10 }" style="height:100% !important">
                <a-tabs v-model:activeKey="activeKey" style="height: calc(100% - 60px);">
                    <a-tab-pane key="1" tab="系统配置">
                        <div class="item mt20">
                            <a-form-item label="系统名称">
                                <a-input v-model:value="visibleForm.system_name" />
                            </a-form-item>
                            <a-form-item label="系统LOGO">
                                <a-input v-model:value="visibleForm.system_logo" />
                            </a-form-item>
                            <a-form-item label="软件版本">
                                <a-input v-model:value="visibleForm.system_version" />
                            </a-form-item>
                            <a-form-item label="电脑端下载地址">
                                <a-input v-model:value="visibleForm.pc_download_link" />
                            </a-form-item>
                            <a-form-item label="操作视频地址">
                                <a-input v-model:value="visibleForm.system_video" />
                            </a-form-item>
                            <a-form-item label="启用客服后台实名认证">
                                <a-switch  :checked="visibleForm.real_name==1?true:false" checked-children="开" un-checked-children="关" @change="switchCheck" />
                            </a-form-item>
                        </div>
                    </a-tab-pane>
                    <a-tab-pane key="2" tab="对象储存">
                        <div class="item mt20">
                            <a-form-item label="储存使用" class="mb16">
                                <a-radio-group v-model:value="visibleForm.oss_storage">
                                    <a-radio value="tencent">腾讯云</a-radio>
                                    <a-radio value="tencent2">腾讯云备用</a-radio>
                                </a-radio-group>
                            </a-form-item>
                            <div class="car-box flex flex-warp">
                                <div class="item mr20" style="width:400px">
                                    <a-card title="腾讯云">
                                        <a-form-item label="appid" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.oss_tencent_appid" />
                                        </a-form-item>
                                        <a-form-item label="域名" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.oss_tencent_domain" />
                                        </a-form-item>
                                        <a-form-item label="SecretId" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.oss_tencent_secret_id" />
                                        </a-form-item>
                                        <a-form-item label="SecretKey" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.oss_tencent_secret_key" />
                                        </a-form-item>
                                        <a-form-item label="区域" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.oss_tencent_region" />
                                        </a-form-item>
                                        <a-form-item label="储存桶" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.oss_tencent_storage" />
                                        </a-form-item>
                                    </a-card>
                                </div>
                                <div class="item" style="width:400px">
                                    <a-card title="腾讯云备用">
                                        <a-form-item label="appid" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.oss_tencent2_appid" />
                                        </a-form-item>
                                        <a-form-item label="域名" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.oss_tencent2_domain" />
                                        </a-form-item>
                                        <a-form-item label="SecretId" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.oss_tencent2_secret_id" />
                                        </a-form-item>
                                        <a-form-item label="SecretKey" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.oss_tencent2_secret_key" />
                                        </a-form-item>
                                        <a-form-item label="区域" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.oss_tencent2_region" />
                                        </a-form-item>
                                        <a-form-item label="储存桶" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.oss_tencent2_storage" />
                                        </a-form-item>
                                    </a-card>
                                </div>
                            </div>
                        </div>
                    </a-tab-pane>
                    <a-tab-pane key="3" tab="支付配置">
                        <div class="item mt20">
                            <a-form-item label="激活码价格">
                                <a-input v-model:value="visibleForm.pay" type="number" style="width:400px" />
                            </a-form-item>
                            <div class="car-box flex flex-warp">
                                <div class="item mr20" style="width:400px">
                                    <a-card title="微信支付">                                        
                                        <a-form-item label="商户ID(mchID)" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.pay_wechat_mch_id" />
                                        </a-form-item>
                                        <a-form-item label="appid" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.pay_wechat_appid" />
                                        </a-form-item>
                                        <a-form-item label="密钥(appKey)" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.pay_wechat_app_key" />
                                        </a-form-item>
                                    </a-card>
                                </div>
                                <div class="item" style="width:400px">
                                    <a-card title="支付宝支付">
                                        <a-form-item label="商户ID(mchID)" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.pay_ali_mch_id" />
                                        </a-form-item>
                                        <a-form-item label="appid" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.pay_ali_appid" />
                                        </a-form-item>
                                        <a-form-item label="密钥(appKey)" :wrapper-col="{ span: 24 }">
                                            <a-input v-model:value="visibleForm.pay_ali_app_key" />
                                        </a-form-item>
                                    </a-card>
                                </div>
                            </div>
                        </div>
                    </a-tab-pane>
                    <a-tab-pane key="4" tab="隐私协议">
                        <wangeditor :detail="visibleForm.ad_default" @editorHtml="editorHtml" class="item mt20"></wangeditor>
                    </a-tab-pane>
                </a-tabs>
                
                <div class="layout tr mt30">
                    <a-button type="primary" @click="onFinish" class="plr30 br12 ml26">保存</a-button>
                    <a-button type="primary" class="br12 bg-tipc bn ml26 plr30" @click="reset()">重置</a-button>
                </div>
            </a-form>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, computed } from "vue";
import wangeditor from '../../components/wangeditor.vue'
import axios from "@/utils/axios";
import { useStore } from "vuex";
export default {
    name: "service",
    components: { wangeditor },
    setup() {
        const config = computed(() => store.state.config) //系统配置信息
        const store = useStore();
        const state = reactive({
            activeKey:'1', //菜单索引
            visibleForm: {},
        });

        onMounted(() => {
            state.visibleForm = JSON.parse(JSON.stringify(config.value))
        })

        //富文本内容
        const editorHtml = (e) => {
            state.visibleForm['ad_default'] = e;
        }

        //实名开关
        const switchCheck = async (status) => {
            state.visibleForm.real_name = status ? 1 : 2;
        }

        //保存
        const onFinish = async () => {
            state.visibleForm.pay = parseInt(state.visibleForm.pay)
            const res = await axios.post("/config/update", state.visibleForm);
            store.dispatch("asyncSetInfo"); //获取配置信息
        }

        //重置
        const reset = () => {
            state.visibleForm = {}
        }

        return {
            ...toRefs(state),
            editorHtml, //富文本内容
            switchCheck, //实名开关
            onFinish, //保存
            reset, //重置
        };
    },
};
</script>

<style lang="less" scoped>
:deep(.ant-tabs-content){
    height: 100%;
}
.item{
    height: calc(100% - 40px);
    overflow: overlay;
}
</style>