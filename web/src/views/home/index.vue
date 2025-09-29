<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/05 21:28:22
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang | https://isme.top
 --------------------------------->

<template>
  <AppPage show-footer>
    <div class="flex">
      <n-card class="min-w-200 w-30%">
        <div class="flex items-center">
          <n-avatar round :size="60" :src="userStore.avatar" class="flex-shrink-0" />
          <div class="ml-20 flex-col">
            <span class="text-20 opacity-80">
              Hello, {{ userStore.nickName ?? userStore.username }}
            </span>
            <span class="mt-4 opacity-50">Current Role: {{ userStore.currentRole?.name }}</span>
          </div>
        </div>

        <p class="mt-28 text-14 opacity-60">
          A person can succeed in almost anything for which he has unlimited enthusiasm.
        </p>
        <p class="mt-12 text-right text-12 opacity-40">
          —— Charles Schwab
        </p>
      </n-card>
      <n-card class="ml-12 w-70%" title="✨ Welcome to Vue Naive Admin 2.0">
        <template #header-extra>
          <a
            class="text-highlight text-14 text-primary hover:underline hover:opacity-80"
            href="https://isme.top"
            target="_blank"
          >
            isme.top
          </a>
        </template>

        <p class="opacity-60">
          This is a minimalist admin template with full-stack solutions. Frontend uses Vite + Vue3 + Pinia + Unocss, backend uses Nestjs + TypeOrm + MySql. Simple and elegant, refined through multiple iterations and detail polishing, full of sincerity!
        </p>
        <footer class="mt-12 flex items-center justify-end">
          <n-button
            type="primary"
            ghost
            tag="a"
            href="https://isme.top"
            target="__blank"
          >
            Development Docs
          </n-button>
          <n-button
            type="primary"
            class="ml-12"
            tag="a"
            href="https://github.com/zclzone/vue-naive-admin/tree/2.x"
            target="__blank"
          >
            Code Repository
          </n-button>
        </footer>
      </n-card>
    </div>
    <div class="mt-12 flex">
      <n-card class="w-50%" title="💯 Features" segmented>
        <template #header-extra>
          <span class="text-highlight opacity-90">👏 After more than ten iterations and detail polishing</span>
        </template>

        <ul class="opacity-90">
          <li class="py-4">
            🆒 Uses mainstream
            <b>Vue3</b>
            technology stack:
            <span class="text-highlight">Vite + Vue3 + Pinia</span>
          </li>
          <li class="py-4">
            🍇 Uses
            <b>Atomic CSS</b>
            framework:
            <span class="text-highlight">Unocss</span>
            , elegant, lightweight, easy to use
          </li>
          <li class="py-4">
            🤹 Uses mainstream
            <span class="text-highlight">iconify + unocss</span>
            icon solution, supports custom icons, supports dynamic rendering
          </li>
          <li class="py-4">
            🎨 Uses Naive UI，
            <span class="text-highlight">extremely concise code style and refreshing page design</span>
            , aesthetic online, theme customization is easy
          </li>
          <li class="py-4">
            👏 Advanced and easy-to-understand file structure design, multiple modules have
            <b>zero coupling</b>
            , deleting a single business module does not affect other modules
          </li>
          <li class="py-4">
            🚀
            <span class="text-highlight">Flat routing</span>
            design, every component can be a page, eliminating the difficulty of implementing KeepAlive for multi-level routes
          </li>

          <li class="py-4">
            🍒
            <span class="text-highlight">Dynamically generate routes based on permissions</span>
            , no need to define routes additionally,
            <span class="text-highlight">403 and 404 can be distinguished</span>
            , rather than jumping to 404 even without permissions
          </li>
          <li class="py-4">
            🔐 Integrated
            <span class="text-highlight">seamless refresh</span>
            with Redis, user login state is controllable, security and experience are both indispensable
          </li>
          <li class="py-4">
            ✨ Encapsulated global tool method based on Naive UI
            <span class="text-highlight">message</span>
            , supports batch reminders, supports cross-page shared instances
          </li>
          <li class="py-4">
            ⚡️ Encapsulated commonly used business components based on Naive UI, including
            <span class="text-highlight">Page</span>
            component、
            <span class="text-highlight">CRUD</span>
            table component and
            <span class="text-highlight">Modal</span>
            component, reducing a lot of repetitive work
          </li>
        </ul>

        <n-divider class="mb-0! mt-12!">
          <p class="text-14 opacity-60">
            👉 Click
            <b class="mx-2 transition hover:text-primary">
              <a href="https://isme.top" target="_blank">more</a>
            </b>
            to view more practical features, continuously developing...
          </p>
        </n-divider>
      </n-card>

      <n-card class="ml-12 w-50%" title="🛠️ Tech Stack" segmented>
        <VChart :option="skillOption" autoresize />
      </n-card>
    </div>

    <n-card class="mt-12" title="⚡️ Trends" segmented>
      <div class="h-400">
        <VChart :option="trendOption" autoresize />
      </div>
    </n-card>
  </AppPage>
</template>

<script setup>
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import * as echarts from 'echarts/core'
import { UniversalTransition } from 'echarts/features'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import { useUserStore } from '@/store'

const userStore = useUserStore()

echarts.use([
  TooltipComponent,
  GridComponent,
  LegendComponent,
  BarChart,
  LineChart,
  CanvasRenderer,
  UniversalTransition,
  PieChart,
])

const trendOption = {
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'cross',
      crossStyle: {
        color: '#999',
      },
    },
  },
  legend: {
    top: '5%',
    data: ['star', 'fork'],
  },
  xAxis: [
    {
      type: 'category',
      data: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
      axisPointer: {
        type: 'shadow',
      },
    },
  ],
  yAxis: [
    {
      type: 'value',
      min: 0,
      max: 3000,
      interval: 500,
      axisLabel: {
        formatter: '{value}',
      },
    },
    {
      type: 'value',
      min: 0,
      max: 500,
      interval: 100,
      axisLabel: {
        formatter: '{value}',
      },
    },
  ],
  series: [
    {
      name: 'star',
      type: 'line',
      data: [200, 320, 520, 550, 600, 805, 888, 950, 1300, 2503, 2702, 2712],
    },
    {
      name: 'fork',
      yAxisIndex: 1,
      type: 'bar',
      data: [40, 72, 110, 115, 121, 175, 180, 201, 260, 398, 423, 455],
    },
  ],
}

const skillOption = {
  tooltip: {
    trigger: 'item',
    formatter({ name, value }) {
      return `${name} ${value}%`
    },
  },
  legend: {
    left: 'center',
  },
  series: [
    {
      top: '12%',
      type: 'pie',
      radius: ['35%', '90%'],
      avoidLabelOverlap: true,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2,
      },
      label: {
        show: false,
        position: 'center',
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 36,
          fontWeight: 'bold',
        },
      },
      labelLine: {
        show: false,
      },
      data: [
        { value: 38.5, name: 'Vue' },
        { value: 37.0, name: 'JavaScript' },
        { value: 6.5, name: 'CSS' },
        { value: 6.2, name: 'HTML' },
        { value: 1.8, name: 'Other' },
      ],
    },
  ],
}
</script>
