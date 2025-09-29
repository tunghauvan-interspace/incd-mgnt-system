<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2024/01/13 17:41:47
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <CommonPage show-footer>
    <n-button type="primary" @click="openModal1">
      Open first modal
    </n-button>
    <MeModal ref="$modal1">
      <n-input v-model:value="text" />
    </MeModal>
  <MeModal ref="$modal2" title="Content submitted by previous modal">
      <h2>{{ text }}</h2>
    </MeModal>
  </CommonPage>
</template>

<script setup>
import { MeModal } from '@/components'
import { useModal } from '@/composables'
import { sleep } from '@/utils'

const text = ref('')
const [$modal1, okLoading1] = useModal()
function openModal1() {
  $modal1.value?.open({
  title: 'First modal',
    width: '600px',
  okText: 'Open another',
  cancelText: 'Close',
    async onOk() {
      if (!text.value) {
        $message.warning('Please enter content')
        return false // Prevent the modal from closing
      }
      okLoading1.value = true
  $message.loading('Submitting...', { key: 'modal1' })
      await sleep(1000)
      okLoading1.value = false
  $message.success('Submitted', { key: 'modal1' })
      openModal2()
  return false // Default behavior is to close the modal; returning false prevents it from closing
    },
    onCancel(message) {
      $message.info(message ?? 'Cancelled')
    },
  })
}

const [$modal2, okLoading2] = useModal()
function openModal2() {
  $modal2.value?.open({
  cancelText: 'Close this',
  okText: 'Close all modals',
    width: '400px',
    async onOk() {
      okLoading2.value = true
  $message.loading('Closing...', { key: 'modal2' })
      await sleep(1000)
      okLoading2.value = false

  // close modal1 as well
  $modal1.value?.close()
  $message.success('Closed', { key: 'modal2' })
    },
  })
}
</script>
