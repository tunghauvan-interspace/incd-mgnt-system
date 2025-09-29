/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:31:02
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import path from 'node:path'
import Vue from '@vitejs/plugin-vue'
import VueJsx from '@vitejs/plugin-vue-jsx'
import Unocss from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import { defineConfig, loadEnv } from 'vite'
import removeNoMatch from 'vite-plugin-router-warn'
import VueDevTools from 'vite-plugin-vue-devtools'
import { pluginIcons, pluginPagePathes } from './build/plugin-isme'

export default defineConfig(({ mode }) => {
  const viteEnv = loadEnv(mode, process.cwd())
  const { VITE_PUBLIC_PATH, VITE_PROXY_TARGET } = viteEnv

  return {
    base: VITE_PUBLIC_PATH || '/',
    plugins: [
      Vue(),
      VueJsx(),
      VueDevTools(),
      Unocss(),
      AutoImport({
        imports: ['vue', 'vue-router'],
        dts: false,
      }),
      Components({
        resolvers: [NaiveUiResolver()],
        dts: false,
      }),
  // Custom plugin: generate page file paths and add them to a virtual module
      pluginPagePathes(),
  // Custom plugin: generate custom icons and add them to a virtual module
      pluginIcons(),
  // Remove unnecessary vue-router dynamic route warnings: "No match found for location with path"
      removeNoMatch(),
    ],
    resolve: {
      alias: {
        '@': path.resolve(process.cwd(), 'src'),
        '~': path.resolve(process.cwd()),
      },
    },
    server: {
      host: '0.0.0.0',
      port: 3200,
      open: false,
      proxy: {
        '/api': {
          target: VITE_PROXY_TARGET,
          changeOrigin: true,
          rewrite: path => path.replace(/^\/api/, ''),
          secure: false,
          configure: (proxy, options) => {
            // Configuring this option adds the real request URL to the response headers
            proxy.on('proxyRes', (proxyRes, req) => {
              proxyRes.headers['x-real-url'] = new URL(req.url || '', options.target)?.href || ''
            })
          },
        },
      },
    },
    build: {
      chunkSizeWarningLimit: 1024, // chunk size warning limit (in KB)
    },
  }
})
