<script setup lang="ts">
  import Loader from '@/components/Loader.vue';
  import { StringConstants } from '@/constants';
  import type { IFood } from '@/services/FoodService';
  import { getFoodList } from '@/services/FoodService';
  import { onMounted, ref, computed } from 'vue';
  import { toast } from 'vue3-toastify';
  import Search from '@/components/Search.vue';
  import BtnCreate from '@/components/BtnCreate.vue';

  interface IState {
    showLoader: boolean
    showResult: boolean
    foodList: IFood[]
    search: string
  }

  const state = ref({
    showLoader: true,
    showResult: false,
    foodList: [],
    search: ""
  } as IState);

  onMounted(() => {
    getFoodList()
      .finally(() => {
        state.value.showLoader = false;
      })
      .then((result) => {
        state.value.showResult = true;
        state.value.foodList = result;
      })
      .catch((error: Error) => {
        toast.error(error.message);
      });
  });

  const filteredFoodList = computed(() => {
    const pattern = state.value.search.toLocaleUpperCase();

    return state.value.foodList.filter((f) => 
      f.name.toLocaleUpperCase().indexOf(pattern) !== -1 ||
      f.brand.toLocaleUpperCase().indexOf(pattern) !== -1 ||
      f.comment.toLocaleUpperCase().indexOf(pattern) !== -1
    );
  });

</script>

<template>
  <h3>{{ StringConstants.FoodList }}</h3>
  
  <div class="row mb-2">
    <div class="col-sm-4">
      <BtnCreate url="/food/create"/>
    </div>
    <div class="col-sm-8 float-end">
      <Search v-model="state.search" />
    </div>
  </div>

  <Loader :show="state.showLoader"/>
  
  <div class="table-responsive" v-if="state.showResult">
    <table class="table table-striped table-bordered table-hover">
      <thead>
        <tr>
          <th class="align-middle col-4">{{ StringConstants.Name }}</th>
          <th class="align-middle col-2">{{ StringConstants.Brand }}</th>
          <th class="align-middle col-1">{{ StringConstants.Cal100 }}</th>
          <th class="align-middle col-2">{{ StringConstants.Comment }}</th>
          <th class="align-middle col-1 text-center"><i className="bi bi-pencil"></i></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="food in filteredFoodList" :key="food.key">
          <td>{{ food.name }}</td>
          <td>{{ food.brand }}</td>
          <td>{{ food.cal100 }}</td>
          <td>{{ food.comment}}</td>
          <td></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
