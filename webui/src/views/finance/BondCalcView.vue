<script setup lang="ts">
import { StringConstants } from '@/constants';
import { ref, computed } from 'vue';

interface State {
  nominal: number,
  price: number
  matDate: string
  ncd: number
  cd: number
  cdCnt: number
  sum: number
  showNominalHelp: boolean
  showPriceHelp: boolean
  showMatDateHelp: boolean
  showCdHelp: boolean
  showCdCntHelp: boolean
  showSumHelp: boolean
  calculated: Array<{indicator: string, value: string}>
}

const state = ref({
  nominal: 1000,
  price: 0.0,
  matDate: getTomorrowDate(),
  ncd: 0.0,
  cd: 0.0,
  cdCnt: 0,
  sum: 0.0,
  showNominalHelp: false,
  showPriceHelp: false,
  showMatDateHelp: false,
  showCdHelp: false,
  showCdCntHelp: false,
  showSumHelp: false,
  calculated: []
} as State)

function bondCalc() {
  if (!isCorrectInput()) {
    return;
  }

  const totalPrice = state.value.price + state.value.ncd;
  const cdSum = state.value.cd * state.value.cdCnt;
  const totalBoughtCnt = Math.round(state.value.sum / totalPrice);
  const totalBoughtSum = totalBoughtCnt * totalPrice;
  const daysToMaturity = getDaysBetween(new Date(state.value.matDate), new Date());
  const ytm = ((state.value.nominal + cdSum - totalPrice) / totalPrice) * (365 / daysToMaturity) * 100;
  const totalSum = totalBoughtCnt * (state.value.nominal + cdSum);
  const diffSum = totalSum - totalBoughtSum;

  state.value.calculated = [{indicator: StringConstants.TotalBoughtCnt, value: totalBoughtCnt.toString()}]
  state.value.calculated.push({indicator: StringConstants.TotalBoughtSum, value: totalBoughtSum.toFixed(2)});
  state.value.calculated.push({indicator: StringConstants.DaysToMaturity, value: daysToMaturity.toString()});
  state.value.calculated.push({indicator: StringConstants.YTM, value: ytm.toFixed(2)});
  state.value.calculated.push({indicator: StringConstants.TotalSum, value: totalSum.toFixed(2)});
  state.value.calculated.push({indicator: StringConstants.DiffSum, value: diffSum.toFixed(2)});
}

function isCorrectInput() {
  let res = true;

  if (state.value.nominal <= 0.0) {
    state.value.showNominalHelp = true;
    res = false;
  } else {
    state.value.showNominalHelp = false;
  }

  if (state.value.price <= 0.0) {
    state.value.showPriceHelp = true;
    res = false;
  } else {
    state.value.showPriceHelp = false;
  }

  if (new Date() > new Date(state.value.matDate)) {
    state.value.showMatDateHelp = true;
    res = false;
  } else {
    state.value.showMatDateHelp = false;
  }

  if (state.value.cd <= 0.0) {
    state.value.showCdHelp = true;
    res = false;
  } else {
    state.value.showCdHelp = false;
  }

  if (state.value.cdCnt <= 0) {
    state.value.showCdCntHelp = true;
    res = false;
  } else {
    state.value.showCdCntHelp = false;
  }

  if (state.value.sum <= 0.0) {
    state.value.showSumHelp = true;
    res = false;
  } else {
    state.value.showSumHelp = false;
  }

  return res;
}

function getTomorrowDate() {
  const today = new Date();
  const dt = new Date(today);
  
  dt.setDate(today.getDate() + 1);
  return `${dt.getFullYear()}-${String(dt.getMonth() + 1).padStart(2, '0')}-${String(dt.getDate()).padStart(2, '0')}`;
}

function getDaysBetween(startDate: Date, endDate: Date) {
  const date1 = new Date(startDate);
  const date2 = new Date(endDate);

  const diffInMilliseconds = Math.abs(date2.getTime() - date1.getTime());
  const oneDayInMilliseconds = 1000 * 60 * 60 * 24;

  const diffInDays = Math.round(diffInMilliseconds / oneDayInMilliseconds);
  return diffInDays;
}

function trunc(v: number) {
  return Math.trunc(v * 100) / 100;
}

</script>

<template>
  <h3>{{ StringConstants.FinanceBondCalc }}</h3>
  <form @submit.prevent="bondCalc">
    <div class="mb-3">
      <label for="nominal" class="form-label">{{ StringConstants.Nominal }}</label>
      <input type="number" id="nominal" min="0" step="0.01" class="form-control" v-model="state.nominal">
      <div id="nominalHelp" v-if="state.showNominalHelp" class="form-text text-danger">{{ StringConstants.ValueG0 }}</div>
    </div>    
    <div class="mb-3">
      <label for="price" class="form-label">{{ StringConstants.BondCurrentPrice }}</label>
      <input type="number" id="price" min="0" step="0.01" class="form-control" v-model="state.price">
      <div id="priceHelp" v-if="state.showPriceHelp" class="form-text text-danger">{{ StringConstants.ValueG0 }}</div>
    </div>
    <div class="mb-3">
      <label for="matDate" class="form-label">{{ StringConstants.BondMaturityDate }}</label>
      <input type="date" id="matDate" class="form-control" v-model="state.matDate">
      <div id="matDateHelp" v-if="state.showMatDateHelp" class="form-text text-danger">{{ StringConstants.DateGToday }}</div>
    </div>
    <div class="mb-3">
      <label for="ncd" class="form-label">{{ StringConstants.BondNcd }}</label>
      <input type="number" id="ncd" min="0" step="0.01" class="form-control" v-model="state.ncd">
    </div>
    <div class="mb-3">
      <label for="cd" class="form-label">{{ StringConstants.BondCd }}</label>
      <input type="number" id="cd" min="0" step="0.01" class="form-control" v-model="state.cd">
      <div id="cdHelp" v-if="state.showCdHelp" class="form-text text-danger">{{ StringConstants.ValueG0 }}</div>
    </div>    
    <div class="mb-3">
      <label for="cdCnt" class="form-label">{{ StringConstants.BondCdCount }}</label>
      <input type="number" id="cdCnt" min="0" step="1" class="form-control" v-model="state.cdCnt">
      <div id="cdCntHelp" v-if="state.showCdCntHelp" class="form-text text-danger">{{ StringConstants.ValueG0 }}</div>
    </div>
    <div class="mb-3">
      <label for="sum" class="form-label">{{ StringConstants.BondSum }}</label>
      <input type="number" id="sum" min="0" step="0.01" class="form-control" v-model="state.sum">
      <div id="sumHelp" v-if="state.showSumHelp" class="form-text text-danger">{{ StringConstants.ValueG0 }}</div>
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