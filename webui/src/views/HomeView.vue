<script setup lang="ts">
  import { ref } from 'vue';
  import Todo from '@/components/Todo.vue';

  const count = ref(0);
</script>

<template>
  <h1>Main page</h1>
  <Todo/>
  <button class="btn btn-primary" @click="count++">Count: {{ count }}</button>
</template>
