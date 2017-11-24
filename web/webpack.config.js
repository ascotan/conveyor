module.exports = {
  entry: './src/main.js',
  output: {
    path: '../build/assets',
    publicPath: "/assets/",
    filename: '[name].js'
  },
  module: {
    loaders: [
      {
        test: /\.js$/,
        loader: 'babel',
        exclude: /node_modules/
      },
      {
        test: /\.vue$/,
        loader: 'vue'
      }
    ]
  },
  vue: {
    loaders: {
      js: 'babel'
    }
  },
  resolve: {
    alias: {
      'vue': 'vue/dist/vue.common.js'
    }
  },
}