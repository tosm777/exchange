<template>
    <div class="main">
        <!-- <transition name="fade" appear> -->
        <transition name="fade" appear>
            <h1>For get flow of Blockchain</h1>
        </transition>
        <transition name="fade" appear>
            <div class="message">
                <p>Connect Our Moments</p>
                <p>Be fantastic</p>
            </div>
        </transition>
        <transition name="fade" appear>
            <q-row>
                <q-col>
                    <router-link :to="{ name: 'wallet', params: { user_id: 1 } }">for users</router-link>
                </q-col>
                <q-col>for miners</q-col>
            </q-row>
        </transition>

    </div>
</template>


<script setup>
import { ref } from 'vue'
import axios from 'axios'

</script>

<style>
.main {
    text-align: center;
}
.fade-enter-from {
    opacity: 0;
}
.fade-enter-active {
    transition: opacity 2.0s ease;
}
.message {
    padding: 10%;
}
.message p {
    padding: 10%;
}
</style>
