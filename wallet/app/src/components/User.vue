<template>
    <div class="main">
        <div style="width: 100%; text-align: right;">
            <p>
                <it-icon name="face" :bottom="0" style="display: inline-flex; vertical-align: bottom;" outlined />
                {{ name }}
            </p>
        </div>

        <h1>Your Info</h1>

        <div>
            <q-form class="user-form">
                <q-form-item label="Public Key" >
                    <q-input v-model="publicKey" :disabled="true" />
                </q-form-item>

                <q-form-item label="PrivateKey Key" >
                    <q-input v-model="privateKey" :disabled="true" />
                </q-form-item>

                <q-form-item label="My Address" >
                    <q-input v-model="userAddress" :disabled="true" />
                </q-form-item>

                <q-form-item style="text-align: center;">
                    <q-button @click="refresh">Refresh</q-button>
                </q-form-item>
            </q-form>
        </div>

        <router-link :to="{ name: 'wallet', param: { user_id: 1} }">
            <q-button theme="link" style="width: 100%; text-align: center;">Go to Wallet</q-button>
        </router-link>
    </div>
</template>


<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import {  mapState } from 'pinia'
import { useUserStore } from "@/store/user";

const userStore = useUserStore()
const name = "Toshi"
const publicKey = ref()
const privateKey = ref()
const userAddress = ref()

onMounted(() => fetch())

function refresh() {
    requestUser()
}

function fetch() {
    if (userStore.hasInfo) {
        setUser()
        return
    }
    requestUser()
}

function requestUser() {
    axios.defaults.headers.get['Content-Type'] = 'application/json;charset=utf-8';
    axios.defaults.headers.get['Access-Control-Allow-Origin'] = '*'

    axios.post('http://net.develop:443/wallet')
        .then(function (response) {
            // handle success(axiosの処理が成功した場合に処理させたいことを記述)
            userAddress.value = response.data['blockchain_address']
            publicKey.value  = response.data['public_key']
            privateKey.value = response.data['private_key']
            console.log(response)

            userStore.update(userAddress.value, publicKey.value, privateKey.value)
            console.log(userStore.getUser)
        })
        .catch(function (error) {
            console.log(error);
        })
        .finally(function () {
            // always executed(axiosの処理結果によらずいつも実行させたい処理を記述)
        });
}

function setUser() {
    const userInfo = userStore.getUser

    userAddress.value = userInfo['userAddress']
    publicKey.value   = userInfo['publicKey']
    privateKey.value  = userInfo['privateKey']
}
</script>

<style>
h1 {
    vertical-align: middle;
}

.user-form {
    margin-top: 5%;
}
</style>
