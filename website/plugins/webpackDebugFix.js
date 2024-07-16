// See: https://github.com/rohit-gohri/redocusaurus/issues/236#issuecomment-1449548972

module.exports = async function webpackDebugFix(context, opts) {
  const webpack = require('webpack');

  return {
    name: 'webpack-fix-plugin',
    configureWebpack(config, isServer, utils, content) {
      return {
        plugins: [
          new webpack.DefinePlugin({
            'process.env.DEBUG': 'process.env.DEBUG',
          }),
        ],
      };
    },
  };
};
