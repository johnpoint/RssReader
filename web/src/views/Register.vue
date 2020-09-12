<template>
  <div>
    <h1>Register</h1>
    <b-card style="width: 350px;margin: auto">
      <b-form v-if="show">
        <b-form-group
            id="input-group-1"
            label="Email:"
            label-align="left"
            label-for="input-1"
            description="We'll never share your email with anyone else."
        >
          <b-form-input
              id="input-1"
              v-model="form.email"
              type="email"
              required
              placeholder="Enter email"
          ></b-form-input>
        </b-form-group>

        <b-form-group
            id="input-group-2"
            label-align="left"
            label="Password:"
            label-for="input-2"
        >
          <b-form-input
              id="input-2"
              v-model="form.password"
              type="password"
              required
              placeholder="Enter password"
          ></b-form-input>
        </b-form-group>

        <b-button @click="onSubmit" variant="primary">Submit</b-button>
        <router-link style="font-size: small;float: right" to="/login"
        >login
        </router-link
        >
      </b-form>
    </b-card>
  </div>
</template>

<script>
import router from "@/router";
import axios from "axios";
import config from "@/config";

export default {
  name: "Register",
  data() {
    return {
      form: {
        email: "",
        password: ""
      },
      show: true
    };
  },
  methods: {
    onSubmit() {
      axios
          .post(config.apiAddress + "/api/register", {
            mail: this.form.email,
            password: this.form.password
          })
          .then(response => {
            alert(response.data.message);
            router.push("/login")
          });
    }
  }
};
</script>

<style scoped></style>
