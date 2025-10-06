<script setup lang="ts">
    import { StringConstants } from '@/constants';
  import { ref, computed } from 'vue';

  interface State {
    gender: string
    weight: number
    height: number
    age: number
    calculated: Array<{indicator: string, value: number}>
  }

  const state = ref({
    gender: "m",
    weight: 0,
    height: 0,
    age: 0,
    calculated: []
  } as State)

  const isCalcEnabled = computed(() => state.value.weight !== 0 && state.value.height !== 0 && state.value.age !== 0);

  const genders = [
    {key: "m", label: StringConstants.GenderMale},
    {key: "f", label: StringConstants.GenderFemale}
  ] as Array<{key: string, label: string}>;

  function calcCal() {
    let ubm = 10 * state.value.weight + 6.25 * state.value.height - 5 * state.value.age
    if (state.value.gender === "m") 
        ubm += 5
    else
        ubm -= 161

    const activities: Array<{name: string, k: number}> = [
      {name: StringConstants.Activity1, k: 1.2},
      {name: StringConstants.Activity2, k: 1.375},
      {name: StringConstants.Activity3, k: 1.55},
      {name: StringConstants.Activity4, k: 1.725},
      {name: StringConstants.Activity5, k: 1.9}
    ];
    
    state.value.calculated = [{indicator: StringConstants.Ubm, value: Math.round(ubm)}]
    activities.forEach(a => {
      state.value.calculated.push({indicator: a.name, value: Math.round(ubm * a.k)});
    });
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
      <input type="number" id="weight" min="0" step="0.1" class="form-control" v-model="state.weight"></input>
    </div>
    <div class="mb-3">
      <label for="height" class="form-label">{{ StringConstants.Height }}</label>
      <input type="number" id="height" min="0" step="0.1" class="form-control" v-model="state.height">
    </div>
    <div class="mb-3">
      <label for="age" class="form-label">{{ StringConstants.Age }}</label>
      <input type="number" id="age" min="0" step="0.1" class="form-control" v-model="state.age">
    </div>
    <button type="submit" :disabled="!isCalcEnabled" class="btn btn-primary">{{ StringConstants.Calculate }}</button>
  </form>
  <table class="table mt-3" v-if="state.calculated.length > 0">
    <thead>
      <tr>
        <th scope="col">{{ StringConstants.Indicator }}</th>
        <th scope="col">{{ StringConstants.ValueCal }}</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="calc in state.calculated">
        <td>{{ calc.indicator }}</td>
        <td>{{ calc.value }}</td>
      </tr>
    </tbody>
  </table>
</template>
