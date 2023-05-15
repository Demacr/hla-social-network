import 'bootstrap/dist/css/bootstrap.css';
import 'primevue/resources/themes/saga-blue/theme.css';
import 'primevue/resources/primevue.min.css';
import 'primeflex/primeflex.min.css';
import 'primeicons/primeicons.css';
import { createApp } from 'vue';
// import Axios from 'axios';
import App from './App.vue';
import router from './router';
import PrimeVue from 'primevue/config';

import Button from 'primevue/button';
import CheckBox from 'primevue/checkbox';
import Column from 'primevue/column';
import DataTable from 'primevue/datatable';
import Dialog from 'primevue/dialog';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import RadioButton from 'primevue/radiobutton';
import ScrollPanel from 'primevue/scrollpanel';

import MyInfo from './components/MyInfo'

import ToastService from 'primevue/toastservice';

// const token = localStorage.getItem('token');

const app = createApp(App);

app.use(router);
app.use(PrimeVue);
app.use(ToastService);
app.component('InputNumber', InputNumber);
app.component('InputText', InputText);
// eslint-disable-next-line
app.component('Button', Button);
// eslint-disable-next-line
app.component('Checkbox', CheckBox);
// eslint-disable-next-line
app.component('Column', Column);
app.component('DataTable', DataTable);
// eslint-disable-next-line
app.component('Dialog', Dialog);
app.component('RadioButton', RadioButton);
app.component('ScrollPanel', ScrollPanel);
app.component('MyInfo', MyInfo);
app.mount('#app');
// app.config.globalProperties.$http = Axios;
// if (token) {
//     app.config.globalProperties.$http.defaults.headers.common.Authorization = token;
// }
