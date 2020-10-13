<template>
  <div>
    <h1>{{ $t("auth.register") }}</h1>
    <span>{{ info }}</span>
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

        <b-button @click="onSubmit" variant="primary">OK</b-button>
        <router-link style="font-size: small;float: right" to="/login"
        >{{ $t("auth.login") }}
        </router-link
        >
      </b-form>
    </b-card>
    <div v-if="showLoading">
      <b-spinner class="loading" variant="success" label="Spinning"></b-spinner>
    </div>
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
      show: true,
      showLoading: false,
      info: "",
    };
  },
  methods: {
    onSubmit: function () {
      this.showLoading = true
      this.info = ""
      axios
          .post(config.apiAddress + "/api/register", {
            mail: this.form.email,
            password: this.form.password
          })
          .then(response => {
                this.showLoading = false
                this.info = response.data.message
                router.push("/login")
              },
              (error) => {
                let errText
                if (error.response == undefined) {
                  errText = "Unable to connect to server";
                } else {
                  errText = error.response.status + " " + error.response.data.message;
                }
                this.info = errText;
              }
          );
    }
  },
  beforeMount() {
    if (window.localStorage.getItem("login") == "true") {
      this.$store.commit("setStatus", true)
      this.$store.commit("setjwt", window.localStorage.getItem("jwt"))
      router.push("/posts")
    }
  }
};
</script>

<style scoped></style>
