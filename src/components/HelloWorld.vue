<script setup lang="ts">
import BaseInput from './BaseInput.vue';
import BaseForm from './BaseForm.vue';
import BaseRadioGroup from './BaseRadioGroup.vue';

const rg = [
  {
    label: "first of one",
    value: "wgbg"
  },
  {
    label: "el segundo",
    value: "wmex"
  }
];

type JSOB = {
  [key: string]: string | number
}
function processForm(obj: JSOB) {
  console.log("User:", obj.email);
  console.log("Password as well");
  console.log(obj);
}

</script>

<template>
  <div class="container">
    <BaseForm :process="processForm" submitText="Login" resetText="Reset">
      <template #default="fromForm">
        <BaseInput label="Email" id="email" inputType="email" required="true" />
        <BaseInput label="Password" id="password" inputType="password" required="true" />
        <BaseInput
          label="Agree To Licence"
          id="license"
          value="agreed-to-licence"
          inputType="checkbox"
        ></BaseInput>
        <BaseRadioGroup
          :radios="rg"
          id="lone-radio"
          legend="Radio Stations"
          name="station"
          required="true"
          :addedClasses="fromForm.addedClasses"
        />
      </template>
    </BaseForm>
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
