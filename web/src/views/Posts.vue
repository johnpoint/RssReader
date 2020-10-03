<template>
  <div class="home">
    <span>{{ info }}</span>
    <div v-if="!showPost" id="list">
      <label @click="getPostList()" class="tab lefttab">{{
          $t("post.update")
        }}</label>
      <label
          @click="
          top = 0;
          backTop();
        "
          class="tab righttab"
      >{{ $t("post.totop") }}</label
      >
      <label
          class="tab"
          :class="showUnread && !showRead ? 'select' : ''"
          @click="
          showUnread = true;
          showRead = false;
        "
      >{{ $t("post.unread") }}
        <a style="color: #42b983">{{ unreadpost }}</a>
      </label>
      |
      <label
          class="tab"
          :class="showRead && showUnread ? 'select' : ''"
          @click="
          showUnread = true;
          showRead = true;
        "
      >{{ $t("post.all") }}</label
      >
      |
      <label
          class="tab"
          :class="showRead && !showUnread ? 'select' : ''"
          @click="
          showUnread = false;
          showRead = true;
        "
      >{{ $t("post.read") }}</label
      >
      |
      <label
          class="tab"
          :class="!showRead && !showUnread ? 'select' : ''"
          @click="
          showUnread = false;
          showRead = false;
        "
      >{{ $t("post.cache") }}</label
      >
      <div v-for="(i,index) in savePostData" :key="i.title" style="text-align: left">
        <div
            class="post"
            v-if="!showRead&&!showUnread"
        >
          <a style="font-size: small; color: rgba(0, 0, 0, 0.7)"
          >{{ i.source }} >>
          </a>
          <a
              style="font-size: large"
              @click="
              setTop();
              showPost = true;
              nowPost=index;
              nowData=savePostData;
              postContent=nowData[index].content
            "
              class="postlisttitle"
          >{{ i.title }}
          </a>
          <b-icon-trash
              class="readbtn"
              style="float: right; margin: 5px;"
              @click="removeCache(index)"
          >saved
          </b-icon-trash>
        </div>
      </div>

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
          <b-icon-box-seam
              style="margin-right: 10px;"
              v-if="savePost.indexOf(String(post[index].id))!==-1"
              @click="save(index)"
          >saved
          </b-icon-box-seam>
          <b-icon-download
              style="margin-right: 10px;"
              v-else
              @click="save(index)"
          >save
          </b-icon-download>
          <a style="font-size: small; color: rgba(0, 0, 0, 0.7)"
          >{{ i.source }} >>
          </a>
          <a
              style="font-size: large"
              @click="
              setTop();
              nowPost = index;
              showPost = true;
              i.read = false;
              nowData=savePostData=post;
              getPostContent(index);
              change(index);
            "
              class="postlisttitle"
          >{{ i.title }}
          </a>
          <b-icon-check-square-fill
              class="readbtn"
              style="float: right; margin: 5px; color: #42b983"
              v-if="i.read"
              @click="change(index)"
          >read
          </b-icon-check-square-fill>
          <b-icon-check-square
              class="readbtn"
              style="float: right; margin: 5px"
              v-else
              @click="change(index)"
          >unread
          </b-icon-check-square>
          <a class="postdate">
            {{ i.date }}
          </a>
        </div>
      </div>
    </div>
    <!--<label class="tab">{{ $t("post.prev") }}</label> | <label class="tab">{{ $t("post.next") }}</label>-->
    <div v-if="showPost" id="postinfo">
      <div>
        <label
            class="tab lefttab"
            @click="
            showPost = false;
            info = '';
            backTop();
          "
        >{{ $t("post.back") }}</label
        >
        <label
            @click="change(nowPost)"
            class="tab righttab"
            v-if="post[nowPost].read!==undefined&&post[nowPost].read"
        >{{ $t("post.setunread") }}</label
        >
        <label
            @click="change(nowPost)"
            class="tab righttab"
            v-if="post[nowPost].read!==undefined&&!post[nowPost].read"
        >{{ $t("post.setread") }}</label
        >
      </div>

      <h1 class="title">{{ nowData[nowPost].title }}</h1>
      <span>{{ nowData[nowPost].source }}</span> |
      <a :href="nowData[nowPost].link">{{ $t("post.link") }}</a>
      <b-card id="postcontext" style="margin: 15px" v-html="postContent">
      </b-card>
    </div>
    <div v-if="showLoading">
      <b-spinner class="loading" variant="success" label="Spinning"></b-spinner>
    </div>
    <span v-if="unreadpost === 0 && !showLoading" class="empty">{{ $t("post.empty") }}</span>
    <span v-else class="empty" style="color:rgba(0,0,0,0)">1</span>
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
      this.$store.commit("setStatus", true);
      this.$store.commit("setjwt", window.localStorage.getItem("jwt"));
    }
    if (!this.$store.state.isLogin) {
      router.push("/login");
    }
    this.getData();
    this.getPostList();
    this.updateCache();
  },
  methods: {
    removeCache: function (index) {
      window.localStorage.removeItem("post" + this.savePostData[index].id);
      this.updateCache();
    },
    updateCache: function () {
      let cachePost = Object.keys(window.localStorage);
      this.savePost = [];
      this.savePostData = [];
      for (let i of cachePost) {
        if (i.indexOf("post") !== -1 && i.indexOf("posts") === -1) {
          this.savePost.push(i.replace("post", ""))
          let data = JSON.parse(window.localStorage.getItem(i))
          data.id = parseInt(i.replace("post", ""))
          this.savePostData.push(data)
        }
      }
    },
    setTop: function () {
      this.top = document.documentElement.scrollTop;
      document.documentElement.scrollTop = 0;
    },
    backTop: function () {
      document.documentElement.scrollTop = this.top;
    },
    save: function (index) {
      this.getPostContent(index);
    },
    change: function (index) {
      this.post[index].read ? this.unread(index) : this.read(index);
      this.post[index].read = !this.post[index].read;
    },
    read: function (index) {
      this.info = "";
      axios
          .post(
              config.apiAddress + "/api/post/read/" + this.post[index].id,
              null,
              {
                headers: {
                  Authorization: "Bearer " + this.$store.state.jwt,
                  Accept: "application/json",
                },
              }
          )
          .then(
              (response) => {
                if (response.data.code !== 200) {
                  this.info = response.data.message;
                  return;
                }
                this.unreadpost -= 1;
              },
              (err) => {
                console.log(err);
                this.info = "请检查网络连接";
              }
          );
    },
    unread: function (index) {
      this.info = "";
      axios
          .post(
              config.apiAddress + "/api/post/unread/" + this.post[index].id,
              null,
              {
                headers: {
                  Authorization: "Bearer " + this.$store.state.jwt,
                  Accept: "application/json",
                },
              }
          )
          .then(
              (response) => {
                if (response.data.code !== 200) {
                  this.info = response.data.message;
                  return;
                }
                console.log(this.unreadpost);
                this.unreadpost += 1;
              },
              (err) => {
                console.log(err);
                this.info = "请检查网络连接";
              }
          );
    },
    saveData: function () {
      window.localStorage.setItem("posts", JSON.stringify(this.post));
    },
    getData: function () {
      if (window.localStorage.getItem("posts") === null) {
        window.localStorage.setItem("posts", JSON.stringify([]));
      }
      this.post = JSON.parse(window.localStorage.getItem("posts"));
    },
    getReadList: function () {
      this.unreadpost = 0;
      this.info = "";
      this.showLoading = true;
      axios
          .get(config.apiAddress + "/api/post/read", {
            headers: {
              Authorization: "Bearer " + this.$store.state.jwt,
              Accept: "application/json",
            },
          })
          .then(
              (response) => {
                if (response.data.code !== 200) {
                  this.info = response.data.message;
                  return;
                }
                this.readPost = JSON.parse(response.data.message);
                this.post = [];
                this.unreadpost = 0;
                this.postList.forEach((item) => {
                  this.post.push({
                    id: item.ID,
                    title: item.Title,
                    source: item.FeedTitle,
                    date: new Date(parseInt(item.Time) * 1000).format(
                        "yyyy-MM-dd hh:mm:ss"
                    ),
                    link: item.Link,
                    read: this.readPost.indexOf(item.ID) !== -1,
                  });
                  this.readPost.indexOf(item.ID) === -1 ? this.unreadpost++ : null;
                });
                this.saveData();
                this.showLoading = false;
              },
              (err) => {
                console.log(err);
                this.info = "请检查网络连接";
              }
          );
    },
    getPostList: function () {
      this.info = "";
      axios
          .get(config.apiAddress + "/api/post/", {
            headers: {
              Authorization: "Bearer " + this.$store.state.jwt,
              Accept: "application/json",
            },
          })
          .then(
              (response) => {
                if (response.data.code !== 200) {
                  this.info = response.data.message;
                  return;
                }
                this.postList = JSON.parse(response.data.message);
                this.getReadList();
              },
              (err) => {
                console.log(err);
                this.info = "请检查网络连接";
              }
          );
    },
    getPostContent: function (index) {
      if (window.localStorage.getItem("post" + this.post[index].id) !== null) {
        this.postContent = JSON.parse(window.localStorage.getItem(
            "post" + this.post[index].id
        )).content;
        this.showLoading = false;
        return;
      }
      this.info = "";
      this.showLoading = true;
      this.postContent = "";
      axios
          .get(config.apiAddress + "/api/post/content/" + this.post[index].id, {
            headers: {
              Authorization: "Bearer " + this.$store.state.jwt,
              Accept: "application/json",
            },
          })
          .then(
              (response) => {
                if (response.data.code !== 200) {
                  this.info = response.data.message;
                  return;
                }
                this.postContent = response.data.message;
                let newPostCache = {"title": this.post[index].title}
                newPostCache.content = this.postContent
                newPostCache.source = this.post[index].source
                newPostCache.link = this.post[index].link
                window.localStorage.setItem(
                    "post" + this.post[index].id,
                    JSON.stringify(newPostCache)
                );
                this.showLoading = false;
                this.updateCache();
              },
              (err) => {
                console.log(err);
                this.info = "请检查网络连接";
              }
          );
    },
  },
  data() {
    return {
      post: [],
      postListPage: 0,
      showPost: false,
      showRead: false,
      showUnread: true,
      nowPost: null,
      readPost: [],
      postList: null,
      postContent: "",
      showLoading: true,
      info: "",
      unreadpost: "-",
      top: 0,
      savePost: [],
      savePostData: [],
      nowData: null,
    };
  },
};
</script>
