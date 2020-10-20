<template>
  <div class="home">
    <span>{{ info }}</span>
    <div id="list">
      <div class="tablist">
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
      </div>

      <label
          class="tab"
          :class="showUnread && !showRead ? 'select' : ''"
          @click="
          showUnread = true;
          showRead = false;
          nowData=post;
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
      >{{ $t("post.readafter") }}</label
      >
      <div v-for="(i,index) in cachePostData" :key="i.title" style="text-align: left">
        <div
            class="post"
            v-if="!showRead&&!showUnread&&i.readafter===true"
        >
          <a style="font-size: small; color: rgba(0, 0, 0, 0.7)"
          >{{ i.Source }} >>
          </a>
          <a
              style="font-size: large"
              @click="
              setTop();
              nowshowpost=cachePostData[index]
            "
              class="postlisttitle"
          >{{ i.Title }}
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
          <!--<b-icon-box-seam
              style="margin-right: 10px;"
              v-if="cachePostList.indexOf(String(post[index].id))!==-1"
          >saved
          </b-icon-box-seam>
          <b-icon-download
              style="margin-right: 10px;"
              v-else
              @click="getPostContent(post[index].id)"
          >save
          </b-icon-download>-->
          <a style="font-size: small; color: rgba(0, 0, 0, 0.7)"
          >{{ i.Source }} >>
          </a>
          <a
              style="font-size: large"
              @click="
              setTop();
              !post[index].read?change(index):'';
              toPost(post[index].ID);
          "
              class="postlisttitle"
          >{{ i.Title }}
          </a>
          <!--<b-icon-clock-history
              class="readbtn"
              style="float: right; margin: 5px;"
              v-if="!i.readafter"
              @click="readAfter(getPostContent(post[index].ID));nowshowpost=loadingShowPost">
          </b-icon-clock-history>-->
          <b-icon-check-square-fill
              class="readbtn"
              style="float: right; margin: 5px; color: #42b983"
              v-if="i.read"
              @click="change(index);"
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
    <div v-if="showLoading">
      <b-spinner class="loading" variant="success" label="Spinning"></b-spinner>
    </div>
    <span
        v-if="((showUnread && !showRead && unreadpost===0) || (!showRead&&!showUnread&&readafter.length===0)) && !showLoading"
        class="empty">{{ $t("post.empty") }}</span>
    <span v-else class="empty" style="color:rgba(0,0,0,0)">1</span>
  </div>
</template>

<script>
import axios from "axios";
import config from "@/config";
import router from "@/router";

