/**********************************
 * @FilePath: common.js
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/04 22:45:46
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import dayjs from 'dayjs'

/**
 * @param {(object | string | number)} time
 * @param {string} format
 * @returns {string | null} Formatted time string
 *
 */
export function formatDateTime(time = undefined, format = 'YYYY-MM-DD HH:mm:ss') {
  return dayjs(time).format(format)
}

export function formatDate(date = undefined, format = 'YYYY-MM-DD') {
  return formatDateTime(date, format)
}

/**
 * @param {Function} fn
 * @param {number} wait
 * @returns {Function}  Throttle function
 *
 */
export function throttle(fn, wait) {
  let context, args
  let previous = 0

  return function (...argArr) {
    const now = +new Date()
    context = this
    args = argArr
    if (now - previous > wait) {
      fn.apply(context, args)
      previous = now
    }
  }
}

/**
 * @param {Function} method
 * @param {number} wait
 * @param {boolean} immediate
 * @return {*} Debounce function
 */
export function debounce(method, wait, immediate) {
  let timeout
  return function (...args) {
    const context = this
    if (timeout) {
      clearTimeout(timeout)
    }
  // Immediate execution requires two conditions: immediate is true, and timeout is not set or is null
    if (immediate) {
  /**
   * If the timer doesn't exist, execute immediately and set a timer that will set timeout to null after wait milliseconds
   * This ensures the immediate execution won't be triggered again within wait milliseconds
   */
      const callNow = !timeout
      timeout = setTimeout(() => {
        timeout = null
      }, wait)
      if (callNow) {
        method.apply(context, args)
      }
    }
    else {
  // If immediate is false, execute the function after wait milliseconds
      timeout = setTimeout(() => {
  /**
   * args is an array-like object, so use method.apply
   * Alternatively, you could write method.call(context, ...args)
   */
  method.apply(context, args)
      }, wait)
    }
  }
}

/**
 * @param {number} time Milliseconds
 * @returns Sleep for a while (helper for delays)
 */
export function sleep(time) {
  return new Promise(resolve => setTimeout(resolve, time))
}

/**
 * @param {HTMLElement} el
 * @param {Function} cb
 * @return {ResizeObserver}
 */
export function useResize(el, cb) {
  const observer = new ResizeObserver((entries) => {
    cb(entries[0].contentRect)
  })
  observer.observe(el)
  return observer
}
