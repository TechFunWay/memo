<template>
  <!-- 主弹窗：精简为一屏展示（mobile-ui 规范 2026-10-01） -->
  <n-modal
    :show="support.show"
    preset="card"
    class="support-modal-card w-[88vw] max-w-[320px]"
    :bordered="false"
    :auto-focus="false"
    @update:show="onMainShowChange"
  >
    <div class="text-center">
      <div class="w-11 h-11 mx-auto mb-2 rounded-xl bg-amber-100 dark:bg-amber-500/20 flex items-center justify-center">
        <span class="text-xl">☕</span>
      </div>
      <h3 class="text-base font-bold text-gray-900 dark:text-gray-100 mb-1.5">请作者喝杯咖啡</h3>

      <p class="text-xs text-gray-500 dark:text-gray-400 leading-relaxed mb-3">
        应用免费无广告、数据保存在本机。若对你有帮助，欢迎赞赏——<strong class="text-gray-900 dark:text-gray-100">金额随意，1 元也是心意</strong>。
      </p>

      <div class="flex justify-center my-2">
        <div>
          <img
            :src="donateQr"
            alt="微信赞赏码"
            width="352"
            height="480"
            decoding="async"
            class="mx-auto w-[150px] h-auto max-w-full object-contain rounded-lg border border-gray-200 dark:border-gray-700 bg-white p-1"
          />
          <span class="block mt-1 text-[11px] text-gray-500 dark:text-gray-400">微信扫码赞赏</span>
        </div>
      </div>

      <p v-if="support.errorText" class="text-xs text-red-500 mt-1.5">{{ support.errorText }}</p>

      <div class="flex justify-center gap-2.5 mt-3.5">
        <n-button size="small" :disabled="support.sending" @click="support.dismiss()">暂不支持</n-button>
        <n-button size="small" type="primary" :disabled="support.sending" :loading="support.sending" @click="openAmount">
          ❤ {{ support.sending ? '发送中…' : '已支持' }}
        </n-button>
      </div>
    </div>
  </n-modal>

  <!-- 金额输入弹窗（点击「已支持」后出现） -->
  <n-modal
    :show="amountDialog"
    preset="card"
    class="w-[88vw] max-w-[300px]"
    :bordered="false"
    :auto-focus="false"
    :title="''"
    @update:show="closeAmount"
  >
    <div class="text-center">
      <h4 class="text-sm font-bold text-gray-900 dark:text-gray-100">填写赞赏金额</h4>
      <p class="mt-1 text-[11px] text-gray-500 dark:text-gray-400">仅用于统计，不做支付核验</p>
      <div class="mt-2.5 flex flex-wrap justify-center gap-1.5">
        <n-button
          v-for="p in presets"
          :key="p"
          size="tiny"
          round
          :type="amount === p ? 'primary' : 'default'"
          :secondary="amount !== p"
          @click="selectPreset(p)"
        >
          {{ p }} 元
        </n-button>
      </div>
      <div class="mt-2.5 flex items-center justify-center gap-1.5">
        <span class="text-xs text-gray-500 dark:text-gray-400">或输入</span>
        <n-input-number
          v-model:value="customAmount"
          :min="0"
          :step="0.01"
          placeholder="金额（元）"
          class="w-28"
          size="tiny"
          @update:value="onCustomInput"
        />
      </div>
      <p v-if="amount === 0" class="mt-1.5 text-[11px] text-amber-500">未填写将以 0 元上报</p>
      <p v-if="support.errorText" class="mt-1.5 text-xs text-red-500">{{ support.errorText }}</p>
      <div class="mt-3 flex justify-center gap-2">
        <n-button size="small" :disabled="support.sending" @click="closeAmount">返回</n-button>
        <n-button size="small" type="primary" :disabled="support.sending" :loading="support.sending" @click="confirmAmount">
          {{ support.sending ? '发送中…' : '确定' }}
        </n-button>
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { NModal, NButton, NInputNumber } from 'naive-ui'
import { useSupportStore } from '@/stores/support'
import donateQr from '@/assets/donate-wechat.png'

const support = useSupportStore()

// 金额输入弹窗状态
const amountDialog = ref(false)
const presets = [1, 5, 10, 50]
const customAmount = ref<number | null>(null)
const amount = ref<number>(0)

// 主弹窗关闭时（含关闭按钮/遮罩/发送成功）一并复位金额弹窗
watch(
  () => support.show,
  (visible) => {
    if (!visible) resetAmount()
  },
)

function onMainShowChange(visible: boolean) {
  if (!visible) support.dismiss()
}

function resetAmount() {
  amountDialog.value = false
  customAmount.value = null
  amount.value = 0
}

function openAmount() {
  resetAmount()
  support.errorText = ''
  amountDialog.value = true
}

function selectPreset(p: number) {
  amount.value = p
  customAmount.value = p
}

function onCustomInput(v: number | null) {
  amount.value = Number(v) || 0
}

function closeAmount() {
  if (support.sending) return
  resetAmount()
}

function confirmAmount() {
  if (support.sending) return
  void support.confirmSupported(amount.value || 0)
}
</script>
