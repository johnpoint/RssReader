import Vue from "vue";
import Vuex from "vuex";
import App from "./App.vue";
import router from "./router";
import { BootstrapVue, IconsPlugin } from "bootstrap-vue";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";
import "@/assets/main.less";

Vue.use(BootstrapVue);
Vue.use(IconsPlugin);

Vue.use(Vuex);

Vue.config.productionTip = false;

const store = new Vuex.Store({
  state: {
    jwt: "",
    isLogin: false,
    userInfo: null
  },
  mutations: {
    setjwt(state, jwt) {
      state.jwt = jwt;
    },
    setStatus(state, bool) {
      state.isLogin = bool;
    },
    setUserInfo(state, data) {
      state.userInfo = data;
    }
  }
});

new Vue({
  router,
  store: store,
  render: h => h(App)
}).$mount("#app");
