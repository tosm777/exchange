import * as vueRouter from "vue-router";

const routes = [
    {
        name: "home",
        path: "/home",
        component: () => import('../src/components/Home.vue')
    },
    {
        name: "user",
        path: "/user/:user_id",
        component: () => import('../src/components/User.vue')
    },
    {
        name: "wallet",
        path: "/user/:user_id/wallet",
        component: () => import('../src/components/Wallet.vue')
    },
];

const router = vueRouter.createRouter({
    history: vueRouter.createWebHistory(),
    routes,
});

export default router;
