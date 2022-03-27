<template>
  <div :class="groupClass" ref="controlBlock">
    <label v-if="inputType !== 'checkbox'" :for="id" class="form-label align-self-start">{{ label }}</label>

    <input
      v-if="inputType !== 'checkbox'"
      :id="id"
      :type="inputType"
      class="form-control"
      :required="required != ''"
      :value="value"
      @invalid="onInvalid"
    />
    <input
      v-else
      :id="id"
      :type="inputType"
      class="form-check ms-3 pb-2"
      :required="required != ''"
      :value="value"
      @invalid="onInvalid"
    />
    <label
      v-if="inputType === 'checkbox'"
      :for="id"
      class="form-label checkbox pt-1 align-self-start"
    >{{ label }}</label>

    <div class="errors text-danger mt-1 ms-2 pb-1 align-self-start d-none"></div>
  </div>
</template>

<script setup lang="ts" >
import { ref, Ref, onMounted, computed } from 'vue';

//const controlInvalid = ref(false);

const onInvalid = (evt: Event) => {
  // Stop the default browser UI from
  // popping up, since it ignores novalidate.
  evt.preventDefault();
  updateErrorBlock(false);
}

const onReset = () => {
  setInvalidClass(false);
  // controlInvalid.value = false;
}


const props = withDefaults(defineProps<{
  label: string,
  id: string,
  required?: string,
  inputType?: string,
  value?: string,
}>(), {
  inputType: "text"
});

const validateItem = () => {
  //controlInvalid.value = false;
  const input = controlBlock.value?.querySelector("input");
  const isValid = input?.checkValidity();
  updateErrorBlock(isValid as boolean);
}

const setInvalidClass = (isInvalid: boolean) => {
  const classOn = controlBlock.value?.classList.contains("invalid");
  if (isInvalid && !classOn) {
    controlBlock.value?.classList.add("invalid");
  } else if (!isInvalid && classOn) {
    controlBlock.value?.classList.remove("invalid");
  }
}

const updateErrorBlock = (valid: boolean) => {
  const input = controlBlock.value?.querySelector("input");
  const errorBlock = controlBlock.value?.querySelector("div.errors");

  if (!valid) {
    setInvalidClass(true);
    const block = errorBlock as HTMLElement;
    block.classList.remove("d-none");
    block.textContent = input?.validationMessage as string;
  } else {
    setInvalidClass(false);
    if (!errorBlock?.classList.contains("d-none")) {
      errorBlock?.classList.add("d-none");
    }
  }
}


const controlBlock: Ref<HTMLElement | null> = ref(null);

const groupClass = computed<string>(() => {
  // const invalid = controlInvalid.value ? "invalid" : "";
  const invalid = ""; // temp fix
  if (props.inputType === "checkbox") {
    return `${invalid} mt-3 d-flex flex-row justify-content-start`;
  } else {
    return `${invalid} mt-3 px-3 d-flex flex-column`;
  }
});

onMounted(() => {
  const input = controlBlock.value?.querySelector("input");
  input?.addEventListener("input", evt => {
    validateItem();
  });
  input?.addEventListener("blur", evt => {
    validateItem();
  });

  // Handle a standard CustomEvent of spReset, since
  // it appears the standard event handling ignores CustomEvents.
  input?.addEventListener("spReset", onReset);
});

</script>

<style scoped>
label.checkbox {
  margin-left: 0.3rem;
}
label:not(.checkbox) {
  font-weight: bold;
  margin-left: 1rem;
}

div.invalid label {
  color: #dc3545;
}

div.invalid input {
  background-color: rgba(230, 139, 166, 0.2);
}

div.errors {
  font-weight: bold;
  font-style: italic;
}
</style>