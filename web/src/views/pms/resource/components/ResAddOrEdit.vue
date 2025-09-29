<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2024/04/01 15:52:31
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <MeModal ref="modalRef">
    <n-form
      ref="modalFormRef"
      label-placement="left"
      require-mark-placement="left"
      :label-width="100"
      :model="modalForm"
    >
      <n-grid :cols="24" :x-gap="24">
  <n-form-item-gi :span="12" label="Parent Menu" path="parentId">
          <n-tree-select
            v-model:value="modalForm.parentId"
            :options="menuOptions"
            :disabled="parentIdDisabled"
            label-field="name"
            key-field="id"
            placeholder="Root menu"
            clearable
          />
        </n-form-item-gi>
        <n-form-item-gi :span="12" path="name" :rule="required">
          <template #label>
            <QuestionLabel label="Name" content="Title" />
          </template>
          <n-input v-model:value="modalForm.name" />
        </n-form-item-gi>
        <n-form-item-gi :span="12" path="code" :rule="required">
          <template #label>
            <QuestionLabel label="Code" content="If menu, corresponds to front-end route name (PascalCase)" />
          </template>
          <n-input v-model:value="modalForm.code" />
        </n-form-item-gi>
        <n-form-item-gi
          v-if="modalForm.type === 'MENU'"
          :span="12"
          path="path"
          :rule="{
            trigger: ['blur', 'change'],
            type: 'string',
            message: 'Must start with /, http or https',
            validator(rule, value) {
              if (value) {
                return /\/|http|https/.test(value)
              }
              return true
            },
          }"
        >
          <template #label>
            <QuestionLabel label="Route path" content="Optional for parent menus" />
          </template>
          <n-input v-model:value="modalForm.path" />
        </n-form-item-gi>
        <n-form-item-gi v-if="modalForm.type === 'MENU'" :span="12" path="icon">
          <template #label>
            <QuestionLabel
              label="Menu icon"
              content="e.g. material-symbols:help, icon library: https://icones.js.org/collection/all"
            />
          </template>
          <n-select v-model:value="modalForm.icon" :options="iconOptions" clearable filterable />
        </n-form-item-gi>
        <n-form-item-gi v-if="modalForm.type === 'MENU'" :span="12" path="layout">
          <template #label>
            <QuestionLabel
              label="Layout"
              content="Corresponds to a folder in layouts; default is 'default'"
            />
          </template>
          <n-select v-model:value="modalForm.layout" :options="layoutOptions" clearable />
        </n-form-item-gi>
        <n-form-item-gi v-if="modalForm.type === 'MENU'" :span="24" path="component">
          <template #label>
            <QuestionLabel
              label="Component path"
              content="Front-end component path, starts with /src; optional for parent menus"
            />
          </template>
          <n-select
            v-model:value="modalForm.component"
            :options="componentOptions"
            clearable
            filterable
            tag
          />
        </n-form-item-gi>

        <n-form-item-gi v-if="modalForm.type === 'MENU'" :span="12" path="show">
          <template #label>
            <QuestionLabel label="Display" content="Control whether shown in menu; does not affect route registration" />
          </template>
          <n-switch v-model:value="modalForm.show">
            <template #checked>
              Show
            </template>
            <template #unchecked>
              Hide
            </template>
          </n-switch>
        </n-form-item-gi>
        <n-form-item-gi :span="12" path="enable">
          <template #label>
            <QuestionLabel
              label="Status"
              content="If menu, disabling will remove it from route table and make the page inaccessible"
            />
          </template>
          <n-switch v-model:value="modalForm.enable">
            <template #checked>
              Enabled
            </template>
            <template #unchecked>
              Disabled
            </template>
          </n-switch>
        </n-form-item-gi>
        <n-form-item-gi v-if="modalForm.type === 'MENU'" :span="12" path="keepAlive">
          <template #label>
            <QuestionLabel
              label="KeepAlive"
              content="When enabling keepAlive, set the component's name to this menu's code"
            />
          </template>
          <n-switch v-model:value="modalForm.keepAlive">
            <template #checked>
              Yes
            </template>
            <template #unchecked>
              No
            </template>
          </n-switch>
        </n-form-item-gi>
        <n-form-item-gi
          v-if="modalForm.type === 'MENU'"
          :span="12"
          label="Sort Order"
          path="order"
          :rule="{
            type: 'number',
            required: true,
            message: 'This field is required',
            trigger: ['blur', 'change'],
          }"
        >
          <n-input-number v-model:value="modalForm.order" />
        </n-form-item-gi>
        
      </n-grid>
    </n-form>
  </MeModal>
</template>

<script setup>
import icons from 'isme:icons'
import pagePathes from 'isme:page-pathes'
import { MeModal } from '@/components'
import { useForm, useModal } from '@/composables'
import api from '../api'
import QuestionLabel from './QuestionLabel.vue'

const props = defineProps({
  menus: {
    type: Array,
    required: true,
  },
})
const emit = defineEmits(['refresh'])

const menuOptions = computed(() => {
  return [{ name: 'Root menu', id: '', children: props.menus || [] }]
})
const componentOptions = pagePathes.map(path => ({ label: path, value: path }))
const iconOptions = icons.map(item => ({
  label: () =>
    h('span', { class: 'flex items-center' }, [h('i', { class: `${item} text-18 mr-8` }), item]),
  value: item,
}))
const layoutOptions = [
  { label: 'Follow system', value: '' },
  { label: 'Simple', value: 'simple' },
  { label: 'Normal', value: 'normal' },
  { label: 'Full', value: 'full' },
  { label: 'Empty', value: 'empty' },
]
const required = {
  required: true,
  message: 'This field is required',
  trigger: ['blur', 'change'],
}

const defaultForm = { enable: true, show: true, layout: '' }
const [modalFormRef, modalForm, validation] = useForm()
const [modalRef, okLoading] = useModal()

const modalAction = ref('')
const parentIdDisabled = ref(false)
function handleOpen(options = {}) {
  const { action, row = {}, ...rest } = options
  modalAction.value = action
  modalForm.value = { ...defaultForm, ...row }
  parentIdDisabled.value = !!row.parentId && row.type === 'BUTTON'
  modalRef.value.open({ ...rest, onOk: onSave })
}

async function onSave() {
  await validation()
  okLoading.value = true
  try {
    let newFormData
    if (!modalForm.value.parentId)
      modalForm.value.parentId = null
    if (modalAction.value === 'add') {
      const res = await api.addPermission(modalForm.value)
      newFormData = res.data
    }
    else if (modalAction.value === 'edit') {
      await api.savePermission(modalForm.value.id, modalForm.value)
    }
    okLoading.value = false
  $message.success('Saved successfully')
    emit('refresh', modalAction.value === 'add' ? newFormData : modalForm.value)
  }
  catch (error) {
    console.error(error)
    okLoading.value = false
    return false
  }
}

defineExpose({
  handleOpen,
})
</script>
