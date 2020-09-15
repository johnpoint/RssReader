<template>
  <div class="home">
    <span>{{ info }}</span>
    <div v-if="!showPost" id="list">
      <label @click="getPostList()" class="tab backtab">更新</label>
      <label
          class="tab"
          :class="showUnread && !showRead ? 'select' : ''"
          @click="
          showUnread = true;
          showRead = false;
        "
      >Unread
        <b-badge pill variant="success">{{ unreadpost }}</b-badge>
      </label
      >
      |
      <label
          class="tab"
          :class="showRead && showUnread ? 'select' : ''"
          @click="
          showUnread = true;
          showRead = true;
        "
      >All</label
      >
      |
      <label
          class="tab"
          :class="showRead && !showUnread ? 'select' : ''"
          @click="
          showUnread = false;
          showRead = true;
        "
      >Read</label
      >

      <div v-for="(i, index) in post" :key="index" style="text-align: left">
        <div
            class="post"
            :class="i.read ? 'read' : 'unread'"
            v-if="
            (showRead && i.read && !showUnread) ||
              (showUnread && !i.read) ||
              (showRead && showUnread)
          "
        >
          <a style="font-size: small;color: rgba(0,0,0,.7)"
          >{{ i.source }} >>
          </a>
          <a
              style="font-size: large"
              @click="
              nowPost = index;
              showPost = true;
              i.read = false;
              getPostContent(index);
              change(index);
            "
              class="postlisttitle"
          >{{ i.title }}
          </a>
          <b-icon-check-square-fill
              style="float: right;margin: 5px;color: rgb(69,123,48)"
              v-if="i.read"
              @click="change(index)"
          >read
          </b-icon-check-square-fill
          >
          <b-icon-check-square
              style="float: right;margin: 5px"
              v-else
              @click="change(index)"
          >unread
          </b-icon-check-square
          >
          <a
              style="font-size: small;color: rgba(0,0,0,.7);float: right;margin: 5px"
              class="postdate"
          >
            {{ i.date }}
          </a>
        </div>
      </div>
    </div>
    <div v-if="showPost" id="postinfo">
      <div>
        <label
            class="tab backtab"
            @click="showPost = false;info=''"
        >Back</label
        >
        <label @click="change(nowPost)" class="readtab" v-if="post[nowPost].read">标为未读</label>
        <label @click="change(nowPost)" class="readtab" v-if="!post[nowPost].read">标为已读</label>
      </div>

      <h1 class="title">{{ post[nowPost].title }}</h1>
      <span>{{ post[nowPost].source }}</span> |
      <a :href="post[nowPost].link">Link</a>
      <b-card id="postcontext" style="margin: 15px" v-html="postContent">
      </b-card>
    </div>
    <div v-if="showLoading">
      <b-spinner class="loading" variant="success" label="Spinning"></b-spinner>
    </div>
  </div>
</template>

<script>
import axios from "axios"
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
    this.getPostList()
  },
  methods: {
    change: function (index) {
      this.post[index].read ? this.unread(index) : this.read(index);
      this.post[index].read = !this.post[index].read;
    },
    read: function (index) {
      this.info = ""
      axios.post(config.apiAddress + "/api/post/read/" + this.post[index].id, null, {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        if (response.data.code !== 200) {
          this.info = response.data.message
          return
        }
      }, err => {
        console.log(err)
        this.info = "请检查网络连接";
      })
    },
    unread: function (index) {
      this.info = ""
      axios.post(config.apiAddress + "/api/post/unread/" + this.post[index].id, null, {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        if (response.data.code !== 200) {
          this.info = response.data.message
          return
        }
      }, err => {
        console.log(err)
        this.info = "请检查网络连接";
      })
    },
    saveData: function () {
      window.localStorage.setItem("posts", JSON.stringify(this.post))
    },
    getData: function () {
      if (window.localStorage.getItem("posts") === null) {
        window.localStorage.setItem("posts", JSON.stringify([]))
      }
      this.post = JSON.parse(window.localStorage.getItem("posts"))
    },
    getReadList: function () {
      this.unreadpost = 0
      this.info = ""
      this.showLoading = true
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
        this.readPost = JSON.parse(response.data.message)
        this.post = []
        this.postList.forEach(item => {
          this.post.push({
            id: item.ID,
            title: item.Title,
            source: item.FeedTitle,
            date: new Date(item.Time).format("yyyy-MM-dd hh:mm:ss"),
            link: item.Link,
            read: this.readPost.indexOf(item.ID) !== -1
          })
          this.readPost.indexOf(item.ID) === -1 ? this.unreadpost++ : null
        })
        this.post.sort(function (a, b) {
          var x = a.date.toLowerCase();
          var y = b.date.toLowerCase();
          if (x < y) {
            return 1;
          }
          if (x > y) {
            return -1;
          }
          return 0;
        })
        this.saveData()
        this.showLoading = false
      }, err => {
        console.log(err)
        this.info = "请检查网络连接";
      })
    },
    getPostList: function () {
      this.info = ""
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
        this.getReadList()
      }, err => {
        console.log(err)
        this.info = "请检查网络连接";
      })
    },
    clearInfo: function () {
      this.info = ""
    },
    getPostContent: function (index) {
      if (window.localStorage.getItem("post" + this.post[index].id) !== null) {
        this.postContent = window.localStorage.getItem("post" + this.post[index].id)
        this.showLoading = false
        return
      }
      this.info = ""
      this.showLoading = true
      this.postContent = ""
      axios.get(config.apiAddress + "/api/post/content/" + this.post[index].id, {
        headers: {
          'Authorization': "Bearer " + this.$store.state.jwt,
          'Accept': 'application/json'
        }
      }).then(response => {
        if (response.data.code !== 200) {
          this.info = response.data.message
          return
        }
        this.postContent = response.data.message
        window.localStorage.setItem("post" + this.post[index].id, response.data.message)
        this.showLoading = false
      }, err => {
        console.log(err)
        this.info = "请检查网络连接";
      })
    }
  },
  data() {
    return {
      post: [],
      showPost: false,
      showRead: false,
      showUnread: true,
      nowPost: null,
      readPost: [],
      postList: null,
      postContent: "",
      showLoading: true,
      info: "",
      unreadpost: 0
    };
  }
};
</script>
