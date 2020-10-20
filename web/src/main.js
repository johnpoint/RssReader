import Vue from "vue";
import Vuex from "vuex";
import i18n from "@/locales"
import App from "./App.vue";
import router from "./router";
import {BootstrapVue, IconsPlugin} from "bootstrap-vue";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";
import "@/assets/main.less";
import '@/registerServiceWorker'

Vue.use(BootstrapVue);
Vue.use(IconsPlugin);
Vue.use(Vuex);

Vue.config.productionTip = false;

Date.prototype.format = function (fmt) {
    var o = {
        "M+": this.getMonth() + 1,                 //月份
        "d+": this.getDate(),                    //日
        "h+": this.getHours(),                   //小时
        "m+": this.getMinutes(),                 //分
        "s+": this.getSeconds(),                 //秒
        "q+": Math.floor((this.getMonth() + 3) / 3), //季度
        "S": this.getMilliseconds()             //毫秒
    };
    if (/(y+)/.test(fmt)) {
        fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    }
    for (var k in o) {
        if (new RegExp("(" + k + ")").test(fmt)) {
            fmt = fmt.replace(RegExp.$1, (RegExp.$1.length === 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
        }
    }
    return fmt;
}

const store = new Vuex.Store({
    state: {
        jwt: "",
        isLogin: false,
        config: {
            postnum: null
        },
        postList: null
    },
    mutations: {
        setjwt(state, jwt) {
            state.jwt = jwt;
        },
        setStatus(state, bool) {
            state.isLogin = bool;
        }
    }
});

new Vue({
    router,
    i18n,
    store: store,
    render: h => h(App)
}).$mount("#app");
