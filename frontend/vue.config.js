const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    proxy: 'http://localhost:8080',
  },
  pluginOptions: {
    webpack: {
      dir: [
        './webpack'
      ]
    }
  }
})
