<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2024/04/01 15:52:40
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <CommonPage>
    <template #action>
      <NButton type="primary" @click="handleAdd()">
        <i class="i-material-symbols:add mr-4 text-18" />
        Add Role
      </NButton>
    </template>

    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :scroll-x="1200"
      :columns="columns"
      :get-data="api.read"
    >
      <MeQueryItem label="Role Name" :label-width="50">
        <n-input v-model:value="queryItems.name" type="text" placeholder="Please enter role name" clearable />
      </MeQueryItem>
  <MeQueryItem label="Status" :label-width="50">
        <n-select
          v-model:value="queryItems.enable"
          clearable
          :options="[
            { label: 'Enabled', value: 1 },
            { label: 'Disabled', value: 0 },
          ]"
        />
      </MeQueryItem>
    </MeCrud>
    <MeModal ref="modalRef" width="520px">
      <n-form
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
      >
        <n-form-item
          label="Role Name"
          path="name"
          :rule="{
            required: true,
            message: 'Please enter role name',
            trigger: ['input', 'blur'],
          }"
        >
          <n-input v-model:value="modalForm.name" />
        </n-form-item>
        <n-form-item
          label="Role Code"
          path="code"
          :rule="{
            required: true,
            message: 'Please enter role code',
            trigger: ['input', 'blur'],
          }"
        >
          <n-input v-model:value="modalForm.code" :disabled="modalAction !== 'add'" />
        </n-form-item>
  <n-form-item label="Permissions" path="permissionIds">
          <n-tree
            key-field="id"
            label-field="name"
            :selectable="false"
            :data="permissionTree"
            :checked-keys="modalForm.permissionIds"
            :on-update:checked-keys="(keys) => (modalForm.permissionIds = keys)"

            checkable check-on-click default-expand-all
            class="cus-scroll max-h-200 w-full"
          />
        </n-form-item>
        <n-form-item label="Status" path="enable">
          <NSwitch v-model:value="modalForm.enable">
            <template #checked>
              Enabled
            </template>
            <template #unchecked>
              Disabled
            </template>
          </NSwitch>
        </n-form-item>
      </n-form>
    </MeModal>
  </CommonPage>
</template>

<script setup>
import { NButton, NSwitch } from 'naive-ui'
import { MeCrud, MeModal, MeQueryItem } from '@/components'
import { useCrud } from '@/composables'
import api from './api'

defineOptions({ name: 'RoleMgt' })

const router = useRouter()

const $table = ref(null)
/** QueryBar filter parameters (optional) */
const queryItems = ref({})

onMounted(() => {
  $table.value?.handleSearch()
})

const { modalRef, modalFormRef, modalAction, modalForm, handleAdd, handleDelete, handleEdit }
  = useCrud({
    name: 'Role',
    doCreate: api.create,
    doDelete: api.delete,
    doUpdate: api.update,
    initForm: { enable: true },
    refresh: (_, keepCurrentPage) => $table.value?.handleSearch(keepCurrentPage),
  })

const columns = [
  { title: 'Role Name', key: 'name' },
  { title: 'Role Code', key: 'code' },
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
          disabled: row.code === 'SUPER_ADMIN',
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
            secondary: true,
            onClick: () =>
              router.push({ path: `/pms/role/user/${row.id}`, query: { roleName: row.name } }),
          },
          {
            default: () => 'Assign Users',
            icon: () => h('i', { class: 'i-fe:user-plus text-14' }),
          },
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-left: 12px;',
            disabled: row.code === 'SUPER_ADMIN',
            onClick: () => handleEdit(row),
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
            disabled: row.code === 'SUPER_ADMIN',
            onClick: () => handleDelete(row.id),
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

async function handleEnable(row) {
  row.enableLoading = true
  try {
    await api.update({ id: row.id, enable: !row.enable })
    row.enableLoading = false
  $message.success('Operation successful')
    $table.value?.handleSearch()
  }
  catch (error) {
    console.error(error)
    row.enableLoading = false
  }
}

const permissionTree = ref([])
api.getAllPermissionTree().then(({ data = [] }) => (permissionTree.value = data))
</script>
