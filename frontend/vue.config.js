module.exports = {
    devServer: {
        proxy: {
            '/api': {
                target: 'http://xxxx/device/', // 对应自己的接口
                changeOrigin: true,
                ws: true,
                pathRewrite: {
                    '^/api': '',
                },
            },
        },
    },
    publicPath: './',
    lintOnSave: false, // 解决编译问题
}