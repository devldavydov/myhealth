import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import FoodListView from '@/views/food/FoodListView.vue'
import FoodJournalView from '@/views/food/FoodJournalView.vue'
import FoodStatsView from '@/views/food/FoodStatsView.vue'
import WeightListView from '@/views/weight/WeightListView.vue'
import WeightStatsView from '@/views/weight/WeightStatsView.vue'
import ActivityJournalView from '@/views/activity/ActivityJournalView.vue'
import ActivityStatsView from '@/views/activity/ActivityStatsView.vue'
import SportListView from '@/views/activity/SportListView.vue'
import CalcCalView from '@/views/settings/CalcCalView.vue'
import UserView from '@/views/settings/UserView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { title: 'Главная' }
    },
    // food
    {
      path: '/food',
      name: 'foodList',
      component: FoodListView,
      meta: { title: 'Управление едой' }
    },
    {
      path: '/food/journal',
      name: 'foodJournal',
      component: FoodJournalView,
      meta: { title: 'Журнал приема пищи' }
    },
    {
      path: '/food/stats',
      name: 'foodStats',
      component: FoodStatsView,
      meta: { title: 'Статистика' }
    },
    // weight
    {
      path: '/weight',
      name: 'weightList',
      component: WeightListView,
      meta: { title: 'Вес тела' }
    },
    {
      path: '/weight/stats',
      name: 'weightStats',
      component: WeightStatsView,
      meta: { title: 'Статистика' }
    },
    // activity
    {
      path: '/activity',
      name: 'sportList',
      component: SportListView,
      meta: { title: 'Управление спортом' }
    },
    {
      path: '/activity/journal',
      name: 'activityJournal',
      component: ActivityJournalView,
      meta: { title: 'Журнал активности' }
    },
    {
      path: '/activity/stats',
      name: 'activityStats',
      component: ActivityStatsView,
      meta: { title: 'Статистика' }
    },
    // settings
    {
      path: '/settings/calccal',
      name: 'settingsCalcCal',
      component: CalcCalView,
      meta: { title: 'Расчет лимита ккал' }
    },
    {
      path: '/settings/user',
      name: 'settingsUser',
      component: UserView,
      meta: { title: 'Пользователь' }
    },
  ],
})

export default router
