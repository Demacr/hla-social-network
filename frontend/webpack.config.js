module.exports = {
    context: __dirname,
    devtool: "source-map",
    entry: "./js/profile.js",
    output: {
      path: __dirname + "/dist",
      filename: "bundle.js"
    }
  }