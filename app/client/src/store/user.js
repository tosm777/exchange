import { defineStore } from 'pinia'

export const useUserStore = defineStore({
    id: 'user',
    state: () => {
        return {
            userAddress: "",
            publicKey: "",
            privateKey: "",
            accessedUser: false,
        }
    },
    getters: {
        hasInfo() {
            return this.accessedUser
        },
        getUser() {
            return {
                userAddress: this.userAddress,
                publicKey:   this.publicKey,
                privateKey:  this.privateKey,
            }
        }
    },
    actions: {
        update(address, publicKey, privateKey) {
            this.userAddress  = address
            this.publicKey    = publicKey
            this.privateKey   = privateKey
            this.accessedUser = true
        }
    },
    persist: {
        storage: window.sessionStorage,
    },
})
