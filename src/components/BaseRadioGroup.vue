<template>
<fieldset>
  <BaseRadio
    v-for="radio in radios"
    :name="name"
    :value="radio.value"
    :label="radio.label"
  />
</fieldset>
</template>

<script setup lang="ts">
import BaseRadio from './BaseRadio.vue';

type RadioData = {
  label: string,
  value: string,
}

const props = defineProps<{
  id: string,
  name: string,
  radios: RadioData[],
  required?: string,
}>();
</script>

<style scoped>
  fieldset {
    padding: 1rem;
  }
</style>