<template>
  <div class="home">
    <check/>
    <div v-if="!addRss" id="list">
      <label class="tab" @click="addRss = true">Add</label>
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
    <div v-else id="postinfo">
      <label class="tab" style="margin: 5px;float: left" @click="addRss = false"
      >Back</label
      >
      <input v-model="searchrss"/>
      <b-button style="margin: 5px" @click="searchRss()">ok</b-button>
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
  </div>
</template>

<script>
import check from "@/components/check";
import axios from "axios";
import config from "@/config";

export default {
  name: "Overview",
  components: {
    check
  },
  beforeMount() {
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
      axios.post(config.apiAddress + "/api/feed/subscribe/" + this.search[index].id, null, {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        console.log(response.data)
        this.getRss()
        this.addRss = false
      })
    },
    removeRss: function (index) {
      axios.post(config.apiAddress + "/api/feed/unsubscribe/" + this.rss[index].id, null, {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        console.log(response.data)
        this.getRss()
      })
    },
    searchRss: function () {
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
        }
      })
    },
    getRss: function () {
      axios.get(config.apiAddress + "/api/feed/list", {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        this.rsslist = JSON.parse(response.data.message);
        axios.get(config.apiAddress + "/api/post/", {
          headers: {
            'Authorization': "Bearer " + this.$store.state.jwt,
            'Accept': 'application/json'
          }
        }).then(response => {
          this.postList = JSON.parse(response.data.message)
          axios.get(config.apiAddress + "/api/post/read", {
            headers: {
              'Authorization': "Bearer " + this.$store.state.jwt,
              'Accept': 'application/json'
            }
          }).then(response => {
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
          })
        })
      })
    }
  },
  data() {
    return {
      rss: [],
      addRss: false,
      delRss: false,
      delRssIndex: 0,
      rsslist: null,
      readPost: null,
      searchrss: "",
      unread: 0,
      search: []
    };
  }
};
</script>
