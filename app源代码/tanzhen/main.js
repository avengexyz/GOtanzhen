import Vue from 'vue'
import App from './App'

Vue.config.productionTip = false

App.mpType = 'app'

// 引入uView主JS库
import uView from "uview-ui";
Vue.use(uView);
// 引入uView主JS库

const app = new Vue({
    ...App
})
app.$mount()
