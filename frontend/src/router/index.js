import {createRouter, createWebHistory, createWebHashHistory} from 'vue-router'
import LayoutAdmin from "@/layout/LayoutAdmin";
import LayoutUser from "@/layout/LayoutUser";
import LayoutNew from "@/layout/LayoutNew";

const routes = [
    {
      path: '/',
      redirect: '/login',
    },
    {
        path: '/front',
        component: LayoutNew,
        redirect: "/front/home",
        children: [
            {
                path: 'basic_behavior_data',
                component: () => import("@/views/new/BasicBehaviorData")
            },
            {
                path: 'behavior_rule',
                component: () => import("@/views/new/BehaviorRule")
            },
            {
                path: 'behavior_rule_data',
                component: () => import("@/views/new/BehaviorRuleData")
            },
            {
                path: 'data_mean',
                component: () => import("@/views/new/DataMean")
            },
            {
                path: 'event_rule',
                component: () => import("@/views/new/EventRule")
            },
            {
                path: 'event_rule_data',
                component: () => import("@/views/new/EventRuleData")
            },
            {
                path: 'home',
                component: () => import("@/views/new/Home")
            },
            {
                path: 'label',
                component: () => import("@/views/new/Label")
            },
            {
                path: 'learning_model',
                component: () => import("@/views/new/LearningModel")
            },
            {
                path: 'predict',
                component: () => import("@/views/new/Predict")
            },
            {
                path: 'statistics_model',
                component: () => import("@/views/new/StatisticsModel")
            },
            {
                path: 'user_record',
                component: () => import("@/views/new/UserRecord")
            }
        ]
    },
    {
        path: '/admin',
        // name: 'Admin',
        component: LayoutAdmin,
        redirect: "/admin/home",
        children: [
            {
                path: 'home',
                // name: 'Home',
                component: () => import("@/views/admin/Home")
            },
            {
                path: 'department',
                // name: 'Department',
                component: () => import("@/views/admin/Department")
            },
            {
                path: 'goods',
                // name: 'Goods',
                component: () => import("@/views/admin/Goods")
            },
            {
                path: 'insurance',
                // name: 'Insurance',
                component: () => import("@/views/admin/Insurance")
            },
            {
                path: 'log',
                // name: 'Log',
                component: () => import("@/views/admin/Log")
            },
            {
                path: 'manager',
                // name: 'Manager',
                component: () => import("@/views/admin/Manager")
            },
            {
                path: 'order',
                // name: 'Order',
                component: () => import("@/views/admin/Order")
            },
            {
                path: 'price',
                // name: 'Price',
                component: () => import("@/views/admin/Price")
            },
            {
                path: 'profile',
                // name: 'Profile',
                component: () => import("@/views/admin/Profile")
            },
            {
                path: 'staff',
                // name: 'Staff',
                component: () => import("@/views/admin/Staff")
            },
            {
                path: 'storehouse',
                // name: 'Storehouse',
                component: () => import("@/views/admin/Storehouse")
            },
            {
                path: 'user',
                // name: 'User',
                component: () => import("@/views/admin/User")
            },
            {
                path: 'vehicle',
                // name: 'Vehicle',
                component: () => import("@/views/admin/Vehicle")
            }
        ]
    },
    {
        path: '/user',
        // name: 'user',
        component: LayoutUser,
        redirect: "/user/home",
        children: [
            {
                path: 'home',
                // name: 'Home',
                component: () => import("@/views/user/Home")
            },
            {
                path: 'profile',
                // name: 'Profile',
                component: () => import("@/views/user/Profile")
            },
            {
                path: 'goods',
                // name: 'Goods',
                component: () => import("@/views/user/Goods")
            },
            {
                path: 'order',
                // name: 'Order',
                component: () => import("@/views/user/Order")
            }
        ]
    },
    {
        path: '/login',
        // name: 'Login',
        component: () => import("@/views/Login")
    },
    {
        path: '/register',
        // name: 'Register',
        component: () => import("@/views/Register")
    },
]

const router = createRouter({
    history: createWebHashHistory(process.env.BASE_URL),
    
    routes
})

function isEmptyStr(s) {
    if (s === undefined || s == null || s === '') {
        return true
    }
    return false
}

// 限制某些页面禁止未登录访问
let limitPagePath = ['/about']
router.beforeEach((to, from, next) => {
    if (limitPagePath.includes(to.path)) {
        // 判断sessionStorage是否保存了用户信息
        let userStr = sessionStorage.getItem("user") || "{}"
        let user = JSON.parse(userStr)
        // 判断sessionStorage是否保存了管理员信息
        let adminStr = sessionStorage.getItem("admin") || "{}"
        let admin = JSON.parse(adminStr)
        console.log(adminStr)
        console.log(admin)
        if (isEmptyStr(user.username) && isEmptyStr(admin.accountName)) {
            // 跳转到登录页面
            next({path: "/login"})
        } else {
            next()
        }
    } else {
        next()
    }

})

export default router
