// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import 'bootstrap/dist/css/bootstrap.css';
import BootstapVue from 'bootstrap-vue';
import Axios from 'axios';
import Vue from 'vue';
import App from './App';
import router from './router';

Vue.prototype.$http = Axios;

const token = localStorage.getItem('token');

Vue.config.productionTip = false;
if (token) {
  Vue.prototype.$http.defaults.headers.common.Authorization = token;
}

Vue.use(BootstapVue);

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>',
});
