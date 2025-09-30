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

import { StringConstants } from '@/constants'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { title: StringConstants.MainPage }
    },
    // food
    {
      path: '/food',
      name: 'foodList',
      component: FoodListView,
      meta: { title: StringConstants.FoodList }
    },
    {
      path: '/food/journal',
      name: 'foodJournal',
      component: FoodJournalView,
      meta: { title: StringConstants.FoodJournal }
    },
    {
      path: '/food/stats',
      name: 'foodStats',
      component: FoodStatsView,
      meta: { title: StringConstants.Statistics }
    },
    // weight
    {
      path: '/weight',
      name: 'weightList',
      component: WeightListView,
      meta: { title: StringConstants.WeightList }
    },
    {
      path: '/weight/stats',
      name: 'weightStats',
      component: WeightStatsView,
      meta: { title: StringConstants.Statistics }
    },
    // activity
    {
      path: '/activity',
      name: 'sportList',
      component: SportListView,
      meta: { title: StringConstants.SportList }
    },
    {
      path: '/activity/journal',
      name: 'activityJournal',
      component: ActivityJournalView,
      meta: { title: StringConstants.ActivityJournal }
    },
    {
      path: '/activity/stats',
      name: 'activityStats',
      component: ActivityStatsView,
      meta: { title: StringConstants.Statistics }
    },
    // settings
    {
      path: '/settings/calccal',
      name: 'settingsCalcCal',
      component: CalcCalView,
      meta: { title: StringConstants.SettingsCalcCal }
    },
    {
      path: '/settings/user',
      name: 'settingsUser',
      component: UserView,
      meta: { title: StringConstants.SettingsUser }
    },
  ],
})

export default router
