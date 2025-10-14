<script setup lang="ts">
  import { StringConstants } from '@/constants';
  import { defineModel } from 'vue';
  
  const model = defineModel();
</script>

<template>
  <div class="input-group">
    <span class="input-group-text"><i class="bi bi-search"></i></span>
    <input type="text" class="form-control" v-model="model" :placeholder="StringConstants.Search">
  </div>
</template>