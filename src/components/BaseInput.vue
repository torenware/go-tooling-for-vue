<template>
  <div  :class="groupClass" ref="controlBlock">
    <label
        v-if="inputType !== 'checkbox'"
        :for=id
        class="form-label align-self-start">{{label}}</label>

    <input
      v-if="inputType !== 'checkbox'"
      :id=id
      :type="inputType"
      class="form-control"
      :required="required != ''"
      @invalid="onInvalid"
    >
    <input
      v-else
      :id=id
      :type="inputType"
      class="form-check ms-3 pb-2"
      :required="required != ''"
      @invalid="onInvalid"
    >
    <label
        v-if="inputType === 'checkbox'"
        :for=id
        class="form-label checkbox pt-1 align-self-start">{{label}}</label>

    <div class="errors text-danger mt-1 ms-2 pb-1  align-self-start d-none"></div>
  </div>
</template>

<script setup lang="ts" >
import { ref, Ref, onMounted, computed } from 'vue';

const onInvalid = () => {
  // console.log(`invalid called for ID=${props.id}`);
  updateErrorBlock(false);
}


const props = withDefaults(defineProps<{
      label: string,
      id: string,
      required?: string,
      inputType?: string,
    }>(), {
  inputType: "text"
});

const validateItem = () => {
  const input = controlBlock.value?.querySelector("input");
  const isValid = input?.checkValidity();
  updateErrorBlock(isValid as boolean);
}

const updateErrorBlock = (valid: boolean) => {
  const input = controlBlock.value?.querySelector("input");
  const errorBlock = controlBlock.value?.querySelector("div.errors");

  if (!valid) {
    const block = errorBlock as HTMLElement;
    block.classList.remove("d-none");
    block.textContent = input?.validationMessage as string;
  } else {
    if (!errorBlock?.classList.contains("d-none")) {
      errorBlock?.classList.add("d-none");
    }
  }
}


const controlBlock: Ref<HTMLElement | null> = ref(null);

const groupClass = computed<string>(() => {
  if (props.inputType === "checkbox") {
    return "mt-3 d-flex flex-row justify-content-start";
  } else {
    return "mt-3 px-3 d-flex flex-column";
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
});

</script>

<style scoped>
label.checkbox {
  margin-left: .3rem;
}
label:not(.checkbox) {
  font-weight: bold;
  margin-left: 1rem;
}
</style>