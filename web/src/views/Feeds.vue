<template>
  <div class="home">
    <div class="tablist">
      <label class="tab lefttab" v-if="!addRss"
             @click="addRss = true">{{ $t("feed.add") }}</label>
      <label class="tab lefttab" v-else style="margin: 5px;" @click="addRss = false">{{ $t("feed.cancel") }}</label>
      <label class="tab righttab"
      >{{ $t("feed.import") }}</label>
    </div>


    <div v-if="addRss" id="postinfo">
      <label>
        <input v-model="searchrss"/>
      </label>
      <b-button style="margin: 5px" @click="searchRss()">{{ $t("feed.search") }}</b-button>
      <br>
      <span>{{ info }}</span>
      <b-container>
        <b-row>
          <b-col v-for="(i, index) in search" :key="index" style="text-align: left">
            <b-card>
              <p>{{ i.title }}</p><span style="font-size: small">{{ i.link }}</span>
              <hr>
              <span>{{ $t("feed.subscriber") }} {{ i.subscriber }}</span>
              <b-icon-plus-square
                  style="margin: 5px;float: right"
                  @click="addSub(index)"
              ></b-icon-plus-square>
            </b-card>
          </b-col>
        </b-row>
      </b-container>
    </div>
    <label v-if="!addRss" style="margin: 5px;font-size: larger">{{ $t("feed.subscribed") }}</label>
    <div v-if="!addRss" id="list">
      <b-container v-for="(i, index) in rss" :key="index" style="text-align: left">
        <b-row class="post">
          <b-col><a
              style="font-size: small;color: rgba(0,0,0,.7);margin: 5px"
          >{{ i.unread }}
          </a>
            <a style="font-size: large">{{ i.title }} </a>
            <span style="font-size: small">{{ i.link }}</span>
          </b-col>
          <b-col cols="3">
            <b-icon-check-square-fill
                style="color: rgb(69,123,48);margin: 5px"
                v-if="i.unread === 0"
            >read
            </b-icon-check-square-fill
            >
            <b-icon-check-square
                v-else
                style="margin: 5px"
                @click="i.unread = 0;readFeed(index)"
            >unread
            </b-icon-check-square
            >
            <b-icon-x
                v-if="delRss && delRssIndex === index"
                style="margin: 5px"
                @click="delRss = false"
            ></b-icon-x>
            <b-icon-check
                v-if="delRss && delRssIndex === index"
                style="margin: 5px"
                @click="removeRss(index)"
            ></b-icon-check>
            <b-icon-trash
                v-else
                style="margin: 5px"
                @click="
              delRss = true;
              delRssIndex = index;
            "
            ></b-icon-trash>
          </b-col>
        </b-row>
      </b-container>
    </div>
    <div v-if="showLoading">
      <b-spinner class="loading" variant="success" label="Spinning"></b-spinner>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import config from "@/config";
import router from "@/router";

export default {
  name: "Overview",
  components: {},
  beforeMount() {
    if (window.localStorage.getItem("login") === "true") {
      this.$store.commit("setStatus", true)
      this.$store.commit("setjwt", window.localStorage.getItem("jwt"))
    }
    if (!this.$store.state.isLogin) {
      router.push("/login");
    }
    this.getData()
    this.getRss()
  },
  methods: {
    saveData: function () {
      window.localStorage.setItem("feeds", JSON.stringify(this.rss))
    },
    getData: function () {
      if (window.localStorage.getItem("feeds") === null) {
        window.localStorage.setItem("feeds", JSON.stringify([]))
      }
      this.rss = JSON.parse(window.localStorage.getItem("feeds"))
    },
    readFeed: function (index) {
      axios.post(config.apiAddress + "/api/feed/read/" + this.rss[index].id, null, {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        console.log(response.data)
      })
    },
    change: function (index) {
      this.post[index].read = !this.post[index].read;
    },
    addSub: function (index) {
      this.info = ""
      this.showLoading = true
      axios.post(config.apiAddress + "/api/feed/subscribe/" + this.search[index].id, null, {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        if (response.data.code === 200) {
          this.getRss()
          this.addRss = false
          this.searchrss = ""
          this.search = []
        } else {
          this.info = response.data.message
        }
        this.showLoading = false
      })
    },
    removeRss: function (index) {
      this.showLoading = true
      axios.post(config.apiAddress + "/api/feed/unsubscribe/" + this.rss[index].id, null, {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(() => {
        this.getRss()
      })
      this.showLoading = false
    },
    searchRss: function () {
      this.showLoading = true
      this.info = ""
      axios.post(config.apiAddress + "/api/feed/search", {
        Url: this.searchrss
      }, {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        if (response.data.code === 200) {
          let data = JSON.parse(response.data.message)
          this.search = []
          this.search.push({
            id: data.ID,
            title: data.Title,
            link: data.Url,
            subscriber: data.Subscriber
          })
          this.showLoading = false
        } else {
          this.info = response.data.message
          this.showLoading = false
        }
      })
    },
    getRss: function () {
      this.info = ""
      axios.get(config.apiAddress + "/api/feed/list", {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        if (response.data.code !== 200) {
          this.info = response.data.message
          return
        }
        this.rsslist = JSON.parse(response.data.message);
        axios.get(config.apiAddress + "/api/post/", {
          headers: {
            'Authorization': "Bearer " + this.$store.state.jwt,
            'Accept': 'application/json'
          }
        }).then(response => {
          if (response.data.code !== 200) {
            this.info = response.data.message
            return
          }
          this.postList = JSON.parse(response.data.message)
          axios.get(config.apiAddress + "/api/post/read", {
            headers: {
              'Authorization': "Bearer " + this.$store.state.jwt,
              'Accept': 'application/json'
            }
          }).then(response => {
            if (response.data.code !== 200) {
              this.info = response.data.message
              return
            }
            this.rss = []
            this.readPost = JSON.parse(response.data.message)
            this.rsslist.forEach(item => {
              let unread = 0
              this.postList.forEach(post => {
                if (post.Feed === item.ID && this.readPost.indexOf(post.ID) === -1) {
                  unread++
                }
              })

              this.rss.push({
                id: item.ID,
                title: item.Title,
                link: item.Url,
                unread: unread
              })
            })
            this.saveData()
            this.showLoading = false
          })
        })
      })
    }
  },
  data() {
    return {
      rss: [],
      info: "",
      addRss: false,
      delRss: false,
      delRssIndex: 0,
      rsslist: null,
      readPost: null,
      searchrss: "",
      unread: 0,
      search: [],
      showLoading: true,
      dismissSecs: 10,
      dismissCountDown: 0,
      showDismissibleAlert: false,
      postList: []
    };
  }
};
</script>
