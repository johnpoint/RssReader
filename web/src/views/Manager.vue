<template>
  <div class="home">
    <b-spinner v-if="showLoading" style="position: fixed;top: 5px;right: 5px" variant="primary"
               label="Spinning"></b-spinner>
    <label class="tab" v-if="!addRss" style="margin: 5px;width: 100%;text-align: left"
           @click="addRss = true">Add</label>
    <label class="tab" v-else style="margin: 5px;width: 100%;text-align: left" @click="addRss = false">Cancel</label>
    <div v-if="addRss" id="postinfo">
      <input v-model="searchrss"/>
      <b-button style="margin: 5px" @click="searchRss()">ok</b-button>
      <br>
      <span>{{ info }}</span>
      <div v-for="(i, index) in search" :key="index" style="text-align: left">
        <div class="post">
          <a>{{ i.title }} </a>
          <b-icon-plus-square
              style="margin: 5px;float: right"
              @click="addSub(index)"
          ></b-icon-plus-square>
          <span style="font-size: small;margin: 5px;float: right"> {{ i.link }}</span>
        </div>
      </div>
    </div>
    <label style="margin: 5px;text-align: left;width: 100%;font-size: larger">Subscribed</label>
    <div id="list">
      <div v-for="(i, index) in rss" :key="index" style="text-align: left">
        <div class="post">
          <a style="font-size: large">{{ i.title }} </a>
          <span style="font-size: small">{{ i.link }}</span>
          <b-icon-x
              v-if="delRss && delRssIndex == index"
              style="float: right;margin: 5px"
              @click="delRss = false"
          ></b-icon-x>
          <b-icon-check
              v-if="delRss && delRssIndex == index"
              style="float: right;margin: 5px"
              @click="removeRss(index)"
          ></b-icon-check>
          <b-icon-trash
              v-else
              style="float: right;margin: 5px;"
              @click="
              delRss = true;
              delRssIndex = index;
            "
          ></b-icon-trash>
          <b-icon-check-square-fill
              style="float: right;margin: 5px;color: rgb(69,123,48)"
              v-if="i.unread == 0"
          >read
          </b-icon-check-square-fill
          >
          <b-icon-check-square
              style="float: right;margin: 5px"
              v-else
              @click="i.unread = 0;readFeed(index)"
          >unread
          </b-icon-check-square
          >
          <a
              style="font-size: small;color: rgba(0,0,0,.7);float: right;margin: 5px"
          >
            {{ i.unread }}
          </a>
        </div>
      </div>
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
    if (window.localStorage.getItem("login")) {
      this.$store.commit("setStatus", true)
      this.$store.commit("setjwt", window.localStorage.getItem("jwt"))
    }
    if (!this.$store.state.isLogin) {
      router.push("/");
    }
    this.getRss()
  },
  methods: {
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
      this.post[index].read = this.post[index].read ? false : true;
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
        console.log(response.data)
        if (response.data.code == 200) {
          this.getRss()
          this.addRss = false
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
      }).then(response => {
        console.log(response.data)
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
        if (response.data.code == 200) {
          let data = JSON.parse(response.data.message)
          this.search = []
          this.search.push({
            id: data.ID,
            title: data.Title,
            link: data.Url
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
        if (response.data.code != 200) {
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
          if (response.data.code != 200) {
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
            if (response.data.code != 200) {
              this.info = response.data.message
              return
            }
            this.rss = []
            this.readPost = JSON.parse(response.data.message)
            this.rsslist.forEach(item => {
              this.unread = 0
              this.postList.forEach(post => {
                if (post.Feed == item.ID && this.readPost.indexOf(post.id) == -1) {
                  this.unread++
                }
              })
              this.rss.push({
                id: item.ID,
                title: item.Title,
                link: item.Url,
                unread: this.unread
              })
            })
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
      showDismissibleAlert: false
    };
  }
};
</script>
