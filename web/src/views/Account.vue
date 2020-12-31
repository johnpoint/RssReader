<template>
  <div class="home">
    <span>{{ info }}</span>
    <b-container class="setting">
      <b-row>
        <b-col>
          <label>
            {{ $t("setting.CachedArticle") }}: {{ localpostnum }}
          </label>
        </b-col>
        <b-col align="end">
          <b-button size="sm" variant="outline-info"
                    @click="showAnalysis();Analysis=!Analysis">
            {{ Analysis ? $t("setting.hide") : $t("setting.info") }}
          </b-button>
          <b-button size="sm" variant="outline-primary" @click="clearCache">
            {{ $t("setting.clearCache") }}
          </b-button>
        </b-col>
      </b-row>
      <b-row v-if="Analysis">
        <b-col>
          <b-container>
            <b-row>
              <b-col cols="0">
                <a style="color: #42b983">#</a>
              </b-col>
              <b-col>
                {{ $t("setting.info") }}
              </b-col>
            </b-row>
          </b-container>
          <b-container v-for="i in cacheAnalysis" :key="i">
            <b-row>
              <b-col cols="0">
                <a style="color: #42b983">{{ i.num }}</a>
              </b-col>
              <b-col>
                {{ i.source }}
              </b-col>
            </b-row>
          </b-container>
        </b-col>
      </b-row>

      <!--<hr>
      <div class="setting">
        <span>{{ $t("setting.autodownload") }}</span>
        <b-button size="sm" style="float: right;margin: 5px" @click="autodownload=!autodownload" :variant='autodownload?"outline-primary":"outline-danger"'>
          {{ autodownload ? "ON" : "OFF" }}
        </b-button>
      </div>
      <div v-if="autodownload" class="setting">
        <span># {{ $t("setting.autodownloadset") }}</span>
      </div>-->
      <b-row>
        <b-col>
          <span>{{ $t("setting.postlistnum") }}</span>
        </b-col>
        <b-col align="end">
          <b-form-select style="width: 100px" v-model="loadpostnum" :options="loadpostnums"></b-form-select>
          <b-button size="sm" @click="savepostnum()" variant="outline-primary">{{
              $t("setting.save")
            }}
          </b-button>
        </b-col>
      </b-row>
      <b-row>
        <b-col>
          <span>{{ $t("setting.password") }}</span>
        </b-col>
        <b-col align="end">
          <input style="margin: 5px" v-if="editpassword&&!showpassword" v-model="newpassword"
                 type="password">
          <input style="margin: 5px" v-if="editpassword&&showpassword" v-model="newpassword" type="text">
          <b-button v-if="newpassword!==''" size="sm" @click="changePassword"
                    variant="outline-primary">{{
              $t("setting.save")
            }}
          </b-button>
          <b-button v-if="editpassword" size="sm" @click="showpassword=!showpassword"
                    variant="outline-primary">{{
              showpassword ? $t("setting.hidepassword") : $t("setting.showpassword")
            }}
          </b-button>
          <b-button size="sm" @click="editpassword=!editpassword;newpassword=''"
                    variant="outline-primary">{{
              editpassword ? $t("feed.cancel") : $t("setting.edit")
            }}
          </b-button>
        </b-col>
      </b-row>
      <b-row>
        <b-col>
          <span>语言 / Language</span>
        </b-col>
        <b-col align="end">
          <b-button size="sm" @click="changeLg()" variant="outline-primary">{{
              $i18n.locale == "zh" ? "EN" : "中文"
            }}
          </b-button>
        </b-col>
      </b-row>
      <!--<div class="setting">
        <span>{{ $t("setting.syncsetting") }}</span>
        <div style="float: right;">
          <span>{{ $t("setting.lastupdate") }} {{ new Date().format("yyyy-MM-dd") }}</span>
          <b-button size="sm" style="margin: 5px" variant="outline-info">
            {{ $t("setting.sync") }}
          </b-button>
        </div>
      </div>
      <hr>-->
      <b-row>
        <b-col align="end">
          <b-button size="sm" variant="outline-danger" @click="logout">
            {{ $t("setting.exit") }}
          </b-button>
        </b-col>
      </b-row>
    </b-container>
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
  name: "Account",
  components: {},
  data() {
    return {
      localpostnum: 0,
      localpost: [],
      showLoading: false,
      cache: [],
      cacheAnalysis: [],
      Analysis: false,
      autodownload: false,
      editpassword: false,
      loadpostnums: [
        {value: 50, text: '50'},
        {value: 75, text: '75'},
        {value: 100, text: '100'},
        {value: 125, text: '125'},
        {value: 150, text: '150'},
        {value: 200, text: '200'},
      ],
      loadpostnum: null,
      newpassword: "",
      showpassword: false,
      info: ""
    };
  },
  beforeMount() {
    if (window.localStorage.getItem("login") === "true") {
      this.$store.commit("setStatus", true);
      this.$store.commit("setjwt", window.localStorage.getItem("jwt"));
      if (window.localStorage.getItem("config") !== null) {
        this.$store.state.config = JSON.parse(window.localStorage.getItem("config"))
      } else {
        window.localStorage.setItem("config", JSON.stringify({"postnum": 50}))
      }
    } else {
      this.$store.commit("setStatus", false);
    }
    if (!this.$store.state.isLogin) {
      router.push("/login");
    }
    if (this.$store.state.config.postnum !== null) {
      this.loadpostnum = this.$store.state.config.postnum
    } else {
      this.loadpostnum = 50
      window.localStorage.setItem("config", JSON.stringify(this.$store.state.config))
    }
    this.getCache()
  },
  methods: {
    changePassword: function () {
      axios.post(config.apiAddress + "/api/user/password", {Password: this.newpassword}, {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(
          (response) => {
            this.info = response.data.message
            if (response.data.code === 200) {
              this.logout();
            }
          },
          (error) => {
            let errText
            if (error.response === undefined) {
              errText = "Unable to connect to server";
            } else {
              if (error.response.status === 401) {
                window.localStorage.setItem("login", false)
                router.push("/login")
              }
              errText = error.response.status + " " + error.response.data.message;
            }
            this.info = errText;
          }
      )
    },
    savepostnum: function () {
      let config = JSON.parse(window.localStorage.getItem("config"))
      config.postnum = this.loadpostnum
      window.localStorage.setItem("config", JSON.stringify(config))
      this.$store.state.config = JSON.parse(window.localStorage.getItem("config"))
    },
    showAnalysis: function () {
      this.cache = []
      for (let i of this.localpost) {
        let keys = Object.keys(this.cache)
        let csource = JSON.parse(window.localStorage.getItem(i)).Source
        if (keys.indexOf(csource) === -1) {
          this.cache[csource] = 1
        } else {
          this.cache[csource] += 1
        }
      }
      this.cacheAnalysis = []
      let keys = Object.keys(this.cache)
      for (let i of keys) {
        this.cacheAnalysis.push(
            {
              "source": i,
              "num": this.cache[i]
            }
        )
      }
      this.cacheAnalysis.sort(function (a, b) {
        return -(a["num"] - b["num"]);
      });
    },
    changeLg: function () {
      this.$i18n.locale = (this.$i18n.locale === "zh" ? "en" : "zh")
      window.localStorage.setItem("i18n", this.$i18n.locale)
    },
    logout: function () {
      this.$store.commit("setStatus", false);
      window.localStorage.removeItem("jwt")
      window.localStorage.setItem("login", false)
      window.localStorage.removeItem("posts")
      window.localStorage.removeItem("feeds")
      router.push("/");
    },
    getCache: function () {
      this.localpost = []
      this.localpostnum = 0
      Object.keys(window.localStorage).forEach(i => {
        if (i.indexOf("post") !== -1 && i.indexOf("posts") === -1) {
          this.localpost.push(i)
          this.localpostnum++
        }
      })
    },
    clearCache: function () {
      for (let i of this.localpost) {
        window.localStorage.removeItem(i)
      }
      this.getCache()
    }
  }
};
</script>