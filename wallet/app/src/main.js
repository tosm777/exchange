import { createApp, defineComponent } from 'vue'
import App from './App.vue'
import Qui from '@qvant/qui-max'
import '@qvant/qui-max/styles'
import { VuesticPlugin } from 'vuestic-ui'
import 'vuestic-ui/dist/vuestic-ui.css'
import Equal from 'equal-vue'
import 'equal-vue/dist/style.css'
import router from '../router'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'


const app = createApp(App);
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

app.use(Qui);
app.use(VuesticPlugin);
app.use(Equal);
app.use(router);
app.use(pinia);

app.mount('#app');
