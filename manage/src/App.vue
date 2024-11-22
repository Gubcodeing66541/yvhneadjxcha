<template>
  <a-config-provider :locale="locale">
    <router-view />
  </a-config-provider>
</template>

<script>
import { defineComponent, reactive, toRefs } from "vue";
import "dayjs/locale/zh-cn";
import zhCN from 'ant-design-vue/lib/locale-provider/zh_CN';

export default defineComponent({
  name: "App",
  setup() {
    const state = reactive({
      locale: zhCN,
    });
    return {
      ...toRefs(state),
    };
  },
});
</script>
<style lang="less">
/********** 内容页主体盒子 **********/
.content {
  height: calc(100vh - 140px);

  .content-box {
    height: 100%;

    .main {
      box-shadow: 0px 3px 6px 1px rgba(211, 229, 233, 0.2);
      border-radius: 42px;
      height: calc(100% - 40px);
      padding: 30px 35px;
      overflow: overlay;
    }
  }
}

.popover-pointer {
  .item {
    cursor: pointer;
    padding: 10px 0;
    border-bottom: 1px solid #eee;

    &:hover,
    span:hover {
      color: @primary-color;
    }

    &:last-child {
      border: none;
    }
  }
}

/********** 弹窗统一样式 **********/
.ant-modal-content {
  border-radius: 20px;
  overflow: hidden;

  .ant-modal-header {
    padding-top: 30px;
    border-bottom: none;

    .ant-modal-title {
      text-align: center;
      font-weight: bold;
    }
  }

  .ant-modal-footer {
    text-align: center;
    border: none;

    .ant-btn {
      border-radius: 12px;
    }
  }
}
.ant-modal-confirm .ant-modal-confirm-btns .ant-btn {
  border-radius: 12px;
}
/********** 输入框统一样式 **********/
.ant-input {
  border-radius: 8px;
}

.input {
  height: 35px;
  border-radius: 8px;
}

.ant-picker-range {
  border-radius: 8px;
}

.ant-select:not(.ant-select-customize-input) .ant-select-selector {
  border-radius: 8px;
}
/********** 表格统一样式 **********/
.tab-box {
  .ant-table-wrapper {
    border: 1px solid #EEEEEE;

    tr {
      &:last-child {
        td {
          border-bottom: none;
        }
      }
    }
  }
}

/******* 组件按钮 ******/
.btn-add {
  background: #31DAC3;
  border-color: #31DAC3;

  &:hover,
  &:focus {
    background: #70e2d3;
    border-color: #70e2d3;
  }
}

.btn-red {
  background: #FF6E6E;
  border-color: #FF6E6E;

  &:hover,
  &:focus {
    background: #fc8989;
    border-color: #fc8989;
  }
}

.btn-green {
  background: #039B84;
  border-color: #039B84;

  &:hover,
  &:focus {
    background: #039B84;
    border-color: #039B84;
  }
}

.btn-tipc{
    background: #bbb;
    border-color: #bbb;
  
    &:hover,
    &:focus {
      background: #bbb;
      border-color: #bbb;
    }
}

.btn-open{
  background: #20bf1e;
  border-color: #20bf1e;

  &:hover{
    background: #20bf1e;
    border-color: #20bf1e;
  }
}
</style>