<script setup lang="ts">
  import Todo from '@/components/Todo.vue';
  import { StringConstants } from '@/constants';
  import { ref } from 'vue';

  interface State {
    gender: string
    weight: number
    height: number
    age: number
  }

  const state = ref({
    gender: "m",
    weight: 0,
    height: 0,
    age: 0
  } as State)

  const genders = [
    {key: "m", label: StringConstants.GenderMale},
    {key: "f", label: StringConstants.GenderFemale}
  ] as Array<{key: string, label: string}>;

  function calcCal() {
    alert(1);
  }
</script>

<template>
  <h3>{{ StringConstants.SettingsCalcCal }}</h3>
  <form @submit.prevent="calcCal">
  <div class="mb-3">
    <label for="gender" class="form-label">{{ StringConstants.Gender }}</label>
    <div class="form-check" v-for="gender in genders" :key="gender.key">
      <input
        class="form-check-input"
        type="radio"
        v-model="state.gender"
        :id="gender.key"        
        :value="gender.key"
        >
      <label class="form-check-label" :for="gender.key">
        {{ gender.label }}
      </label>
    </div>
  </div>
  <div class="mb-3">
    <label for="weight" class="form-label">{{ StringConstants.Weight }}</label>
    <input type="number" min="0" step="0.1" class="form-control" v-model="state.weight">
  </div>  
  <button type="submit" class="btn btn-primary">Submit</button>
</form>
</template>
