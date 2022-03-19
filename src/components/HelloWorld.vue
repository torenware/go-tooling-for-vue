<script setup lang="ts">
import { ref, Ref, onMounted } from 'vue'
import BaseInput from './BaseInput.vue';

type FormControl = HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement;

const processSubmit = (form: HTMLFormElement) => {
  const values: {[key: string]: string | number} = {};

  // First check validity
  if (form.reportValidity()) {
    const elems = form.querySelectorAll("input, select, textarea");
    elems.forEach((elem: Element) => {
      values[elem.id] = (elem as FormControl).value
    });
  }
  console.log(values);
}

onMounted(() => {
  const ourForm = form.value as HTMLFormElement;
  const btn = ourForm.querySelector("button");
  btn?.addEventListener("click", () => {
    const form = btn.form as HTMLFormElement;
    processSubmit(form);
  });

  const reset = ourForm.querySelector("#reset-form") as HTMLFormElement;
  reset.addEventListener("click", () => {
    reset.form.reset();
    const elems = reset.form.querySelectorAll("input, select, textarea");
    elems.forEach((elem: Element) => {
      const parent = elem.parentElement;
      const block = parent?.querySelector("div.errors") as HTMLElement;
      if (!block.classList.contains("d-none")) {
        block.classList.add("d-none");
      }
      block.innerText = "";
    });
  });

});

const form: Ref<HTMLFormElement | null> = ref(null)
</script>

<template>
<div class="container">
  <form
    id="test-form"
    class="d-flex flex-column"
    novalidate
    @submit.prevent
    ref="form"
  >
    <BaseInput label="Email" id="test-id" inputType="email" required="true"/>
    <BaseInput label="Password" id="test-pw" inputType="password" required="true" />
    <div class="d-flex justify-content-start">
      <button id="submitter" class="btn btn-primary mt-3 ms-3">
        Submit Me!
      </button>
      <button id="reset-form" class="btn btn-secondary mt-3 ms-3 me-auto">Reset</button>
    </div>
  </form>
</div>
</template>

<style scoped>
a {
  color: #42b983;
}

label {
  margin: 0 0.5em;
  font-weight: bold;
}

code {
  background-color: #eee;
  padding: 2px 4px;
  border-radius: 4px;
  color: #304455;
}
</style>
