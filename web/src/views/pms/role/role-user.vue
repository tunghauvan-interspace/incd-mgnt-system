<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/05 21:29:43
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <CommonPage back>
    <template #title-suffix>
      <NTag class="ml-12" type="warning">
        {{ route.query.roleName }}
      </NTag>
    </template>
    <template #action>
      <div class="flex items-center">
        <NButton :disabled="!userIds.length" type="error" @click="handleBatchRemove()">
          <i v-if="userIds.length" class="i-material-symbols:delete-outline mr-4 text-18" />
          Batch Remove Authorization
        </NButton>
        <NButton
          class="ml-12"
          :disabled="!userIds.length"
          type="primary"
          @click="handleBatchAdd()"
        >
          <i v-if="userIds.length" class="i-line-md:confirm-circle mr-4 text-18" />
          Batch Authorize
        </NButton>
      </div>
    </template>

    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :scroll-x="1200"
      :columns="columns"
      :get-data="api.getAllUsers"
      @on-checked="onChecked"
    >
      <MeQueryItem label="Username" :label-width="50">
        <n-input
          v-model:value="queryItems.username"
          type="text"
          placeholder="Please enter username"
          clearable
        />
      </MeQueryItem>

      <MeQueryItem label="Gender" :label-width="50">
        <n-select v-model:value="queryItems.gender" clearable :options="genders" />
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
  </CommonPage>
</template>

<script setup>
import { NAvatar, NButton, NSwitch, NTag } from 'naive-ui'
import { h } from 'vue'
import { MeCrud, MeQueryItem } from '@/components'
import { formatDateTime } from '@/utils'
import api from './api'

defineOptions({ name: 'RoleUser' })
const route = useRoute()

const $table = ref(null)
/** QueryBar filter parameters (optional) */
const queryItems = ref({})

onMounted(() => {
  $table.value?.handleSearch()
})

const genders = [
  { label: 'Male', value: 1 },
  { label: 'Female', value: 2 },
]

const columns = [
  { type: 'selection', fixed: 'left' },
  {
    title: 'Avatar',
    key: 'avatar',
    width: 80,
    render: ({ avatar }) =>
      h(NAvatar, {
        size: 'medium',
        src: avatar,
      }),
  },
  { title: 'Username', key: 'username', width: 150, ellipsis: { tooltip: true } },
  {
  title: 'Roles',
    key: 'roles',
    width: 200,
    ellipsis: { tooltip: true },
    render: ({ roles }) => {
      if (roles?.length) {
        return roles.map((item, index) =>
          h(
            NTag,
            { type: 'success', style: index > 0 ? 'margin-left: 8px;' : '' },
            { default: () => item.name },
          ),
        )
      }
      return 'No roles'
    },
  },
  {
    title: 'Gender',
    key: 'gender',
    width: 80,
    render: ({ gender }) => genders.find(item => gender === item.value)?.label ?? '',
  },
  {
    title: 'Created At',
    key: 'createDate',
    width: 180,
    render(row) {
      return h('span', formatDateTime(row.createTime))
    },
  },
  {
  title: 'Status',
    key: 'enable',
    width: 100,

    render: row =>
      h(
        NSwitch,
        {
          size: 'small',
          rubberBand: false,
          value: row.enable,
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
    width: 100,
    align: 'right',
    fixed: 'right',
    hideInExcel: true,
    render(row) {
      return row.roles?.some(item => item.id === +route.params.roleId)
        ? h(
            NButton,
            {
              size: 'small',
              type: 'error',
              secondary: true,
              onClick: () => handleBatchRemove([row.id]),
            },
            {
              default: () => 'Remove Authorization',
              icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
            },
          )
        : h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              secondary: true,
              onClick: () => handleBatchAdd([row.id]),
            },
            {
              default: () => 'Authorize',
              icon: () => h('i', { class: 'i-line-md:confirm-circle text-14' }),
            },
          )
    },
  },
]

const userIds = ref([])
function onChecked(rowKeys) {
  userIds.value = rowKeys || []
}

function handleBatchAdd(ids = userIds.value) {
  const roleId = route.params.roleId
  if (!roleId)
    return $message.error('Role error, please re-select the role')
  if (!ids.length)
    return $message.error('Please select users first')
  $dialog.confirm({
    content: `Confirm assign 【${route.query.roleName}】?`,
    async confirm() {
      await api.addRoleUsers(roleId, { userIds: ids })
      $table.value?.handleSearch()
    },
  })
}
function handleBatchRemove(ids = userIds.value) {
  const roleId = route.params.roleId
  if (!roleId)
    return $message.error('Role error, please re-select the role')
  if (!ids.length)
    return $message.error('Please select users first')
  $dialog.confirm({
    content: `Confirm remove assignment of 【${route.query.roleName}】?`,
    async confirm() {
      await api.removeRoleUsers(roleId, { userIds: ids })
      $table.value?.handleSearch()
    },
  })
}
</script>
