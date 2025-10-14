<script setup lang="ts">
  import { StringConstants } from '@/constants';
  defineProps(["url"]);
</script>

<template>
  <a class="btn btn-primary" :href="url" role="button"><i class="bi-plus-square"></i> {{ StringConstants.Create }}</a>
</template>