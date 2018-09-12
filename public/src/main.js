import Vue from 'vue'
import App from './App.vue'
import router from "./router";
import axios from "axios";
import globalConfig from "../appConfig.json";

Vue.prototype.$globalConfig = globalConfig;

Vue.config.productionTip = false
Vue.prototype.$http = axios;

export const EventBus = new Vue(); // To pass messages between components

Vue.mixin({
    methods: {
        paramsToUrl: function(formValues) {
            var queryParams = [];
            for (var param in formValues) {
                queryParams.push(
                    `${param}=${encodeURIComponent(formValues[param])}`
                );
            }
            return queryParams.join("&");
        }
    }
});

new Vue({
    router: router,
    render: h => h(App)
}).$mount('#app')