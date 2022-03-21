<template>
  <form
    class="d-flex flex-column"
    novalidate
    @submit.prevent
    ref="form"
  >
  <slot>

  </slot>

  <div class="d-flex justify-content-start">
      <button id="submitter" class="btn btn-primary mt-3 ms-3">
        Submit Me!
      </button>
      <button id="reset-form" class="btn btn-secondary mt-3 ms-3 me-auto">Reset</button>
    </div>
  </form>

</template>

<script setup lang="ts">
import { ref, Ref, onMounted } from 'vue'

type JSPO = {
  [key: string]: string | number
}

type ProcessSubmitFunc = (obj: JSPO) => void

const props = withDefaults(defineProps<{
  process?: ProcessSubmitFunc
    }>(), {
      process: (obj: JSPO) => {
        console.log("supply @prop process func taking a JS object to get results of form");
        console.log(obj);
      }
});


type FormControl = HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement;

const processSubmit = (form: HTMLFormElement) => {
  const values: JSPO = {};

  // First check validity
  if (form.reportValidity()) {
    const elems = form.querySelectorAll("input, select, textarea");
    elems.forEach((elem) => {
      if (elem.tagName === 'INPUT') {
        const input = elem as HTMLInputElement;
        if (input.type === "radio") {
          values[input.name] = input.value;
          return;
        }
      }
      // anything else
      if (elem.id) {
        values[elem.id] = (elem as FormControl).value
      }
    });
    props.process(values);
  }
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

    // make the button unselect itself
    reset.blur();

    const elems = reset.form.querySelectorAll("input, select, textarea");
    elems.forEach((elem: Element) => {
      const parent = elem.parentElement;
      const block = parent?.querySelector("div.errors") as HTMLElement;
      if (block) {
        // groups will not have their own error block.
        if (!block.classList.contains("d-none")) {
          block.classList.add("d-none");
        }
        block.innerText = "";
      }
    });
  });

});

const form: Ref<HTMLFormElement | null> = ref(null)

</script>

<style scoped>

label {
  margin: 0 0.5em;
  font-weight: bold;
}
</style>
