<template>
  <div>
    <Toast />
    <Dialog header="Registration" :visible="visible" :modal="true" @update:visible="onClose">
      <Form @submit="onSubmit" class="flex flex-column gap-2">
        <div>
          <label class="col-fixed" style="width:100px" for="name">Name</label>
          <InputText id="name" type="text" v-model.trim="registrateForm.name" autofocus></InputText>
        </div>
        <div>
          <label class="col-fixed" style="width:100px" for="surname">Surname</label>
          <InputText id="surname" type="text" v-model.trim="registrateForm.surname"></InputText>
        </div>
        <div>
          <label class="col-fixed" style="width:100px" for="age">Age</label>
          <InputNumber id="age" type="text" v-model.number="registrateForm.age" />
        </div>
        <div>
          <label class="col-fixed" style="width:100px" for="sex">Sex</label>
          <InputText id="sex" type="text" v-model.trim="registrateForm.sex"></InputText>
        </div>
        <div>
          <label class="col-fixed" style="width:100px" for="interests">Interests</label>
          <InputText id="interests" type="text" v-model.trim="registrateForm.interests"></InputText>
        </div>
        <div>
          <label class="col-fixed" style="width:100px" for="city">City</label>
          <InputText id="city" type="text" v-model.trim="registrateForm.city"></InputText>
        </div>
        <div>
          <label class="col-fixed" style="width:100px" for="email">Email</label>
          <Field id="email" name="email" type="email" v-model.trim="registrateForm.email"
            :class="['p-inputtext p-component']" :rules="validateEmail">
          </Field>
        </div>
        <ErrorMessage class="text-xs text-red-400" name="email" />
        <div>
          <label class="col-fixed" style="width:100px" for="password">Password</label>
          <Field id="password" name="password" type="password" v-model="registrateForm.password"
            :class="['p-inputtext p-component']" :rules="validatePassword">
          </Field>
        </div>
        <ErrorMessage class="text-xs text-red-400" name="password" />
        <div class="flex flex-row" style="justify-content: flex-end;">
          <Button label="Cancel" icon="pi pi-times" class="p-button-text" @click="onClose" />
          <Button label="Registrate" icon="pi pi-check" type="submit" />
        </div>
      </Form>
    </Dialog>
  </div>
</template>

<script>
import { Form, Field, ErrorMessage } from 'vee-validate';
import Toast from "primevue/toast";
import axios from 'axios';


export default {
  name: 'registration-form',
  components: {
    Toast,
    Form,
    Field,
    ErrorMessage,
  },
  data() {
    return {
      registrateForm: {},
    }
  },
  methods: {
    validateEmail(value) {
      if (!value || value === "") {
        return 'This field is required';
      }

      const regex = /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$/i;
      if (!regex.test(value)) {
        return 'This field must be a valid email';
      }

      return true;
    },
    validatePassword(value) {
      if (!value) {
        return 'This field is required';
      }

      return true;
    },
    onSubmit() {
      const path = `${process.env.VUE_APP_API_HOST || ''}/api/registrate`;
      axios.post(path, this.registrateForm)
        .then(() => {
          this.onClose()
        })
        .catch((err) => {
          // console.log(err);
          this.$toast.add({ severity: 'error', summary: 'Error', detail: err.message, life: 5000 });
        });
    },
    onClose() {
      this.registrateForm = {}
      this.$emit('update:visible', false)
    }
  },
  props: {
    visible: {
      type: Boolean,
      default: false
    },
  }
}
</script>

<style scoped></style>