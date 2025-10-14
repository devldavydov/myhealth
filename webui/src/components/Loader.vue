<script setup lang="ts">
  defineProps(['show']);
</script>

<template>
  <div class="spinner-border" role="status" v-if="show">
    <span class="visually-hidden">Loading...</span>
  </div>
</template>