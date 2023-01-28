import { createWebHistory, createRouter } from "vue-router";
import Index from '@/components/IndexPage';
import Home from '@/components/HomePage';
import People from '@/components/PeoplePage';

// Vue.use(Router);

const routes = [
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
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router;
