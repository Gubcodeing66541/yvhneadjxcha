<template>
    <div class="content-box">
        <div class="main bg-white">
            <!-- <div class="lh20">
                <span class="f20 t-them">余额：4567.2</span>
            </div> -->
            <div class="flex align-center">
                <div class="left layout mt54 tel-mt30">
                    <a-form layout="vertical" :model="form" name="basic" autocomplete="off" @finish="onFinish" :wrapper-col="{ span: 16 }">
                        <a-form-item label="自定义充值额度" name="money" :rules="[{ required: true, message: ' ' }]">
                            <a-input type="number" v-model:value="form.money" placeholder="请输入充值额度" />
                        </a-form-item>
                    
                        <a-form-item label="支付方式" name="payment" :rules="[{ required: true, message: ' ' }]">
                            <div class="payment-box flex">
                                <p class="item flex align-center just-center" :class="form.payment == 1?'active':''" @click="form.payment = 1">
                                    <i class="ico icon_zfb"></i>
                                </p>
                                <p class="item flex align-center just-center ml32" :class="form.payment == 2?'active':''" @click="form.payment = 2">
                                    <i class="ico icon_wx"></i>
                                </p>
                            </div>
                        </a-form-item>

                        <a-form-item>
                            <!-- <span>温馨提示：遇到问题请联系我们，QQ：2344455555</span> -->
                        </a-form-item>
                    
                        <a-form-item class="mb0">
                            <a-button type="primary" html-type="submit" class="btn f16">确认支付</a-button>
                        </a-form-item>
                    </a-form>
                </div>
                <img src="../../assets/image/icon_cz.png" class="w44 mr100 pc">
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs } from "vue";
import { notification } from 'ant-design-vue';
import { message } from 'ant-design-vue';
import axios from "@/utils/axios";
export default {
    name: "recharge",
    setup() {
        const state = reactive({
            form: {
                money: '',
                payment: 1,
            },
        });

        //确认支付
        const onFinish = (values) => {
            return message.info("暂无通道")
            notification['success']({
                message: '温馨提示',
                placement: 'bottomRight',
                description: "【xx客服系统】尊敬的用户：您的账户3456787999已于2022年10月10日22: 23 在xx客服系统成功充值666元。感谢您的信赖！",
                onClick: () => {
                    console.log('Notification Clicked!');
                },
            });
        };

        return {
            ...toRefs(state),
            onFinish,
        };
    },
};
</script>

<style lang="less" scoped>
:deep(.ant-form-vertical .ant-form-item-label>label){
    font-size: 16px;
}
.ant-input{
    height: 61px;
    border-radius: 6px;
    border: 1px solid #D6D6D6;
}
.payment-box{
    .item{
        width: 188px;
        height: 76px;
        border-radius: 6px;
        border: 1px solid #D6D6D6;
        cursor: pointer;
    }
    .active{
        border: 2px solid #6971F8;
    }
}
.btn{
    width: 327px;
    height: 53px;
    background: #6971F8;
    border-radius: 12px;
    margin-top: 100px;
}
</style>