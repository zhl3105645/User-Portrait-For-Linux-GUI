import axios from 'axios'
import router from "@/router";

const request = axios.create({
    baseURL: "/",
    timeout: 5000
})

// 请求白名单，如果请求在白名单里面，将不会被拦截校验权限
const whiteUrls = ["/login", '/register', "/applist"]
// 上传文件url
const uploadUrls = ["/api/user/upload/"]
// 下载文件url
const downloadUrls = ["/api/seq_mining_result/"]

// request 拦截器
// 可以自请求发送前对请求做一些处理
// 比如统一加token，对请求参数统一加密
request.interceptors.request.use(config => {
    config.headers['Content-Type'] = 'application/json;charset=utf-8'
    uploadUrls.forEach((val, index) => {
        if (config.url.includes(val)) {
            config.headers['Content-Type'] = 'multipart/form-data'
        }
    })
    downloadUrls.forEach((val, index) => {
        if (config.url.includes(val)) {
            console.log("下载")
            config.responseType = 'blob'
            config.headers['content-type'] = 'application/x-download;charset=utf-8'
            config.headers['content-disposition'] = 'attachment;filename=*'
        }
    })
    
    if (!whiteUrls.includes(config.url)) { // 校验请求白名单
        let token = sessionStorage.getItem("token")
        console.log("token=", token)
        if (isEmptyStr(token)) {
            router.push("/login")
        } else {
            console.log("加token头:" + token)
            config.headers['Authorization'] = "Bearer " + token;  // 设置请求头 ，将用户的token保存在请求头中
        }
    }
    return config
}, error => {
    return Promise.reject(error)
});

// response 拦截器
// 可以在接口响应后统一处理结果
request.interceptors.response.use(
    response => {
        let res = response.data;
        // 如果是返回的文件
        if (response.config.responseType === 'blob') {
            return res
        }
        // 兼容服务端返回的字符串数据
        if (typeof res === 'string') {
            res = res ? JSON.parse(res) : res
        }
        // 验证token
        if (res.status_code === 401) {
            console.error("token过期，重新登录")
            router.push("/login")
        }
        return res;
    },
    error => {
        console.log('err' + error) // for debug
        return Promise.reject(error)
    }
)
function isEmptyStr(s) {
    return s === undefined || s === null || s === '';
}

export default request

