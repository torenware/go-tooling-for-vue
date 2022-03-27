<template>
  <div class="d-flex radio align-items-start mb-3" :class="invalid">
    <input :id="id" type="radio" :name="name" :value="value" />
    <label :for="id">{{ label }}</label>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  id: string,
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

label {
  margin-top: -0.4rem;
}

div.invalid label {
  color: #dc3545;
}
</style>