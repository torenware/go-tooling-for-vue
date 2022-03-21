<template>

  <div class="d-flex radio">
    <input
      type="radio"
      :name="name"
      :value="value"
    >
    <label>{{label}}</label>
  </div>

</template>

<script setup lang="ts">

const props = defineProps<{
  label: string,
  value: string,
  name: string,
}>();

</script>

<style scoped>
input[type=radio] {
  margin-right: .3rem;
}
</style>