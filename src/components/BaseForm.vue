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

enum FieldSetState {
  REQUIRED_MISSING = -1
}

type JSPO = {
  [key: string]: string | number
}

type ProcessSubmitFunc = (obj: JSPO) => void;
type GatherValueFunc = (form: HTMLFormElement) => JSPO;

const props = withDefaults(defineProps<{
  process?: ProcessSubmitFunc
    }>(), {
      process: (obj: JSPO) => {
        console.log("supply @prop process func taking a JS object to get results of form");
        console.log(obj);
      }
});


type FormControl = HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement;

const validateRadioGroups = (form: HTMLFormElement) => {
  const values = checkRadioGroups(form);
  // We assume that our structure is similar for fieldsets as it is
  // for other control types.
  let valid = true; // optimism.

  for (let fsID of Object.keys(values)) {
    if (values[fsID] === FieldSetState.REQUIRED_MISSING) {
      const fs = document.getElementById(fsID);
      const errDiv = fs?.querySelector("div.errors");
      if (errDiv) {
        const block = errDiv as HTMLDivElement;
        if (block.classList.contains("d-none")) {
          block.classList.remove("d-none");
        }
        block.textContent = "At least one radio button must be selected."
        valid = false;
      }
    }
  }
  return valid;
}

const checkRadioGroups: GatherValueFunc = (form: HTMLFormElement) => {
  // We require that radio buttons be inside a fieldset.
  const fsets = form.querySelectorAll("fieldset");
  let values: JSPO = {};
  fsets.forEach(fs => {
    const isRequired = fs.getAttribute("required");
    const fsID: string | number = fs.id;
    values[fsID] = "" // assume no value

    if (isRequired && isRequired.length) {
      // ID is required
      if (!fsID) {
        throw "required fieldset *must* have an ID";
      }
      const inputs = fs.querySelectorAll("input[type=radio]");
      inputs.forEach(elem => {
        const radio = elem as HTMLInputElement;
        if (radio.checked) {
          values[fsID] = radio.value;
        }
      });
      // If the value for this ID is still unchanged,
      // we are empty where a value is required.
      if (values[fsID] === "") {
        values[fsID] = FieldSetState.REQUIRED_MISSING;
      }
    };
  });
  return values;
}

const processSubmit = (form: HTMLFormElement) => {
  const values: JSPO = {};

  // First check validity
  if (form.reportValidity() && validateRadioGroups(form)) {
    const elems = form.querySelectorAll("input, select, textarea");
    elems.forEach((elem) => {
      if (elem.tagName === 'INPUT') {
        const input = elem as HTMLInputElement;
        if (input.type === "radio" && input.checked) {
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

      let parent = elem.parentElement;
      // radio button groups are one level more nested.
      if (elem.tagName.toLowerCase() === "input") {
        const input = elem as HTMLInputElement;
        if (input.type.toLowerCase() === "radio") {
          parent = (parent as HTMLElement).parentElement;
        }
      }

      const block = parent?.querySelector("div.errors") as HTMLElement;
      if (block) {
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
