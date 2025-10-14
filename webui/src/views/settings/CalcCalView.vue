<script setup lang="ts">
  import { StringConstants } from '@/constants';
  import { ref } from 'vue';

  interface State {
    gender: string
    weight: number
    height: number
    age: number
    showWeightHelp: boolean
    showHeightHelp: boolean
    showAgeHelp: boolean
    calculated: Array<{indicator: string, value: number}>
  }

  const state = ref({
    gender: "m",
    weight: 0,
    height: 0,
    age: 0,
    showWeightHelp: false,
    showHeightHelp: false,
    showAgeHelp: false,
    calculated: []
  } as State);

  const genders = [
    {key: "m", label: StringConstants.GenderMale},
    {key: "f", label: StringConstants.GenderFemale}
  ] as Array<{key: string, label: string}>;

  function calcCal() {
    if (!isCorrectInput()) {
      return;
    }

    let ubm = 10 * state.value.weight + 6.25 * state.value.height - 5 * state.value.age;
    if (state.value.gender === "m") 
        ubm += 5;
    else
        ubm -= 161;

    const activities: Array<{name: string, k: number}> = [
      {name: StringConstants.Activity1, k: 1.2},
      {name: StringConstants.Activity2, k: 1.375},
      {name: StringConstants.Activity3, k: 1.55},
      {name: StringConstants.Activity4, k: 1.725},
      {name: StringConstants.Activity5, k: 1.9}
    ];
    
    state.value.calculated = [{indicator: StringConstants.Ubm, value: Math.round(ubm)}];
    activities.forEach(a => {
      state.value.calculated.push({indicator: a.name, value: Math.round(ubm * a.k)});
    });
  }

  function isCorrectInput() {
    let res = true;

    if (state.value.weight <= 0.0) {
      state.value.showWeightHelp = true;
      res = false;
    } else {
      state.value.showWeightHelp = false;
    }

    if (state.value.height <= 0.0) {
      state.value.showHeightHelp = true;
      res = false;
    } else {
      state.value.showHeightHelp = false;
    }

    if (state.value.age <= 0.0) {
      state.value.showAgeHelp = true;
      res = false;
    } else {
      state.value.showAgeHelp = false;
    }

    return res;
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
      <div id="weightHelp" v-if="state.showWeightHelp" class="form-text text-danger">{{ StringConstants.ValueG0 }}</div>
    </div>
    <div class="mb-3">
      <label for="height" class="form-label">{{ StringConstants.Height }}</label>
      <input type="number" id="height" min="0" step="0.1" class="form-control" v-model="state.height">
      <div id="heightHelp" v-if="state.showHeightHelp" class="form-text text-danger">{{ StringConstants.ValueG0 }}</div>
    </div>
    <div class="mb-3">
      <label for="age" class="form-label">{{ StringConstants.Age }}</label>
      <input type="number" id="age" min="0" step="0.1" class="form-control" v-model="state.age">
      <div id="ageHelp" v-if="state.showAgeHelp" class="form-text text-danger">{{ StringConstants.ValueG0 }}</div>
    </div>
    <button type="submit" class="btn btn-primary">{{ StringConstants.Calculate }}</button>
  </form>
  <table class="table mt-3" v-if="state.calculated.length > 0">
    <tbody>
      <tr v-for="calc in state.calculated">
        <th>{{ calc.indicator }}</th>
        <td>{{ calc.value }}</td>
      </tr>
    </tbody>
  </table>
</template>
