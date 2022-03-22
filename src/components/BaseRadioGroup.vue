<template>
<fieldset :id="id" :required="required" class="d-flex flex-column" ref="fsref">
  <BaseRadio
    v-for="radio in radios"
    :name="name"
    :value="radio.value"
    :label="radio.label"
    :key="name + '-' + radio.value"
  />

  <div class="errors text-danger mt-1 ms-2 pb-1  align-self-start d-none"></div>

</fieldset>
</template>

<script setup lang="ts">
import {ref, Ref, onMounted} from 'vue';
import BaseRadio from './BaseRadio.vue';

const fsref: Ref<HTMLFieldSetElement | null> = ref(null);

onMounted(() => {
  if (!fsref.value) {
    return;
  }
  const fs = fsref.value;
  const radioElems = fs.querySelectorAll("input");
  radioElems.forEach(elem => {
    const radio = elem as HTMLInputElement;
    radio.addEventListener("change", () => {
      if (radio.checked) {
        const errBlock = fsref.value?.querySelector("div.errors");
        if (errBlock) {
          const block = errBlock as HTMLDivElement;
          if (!block.classList.contains("d-none")) {
            block.classList.add("d-none");
          }
        }
      }
    });
  })
});

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