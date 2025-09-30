<script setup lang="ts">
  import Todo from '@/components/Todo.vue';
  import { StringConstants } from '@/constants';
</script>

<template>
  <h1>{{ StringConstants.Statistics }}</h1>
  <Todo/>
</template>
