import {createRouter, createWebHistory, createWebHashHistory} from 'vue-router'
import LayoutNew from "@/layout/LayoutNew";

const routes = [
    {
      path: '/',
      redirect: '/login',
    },
    {
        path: '/front',
        component: LayoutNew,
        redirect: "/front/label",
        children: [
            {
                path: 'permission',
                component: () => import("@/views/new/Permission")
            },
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
                path: 'seq_mining',
                component: () => import("@/views/new/SeqMining")
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
                path: 'crowd',
                component: () => import("@/views/new/Crowd")
            },
            {
                path: 'learning_model',
                component: () => import("@/views/new/LearningModel")
            },
            {
                path: 'statistics_model',
                component: () => import("@/views/new/StatisticsModel")
            },
            {
                path: 'user_record',
                component: () => import("@/views/new/UserRecord")
            },
            {
                path: 'profile',
                component: () => import("@/views/new/Profile")
            },
            {
                path: 'group_profile',
                component: () => import("@/views/new/GroupProfile")
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

export default router
