<script setup lang="ts">
  import { ref } from 'vue';
  import Todo from '@/components/Todo.vue';
  import { StringConstants } from '@/constants';

  const count = ref(0);
</script>

<template>
  <h1>{{ StringConstants.MainPage }}</h1>
  <Todo/>
  <button class="btn btn-primary" @click="count++">Count: {{ count }}</button>
</template>
