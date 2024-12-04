import { createRouter, createWebHistory } from 'vue-router';
import HomePage from '../views/HomePage.vue';
import AboutPage from '../views/AboutPage.vue';
import NotFoundPage from '../views/NotFoundPage.vue';

import UserList from '../views/UserListWebPage.vue';
import Login from '../views/LoginPage.vue';

const routes = [
  { path: '/', name: 'home', component: HomePage },
  { path: '/about', name: 'about', component: AboutPage },
  { path: '/:pathMatch(.*)*', name: 'not-found', component: NotFoundPage },
  { path: '/users', name: 'users', component: UserList },
  { path: '/login', name: 'login', component: Login },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;

