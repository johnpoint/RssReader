<template>
  <div>
    <h1>Login</h1>
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
        <router-link style="font-size: small;float: right" to="/register"
        >register
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
  name: "Login",
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
          .post(config.apiAddress + "/api/login", {
            mail: this.form.email,
            password: this.form.password
          })
          .then(response => {
            if (response.data.code != 200) {
              alert(response.data.message);
            } else {
              this.$store.commit("setStatus", true);
              this.$store.commit("setjwt", response.data.message);
              window.localStorage.setItem("jwt", response.data.message)
              window.localStorage.setItem("login", true)
              router.push("/posts");
            }
          });
    }
  }
};
</script>

<style scoped></style>
