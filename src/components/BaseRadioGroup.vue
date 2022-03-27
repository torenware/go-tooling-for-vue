<template>
  <fieldset
    :id="id"
    :required="required"
    class="d-flex flex-column justify-content-start"
    :class="extraClasses"
    ref="fsref"
  >
    <legend v-if="legend" class="mb-3">{{ legend }}</legend>
    <BaseRadio
      v-for="radio in radios"
      :id="id + '-' + radio.value"
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

const extraClasses = computed(() => {
  if (!fsref.value) {
    return "";
  }
  if (props.addedClasses && props.addedClasses.value.has(fsref.value.id)) {
    return props.addedClasses.value.get(fsref.value.id);
  }
});

const sendUserEvent = (form: HTMLFormElement) => {
  const evt = new CustomEvent("userChange");
  const rslt = form.dispatchEvent(evt);
};

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
        sendUserEvent(fs.form!);
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


</script>

<style scoped>
fieldset {
  padding: 1rem;
}

legend {
  font-size: medium;
  font-weight: bold;
  text-align: left;
}

fieldset.invalid label {
  color: red;
}

fieldset.invalid legend {
  color: #dc3545;
}

div.errors {
  font-weight: bold;
  font-style: italic;
}
</style>