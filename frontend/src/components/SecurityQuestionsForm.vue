<template>
  <div class="flex flex-col gap-3">
    <div v-for="(_, i) in SECURITY_QUESTION_COUNT" :key="i" class="flex flex-col gap-1.5">
      <n-select
        v-model:value="questions[i]"
        :options="optionsFor(i)"
        :placeholder="`安全问题 ${i + 1}`"
        size="large"
        filterable
        tag
        clearable
      />
      <n-input
        v-model:value="answers[i]"
        :placeholder="`问题 ${i + 1} 的答案`"
        size="large"
        clearable
        @keyup.enter="emit('enter')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NSelect, NInput } from 'naive-ui'
import { SECURITY_QUESTIONS, SECURITY_QUESTION_COUNT } from '@/constants/security'

const questions = defineModel<string[]>('questions', { required: true })
const answers = defineModel<string[]>('answers', { required: true })

const emit = defineEmits<{ enter: [] }>()

const baseOptions = computed(() =>
  SECURITY_QUESTIONS.map((q) => ({ label: q, value: q }))
)

// 排除其它行已选择的问题，避免重复
function optionsFor(index: number) {
  const chosen = questions.value.filter((q, i) => i !== index && q)
  return baseOptions.value.filter((opt) => !chosen.includes(opt.value))
}
</script>
