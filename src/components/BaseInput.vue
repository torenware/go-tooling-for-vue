<template>
  <div class="mt-3 px-3 d-flex flex-column" ref="controlBlock">
    <label :for=id class="form-label align-self-start">{{label}}</label>
    <input
      :id=id
      :type="inputType"
      class="form-control"
      :required="required != ''"
      @invalid="onInvalid"
    >
    <div class="errors text-danger mt-1  align-self-start d-none"></div>
  </div>
</template>

<script setup lang="ts" >
import { ref, Ref, onMounted } from 'vue';

// type ReportValueFunction = (elem: HTMLElement) => string | number;
const onInvalid = () => {
  console.log(`invalid called for ID=${props.id}`);
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
  // const errorBlock = controlBlock.value?.querySelector("div.errors");
  // if (!isValid) {
  //   const block = errorBlock as HTMLElement;
  //   block.classList.remove("d-none");
  //   block.textContent = input?.validationMessage as string;
  // } else {
  //   if (!errorBlock?.classList.contains("d-none")) {
  //     errorBlock?.classList.add("d-none");
  //   }
  // }
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

onMounted(() => {
  const input = controlBlock.value?.querySelector("input");
  input?.addEventListener("input", evt => {
    console.log(input.value);
    validateItem();
  });
  input?.addEventListener("blur", evt => {
    validateItem();
  });
});

</script>

<style scoped>
label {
  font-weight: bold;
  margin-left: 1rem;
}
</style>