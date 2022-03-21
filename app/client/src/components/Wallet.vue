<template>
    <div class="main">
        <div style="width: 100%; text-align: right;">
            <router-link :to="{ name: 'home' }">
                <va-icon name= "home" size="large" class="mr-2" color="pink" />
            </router-link>
            <router-link :to="{ name: 'user', params: { user_id: 1 } }">
                <it-avatar color="#3051ff" text="Toshi" size="32px" />
            </router-link>
        </div>

        <h1>Wallet</h1>

        <div>
            <q-form class="wallet-form">
                <q-form-item label="User Address" >
                    <q-input :disabled="true" v-model="userAddress" />
                </q-form-item>

                <q-form-item label="User Amount" >
                    <div style="display: inline-block; width: 50%; margin-right: 5%;" >
                        <q-input :disabled="true" v-model="userAmount" />
                    </div>
                    <q-button @click="GetUserAmount" :loading="false" type="icon" icon="q-icon-refresh-fill"  size="medium" />
                </q-form-item>


                <q-form-item label="Send Address" >
                    <q-input v-model="recipientAddress" />
                </q-form-item>

                <q-form-item label="Send Amount" >
                    <q-input v-model="sendAmount" style="width: 50%;" />
                </q-form-item>

                <q-form-item>
                    <q-button @click="$refs.modal.show()">Send</q-button>
                    <va-modal ref="modal" @ok="Send" stateful message="Send Yours, OK?" />
                    <q-button @click="F" theme="secondary">Reset</q-button>
                </q-form-item>
            </q-form>
        </div>
    </div>
</template>


<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import {  mapGetters, mapState } from 'pinia'
import { useUserStore } from "@/store/user";
import { useRecipientStore } from "@/store/recipient"
import { useRoute, useRouter } from 'vue-router'

const userStore        = useUserStore()
const recipientStore   = useRecipientStore()

const userStoreState        = mapState(useUserStore, ['userAddress', 'publicKey', 'privateKey'])
const recipientStoreState   = mapState(useRecipientStore, ['recipientAddress', 'amount'])

const userAddress      = computed(() => userStoreState.userAddress())
const publicKey        = computed(() => userStoreState.publicKey())
const privateKey       = computed(() => userStoreState.privateKey())
const recipientAddress = ref()
const userAmount       = ref(0.0)
const sendAmount       = ref(0.0)

// setInterval(() => {
//     GetUserAmount()
// }, 3000);

onMounted(() => {
    if (!userStore.hasInfo) {
        const router = useRouter()
        router.push({ name: 'user', params: { user_id: 1 } })
        return
    }

    GetUserAmount()

    recipientAddress.value = recipientStoreState.recipientAddress()
    // amount.value           = recipientStoreState.amount()
})

function F() {
    console.log(userStore.getUser)
}

function GetUserAmount() {

    axios.defaults.headers.get['Content-Type'] = 'application/json;charset=utf-8'
    axios.defaults.headers.get['Access-Control-Allow-Origin'] = '*'

    const params = {
        blockchain_address:    userAddress.value,
    }

    axios.get("http://net.develop:443/wallet/amount", {params: params})
        .then(function (response) {
            // handle success(axiosの処理が成功した場合に処理させたいことを記述)
            console.log(response)
            userAmount.value = response.data['amount']
        })
        .catch(function (error) {
            console.log(error)
        })
        .finally(function () {
        });
}

function Send() {

    recipientStore.update(recipientAddress.value, sendAmount.value)

    const params = {
        sender_public_key:    publicKey.value,
        sender_private_key:   privateKey.value,
        sender_address:       userAddress.value,
        recipient_address:    recipientAddress.value,
        value:                sendAmount.value,
    }

    axios.defaults.headers.get['Content-Type'] = 'application/json;charset=utf-8'
    axios.defaults.headers.get['Access-Control-Allow-Origin'] = '*'

    axios.post('http://net.develop:443/transaction', params)
        .then(function (response) {
            // handle success(axiosの処理が成功した場合に処理させたいことを記述)
            console.log(response)
        })
        .catch(function (error) {
            console.log(error)
        })
        .finally(function () {
        });
}
</script>

<style>
h1 {
    vertical-align: middle;
}
.wallet-form {
    margin-top: 5%;
}
</style>
