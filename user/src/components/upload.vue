<template>
    <div class="upload flex align-center flex-warp">
        <van-uploader accept="image/*" :max-count="1" :before-read="beforeUpload" name="jpg" :after-read="afterRead">
            <div class="ico-box br16 bg-white tc ww70 hh70 lh70"><i class="iconfont icon-tianjiatupian f30"></i></div>
            <p class="tc mt12">照片</p>
        </van-uploader>
        <van-uploader accept="video/*" :max-count="1" class="ml30" name="mp4" :after-read="afterRead">
            <div class="ico-box br16 bg-white tc ww70 hh70 lh70"><i class="iconfont icon-shipin f30"></i></div>
            <p class="tc mt12">视频</p>
        </van-uploader>
    </div>
</template>

<script>
import { reactive, toRefs } from "vue";
import { beforeUpload, videoBeforeUpload, uploadFile } from '@/utils';
import axios from "@/utils/axios";
import { playMp3 } from "@/utils";
import { readyState } from "@/socket";
import { Toast } from 'vant';
export default {
    name: "upload",
    setup() {
        const state = reactive({});

        const afterRead = async(files, detail) => {
            const socketState = readyState();
            if (socketState) {
                const res = await uploadFile(files.file, detail.name);
                await axios.post('/auth/send', { type: detail.name == 'jpg' ? 'image' : 'video', content: res }, false);
                await playMp3(); //消息音频播放
            } else {
                tipcSocket();
            }
        }

        //发送消息时检测用户是否连接websocket
        const tipcSocket = () => {
            Toast.loading({
                message: '连接超时正在重新接入...',
                forbidClick: true,
            });
            setTimeout(() => {
                window.location.reload();
            }, 1500);
        }

        return {
            ...toRefs(state),
            beforeUpload,
            videoBeforeUpload,
            afterRead,
        };
    },
};
</script>

<style lang="less" scoped>

</style>