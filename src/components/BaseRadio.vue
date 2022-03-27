<template>
  <div class="d-flex radio" :class="invalid">
    <input type="radio" :name="name" :value="value" />
    <label>{{ label }}</label>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  label: string,
  value: string,
  name: string,
  invalid?: string,
}>();

</script>

<style scoped>
input[type="radio"] {
  margin-right: 0.3rem;
}

div.invalid label {
  color: red;
}
</style>