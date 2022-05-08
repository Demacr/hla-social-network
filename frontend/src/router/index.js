import Vue from 'vue';
import Router from 'vue-router';
import Index from '@/components/Index';
import Home from '@/components/Home';
import People from '@/components/People';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Index',
      component: Index,
    },
    {
      path: '/home',
      name: 'Home',
      component: Home,
    },
    {
      path: '/people',
      name: 'People',
      component: People,
    },
  ],
  mode: 'history',
});
