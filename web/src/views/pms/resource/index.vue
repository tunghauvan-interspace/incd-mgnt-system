<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/05 21:28:53
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <CommonPage>
    <div class="flex">
      <n-spin size="small" :show="treeLoading">
        <MenuTree
          v-model:current-menu="currentMenu"
          class="w-320 shrink-0"
          :tree-data="treeData"
          @refresh="initData"
        />
      </n-spin>

      <div class="ml-40 w-0 flex-1">
        <template v-if="currentMenu">
          <div class="flex justify-between">
            <h3 class="mb-12">
              {{ currentMenu.name }}
            </h3>
            <NButton size="small" type="primary" @click="handleEdit(currentMenu)">
              <i class="i-material-symbols:edit-outline mr-4 text-14" />
              Edit
            </NButton>
          </div>
          <n-descriptions label-placement="left" bordered :column="2">
            <n-descriptions-item label="Code">
              {{ currentMenu.code }}
            </n-descriptions-item>
            <n-descriptions-item label="Name">
              {{ currentMenu.name }}
            </n-descriptions-item>
            <n-descriptions-item label="Route Path">
              {{ currentMenu.path ?? '--' }}
            </n-descriptions-item>
            <n-descriptions-item label="Component Path">
              {{ currentMenu.component ?? '--' }}
            </n-descriptions-item>
            <n-descriptions-item label="Menu Icon">
              <span v-if="currentMenu.icon" class="flex items-center">
                <i :class="`${currentMenu.icon}?mask text-22 mr-8`" />
                <span class="opacity-50">{{ currentMenu.icon }}</span>
              </span>
              <span v-else>None</span>
            </n-descriptions-item>
            <n-descriptions-item label="Layout">
              {{ currentMenu.layout || 'Follow system' }}
            </n-descriptions-item>
            <n-descriptions-item label="Visible">
              {{ currentMenu.show ? 'Yes' : 'No' }}
            </n-descriptions-item>
            <n-descriptions-item label="Enabled">
              {{ currentMenu.enable ? 'Yes' : 'No' }}
            </n-descriptions-item>
            <n-descriptions-item label="KeepAlive">
              {{ currentMenu.keepAlive ? 'Yes' : 'No' }}
            </n-descriptions-item>
            <n-descriptions-item label="Order">
              {{ currentMenu.order ?? '--' }}
            </n-descriptions-item>
          </n-descriptions>

          <div class="mt-32 flex justify-between">
            <h3 class="mb-12">
                Buttons
              </h3>
            <NButton size="small" type="primary" @click="handleAddBtn">
              <i class="i-fe:plus mr-4 text-14" />
              Add
            </NButton>
          </div>

          <MeCrud
            ref="$table"
            :columns="btnsColumns"
            :scroll-x="-1"
            :get-data="api.getButtons"
            :query-items="{ parentId: currentMenu.id }"
          />
        </template>
  <n-empty v-else class="h-450 f-c-c" size="large" description="Please select a menu to view details" />
      </div>
    </div>
    <ResAddOrEdit ref="modalRef" :menus="treeData" @refresh="initData" />
  </CommonPage>
</template>

<script setup>
import { NButton, NSwitch } from 'naive-ui'
import { MeCrud } from '@/components'
import api from './api'
import MenuTree from './components/MenuTree.vue'
import ResAddOrEdit from './components/ResAddOrEdit.vue'

const treeData = ref([])
const treeLoading = ref(false)
const $table = ref(null)
const currentMenu = ref(null)
async function initData(data) {
  if (data?.type === 'BUTTON') {
    $table.value.handleSearch()
    return
  }
  treeLoading.value = true
  const res = await api.getMenuTree()
  treeData.value = res?.data || []
  treeLoading.value = false

  if (data)
    currentMenu.value = data
}
initData()

const modalRef = ref(null)
  function handleEdit(item = {}) {
  modalRef.value?.handleOpen({
    action: 'edit',
    title: `Edit Menu - ${item.name}`,
    row: item,
    okText: 'Save',
  })
}

  const btnsColumns = [
  { title: 'Name', key: 'name' },
  { title: 'Code', key: 'code' },
  {
  title: 'Status',
    key: 'enable',
    render: row =>
      h(
        NSwitch,
        {
          size: 'small',
          rubberBand: false,
          value: row.enable,
          loading: !!row.enableLoading,
          onUpdateValue: () => handleEnable(row),
        },
        {
          checked: () => 'Enabled',
          unchecked: () => 'Disabled',
        },
      ),
  },
  {
  title: 'Actions',
    key: 'actions',
    width: 320,
    align: 'right',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-left: 12px;',
            onClick: () => handleEditBtn(row),
          },
          {
            default: () => 'Edit',
            icon: () => h('i', { class: 'i-material-symbols:edit-outline text-14' }),
          },
        ),

        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            style: 'margin-left: 12px;',
            onClick: () => handleDeleteBtn(row.id),
          },
          {
            default: () => 'Delete',
            icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
          },
        ),
      ]
    },
  },
]

watch(
  () => currentMenu.value,
  async (v) => {
    await nextTick()
    if (v)
      $table.value.handleSearch()
  },
)

  function handleAddBtn() {
  modalRef.value?.handleOpen({
    action: 'add',
    title: 'Add Button',
    row: { type: 'BUTTON', parentId: currentMenu.value.id },
    okText: 'Save',
  })
}

  function handleEditBtn(row) {
  modalRef.value?.handleOpen({
    action: 'edit',
    title: `Edit Button - ${row.name}`,
    row,
    okText: 'Save',
  })
}

function handleDeleteBtn(id) {
  const d = $dialog.warning({
  content: 'Confirm delete?',
  title: 'Warning',
  positiveText: 'Confirm',
  negativeText: 'Cancel',
    async onPositiveClick() {
      try {
        d.loading = true
  await api.deletePermission(id)
  $message.success('Deleted successfully')
        $table.value.handleSearch()
        d.loading = false
      }
      catch (error) {
        console.error(error)
        d.loading = false
      }
    },
  })
}

async function handleEnable(item) {
  try {
    item.enableLoading = true
    await api.savePermission(item.id, {
      enable: !item.enable,
    })
  $message.success('Operation successful')
    $table.value?.handleSearch()
    item.enableLoading = false
  }
  catch (error) {
    console.error(error)
    item.enableLoading = false
  }
}
</script>
