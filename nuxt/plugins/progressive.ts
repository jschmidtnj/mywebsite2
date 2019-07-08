/**
 * progressive plugin required for progressive images
 */
if (process.browser) {
  // @ts-ignore
  window.onNuxtReady(() => {
    require('progressive-image.js')
  })
}