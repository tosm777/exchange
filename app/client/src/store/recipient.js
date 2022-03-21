import { defineStore } from 'pinia'

export const useRecipientStore = defineStore({
    id: 'recipient',
    state: () => {
        return {
            recipientAddress: "",
            amount: 0.0
        }
    },
    getters: {},
    actions: {
        update(address, amount) {
            this.recipientAddress  = address
            this.amount            = amount
        }
    },
    persist: {
        storage: window.sessionStorage,
    },
})
