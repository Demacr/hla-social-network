import { createWebHistory, createRouter } from "vue-router";
import Login from '@/views/LoginPage';
import Home from '@/views/HomePage';
import People from '@/views/PeoplePage';
import Dialog from '@/views/DialogPage';
import DialogsList from '@/views/DialogsListPage';

// Vue.use(Router);

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
  },
  {
    path: '/people',
    name: 'People',
    component: People,
  },
  {
    path: '/dialogs',
    name: 'DialogsList',
    component: DialogsList,
  },
  {
    path: '/dialog/:id',
    name: 'DialogPage',
    props: true,
    component: Dialog,
  },
  {
    path: '/feed',
    name: 'Feed',
    component: () => import('@/views/FeedPage'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router;