export default {
  name: "Overview",
  components: {
  },
  beforeMount() {
    // check power
    if (window.localStorage.getItem("login") === "true") {
      this.$store.commit("setStatus", true);
      this.$store.commit("setjwt", window.localStorage.getItem("jwt"));
      if (window.localStorage.getItem("config") !== null) {
        this.$store.state.config = JSON.parse(window.localStorage.getItem("config"))
      } else {
        window.localStorage.setItem("config", JSON.stringify({"postnum": 50}))
      }
    }
    if (!this.$store.state.isLogin) {
      router.push("/login");
    }
    // get post list cache
    if (window.localStorage.getItem("posts") === null) {
      window.localStorage.setItem("posts", JSON.stringify([]));
    }
    this.post = JSON.parse(window.localStorage.getItem("posts"));
    this.$store.state.postList = JSON.parse(window.localStorage.getItem("posts"));
    // something
    this.getPostList();
    this.updateCache();
  },
  methods: {
    removeCache: function (index) {
      // console.log("removeCache")
      window.localStorage.removeItem("post" + this.cachePostData[index].id);
      this.updateCache();
    },
    updateCache: function () {
      let cachePost = Object.keys(window.localStorage);
      this.cachePostList = [];
      this.cachePostData = [];
      for (let i1 = 0; i1 < cachePost.length; i1++) {
        let i = cachePost[i1];
        if (i.indexOf("post") !== -1 && i.indexOf("posts") === -1) {
          this.cachePostList.push(i.replace("post", ""))
          let data = JSON.parse(window.localStorage.getItem(i))
          data.id = parseInt(i.replace("post", ""))
          this.cachePostData.push(data)
        }
      }
      this.readafter = []
      for (let i of this.cachePostData) {
        if (i.readafter) {
          this.readafter.push(i)
        }
      }
      this.readafter.sort(function (m, n) {
        return -(m["time"] - n["time"])
      })
    },
    setTop: function () {
      // console.log("setTop")
      this.top = document.documentElement.scrollTop;
      document.documentElement.scrollTop = 0;
    },
    backTop: function () {
      // console.log("backTop")
      document.documentElement.scrollTop = this.top;
    },
    toPost: function (id) {
      this.$router.push("/post/" + id)
    },
    change: function (index) {
      // console.log("change")
      this.post[index].read ? this.unread(this.post[index].ID) : this.read(this.post[index].ID);
      this.post[index].read = !this.post[index].read;
    },
    readAfter: function (data) {
      // console.log("readAfter")
      if (data.read !== true) {
        let j = 0;
        for (let i of this.post) {
          if (i.ID === data.ID) {
            this.change(j)
            break
          }
          j++;
        }
      }
      let now = JSON.parse(window.localStorage.getItem("post" + data.ID))
      now.readafter = true
      now.time = Date.parse(new Date())
      window.localStorage.setItem("post" + data.ID, JSON.stringify(now))
      this.showPost = false
      this.updateCache()
    },
    read: function (id) {
      // console.log("read")
      if (id === 0) {
        id = this.nowshowpost.id
        if (this.nowshowpost.read === true) {
          return
        }
      }
      this.info = "";
      axios
          .post(
              config.apiAddress + "/api/post/read/" + id,
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
                this.savePostListCache(JSON.stringify(this.post));
              },
              (error) => {
                let errText
                if (error.response === undefined) {
                  errText = "Unable to connect to server";
                } else {
                  errText = error.response.status + " " + error.response.data.message;
                }
                this.info = errText;
              }
          );
    },
    unread: function (id) {
      // console.log("unread")
      if (id === undefined) {
        id = this.nowshowpost.id
      }
      this.info = "";
      axios
          .post(
              config.apiAddress + "/api/post/unread/" + id,
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
                this.unreadpost += 1;
                this.savePostListCache(JSON.stringify(this.post));
              },
              (error) => {
                let errText
                if (error.response === undefined) {
                  errText = "Unable to connect to server";
                } else {
                  errText = error.response.status + " " + error.response.data.message;
                }
                this.info = errText;
              }
          );
    },
    savePostListCache: function (data) {
      // console.log("savePostListCache")
      window.localStorage.setItem("posts", data);
    },
    getReadList: function (postList) {
      // console.log("getReadList")
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
                postList.forEach((item) => {
                  this.post.push({
                    ID: item.ID,
                    Title: item.Title,
                    Source: item.FeedTitle,
                    date: new Date(parseInt(item.Time) * 1000).format(
                        "yyyy-MM-dd hh:mm:ss"
                    ),
                    link: item.Link,
                    read: this.readPost.indexOf(item.ID) !== -1,
                    readAfter: false
                  });
                  this.readPost.indexOf(item.ID) === -1 ? this.unreadpost++ : null;
                });
                this.savePostListCache(JSON.stringify(this.post));
                this.showLoading = false;
              },
              (error) => {
                let errText
                if (error.response === undefined) {
                  errText = "Unable to connect to server";
                } else {
                  errText = error.response.status + " " + error.response.data.message;
                }
                this.info = errText;
              }
          );
    },
    getPostList: function () {
      // console.log("getPostList")
      this.info = "";
      let postnum = ""
      if (this.$store.state.config.postnum === undefined) {
        postnum = ""
      } else {
        postnum = this.$store.state.config.postnum
      }
      axios
          .get(config.apiAddress + "/api/post/" + postnum, {
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
                let postList = JSON.parse(response.data.message);
                this.getReadList(postList);
              },
              (error) => {
                let errText
                if (error.response === undefined) {
                  errText = "Unable to connect to server";
                } else {
                  errText = error.response.status + " " + error.response.data.message;
                }
                this.info = errText;
              }
          );
    },
    getPostContent: function (id) {
      // console.log("getPostContent")
      if (window.localStorage.getItem("post" + id) !== null) {
        this.nowshowpost = JSON.parse(window.localStorage.getItem(
            "post" + id
        ));
        this.showLoading = false;
        return;
      }
      this.info = "";
      this.showLoading = true;
      axios
          .get(config.apiAddress + "/api/post/content/" + id, {
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
                let newPostCache = JSON.parse(response.data.message)
                for (let i of this.post) {
                  if (i.ID === newPostCache.ID) {
                    newPostCache.Source = i.Source
                    newPostCache.read = i.read
                    break
                  }
                }
                newPostCache.readafter = false
                window.localStorage.setItem(
                    "post" + id,
                    JSON.stringify(newPostCache)
                );
                this.showLoading = false;
                this.nowshowpost = newPostCache
                // console.log(newPostCache)
                this.updateCache();
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
    },
  },
  data() {
    return {
      post: [],//文章列表
      showRead: false,// tab控制开关
      showUnread: true,// tab控制开关
      readPost: [], // 已阅读文章列表
      showLoading: true,// 加载图标显示开关
      info: "",// 提示信息
      unreadpost: "-",// 未读文章计数
      top: 0,// 页面所处高度
      cachePostList: [],// 缓存文章key
      cachePostData: [],// 缓存文章数据
      readafter: [],// 稍后阅读列表
      empty: false,// 判断文章列表是否为空
      nowshowpost: null,// 显示的文章
      loadingShowPost: {
        Title: "Loading Title",
        Url: "/",
        Source: "Loading",
        Description: "Loading"
      },
    };
  },
};
</script>
