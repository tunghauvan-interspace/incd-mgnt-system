<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/05 21:30:11
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <AppPage show-footer>
    <n-card>
      <n-space align="center">
        <n-avatar round :size="100" :src="userStore.avatar" />
        <div class="ml-20">
          <div class="flex items-center text-16">
            <span>Username:</span>
            <span class="ml-12 opacity-80">{{ userStore.username }}</span>
            <n-button class="ml-32" type="primary" text @click="pwdModalRef.open()">
              <i class="i-fe:edit mr-4" />
              Change Password
            </n-button>
          </div>
          <div class="mt-16 flex items-center">
            <n-button type="primary" ghost @click="avatarModalRef.open()">
              Change Avatar
            </n-button>
            <span class="ml-12 opacity-60">
              Avatar change only supports an online URL. Upload is not provided — integrate your own upload if needed.
            </span>
          </div>
        </div>
      </n-space>
    </n-card>

    <n-card class="mt-20" title="Profile Information">
      <template #header-extra>
        <n-button type="primary" text @click="profileModalRef.open()">
          <i class="i-fe:edit mr-4" />
          Edit Profile
        </n-button>
      </template>

      <n-descriptions
        label-placement="left"
        :label-style="{ width: '200px', textAlign: 'center' }"
        :column="1"
        bordered
      >
        <n-descriptions-item label="Nickname">
          {{ userStore.nickName }}
        </n-descriptions-item>
        <n-descriptions-item label="Gender">
          {{ genders.find((item) => item.value === userStore.userInfo?.gender)?.label ?? 'Unknown' }}
        </n-descriptions-item>
        <n-descriptions-item label="Address">
          {{ userStore.userInfo?.address }}
        </n-descriptions-item>
        <n-descriptions-item label="Email">
          {{ userStore.userInfo?.email }}
        </n-descriptions-item>
      </n-descriptions>
    </n-card>

    <MeModal ref="avatarModalRef" width="420px" title="Change Avatar" @ok="handleAvatarSave()">
      <n-input v-model:value="newAvatar" />
    </MeModal>

    <MeModal ref="pwdModalRef" title="Change Password" width="420px" @ok="handlePwdSave()">
      <n-form
        ref="pwdFormRef"
        :model="pwdForm"
        label-placement="left"
        require-mark-placement="left"
      >
        <n-form-item label="Old Password" path="oldPassword" :rule="required">
          <n-input v-model:value="pwdForm.oldPassword" type="password" placeholder="Please enter the old password" show-password-on="mousedown" />
        </n-form-item>
        <n-form-item label="New Password" path="newPassword" :rule="required">
          <n-input v-model:value="pwdForm.newPassword" type="password" placeholder="Please enter the new password" show-password-on="mousedown" />
        </n-form-item>
      </n-form>
    </MeModal>

    <MeModal ref="profileModalRef" title="Edit Profile" width="420px" @ok="handleProfileSave()">
      <n-form ref="profileFormRef" :model="profileForm" label-placement="left">
        <n-form-item label="Nickname" path="nickName">
          <n-input v-model:value="profileForm.nickName" placeholder="Please enter nickname" />
        </n-form-item>
        <n-form-item label="Gender" path="gender">
          <n-select
            v-model:value="profileForm.gender"
            :options="genders"
            placeholder="Please select gender"
          />
        </n-form-item>
        <n-form-item label="Address" path="address">
          <n-input v-model:value="profileForm.address" placeholder="Please enter address" />
        </n-form-item>
        <n-form-item label="Email" path="email">
          <n-input v-model:value="profileForm.email" placeholder="Please enter email" />
        </n-form-item>
      </n-form>
    </MeModal>
  </AppPage>
</template>

<script setup>
import { MeModal } from '@/components'
import { useForm, useModal } from '@/composables'
import { useUserStore } from '@/store'
import { getUserInfo } from '@/store/helper'
import api from './api'

const userStore = useUserStore()
const required = {
  required: true,
  message: 'This field is required',
  trigger: ['blur', 'change'],
}

const [pwdModalRef] = useModal()
const [pwdFormRef, pwdForm, pwdValidation] = useForm()

async function handlePwdSave() {
  await pwdValidation()
  await api.changePassword(pwdForm.value)
  $message.success('Password changed successfully')
  refreshUserInfo()
}

const newAvatar = ref(userStore.avatar)
const [avatarModalRef] = useModal()
async function handleAvatarSave() {
  if (!newAvatar.value) {
    $message.error('Please enter avatar URL')
    return false
  }
  await api.updateProfile({ id: userStore.userId, avatar: newAvatar.value })
  $message.success('Avatar updated successfully')
  refreshUserInfo()
}

const genders = [
  { label: 'Prefer not to say', value: 0 },
  { label: 'Male', value: 1 },
  { label: 'Female', value: 2 },
]
const [profileModalRef] = useModal()
const [profileFormRef, profileForm, profileValidation] = useForm({
  id: userStore.userId,
  nickName: userStore.nickName,
  gender: userStore.userInfo?.gender ?? 0,
  address: userStore.userInfo?.address,
  email: userStore.userInfo?.email,
})
async function handleProfileSave() {
  await profileValidation()
  await api.updateProfile(profileForm.value)
  $message.success('Profile updated successfully')
  refreshUserInfo()
}

async function refreshUserInfo() {
  const user = await getUserInfo()
  userStore.setUser(user)
}
</script>
