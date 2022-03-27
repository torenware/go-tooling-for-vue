<template>
  <fieldset
    :id="id"
    :required="required"
    class="d-flex flex-column"
    :class="extraClasses"
    ref="fsref"
  >
    <legend v-if="legend">{{ legend }}</legend>
    <BaseRadio
      v-for="radio in radios"
      :name="name"
      :value="radio.value"
      :label="radio.label"
      :invalid="extraClasses"
      :key="name + '-' + radio.value"
    />

    <div class="errors text-danger mt-1 ms-2 pb-1 align-self-start d-none"></div>
  </fieldset>
</template>

<script setup lang="ts">
import { ref, Ref, onMounted, onUpdated, computed } from 'vue';
import BaseRadio from './BaseRadio.vue';

const fsref: Ref<HTMLFieldSetElement | null> = ref(null);

const extraClasses = computed(() => {
  if (!fsref.value) {
    return "";
  }
  if (props.addedClasses && props.addedClasses.value.has(fsref.value.id)) {
    return props.addedClasses.value.get(fsref.value.id);
  }
});


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
  legend?: string,
  radios: RadioData[],
  required?: string,
  addedClasses?: Ref<Map<string, string>>,
}>();


</script>

<style scoped>
fieldset {
  padding: 1rem;
}

legend {
  font-size: medium;
  font-weight: bold;
}

fieldset.invalid label {
  color: red;
}

fieldset.invalid legend {
  color: red;
}

div.errors {
  font-weight: bold;
  font-style: italic;
}
</style>